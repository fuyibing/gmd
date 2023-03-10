// author: wsfuyibing <websearch@163.com>
// date: 2023-02-12

package task

import (
	"context"
	"fmt"
	"github.com/fuyibing/db/v8"
	"github.com/fuyibing/gmd/app"
	"github.com/fuyibing/gmd/app/md"
	"github.com/fuyibing/gmd/app/models"
	"github.com/fuyibing/gmd/app/services"
	"github.com/fuyibing/log/v8"
	"github.com/fuyibing/util/v8/web/request"
	"github.com/fuyibing/util/v8/web/response"
	"github.com/kataras/iris/v12"
)

type (
	EditSucceed struct {
		request  *EditSubscriber
		response *EditResponse
	}
)

func NewEditSucceed() *EditSucceed {
	return &EditSucceed{
		request:  &EditSubscriber{},
		response: &EditResponse{},
	}
}

func (o *EditSucceed) Run(ctx context.Context, i iris.Context) (res interface{}) {
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
	c := log.NewChild(ctx)
	if code, err = o.Send(c); err != nil {
		return response.With.ErrorCode(err, code)
	}

	// Return succeed response.
	return response.With.Data(o.response)
}

func (o *EditSucceed) Send(ctx context.Context) (code int, err error) {
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
	o.request.OverrideSucceed(bean)
	req := &models.Task{
		Id:                  bean.Id,
		Succeed:             *o.request.Handler,
		SucceedTimeout:      *o.request.Timeout,
		SucceedMethod:       *o.request.Method,
		SucceedCondition:    *o.request.Condition,
		SucceedResponseType: *o.request.ResponseType,
		SucceedIgnoreCodes:  *o.request.IgnoreCodes,
	}

	// Send update service.
	if affects, err = service.SetSubscriberForSucceed(req); err != nil {
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
