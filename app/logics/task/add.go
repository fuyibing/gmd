// author: wsfuyibing <websearch@163.com>
// date: 2023-02-12

package task

import (
	"context"
	"fmt"
	"github.com/fuyibing/db/v3"
	"github.com/fuyibing/gmd/app"
	"github.com/fuyibing/gmd/app/models"
	"github.com/fuyibing/gmd/app/services"
	"github.com/fuyibing/log/v3"
	"github.com/fuyibing/log/v3/trace"
	"github.com/fuyibing/util/v2/web/request"
	"github.com/fuyibing/util/v2/web/response"
	"github.com/kataras/iris/v12"
)

type (
	Add struct {
		request  *AddRequest
		response *AddResponse
	}

	AddRequest struct {
		DelaySeconds int    `json:"delay_seconds" validate:"gte=0,lte=86400" mock:"0" label:"Delay seconds" desc:"When this configuration is greater than 0, the message sent by the producer needs to wait for the specified seconds before consumption. <br />Unit: Second.<br />Default: 0 (not delay)"`
		TopicName    string `json:"topic_name" validate:"required,gte=2,lte=30" mock:"orders" label:"Topic name"`
		TopicTag     string `json:"topic_tag" validate:"required,gte=2,lte=60" mock:"created" label:"Topic tag"`
		Handler      string `json:"handler" validate:"required,url" mock:"https://example.com/orders/expired/remove" label:"Callback address"`
		Title        string `json:"title" validate:"required,lte=80" mock:"Example task" label:"Task name"`
		Remark       string `json:"remark" mock:"Task remark" label:"Description about task"`
	}

	AddResponse struct {
		DelaySeconds int    `json:"delay_seconds" mock:"0" label:"Delay seconds"`
		Id           int    `json:"id" mock:"1" label:"Task id"`
		Title        string `json:"title" mock:"Example task" label:"Task name"`
		TopicName    string `json:"topic_name" mock:"orders" label:"Topic name"`
		TopicTag     string `json:"topic_tag" mock:"created" label:"Topic tag"`
	}
)

func NewAdd() *Add {
	return &Add{
		request:  &AddRequest{},
		response: &AddResponse{},
	}
}

func (o *Add) Run(ctx context.Context, i iris.Context) (res interface{}) {
	var (
		code int
		err  error
	)

	// Read payload json string
	// then assign to request fields.
	if i.ReadJSON(o.request) != nil {
		err = fmt.Errorf("invalid json payload")
		return response.With.ErrorCode(err, app.CodeInvalidPayloadFormat)
	}

	// Validate
	// requested payload params.
	if err = o.request.Validate(); err == nil {
		err = request.Validate.Struct(o.request)
	}
	if err != nil {
		return response.With.ErrorCode(err, app.CodeInvalidPayloadFields)
	}

	// Call send to do main process.
	log.Infofc(ctx, "logic send request: topic-name=%s, topic-tag=%s", o.request.TopicName, o.request.TopicTag)
	c := trace.Child(ctx)
	if code, err = o.Send(c); err != nil {
		return response.With.ErrorCode(err, code)
	}

	// Return succeed response.
	return response.With.Data(o.response)
}

func (o *Add) Send(ctx context.Context) (code int, err error) {
	var (
		br *models.Registry
		bt *models.Task
		ss = db.Connector.GetMasterWithContext(ctx)
		sr = services.NewRegistryService(ss)
		st = services.NewTaskService()
	)

	// Return error
	// if read registry failed.
	if br, err = sr.GetByNames(o.request.TopicName, o.request.TopicTag); err != nil {
		code = app.CodeServiceReadError
		return
	}

	// Create registry if not found.
	if br == nil {
		if br, err = sr.AddByNames(o.request.TopicName, o.request.TopicTag); br.Id == 0 {
			code = app.CodeServiceWriteError
			return
		}
	}

	// Fill filter tag
	// if undefined.
	if br.FilterTag == "" {
		if _, err = sr.SetFilterTag(br.Id); err != nil {
			code = app.CodeServiceWriteError
			return
		}
	}

	// Exists check.
	if bt, err = st.GetByHandler(br.Id, o.request.Handler); err != nil {
		code = app.CodeServiceReadError
		return
	}

	// Add task if not exists.
	if bt == nil {
		// Create task.
		if bt, err = st.Add(&models.Task{
			Title:        o.request.Title,
			Remark:       o.request.Remark,
			DelaySeconds: o.request.DelaySeconds,
			RegistryId:   br.Id,
			Handler:      o.request.Handler,
		}); err != nil {
			code = app.CodeServiceWriteError
			return
		}
	}

	// Set response result.
	o.response.Id = bt.Id
	o.response.Title = bt.Title
	o.response.DelaySeconds = bt.DelaySeconds
	o.response.TopicName = br.TopicName
	o.response.TopicTag = br.TopicTag
	return
}

// /////////////////////////////////////////////////////////////
// Basic edit request
// /////////////////////////////////////////////////////////////

func (o *AddRequest) Override() {}

func (o *AddRequest) Validate() error { return nil }
