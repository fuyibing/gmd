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
	"github.com/fuyibing/gmd/v8/app/logics"
	"github.com/fuyibing/gmd/v8/app/logics/task"
	"github.com/kataras/iris/v12"
)

type (
	// Controller
	// 订阅任务.
	//
	// @RoutePrefix(/task)
	Controller struct{}
)

// PostEnable
// 启用订阅任务.
//
// @Request(app/logics/task.StatusRequest)
// @Response(app/logics/task.StatusResponse)
func (o *Controller) PostEnable(i iris.Context) interface{} {
	return logics.New(i, task.NewEnable().Run)
}