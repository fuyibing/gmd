// author: wsfuyibing <websearch@163.com>
// date: 2023-02-12

package task

import (
	"context"
	"fmt"
	"github.com/fuyibing/db/v3"
	"github.com/fuyibing/gmd/app"
	"github.com/fuyibing/gmd/app/md"
	"github.com/fuyibing/gmd/app/models"
	"github.com/fuyibing/gmd/app/services"
	"github.com/fuyibing/log/v3"
	"github.com/fuyibing/log/v3/trace"
	"github.com/fuyibing/util/v2/web/request"
	"github.com/fuyibing/util/v2/web/response"
	"github.com/kataras/iris/v12"
)

type (
	Edit struct {
		request  *EditRequest
		response *EditResponse
	}

	EditRequest struct {
		Id           int     `json:"id" validate:"required,gte=1" mock:"1" label:"Task id"`
		DelaySeconds *int    `json:"delay_seconds" validate:"required,gte=0,lte=86400" mock:"0" label:"Delay seconds" desc:"When this configuration is greater than 0, the message sent by the producer needs to wait for the specified seconds before consumption. <br />Unit: Second.<br />Default: 0 (not delay)"`
		Parallels    *int    `json:"parallels" validate:"required,gte=0,lte=5" mock:"1" label:"Max consumers" desc:"Start consumers count per node. <br />Default: 1"`
		Concurrency  *int32  `json:"concurrency" validate:"required,gte=0" mock:"10" label:"Max concurrency" desc:"Max consuming message per consumer.<br />Default: 10.<br />Total: Nodes x Parallels * Concurrency.<br />Attention: If this value is set too large, the subscription service will be killed when there are too many messages in the queue (similar to DDOS)"`
		MaxRetry     *int    `json:"max_retry" validate:"required,gte=0" mock:"3" label:"Max consume times" desc:"Max consume times if failed returned.<br />Default: 3."`
		Broadcasting *int    `json:"broadcasting"  mock:"0" label:"Broadcast enabled" desc:"When enabled, all consumers of each deployment node will consume.<br />0: Disabled<br />1: Enabled"`
		Title        *string `json:"title" mock:"Example task" label:"Task name"`
		Remark       *string `json:"remark" mock:"Description about task" label:"Task remark"`
	}

	EditResponse struct {
		Affects int64  `json:"affects" mock:"1" label:"Updated count"`
		Id      int    `json:"id" mock:"1" label:"Task id"`
		Title   string `json:"title" mock:"Example task" label:"Task name"`
	}

	EditStatus struct {
		Id int `json:"id" validate:"required,gte=1" mock:"1" label:"任务ID"`
	}

	EditSubscriber struct {
		Id           int     `json:"id" validate:"required,gte=1" mock:"1" label:"Task id"`
		Handler      *string `json:"handler" mock:"http://example.com/path/route?key=value" label:"Callback address" desc:"Where is the message delivered.<br />Protocol: http, https, tcp, rpc, ws, wss."`
		Condition    *string `json:"condition" label:"Condition filter" desc:"Consume when the consumption content meets the filtering conditions, otherwise ignore the message."`
		IgnoreCodes  *string `json:"ignore_codes" mock:"1234,1234" label:"Ignore logic code" desc:"When the code returned by the business party is within the specified range, the consumption is considered successful.<br />Description: multiple codes are separated by commas"`
		Method       *string `json:"method" label:"Deliver method" desc:"Request method when delivering message. <br />Default: POST"`
		ResponseType *int    `json:"response_type" label:"Response type" desc:"How to identify the return results of business parties.<br />0: https status code is 200.<br />1: Return json string and errno field value is zero string or integer."`
		Timeout      *int    `json:"timeout" mock:"10" label:"Timeout" desc:"If response not returned within specified seconds."`
	}
)

func NewEdit() *Edit {
	return &Edit{
		request:  &EditRequest{},
		response: &EditResponse{},
	}
}

func (o *Edit) Run(ctx context.Context, i iris.Context) (res interface{}) {
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
	log.Infofc(ctx, "logic send request: task-id=%d", o.request.Id)
	c := trace.Child(ctx)
	if code, err = o.Send(c); err != nil {
		return response.With.ErrorCode(err, code)
	}

	// Return succeed response.
	return response.With.Data(o.response)
}

func (o *Edit) Send(ctx context.Context) (code int, err error) {
	var (
		affects int64
		bean    *models.Task
		sess    = db.Connector.GetMasterWithContext(ctx)
		service = services.NewTaskService(sess)
	)

	// Read task
	// bean from database.
	if bean, err = service.GetById(o.request.Id); err != nil {
		code = app.CodeServiceReadError
		return
	}

	// Return error
	// if task not found.
	if bean == nil {
		code = app.CodeServiceReadNotFound
		err = fmt.Errorf("task not found")
		return
	}

	// Prepare request param.
	o.request.Override(bean)
	req := &models.Task{
		Id:           bean.Id,
		Parallels:    *o.request.Parallels,
		Concurrency:  *o.request.Concurrency,
		MaxRetry:     *o.request.MaxRetry,
		DelaySeconds: *o.request.DelaySeconds,
		Broadcasting: *o.request.Broadcasting,
		Title:        *o.request.Title,
		Remark:       *o.request.Remark,
	}

	// Send update service.
	if affects, err = service.SetBasicFields(req); err != nil {
		code = app.CodeServiceWriteError
		return
	}

	// Set response result.
	o.response.Affects = affects
	o.response.Id = bean.Id
	o.response.Title = bean.Title

	// Call consumer container reload access.
	if affects > 0 && bean.IsEnabled() {
		md.Boot.Consumer().Reload()
	}

	return
}

// /////////////////////////////////////////////////////////////
// Basic edit request
// /////////////////////////////////////////////////////////////

func (o *EditRequest) Override(x *models.Task) {
	if o.Broadcasting != nil {
		o.Broadcasting = &x.Broadcasting
	}
	if o.Title != nil {
		o.Title = &x.Title
	}
	if o.Remark != nil {
		o.Remark = &x.Remark
	}
}

func (o *EditRequest) Validate() error { return nil }

// /////////////////////////////////////////////////////////////
// Subscriber edit request
// /////////////////////////////////////////////////////////////

func (o *EditSubscriber) OverrideFailed(x *models.Task) {
	if o.Handler == nil {
		o.Handler = &x.Failed
	}
	if o.Condition == nil {
		o.Condition = &x.FailedCondition
	}
	if o.Method == nil {
		o.Method = &x.FailedMethod
	}
	if o.Timeout == nil {
		o.Timeout = &x.FailedTimeout
	}
	if o.ResponseType == nil {
		o.ResponseType = &x.FailedResponseType
	}
	if o.IgnoreCodes == nil {
		o.IgnoreCodes = &x.FailedIgnoreCodes
	}
}

func (o *EditSubscriber) OverrideHandler(x *models.Task) {
	if o.Handler == nil {
		o.Handler = &x.Handler
	}
	if o.Condition == nil {
		o.Condition = &x.HandlerCondition
	}
	if o.Method == nil {
		o.Method = &x.HandlerMethod
	}
	if o.Timeout == nil {
		o.Timeout = &x.HandlerTimeout
	}
	if o.ResponseType == nil {
		o.ResponseType = &x.HandlerResponseType
	}
	if o.IgnoreCodes == nil {
		o.IgnoreCodes = &x.HandlerIgnoreCodes
	}
}

func (o *EditSubscriber) OverrideSucceed(x *models.Task) {
	if o.Handler == nil {
		o.Handler = &x.Succeed
	}
	if o.Condition == nil {
		o.Condition = &x.SucceedCondition
	}
	if o.Method == nil {
		o.Method = &x.SucceedMethod
	}
	if o.Timeout == nil {
		o.Timeout = &x.SucceedTimeout
	}
	if o.ResponseType == nil {
		o.ResponseType = &x.SucceedResponseType
	}
	if o.IgnoreCodes == nil {
		o.IgnoreCodes = &x.SucceedIgnoreCodes
	}
}

func (o *EditSubscriber) Validate() error { return nil }
