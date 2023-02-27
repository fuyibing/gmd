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

package index

import (
	"github.com/fuyibing/gmd/v8/app"
	"github.com/fuyibing/log/v5/cores"
	"github.com/fuyibing/util/v8/web/response"
	"github.com/kataras/iris/v12"
	"time"
)

type (
	// Home
	// 默认页.
	Home struct {
		response *HomeResponse
	}

	// HomeResponse
	// 默认页结果.
	HomeResponse struct {
		Time time.Time `json:"time" label:"启动时间"`
	}
)

func NewHome() *Home {
	return &Home{
		response: &HomeResponse{
			Time: app.Config.GetStartedTime(),
		},
	}
}

func (o *Home) Run(_ cores.Span, _ iris.Context) interface{} {
	return response.With.Data(o.response)
}
