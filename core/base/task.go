// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package base

import (
	"github.com/fuyibing/gmd/v8/app/models"
)

// Task
// 订阅任务.
type Task struct {
	Id      int
	Title   string
	Updated int64

	// Registry fields.

	RegistryId                     int
	TopicName, TopicTag, FilterTag string

	// Task fields.

	Broadcasting bool
	Concurrency  int32
	DelaySeconds int
	Parallels    int
	MaxRetry     int

	// Execution fields.

	BasicSubscriber   *Subscriber
	FailureNotify     bool
	FailureSubscriber *Subscriber
	SuccessNotify     bool
	SuccessSubscriber *Subscriber
}

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

	o.initFields(m)
	o.initSubscriber(m)
	o.initStatus()
	return o
}

func (o *Task) initFields(m *models.Task) {
	if o.Parallels = m.Parallels; o.Parallels == 0 {
		o.Parallels = models.DefaultParallels
	}

	if o.Concurrency = m.Concurrency; o.Concurrency == 0 {
		o.Concurrency = models.DefaultConcurrency
	}

	if o.MaxRetry = m.MaxRetry; o.MaxRetry == 0 {
		o.MaxRetry = models.DefaultMaxRetry
	}

	if u := m.GmtUpdated.Time().Unix(); u > 0 {
		o.Updated = u
	}
}

func (o *Task) initSubscriber(_ *models.Task) {}

func (o *Task) initStatus() {}
