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

package logics

import (
	"encoding/json"
	"fmt"
	"github.com/fuyibing/log/v5"
	"github.com/fuyibing/log/v5/tracers"
	"github.com/fuyibing/util/v8/web/response"
	"github.com/kataras/iris/v12"
	"net/http"
	"regexp"
)

type (
	// Logic
	// 业务逻辑接口.
	Logic func(s tracers.Span, i iris.Context) interface{}
)

// New
// 执行业务过程.
//
//   logics.New(i)
//   logics.New(i, example.Fn)
//   logics.New(i, example.NewExample().Run)
func New(i iris.Context, logics ...Logic) (res interface{}) {
	var (
		req = i.Request()
		spa = log.NewSpanFromRequest(i.Request(),
			fmt.Sprintf("%s{%s}", req.Method, req.URL.Path),
		)
	)

	// 请求结束.
	defer func() {
		// 捕获异常.
		if r := recover(); r != nil {
			err := fmt.Errorf("%v", r)
			spa.Logger().Fatal("request fatal: %v", err)
			res = response.With.ErrorCode(err, http.StatusInternalServerError)
		}

		// 记录结果.
		buf, _ := json.Marshal(res)
		spa.Logger().Info("http.response.body: %s", buf)
		spa.End()

		// 链路下探.
		i.Next()
	}()

	// 开始请求.
	func() {
		switch i.Request().Method {
		case http.MethodPost, http.MethodPut:
			if b, be := i.GetBody(); be == nil {
				if s := fmt.Sprintf("%s", b); s != "" {
					spa.Logger().Info("http.request.body: %s", regexp.MustCompile(`\n\s*`).ReplaceAllString(s, ""))
				}
			}
		}
	}()

	// 调用逻辑.
	if len(logics) > 0 {
		res = logics[0](spa, i)
		return
	}

	// 内部错误.
	res = response.With.ErrorCode(fmt.Errorf("logic not specified"), http.StatusInternalServerError)
	return
}
