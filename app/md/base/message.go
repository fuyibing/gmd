// author: wsfuyibing <websearch@163.com>
// date: 2023-02-08

package base

import (
	"context"
	"fmt"
	"github.com/fuyibing/db/v8"
	"github.com/fuyibing/gmd/app/md/conf"
	"github.com/fuyibing/gmd/app/models"
	"github.com/fuyibing/gmd/app/services"
	"github.com/fuyibing/log/v8"
)

type (
	// Message
	// struct for message consumed properties.
	Message struct {
		body     []byte
		c        context.Context
		duration float64
		err      error
		ignored  bool

		Dequeue          int
		Keyword          string
		MessageBody      string
		MessageId        string
		MessageTime      int64
		PayloadMessageId string
		TaskId           int
	}
)

func (o *Message) GetContext() context.Context           { return o.c }
func (o *Message) GetError() error                       { return o.err }
func (o *Message) GetIgnored() bool                      { return o.ignored }
func (o *Message) Release()                              { Pool.ReleaseMessage(o) }
func (o *Message) SetBody(b []byte) *Message             { o.body = b; return o }
func (o *Message) SetContext(c context.Context) *Message { o.c = c; return o }
func (o *Message) SetDuration(d float64) *Message        { o.duration = d; return o }
func (o *Message) SetError(e error) *Message             { o.err = e; return o }
func (o *Message) SetIgnored(i bool) *Message            { o.ignored = i; return o }

// /////////////////////////////////////////////////////////////
// Access methods.
// /////////////////////////////////////////////////////////////

func (o *Message) after() {
	// Call save
	// if enabled.
	if o.ignored {
		if *conf.Config.Consumer.StoreDispatchIgnored {
			o.save()
		}
	} else {
		if o.err != nil {
			if *conf.Config.Consumer.StoreDispatchFailed {
				o.save()
			}
		} else {
			if *conf.Config.Consumer.StoreDispatchSucceed {
				o.save()
			}
		}
	}

	// Reset
	// access properties.
	o.body = nil
	o.c = nil
	o.duration = 0
	o.err = nil

	// Reset
	// data properties.
	o.Dequeue = 0
	o.Keyword = ""
	o.MessageBody = ""
	o.MessageId = ""
	o.MessageTime = 0
	o.PayloadMessageId = ""
	o.TaskId = 0
}

func (o *Message) before() {
	o.ignored = false
}

func (o *Message) init() *Message {
	return o
}

func (o *Message) save() {
	log.Infofc(o.c, "consumed message: store into database")

	var (
		affects int64
		bean    *models.Message
		beanId  int64
		ctx     = log.NewChild(o.c)
		err     error
		sess    = db.Connector.GetMasterWithContext(ctx, models.ConnectionName)
		service = services.NewMessageService(sess)
	)

	// Called when end.
	defer func() {
		// Catch panic.
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}

		// Assign bean id.
		if bean != nil && beanId == 0 {
			beanId = bean.Id
		}

		// Logger dispatcher result.
		if err != nil {
			log.Errorfc(ctx, "store error, bean-id=%d, affects=%d, %v", beanId, affects, err)
		} else {
			log.Infofc(ctx, "store finish, bean-id=%d, affects=%d", beanId, affects)
		}
	}()

	// Get history record
	// from message table of database.
	if bean, err = service.GetByMessageId(o.TaskId, o.MessageId); err != nil {
		return
	}

	// Create
	// if history not found.
	if bean == nil {
		req := &models.Message{
			Duration:         o.duration,
			TaskId:           o.TaskId,
			PayloadMessageId: o.PayloadMessageId,
			MessageDequeue:   o.Dequeue,
			MessageTime:      o.MessageTime,
			MessageId:        o.MessageId,
			MessageBody:      o.MessageBody,
			ResponseBody:     string(o.body),
		}

		// Add record.
		if o.ignored {
			bean, err = service.AddIgnored(req)
		} else {
			if o.err != nil {
				bean, err = service.AddFailed(req)
			} else {
				bean, err = service.AddSucceed(req)
			}
		}

		// Completed.
		if bean != nil {
			affects = 1
		}
		return
	}

	// Update status
	// if saved already.
	if o.ignored {
		affects, err = service.SetStatusAsIgnored(bean.Id)
	} else {
		if o.err != nil {
			affects, err = service.SetStatusAsFailed(bean.Id, o.duration, string(o.body))
		} else {
			affects, err = service.SetStatusAsSucceed(bean.Id, o.duration, string(o.body))
		}
	}
}
