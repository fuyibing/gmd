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
		span       = log.NewSpanFromContext(message.GetContext(), "message.consume")
	)

	span.Kv().Add("consume.message.id", message.MessageId)
	message.TaskId = task.Id

	// 结束消费.
	// 消息消费结束后发送结果通知/释放实例回池.
	defer func() {
		span.End()
		atomic.AddInt32(&o.consuming, -1)

		// 后续处理.
		if retry = !ignored && err != nil && message.Dequeue < task.MaxRetry; retry {
			o.release(message)
		} else {
			o.notify(task, message, o.release)
		}
	}()

	// 检查订阅.
	if source, subscriber, message.DispatcherBody, err = o.check(task, message); err != nil {
		span.Logger().Error("consume check error: %v", err)
		return
	}

	// 条件校验.
	if subscriber.HasCondition() {
		ignored, err = subscriber.GetCondition().Validate(span.Context(), task, source, message)

		// 条件错误.
		if err != nil {
			span.Logger().Error("condition error parse: %v", err)
			return
		}

		// 忽略条件.
		if ignored {
			span.Logger().Error("condition parse ignored")
			return
		}
	}

	// 投递过程.
	if subscriber.HasDispatcher() {
		t := time.Now()
		body, err = subscriber.GetDispatcher().Dispatch(span.Context(), task, source, message)
		message.SetDuration(time.Now().Sub(t)).SetError(err).SetResponseBody(body)

		// 投递出错.
		if err != nil {
			if body == nil {
				message.SetResponseBody([]byte(err.Error()))
			}

			span.Logger().Error("dispatch error: %v", err)
			return
		}
	}

	// 结果校验.
	if subscriber.HasResult() {
		_, err = subscriber.GetResult().Validate(span.Context(), task, source, body)
	}
	return
}

func (o *consume) Idle() bool {
	return atomic.LoadInt32(&o.consuming) == 0 &&
		atomic.LoadInt32(&o.releasing) == 0 &&
		atomic.LoadInt32(&o.notifying) == 0
}

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *consume) check(task *base.Task, message *base.Message) (source *base.Task, subscriber base.Subscriber, dispatchBody string, err error) {
	defer func() {
		if err == nil && subscriber == nil {
			err = fmt.Errorf("subscriber not configured")
		}
	}()

	// 失败通知.
	if task.IsNotificationFailed() {
		v := base.Pool.AcquireNotification()
		defer v.Release()

		if source, err = v.Decoder(message); err == nil {
			dispatchBody = v.MessageBody
			subscriber = source.SubscriberFailed
		}
		return
	}

	// 成功通知.
	if task.IsNotificationSucceed() {
		v := base.Pool.AcquireNotification()
		defer v.Release()

		if source, err = v.Decoder(message); err == nil {
			dispatchBody = v.MessageBody
			subscriber = source.SubscriberSucceed
		}
		return
	}

	// 通用订阅.
	dispatchBody = message.MessageBody
	subscriber = task.Subscriber
	return
}

func (o *consume) init() *consume {
	return o
}

// 结果通知.
func (o *consume) notify(task *base.Task, message *base.Message, releaser func(*base.Message)) {
	atomic.AddInt32(&o.notifying, 1)
	go func() {
		var (
			ntf *base.Notification
			tag = ""
		)

		// 结束通知.
		defer func() {
			if ntf != nil {
				ntf.Release()
			}

			atomic.AddInt32(&o.notifying, -1)
			releaser(message)
		}()

		// 结果通知.
		if message.GetError() != nil {
			if task.IsNotificationFailed() {
				tag = app.Config.GetProducer().GetNotifyTagFailed()
			} else {
				return
			}
		} else {
			// 成功通知.
			if task.IsNotificationSucceed() {
				tag = app.Config.GetProducer().GetNotifyTagSucceed()
			} else {
				return
			}
		}

		ntf = base.Pool.AcquireNotification()
		ntf.MessageBody = message.GetResponseBody()
		ntf.MessageId = message.MessageId
		ntf.TaskId = task.Id

		// 消息结构.
		p := base.Pool.AcquirePayload().SetContext(message.GetContext())
		p.GenHash()

		p.TopicName = app.Config.GetProducer().GetNotifyTopic()
		p.TopicTag = tag
		p.FilterTag = tag
		p.MessageBody = ntf.String()

		// todo: 调用通知发布
	}()
}

// 释放消息.
func (o *consume) release(message *base.Message) {
	atomic.AddInt32(&o.releasing, 1)
	go func() {
		defer atomic.AddInt32(&o.releasing, -1)
		message.Release()
	}()
}
