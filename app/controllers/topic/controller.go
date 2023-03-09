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

package topic

import (
	"github.com/fuyibing/gmd/v8/app/logics"
	"github.com/fuyibing/gmd/v8/app/logics/topic"
	"github.com/kataras/iris/v12"
)

// Controller
// 主题操作.
//
// @RoutePrefix(/topic)
type Controller struct{}

// func (o *Controller) PostBatch(i iris.Context) interface{} {}

// PostBatch
// 批量发布.
//
// @Request(app/logics/topic.PostBatchRequest)
// @Response(app/logics/topic.PostBatchResponse)
func (o *Controller) PostBatch(i iris.Context) interface{} {
	return logics.New(i, topic.NewPostBatch().Run)
}

// PostPublish
// 发布消息.
//
// @Request(app/logics/topic.PostPublishRequest)
// @Response(app/logics/topic.PostPublishResponse)
func (o *Controller) PostPublish(i iris.Context) interface{} {
	return logics.New(i, topic.NewPostPublish().Run)
}
