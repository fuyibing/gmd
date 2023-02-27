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
// date: 2023-02-27

package task

import (
	"context"
	"fmt"
	"github.com/fuyibing/db/v5"
	"github.com/fuyibing/gmd/v8/app"
	"github.com/fuyibing/gmd/v8/app/models"
	"github.com/fuyibing/gmd/v8/app/services"
	"github.com/fuyibing/log/v5/cores"
	"github.com/fuyibing/util/v8/web/request"
	"github.com/fuyibing/util/v8/web/response"
	"github.com/kataras/iris/v12"
)

type (
	// Enable
	// 启动订阅任务.
	Enable struct {
		request  *EnableRequest
		response *EnableResponse
	}

	// EnableRequest
	// 启用入参.
	EnableRequest struct {
		Id int `json:"id" validate:"required,gte=1" mock:"1" label:"任务ID"`
	}

	// EnableResponse
	// 启用出参.
	EnableResponse struct {
		Id int `json:"id" validate:"required,gte=1" mock:"1" label:"任务ID"`
	}
)

func NewEnable() *Enable {
	return &Enable{
		request:  &EnableRequest{},
		response: &EnableResponse{},
	}
}

func (o *Enable) Run(s cores.Span, i iris.Context) (res interface{}) {
	var (
		code int
		err  error
	)

	// 校验入参.
	if i.ReadJSON(o.request) != nil {
		err = fmt.Errorf("invalid json payload")
		return response.With.ErrorCode(err, app.CodeInvalidPayloadFormatter)
	}

	// 校验入参.
	if err = request.Validate.Struct(o.request); err != nil {
		return response.With.ErrorCode(err, app.CodeInvalidPayloadFields)
	}

	s.GetAttr().Add("task-id", o.request.Id)
	if code, err = o.Send(s.GetContext()); err != nil {
		return response.With.ErrorCode(err, code)
	}
	return response.With.Data(o.response)
}

func (o *Enable) Send(ctx context.Context) (code int, err error) {
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

	// Send update service.
	if affects, err = service.SetStatusAsEnabled(bean.Id); err != nil {
		code = app.CodeServiceWriteError
		return
	}

	// Set response result.
	// o.response.Affects = affects
	o.response.Id = bean.Id
	// o.response.Title = bean.Title

	// Call consumer container reload access.
	if affects > 0 {
		// md.Boot.Consumer().Reload()
	}
	return
}
