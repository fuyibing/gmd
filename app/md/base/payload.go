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
	// Payload
	// struct for message publish properties.
	Payload struct {
		c         context.Context
		duration  float64
		err       error
		ignored   bool
		messageId string

		FilterTag        string
		Hash             string
		Keyword          string
		MessageBody      string
		MessageMessageId string
		MessageTaskId    int
		Offset           int
		RegistryId       int
		TopicName        string
		TopicTag         string
	}
)

func (o *Payload) GetContext() context.Context           { return o.c }
func (o *Payload) GetError() error                       { return o.err }
func (o *Payload) GetIgnored() bool                      { return o.ignored }
func (o *Payload) GetMessageId() string                  { return o.messageId }
func (o *Payload) Release()                              { Pool.ReleasePayload(o) }
func (o *Payload) SetContext(c context.Context) *Payload { o.c = c; return o }
func (o *Payload) SetDuration(d float64) *Payload        { o.duration = d; return o }
func (o *Payload) SetError(e error) *Payload             { o.err = e; return o }
func (o *Payload) SetIgnored(b bool) *Payload            { o.ignored = b; return o }
func (o *Payload) SetMessageId(s string) *Payload        { o.messageId = s; return o }

// /////////////////////////////////////////////////////////////
// Access methods.
// /////////////////////////////////////////////////////////////

func (o *Payload) after() {
	// Call save
	// if enabled.
	if o.ignored {
		if *conf.Config.Producer.StorePublishIgnored {
			o.save()
		}
	} else {
		if o.err != nil {
			if *conf.Config.Producer.StorePublishFailed {
				o.save()
			}
		} else {
			if *conf.Config.Producer.StorePublishSucceed {
				o.save()
			}
		}
	}

	// Reset
	// access properties.
	o.c = nil
	o.duration = 0
	o.err = nil
	o.messageId = ""

	// Reset
	// data properties.
	o.MessageMessageId = ""
	o.MessageTaskId = 0
	o.Hash = ""
	o.Offset = 0
	o.RegistryId = 0
	o.TopicName = ""
	o.TopicTag = ""
	o.FilterTag = ""
	o.Keyword = ""
	o.MessageBody = ""
}

func (o *Payload) before() {
	o.ignored = false
}

func (o *Payload) init() *Payload {
	return o
}

func (o *Payload) save() {
	log.Infofc(o.c, "produced payload: store into database")

	var (
		affects int64
		bean    *models.Payload
		beanId  int64
		ctx     = log.NewChild(o.c)
		err     error
		sess    = db.Connector.GetMasterWithContext(ctx, models.ConnectionName)
		service = services.NewPayloadService(sess)
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
	// from payload table of database.
	if bean, err = service.GetByHash(o.Hash, o.Offset); err != nil {
		return
	}

	// Create
	// if history not found.
	if bean == nil {
		req := &models.Payload{
			Duration:         o.duration,
			MessageTaskId:    o.MessageTaskId,
			MessageMessageId: o.MessageMessageId,
			Hash:             o.Hash,
			Offset:           o.Offset,
			RegistryId:       o.RegistryId,
			MessageId:        o.messageId,
			MessageBody:      o.MessageBody,
		}

		// Assign response body.
		if o.err != nil {
			req.ResponseBody = o.err.Error()
		}

		// Add record.
		if o.ignored {
			bean, err = service.AddIgnored(req)
		} else {
			if o.err != nil {
				if conf.Config.Producer.MaxRetry > 1 {
					bean, err = service.AddWaiting(req)
				} else {
					bean, err = service.AddFailed(req)
				}
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
			if (bean.Retry + 1) < conf.Config.Producer.MaxRetry {
				affects, err = service.SetStatusAsWaiting(bean.Id, o.duration, o.err.Error())
			} else {
				affects, err = service.SetStatusAsFailed(bean.Id, o.duration, o.err.Error())
			}
		} else {
			affects, err = service.SetStatusAsSucceed(bean.Id, o.duration, o.messageId)
		}
	}
}
