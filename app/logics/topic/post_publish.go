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
	"encoding/json"
	"fmt"
	"github.com/fuyibing/gmd/v8/app"
	"github.com/fuyibing/gmd/v8/md"
	"github.com/fuyibing/gmd/v8/md/base"
	"github.com/fuyibing/log/v5/tracers"
	"github.com/fuyibing/util/v8/web/request"
	"github.com/fuyibing/util/v8/web/response"
	"github.com/kataras/iris/v12"
	"strings"
)

type (
	// PostPublish
	// 发布消息.
	PostPublish struct {
		request  *PostPublishRequest
		response *PostPublishResponse
	}

	// PostPublishRequest
	// 消息入参.
	PostPublishRequest struct {
		Topic       string      `json:"topic" label:"主题名" mock:"Topic" validate:"required,gte=2,lte=30"`
		Tag         string      `json:"tag" label:"标签名" mock:"tag" validate:"required,gte=2,lte=60"`
		Message     interface{} `json:"message" label:"消息正文" validate:"required" desc:"接JSON字符串或JSON对象"`
		MessageBody string      `json:"-" label:"消息正文" validate:"required,gte=1,lte=65536"`
	}

	// PostPublishResponse
	// 发布结果.
	PostPublishResponse struct {
		Hash string `json:"hash" label:"哈希码" mock:"CFD44CFBCB0D451A90E7EE193785F289"`
	}
)

func NewPostPublish() *PostPublish {
	return &PostPublish{
		request:  &PostPublishRequest{},
		response: &PostPublishResponse{},
	}
}

// +---------------------------------------------------------------------------+
// + Logic runner                                                              |
// +---------------------------------------------------------------------------+

func (o *PostPublish) Run(span tracers.Span, i iris.Context) (res interface{}) {
	var (
		code     int
		err      error
		exists   bool
		payload  *base.Payload
		registry *base.Registry
	)

	// 发布结束.
	defer func() {
		span.Kv().
			Add("payload.created.topic.tag", o.request.Tag).
			Add("payload.created.topic.name", o.request.Topic)

		// 覆盖结果.
		if err != nil {
			res = response.With.ErrorCode(err, code)
		}
	}()

	// 校验入参.
	if err = i.ReadJSON(o.request); err == nil {
		if err = o.request.Validate(); err == nil {
			err = request.Validate.Struct(o.request)
		}
	}
	if err != nil {
		code = app.CodePayloadInvalid
		return
	}

	// 注册关系.
	if registry, exists = base.Memory.GetRegistryByNames(o.request.Topic, o.request.Tag); !exists {
		code = app.CodeRegistryNotFound
		err = fmt.Errorf("registry not found")
		return
	}

	// 消息结构.
	payload = base.Pool.AcquirePayload().SetContext(span.Context())
	payload.GenHash()

	o.response.Hash = payload.Hash

	// 基础字段.
	payload.FilterTag = registry.FilterTag
	payload.MessageBody = o.request.MessageBody
	payload.RegistryId = registry.Id
	payload.TopicName = registry.TopicName
	payload.TopicTag = registry.TopicTag

	span.Kv().
		Add("payload.created.hash", payload.Hash).
		Add("payload.created.offset", payload.Offset)

	// 发布消息.
	if err = md.Manager.Producer().Publish(payload); err != nil {
		return
	}

	// 完成发布.
	res = response.With.Data(o.response)
	return
}

// +---------------------------------------------------------------------------+
// + Request validate                                                          |
// +---------------------------------------------------------------------------+

func (o *PostPublishRequest) Validate() error {
	if s, ok := o.Message.(string); ok {
		o.MessageBody = strings.TrimSpace(s)
	} else {
		if buf, err := json.Marshal(o.Message); err == nil {
			o.MessageBody = strings.TrimSpace(string(buf))
		}
	}
	return nil
}
