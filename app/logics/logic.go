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
	// get instance from pool.
	Logic func() LogicHandler

	// LogicHandler
	// instance for api request.
	LogicHandler interface {
		// Release
		// instance into pool.
		Release()

		// Run
		// logic executor.
		Run(span tracers.Span, i iris.Context) interface{}
	}
)

// New
// get logic instance from pool then run it, release when executed
// completed.
func New(i iris.Context, logics ...Logic) (res interface{}) {
	var (
		req = i.Request()

		// Create tracer span from http request.
		span = log.NewSpanFromRequest(i.Request(),
			fmt.Sprintf("%s %s", req.Method, req.URL.Path),
		)
	)

	// End request.
	defer func() {
		// Recover runtime fatal.
		if r := recover(); r != nil {
			span.Logger().Fatal("http request fatal: %v", r)
			res = response.With.ErrorCode(fmt.Errorf("%v", r), app.CodePanicOccurred)
		}

		// Store request result.
		buf, _ := json.Marshal(res)
		span.Logger().Info("http request end: %s", buf)

		// End span then call next middlewares.
		span.End()
		i.Next()
	}()

	// Begin request.
	func() {
		span.Logger().Info("http request begin: %s %s %s", req.Proto, req.Method, req.RequestURI)

		// Store request body.
		switch i.Request().Method {
		case http.MethodPost, http.MethodPut:
			if b, be := i.GetBody(); be == nil {
				if s := fmt.Sprintf("%s", b); s != "" {
					span.Logger().Info("http request body: %s", regexp.MustCompile(`\n\s*`).ReplaceAllString(s, ""))
				}
			}
		}
	}()

	// Return error if logic not specified.
	if len(logics) == 0 {
		res = response.With.ErrorCode(fmt.Errorf("logic not specified"), app.CodeLogicUndefined)
		return
	}

	// Return logic processor result.
	res = logics[0]().Run(span, i)
	return
}
