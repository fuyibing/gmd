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
	"github.com/fuyibing/log/v5/tracers"
	"github.com/kataras/iris/v12"
)

type (
	// PostBatch
	// 发布批量消息.
	PostBatch struct {
		request  *PostBatchRequest
		response *PostBatchResponse
	}

	// PostBatchRequest
	// 消息入参.
	PostBatchRequest struct {
	}

	// PostBatchResponse
	// 发布结果.
	PostBatchResponse struct {
	}
)

func NewPostBatch() *PostBatch {
	return &PostBatch{
		request:  &PostBatchRequest{},
		response: &PostBatchResponse{},
	}
}

func (o *PostBatch) Run(span tracers.Span, i iris.Context) interface{} {
	return "/topic/batch"
}
