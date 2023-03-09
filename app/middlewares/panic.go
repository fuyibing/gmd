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
	"github.com/kataras/iris/v12"
	"net/http"
)

// Panic
// 捕获异常.
func Panic(i iris.Context) {
	defer func() {
		v := recover()

		// 取消处理.
		if v == nil || i.IsStopped() {
			return
		}

		// 发送错误.
		sendError(i, http.StatusInternalServerError, v)
	}()

	i.Next()
}
