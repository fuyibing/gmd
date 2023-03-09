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

package middlewares

import (
	"fmt"
	"github.com/fuyibing/util/v8/web/response"
	"github.com/kataras/iris/v12"
	"net/http"
)

// ErrCode
// 处理错误码.
func ErrCode(i iris.Context) {
	sendError(i, i.GetStatusCode(), nil)
	i.Next()
}

// 发送错误.
func sendError(i iris.Context, code int, v interface{}) {
	var err error
	if v != nil {
		err = fmt.Errorf("%v", v)
	} else {
		err = fmt.Errorf("HTTP %d %s", code, http.StatusText(code))
	}

	// 发送内容.
	_, _ = i.JSON(response.With.ErrorCode(err, code))
}
