// author: wsfuyibing <websearch@163.com>
// date: 2023-02-08

package base

import (
	"github.com/fuyibing/gmd/app/md/conf"
	"github.com/fuyibing/gmd/app/models"
)

type (
	// Task
	// memory subscription task.
	Task struct {
		Id      int
		Title   string
		Updated int64

		Parallels    int
		Concurrency  int32
		MaxRetry     int
		DelaySeconds int
		Broadcasting bool

		RegistryId int
		TopicName  string
		TopicTag   string
		FilterTag  string

		HandlerSubscriber *Subscriber
		FailedSubscriber  *Subscriber
		SucceedSubscriber *Subscriber

		enNotificationFailed  bool
		enNotificationSucceed bool

		isNotification        bool
		isNotificationFailed  bool
		isNotificationSucceed bool
	}
)

func (o *Task) EnNotificationFailed() bool  { return o.enNotificationFailed }
func (o *Task) EnNotificationSucceed() bool { return o.enNotificationSucceed }

func (o *Task) IsNotification() bool        { return o.isNotification }
func (o *Task) IsNotificationFailed() bool  { return o.isNotificationFailed }
func (o *Task) IsNotificationSucceed() bool { return o.isNotificationSucceed }

// /////////////////////////////////////////////////////////////
// Access methods.
// /////////////////////////////////////////////////////////////

func (o *Task) bind(r *Registry) *Task {
	o.FilterTag = r.FilterTag
	o.RegistryId = r.Id
	o.TopicName = r.TopicName
	o.TopicTag = r.TopicTag
	return o
}

func (o *Task) init(m *models.Task) *Task {
	o.Id = m.Id
	o.Title = m.Title
	o.DelaySeconds = m.DelaySeconds
	o.Broadcasting = m.Broadcasting == models.StatusEnabled

	if o.Parallels = m.Parallels; o.Parallels == 0 {
		o.Parallels = conf.Config.Consumer.Parallels
	}
	if o.Concurrency = m.Concurrency; o.Concurrency == 0 {
		o.Concurrency = conf.Config.Consumer.Concurrency
	}
	if o.MaxRetry = m.MaxRetry; o.MaxRetry == 0 {
		o.MaxRetry = conf.Config.Consumer.MaxRetry
	}
	if n := m.GmtUpdated.Time().Unix(); n > 0 {
		o.Updated = n
	}

	o.initSubscriber(m)
	o.initStatus()
	return o
}

func (o *Task) initSubscriber(m *models.Task) {
	o.HandlerSubscriber = NewSubscriber(m, SubscriberTypeHandler)
	o.FailedSubscriber = NewSubscriber(m, SubscriberTypeFailed)
	o.SucceedSubscriber = NewSubscriber(m, SubscriberTypeSucceed)
}

func (o *Task) initStatus() {
	o.enNotificationFailed = o.FailedSubscriber != nil
	o.enNotificationSucceed = o.SucceedSubscriber != nil

	o.isNotification = o.TopicName == conf.Config.Producer.NotificationTopic
	o.isNotificationFailed = o.isNotification && o.TopicTag == conf.Config.Producer.NotificationTagFailed
	o.isNotificationSucceed = o.isNotification && o.TopicTag == conf.Config.Producer.NotificationTagSucceed
}
