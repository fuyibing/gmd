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
	"github.com/fuyibing/log/v5/tracers"
	"github.com/fuyibing/util/v8/web/response"
	"github.com/kataras/iris/v12"
	"runtime"
)

type (
	GetPing struct {
		response *GetPingResponse
	}

	GetPingResponse struct {
		Goroutines int `json:"goroutines" label:"协程数" mock:"12"`
	}
)

func NewGetPing() *GetPing {
	return &GetPing{
		response: &GetPingResponse{},
	}
}

func (o *GetPing) Run(span tracers.Span, i iris.Context) interface{} {

	o.response.Goroutines = runtime.NumGoroutine()

	return response.With.Data(o.response)
}
