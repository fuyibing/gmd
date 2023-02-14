// author: wsfuyibing <websearch@163.com>
// date: 2023-01-18

package topic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fuyibing/gmd/app"
	"github.com/fuyibing/gmd/app/md"
	"github.com/fuyibing/gmd/app/md/base"
	"github.com/fuyibing/log/v3"
	"github.com/fuyibing/log/v3/trace"
	"github.com/fuyibing/util/v2/web/request"
	"github.com/fuyibing/util/v2/web/response"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"strings"
)

type (
	Batch struct {
		registry *base.Registry
		request  *BatchRequest
		response *BatchResponse
	}

	BatchRequest struct {
		TopicName     string        `json:"topic_name" label:"Topic name" validate:"required,min=2,max=30"`
		TopicTag      string        `json:"topic_tag" label:"Topic tag" validate:"required,min=2,max=60"`
		Messages      []interface{} `json:"messages" label:"Message list" desc:"Accept json string or json object in list"`
		MessageBodies []string      `json:"-" validate:"required,min=1,max=100" label:"Message list"`
	}

	BatchResponse struct {
		Count      int    `json:"count" label:"Message count" mock:"3"`
		Hash       string `json:"hash" label:"Message hash" mock:"C0837A1B5E264F19826F31457D51546D"`
		RegistryId int    `json:"registry_id" label:"Registry id" mock:"1"`
	}
)

func NewBatch() *Batch {
	return &Batch{
		request:  &BatchRequest{},
		response: &BatchResponse{},
	}
}

func (o *Batch) Run(ctx context.Context, i iris.Context) (res interface{}) {
	var (
		c   context.Context
		err error
	)

	// Read payload json string
	// then assign to request fields.
	if i.ReadJSON(o.request) != nil {
		err = fmt.Errorf("invalid json payload")
		return response.With.ErrorCode(err, app.CodeInvalidPayloadFormat)
	}

	// Validate
	// requested payload params.
	if err = o.request.Validate(); err == nil {
		err = request.Validate.Struct(o.request)
	}
	if err != nil {
		return response.With.ErrorCode(err, app.CodeInvalidPayloadFields)
	}

	// Return error
	// if registry not found in memory.
	if o.registry = base.Memory.GetRegistryByName(o.request.TopicName, o.request.TopicTag); o.registry == nil {
		return response.With.ErrorCode(
			fmt.Errorf("registry not found"),
			app.CodeServiceReadNotFound,
		)
	}

	// Init key fields
	// for response.
	o.response.Hash = strings.ToUpper(strings.ReplaceAll(uuid.NewString(), "-", ""))
	o.response.RegistryId = o.registry.Id
	o.response.Count = len(o.request.MessageBodies)

	// Message send progress.
	log.Infofc(ctx, "logic call producer manager: topic=%s, tag=%s, filter=%s, hash=%s, total=%d", o.registry.TopicName, o.registry.TopicTag, o.registry.FilterTag, o.response.Hash, o.response.Count)
	c = trace.Child(ctx)
	if err = o.Send(c); err != nil {
		return response.With.ErrorCode(
			fmt.Errorf("message publish failed"),
			app.CodeAdapterError,
		)
	}
	return response.With.Data(o.response)
}

func (o *Batch) Send(ctx context.Context) error {
	payloads := make([]*base.Payload, 0)

	// Iterate message list into buffer.
	for i0, s0 := range o.request.MessageBodies {
		log.Infofc(ctx, "batch item: offset=%d, item=%d-%d", i0, o.response.Count, i0+1)
		c0 := trace.Child(ctx)

		// Append to buffers.
		payloads = append(payloads, func(c1 context.Context, o1 int, s1 string) *base.Payload {
			p := base.Pool.AcquirePayload().SetContext(c1)
			p.Hash = o.response.Hash
			p.Offset = o1
			p.RegistryId = o.registry.Id
			p.TopicName = o.registry.TopicName
			p.TopicTag = o.registry.TopicTag
			p.FilterTag = o.registry.FilterTag
			p.MessageBody = s1
			return p
		}(c0, i0, s0))
	}

	// Send message progress.
	return md.Boot.Producer().Publish(payloads...)
}

func (o *BatchRequest) Validate() error {
	o.MessageBodies = make([]string, 0)

	for _, v := range o.Messages {
		if s, ok := v.(string); ok {
			if s = strings.TrimSpace(s); s != "" {
				o.MessageBodies = append(o.MessageBodies, s)
			}
			continue
		}

		b, _ := json.Marshal(v)
		if s := strings.TrimSpace(string(b)); s != "" {
			o.MessageBodies = append(o.MessageBodies, s)
		}
	}

	return nil
}
