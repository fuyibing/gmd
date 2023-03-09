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

package logics

import (
	"encoding/json"
	"fmt"
	"github.com/fuyibing/gmd/v8/app"
	"github.com/fuyibing/log/v5"
	"github.com/fuyibing/log/v5/tracers"
	"github.com/fuyibing/util/v8/web/response"
	"github.com/kataras/iris/v12"
	"net/http"
	"regexp"
)

type (
	// Logic
	// 逻辑回调.
	Logic func(span tracers.Span, i iris.Context) interface{}
)

// New
// 执行逻辑.
func New(i iris.Context, logics ...Logic) (res interface{}) {
	var (
		req = i.Request()

		// 请求链.
		span = log.NewSpanFromRequest(i.Request(), fmt.Sprintf("%s %s",
			req.Method, req.URL.Path,
		))
	)

	// 结束请求.
	defer func() {
		// 捕获异常.
		if r := recover(); r != nil {
			span.Logger().Fatal("http request fatal: %v", r)
			res = response.With.ErrorCode(fmt.Errorf("%v", r), app.CodePanicOccurred)
		}

		// 记录结果.
		buf, _ := json.Marshal(res)
		span.Logger().Info("http request end: %s", buf)

		// 结束链路.
		span.End()
		i.Next()
	}()

	// 请求开始.
	func() {
		span.Logger().Info("http request begin: %s %s %s", req.Proto, req.Method, req.RequestURI)

		// 记录入参.
		switch i.Request().Method {
		case http.MethodPost, http.MethodPut:
			if b, be := i.GetBody(); be == nil {
				if s := fmt.Sprintf("%s", b); s != "" {
					span.Logger().Info("http request body: %s", regexp.MustCompile(`\n\s*`).ReplaceAllString(s, ""))
				}
			}
		}
	}()

	// 请求过程.
	if len(logics) == 0 {
		res = response.With.ErrorCode(
			fmt.Errorf("logic not specified"),
			app.CodeLogicUndefined,
		)
	} else {
		res = logics[0](span, i)
	}
	return
}
