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
// date: 2023-03-06

package topic

import (
	"encoding/json"
	"fmt"
	"github.com/fuyibing/gmd/v8/app"
	"github.com/fuyibing/gmd/v8/core/base"
	"github.com/fuyibing/gmd/v8/core/managers"
	"github.com/fuyibing/log/v5/tracers"
	"github.com/fuyibing/util/v8/web/request"
	"github.com/fuyibing/util/v8/web/response"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"strings"
)

type (
	// Publish
	// 发布单条消息.
	Publish struct {
		request  *PublishRequest
		response *PublishResponse
	}

	// PublishRequest
	// 发布入参.
	PublishRequest struct {
		TopicName string      `json:"topicName" label:"主题名" validate:"required,gte=2,lte=30"`
		TopicTag  string      `json:"topicTag" label:"标签名" validate:"required,gte=2,lte=60"`
		Message   interface{} `json:"message"`

		MessageBody string `json:"-" label:"消息内容" validate:"required,lte=65536"`
	}

	// PublishResponse
	// 发布结果.
	PublishResponse struct {
		Hash string `json:"hash"`
	}
)

func NewPublish() *Publish {
	return &Publish{
		request:  &PublishRequest{},
		response: (&PublishResponse{}).init(),
	}
}

func (o *Publish) Run(span tracers.Span, i iris.Context) (res interface{}) {
	var (
		code int
		err  error
	)

	// 记录结果.
	defer func() {
		span.Kv().
			Add("publish.topic.name", o.request.TopicName).
			Add("publish.topic.tag", o.request.TopicTag)

		if err != nil {
			span.Kv().Add("publish.result.error", err)
			res = response.With.ErrorCode(err, code)
		} else {
			span.Kv().Add("publish.result.hash", o.response.Hash)
			res = response.With.Data(o.response)
		}
	}()

	// 入参格式.
	if err = i.ReadJSON(o.request); err != nil {
		code = app.CodeInvalidPayloadFormatter
		return
	}

	// 入参校验
	if err = o.request.Validate(); err == nil {
		err = request.Validate.Struct(o.request)
	}
	if err != nil {
		code = app.CodeInvalidPayloadFields
		return
	}

	// 校验注册.
	registry := base.Memory.GetRegistryByName(o.request.TopicName, o.request.TopicTag)
	if registry == nil {
		code = app.CodeRegistryNotFound
		err = fmt.Errorf("registry not found")
		return
	}

	// 创建消息.
	v := base.Pool.AcquirePayload()
	v.SetContext(span.Context())
	v.Hash = o.response.Hash
	v.RegistryId = registry.Id
	v.TopicName = registry.TopicName
	v.TopicTag = registry.TopicTag
	v.FilterTag = registry.FilterTag
	v.MessageBody = o.request.MessageBody

	// 发布消息.
	err = managers.Boot.Producer().Publish(v)
	return
}

// Request

func (o *PublishRequest) Validate() (err error) {
	if s, ok := o.Message.(string); ok {
		o.MessageBody = s
		return
	}

	var buf []byte
	buf, err = json.Marshal(o.Message)
	o.MessageBody = string(buf)
	return
}

// Response

func (o *PublishResponse) init() *PublishResponse {
	o.Hash = strings.ReplaceAll(uuid.NewString(), "-", "")
	return o
}
