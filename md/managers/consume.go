// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// author: wsfuyibing <websearch@163.com>
// date: 2023-03-08

package managers

import (
	"fmt"
	"github.com/fuyibing/gmd/v8/app"
	"github.com/fuyibing/gmd/v8/md/base"
	"github.com/fuyibing/log/v5"
	"sync/atomic"
	"time"
)

type (
	// ConsumeExecutor
	// 消费管理.
	ConsumeExecutor interface {
		// Do
		// 消费过程.
		Do(task *base.Task, message *base.Message) (retry bool, err error)

		// Idle
		// 是否空闲.
		Idle() bool
	}

	consume struct {
		consuming, notifying, releasing int32
	}
)

func NewConsume() ConsumeExecutor { return (&consume{}).init() }

// +---------------------------------------------------------------------------+
// + Interface methods                                                         |
// +---------------------------------------------------------------------------+

func (o *consume) Do(task *base.Task, message *base.Message) (retry bool, err error) {
	atomic.AddInt32(&o.consuming, 1)

	var (
		body       []byte
		ignored    bool
		source     *base.Task
		subscriber base.Subscriber
	)

	// 结束消费.
	defer func() {
		atomic.AddInt32(&o.consuming, -1)

		// 结果处理.
		if retry = !ignored && err != nil && message.Dequeue < task.MaxRetry; retry {
			// 直接释放.
			o.doRelease(message)
		} else {
			// 结果通知并释放.
			o.doNotification(task, message)
		}
	}()

	// 校验出错.
	//
	// 1. 通用订阅 (handler_dispatcher_addr) 未定义.
	// 2. 结果通知 (<failed_dispatcher_addr>, <succeed_dispatcher_addr>) 未定义.
	// 3. 结果格式不合法.
	if source, subscriber, err = o.checkTask(task, message); err != nil {
		o.doIllegal(message, err)
		return
	}

	// 条件校验.
	//
	// 1. 规则不合法.
	// 2. 规则不匹配.
	if ignored, err = o.doCondition(subscriber, task, message); err != nil || ignored {
		return
	}

	// 投递过程.
	t := time.Now()
	body, err = o.doDispatch(subscriber, task, source, message)

	// 保存结果.
	message.SetDuration(time.Now().Sub(t)).SetError(err).SetResponseBody(body)

	// 校验结果.
	if err == nil {
		_, err = o.doResult(subscriber, task, source, message, body)
	} else {
		if body == nil {
			message.SetResponseBody([]byte(err.Error()))
		}
	}
	return
}

func (o *consume) Idle() bool {
	return atomic.LoadInt32(&o.consuming) == 0 &&
		atomic.LoadInt32(&o.releasing) == 0 &&
		atomic.LoadInt32(&o.notifying) == 0
}

// +---------------------------------------------------------------------------+
// + Access methods                                                            |
// +---------------------------------------------------------------------------+

// 检查任务.
func (o *consume) checkTask(task *base.Task, message *base.Message) (source *base.Task, subscriber base.Subscriber, err error) {
	message.TaskId = task.Id

	// 失败通知.
	// 消息投递失败时, 投递结果通知订阅方.
	if task.IsNotificationFailed() {
		v := base.Pool.AcquireNotification()
		defer v.Release()

		if source, err = v.Decoder(message); err == nil {
			message.DispatcherBody = v.MessageBody

			if subscriber = source.SubscriberFailed; subscriber == nil {
				err = fmt.Errorf("subscriber undefined for failed handler")
			}
		}
		return
	}

	// 成功通知.
	// 消息投递成功时, 投递结果通知订阅方.
	if task.IsNotificationSucceed() {
		v := base.Pool.AcquireNotification()
		defer v.Release()

		if source, err = v.Decoder(message); err == nil {
			message.DispatcherBody = v.MessageBody

			if subscriber = source.SubscriberSucceed; subscriber == nil {
				err = fmt.Errorf("subscriber undefined for succeed handler")
			}
		}
		return
	}

	// 通用订阅.
	source = task
	if subscriber = source.Subscriber; subscriber == nil {
		message.DispatcherBody = message.MessageBody
		err = fmt.Errorf("subscriber undefined for normal handler")
	}
	return
}

// 条件校验.
func (o *consume) doCondition(subscriber base.Subscriber, task *base.Task, message *base.Message) (ignored bool, err error) {
	if !subscriber.HasCondition() {
		return
	}
	return subscriber.GetCondition().Validate(task, message)
}

// 投递过程.
func (o *consume) doDispatch(subscriber base.Subscriber, task, source *base.Task, message *base.Message) (body []byte, err error) {
	// 未定义投递规则.
	if !subscriber.HasDispatcher() {
		span := log.NewSpanFromContext(message.GetContext(), "message.dispatch.undefined")
		defer span.End()

		err = fmt.Errorf("undefined in subscriber")
		span.Logger().Error("message dispatch: %v", err)
		return
	}

	// 投递过程.
	return subscriber.GetDispatcher().Dispatch(task, source, message)
}

// 无效任务.
func (o *consume) doIllegal(message *base.Message, err error) {
	span := log.NewSpanFromContext(message.GetContext(), "message.illegal")
	span.Logger().Error("message illegal: %v", err)
	span.End()
}

// 发送通知.
func (o *consume) doNotification(task *base.Task, message *base.Message) {
	atomic.AddInt32(&o.notifying, 1)
	go func() {
		var (
			exists       bool
			notification *base.Notification
			payload      *base.Payload
			registry     *base.Registry
			span         = log.NewSpanFromContext(message.GetContext(), "message.notification")
			topic        = app.Config.GetProducer().GetNotifyTopic()
			tag          string
		)

		// 监听结束.
		defer func() {
			// 释放通知.
			if notification != nil {
				notification.Release()
			}

			// 结束通知.
			span.End()
			atomic.AddInt32(&o.notifying, -1)

			// 释放消息.
			message.Release()
		}()

		// 通知类型.
		if message.GetError() != nil {
			// 失败通知.
			if task.IsSubscriberFailed() {
				tag = app.Config.GetProducer().GetNotifyTagFailed()
			}
		} else {
			// 成功通知.
			if task.IsSubscriberSucceed() {
				tag = app.Config.GetProducer().GetNotifyTagSucceed()
			}
		}
		if tag == "" {
			span.Logger().Info("notification not enabled")
			return
		}

		// 注册组合.
		span.Kv().Add("message.notification.topic.tag", tag).
			Add("message.notification.topic.name", topic)
		if registry, exists = base.Memory.GetRegistryByNames(topic, tag); !exists {
			span.Logger().Info("registry not found: topic-name=%s, topic-tag=%s", topic, tag)
			return
		}

		// 通知消息.
		notification = base.Pool.AcquireNotification()
		notification.MessageId = message.MessageId
		notification.MessageBody = message.GetResponseBody()
		notification.TaskId = task.Id

		// 消息正文.
		payload = base.Pool.AcquirePayload().SetContext(span.Context())
		payload.GenHash()
		payload.FilterTag = registry.FilterTag
		payload.MessageBody = notification.String()
		payload.RegistryId = registry.Id
		payload.TopicName = registry.TopicName
		payload.TopicTag = registry.TopicTag

		// 即时发布.
		if err := Boot.Producer().PublishSync(payload); err != nil {
			span.Logger().Error("notification send error: %v", err)
		}
	}()
}

// 释放回池.
func (o *consume) doRelease(message *base.Message) {
	atomic.AddInt32(&o.releasing, 1)
	go func() {
		defer atomic.AddInt32(&o.releasing, -1)
		message.Release()
	}()
}

// 校验结果.
func (o *consume) doResult(subscriber base.Subscriber, task, source *base.Task, message *base.Message, body []byte) (code int, err error) {
	if subscriber.HasResult() {
		code, err = subscriber.GetResult().Validate(task, source, message, body)
	}
	return
}

func (o *consume) init() *consume { return o }
