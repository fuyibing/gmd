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
// date: 2023-03-09

package index

import (
	"github.com/fuyibing/gmd/v8/app"
	"github.com/fuyibing/gmd/v8/app/models"
	"github.com/fuyibing/log/v5/tracers"
	"github.com/fuyibing/util/v8/web/response"
	"github.com/kataras/iris/v12"
)

type (
	// GetHome
	// 首页逻辑.
	GetHome struct {
		response *GetHomeResponse
	}

	// GetHomeResponse
	// 首页返回值.
	GetHomeResponse struct {
		Service string `json:"service" label:"服务名" mock:"gmd"`
		Started string `json:"started" label:"启动时间" mock:"2023-03-01 09:10:11"`
		Version string `json:"version" label:"版本号" mock:"1.2.3"`
	}
)

func NewGetHome() *GetHome {
	return &GetHome{
		response: &GetHomeResponse{},
	}
}

func (o *GetHome) Run(_ tracers.Span, _ iris.Context) interface{} {
	o.response.Service = app.Config.GetName()
	o.response.Started = app.Config.GetStartTime().Format(models.DefaultDatetimeLayout)
	o.response.Version = app.Config.GetVersion()
	return response.With.Data(o.response)
}
