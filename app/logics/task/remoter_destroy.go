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
	RemoteDestroy struct {
		request  *RemoteDestroyRequest
		response *RemoteDestroyResponse
	}

	RemoteDestroyRequest struct {
		Id int `json:"id" validate:"required,gte=1" mock:"1" label:"Task ID"`
	}

	RemoteDestroyResponse struct {
		Id    int    `json:"id" mock:"1" label:"Task ID"`
		Title string `json:"title" mock:"Task name" label:"Task name"`
	}
)

func NewRemoteDestroy() *RemoteDestroy {
	return &RemoteDestroy{
		request:  &RemoteDestroyRequest{},
		response: &RemoteDestroyResponse{},
	}
}

func (o *RemoteDestroy) Run(ctx context.Context, i iris.Context) (res interface{}) {
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
	if err = request.Validate.Struct(o.request); err != nil {
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

func (o *RemoteDestroy) Send(ctx context.Context) (code int, err error) {
	var (
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
	if err = md.Boot.Remoter().Adapter().DestroyById(ctx, bean.Id); err != nil {
		code = app.CodeAdapterError
		return
	}

	// Set response result.
	o.response.Id = bean.Id
	o.response.Title = bean.Title
	return
}
