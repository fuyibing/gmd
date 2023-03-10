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
	"github.com/fuyibing/log/v5"
	"github.com/fuyibing/log/v5/tracers"
	"github.com/fuyibing/util/v8/web/request"
	"github.com/fuyibing/util/v8/web/response"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"strings"
)

type (
	// PostBatch
	// 发布消息.
	PostBatch struct {
		request  *PostBatchRequest
		response *PostBatchResponse
	}

	// PostBatchRequest
	// 消息入参.
	PostBatchRequest struct {
		Topic         string        `json:"topic" label:"主题名" validate:"required,gte=2,lte=30"`
		Tag           string        `json:"tag" label:"标签名" validate:"required,gte=2,lte=60"`
		Messages      []interface{} `json:"messages" label:"消息列表"`
		MessageBodies []string      `json:"-" label:"消息列表" validate:"required,gte=1,lte=100"`
	}

	// PostBatchResponse
	// 发布结果.
	PostBatchResponse struct {
		Hash  string `json:"hash" label:"哈希码"`
		Count int    `json:"count" label:"消息数"`
	}
)

func NewPostBatch() *PostBatch {
	return &PostBatch{
		request:  &PostBatchRequest{},
		response: &PostBatchResponse{},
	}
}

// +---------------------------------------------------------------------------+
// + Logic runner                                                              |
// +---------------------------------------------------------------------------+

func (o *PostBatch) Run(span tracers.Span, i iris.Context) (res interface{}) {
	var (
		code     int
		err      error
		exists   bool
		registry *base.Registry
	)

	// 发布结束.
	defer func() {
		span.Kv().
			Add("payload.created.batch.tag", o.request.Tag).
			Add("payload.created.batch.name", o.request.Topic)

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

	// 哈希表.
	o.response.Hash = strings.ToUpper(strings.ReplaceAll(uuid.NewString(), "-", ""))
	o.response.Count = len(o.request.MessageBodies)

	// 遍历消息.
	for mi, mb := range o.request.MessageBodies {
		o.Range(span, registry, mb, o.response.Hash, mi)
	}

	// 完成发布.
	res = response.With.Data(o.response)
	return
}

func (o *PostBatch) Range(parent tracers.Span, registry *base.Registry, messageBody, hash string, offset int) {
	var (
		span    = log.NewSpanFromContext(parent.Context(), "payload.created.batch")
		payload = base.Pool.AcquirePayload().SetContext(span.Context())
	)

	span.Kv().
		Add("payload.created.batch.hash", hash).
		Add("payload.created.batch.offset", offset)

	defer span.End()

	payload.Hash = hash
	payload.Offset = offset

	payload.FilterTag = registry.FilterTag
	payload.TopicTag = registry.TopicTag
	payload.TopicName = registry.TopicName
	payload.RegistryId = registry.Id
	payload.MessageBody = messageBody

	// 发布消息.
	if err := md.Manager.Producer().Publish(payload); err != nil {
		span.Logger().Error("payload batch error: %v", err)
	} else {
		span.Logger().Info("payload batch completed")
	}
}

// +---------------------------------------------------------------------------+
// + Request validate                                                          |
// +---------------------------------------------------------------------------+

func (o *PostBatchRequest) Validate() error {
	if o.MessageBodies == nil {
		o.MessageBodies = make([]string, 0)
	}
	for _, v := range o.Messages {
		if s, ok := v.(string); ok {
			if s = strings.TrimSpace(s); s != "" {
				o.MessageBodies = append(o.MessageBodies, s)
			}
		} else {
			if buf, err := json.Marshal(v); err == nil {
				if str := strings.TrimSpace(string(buf)); str != "" {
					o.MessageBodies = append(o.MessageBodies, s)
				}
			}
		}
	}
	return nil
}
