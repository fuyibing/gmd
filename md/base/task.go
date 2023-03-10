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
// date: 2023-03-07

package base

import (
	"github.com/fuyibing/gmd/v8/app"
	"github.com/fuyibing/gmd/v8/app/models"
)

type (
	// Task
	// 订阅任务.
	Task struct {
		Id      int
		Title   string
		Updated int64

		// +-------------------------------------------------------------------+
		// + Subscription settings                                             |
		// +-------------------------------------------------------------------+

		Broadcasting bool
		Parallels    int
		Concurrency  int32
		MaxRetry     int
		DelaySeconds int

		// +-------------------------------------------------------------------+
		// + Registry settings                                                 |
		// +-------------------------------------------------------------------+

		RegistryId int
		TopicName  string
		TopicTag   string
		FilterTag  string

		// +-------------------------------------------------------------------+
		// + Calculation properties for subscriber & notification              |
		// +-------------------------------------------------------------------+

		Subscriber,
		SubscriberFailed,
		SubscriberSucceed Subscriber

		isSubscriber,
		isSubscriberFailed,
		isSubscriberSucceed bool

		isNotification,
		isNotificationFailed,
		isNotificationSucceed bool
	}
)

// IsNotification
// 本任务: 是否为通知主题.
func (o *Task) IsNotification() bool { return o.isNotification }

// IsNotificationFailed
// 本任务: 是否为失败通知.
func (o *Task) IsNotificationFailed() bool { return o.isNotificationFailed }

// IsNotificationSucceed
// 本任务: 是否为成功通知.
func (o *Task) IsNotificationSucceed() bool { return o.isNotificationSucceed }

// IsSubscriber
// 通用订阅.
//
// 订阅任务表 (scheme.task.handler_dispatcher_addr) 已配置回调地址.
func (o *Task) IsSubscriber() bool { return o.isSubscriber }

// IsSubscriberFailed
// 订阅投递失败通知.
//
// 订阅任务表 (scheme.task.failed_dispatcher_addr) 已配置回调地址.
func (o *Task) IsSubscriberFailed() bool { return o.isSubscriberFailed }

// IsSubscriberSucceed
// 订阅投递成功通知.
//
// 订阅任务表 (scheme.task.succeed_dispatcher_addr) 已配置回调地址.
func (o *Task) IsSubscriberSucceed() bool { return o.isSubscriberSucceed }

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *Task) bind(r *Registry) *Task {
	o.RegistryId = r.Id
	o.TopicName = r.TopicName
	o.TopicTag = r.TopicTag
	o.FilterTag = r.FilterTag
	return o
}

func (o *Task) init(m *models.Task) *Task {
	o.Id = m.Id
	o.Title = m.Title
	o.Updated = m.GmtUpdated.Time().Unix()
	o.Broadcasting = m.Broadcasting == models.StatusEnabled
	o.DelaySeconds = m.DelaySeconds

	if o.Parallels = m.Parallels; o.Parallels == 0 {
		o.Parallels = models.DefaultTaskParallels
	}
	if o.Concurrency = m.Concurrency; o.Concurrency == 0 {
		o.Concurrency = models.DefaultTaskConcurrency
	}
	if o.MaxRetry = m.MaxRetry; o.MaxRetry == 0 {
		o.MaxRetry = models.DefaultTaskMaxRetry
	}

	o.initSubscriberFailed(m)
	o.initSubscriberNormal(m)
	o.initSubscriberSucceed(m)

	o.initState()
	return o
}

func (o *Task) initState() {
	o.isNotification = o.TopicName == app.Config.GetProducer().GetNotifyTopic()
	o.isNotificationFailed = o.isNotification && o.TopicTag == app.Config.GetProducer().GetNotifyTagFailed()
	o.isNotificationSucceed = o.isNotification && o.TopicTag == app.Config.GetProducer().GetNotifyTagSucceed()
}

func (o *Task) initSubscriberFailed(m *models.Task) {
	if m.FailedDispatcherAddr != "" {
		x := NewSubscriber().
			SetCondition(m.FailedConditionKind, m.FailedConditionFilter).
			SetDispatcher(m.FailedDispatcherKind, m.FailedDispatcherAddr, m.FailedDispatcherMethod, m.FailedDispatcherTimeout).
			SetResult(m.FailedResultKind, m.FailedResultIgnoreCodes)
		if x.HasDispatcher() {
			o.isSubscriberFailed = true
			o.SubscriberFailed = x
		}
	}
}

func (o *Task) initSubscriberNormal(m *models.Task) {
	if m.HandlerDispatcherAddr != "" {
		x := NewSubscriber().
			SetCondition(m.HandlerConditionKind, m.HandlerConditionFilter).
			SetDispatcher(m.HandlerDispatcherKind, m.HandlerDispatcherAddr, m.HandlerDispatcherMethod, m.HandlerDispatcherTimeout).
			SetResult(m.HandlerResultKind, m.HandlerResultIgnoreCodes)
		if x.HasDispatcher() {
			o.isSubscriber = true
			o.Subscriber = x
		}
	}
}

func (o *Task) initSubscriberSucceed(m *models.Task) {
	if m.SucceedDispatcherAddr != "" {
		x := NewSubscriber().
			SetCondition(m.SucceedConditionKind, m.SucceedConditionFilter).
			SetDispatcher(m.SucceedDispatcherKind, m.SucceedDispatcherAddr, m.SucceedDispatcherMethod, m.SucceedDispatcherTimeout).
			SetResult(m.SucceedResultKind, m.SucceedResultIgnoreCodes)
		if x.HasDispatcher() {
			o.isSubscriberSucceed = true
			o.SubscriberSucceed = x
		}
	}
}
