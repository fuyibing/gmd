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
	Publish struct {
		registry *base.Registry
		request  *PublishRequest
		response *PublishResponse
	}

	PublishRequest struct {
		TopicName string      `json:"topic_name" label:"Topic name" validate:"required,min=2,max=30"`
		TopicTag  string      `json:"topic_tag" label:"Topic tag" validate:"required,min=2,max=60"`
		Message   interface{} `json:"message" label:"Message content" desc:"Accept json string or json object"`

		MessageBody string `json:"-" label:"Message content" validate:"required,min=2,max=65536"`
	}

	PublishResponse struct {
		Hash       string `json:"hash" label:"Message hash" mock:"C0837A1B5E264F19826F31457D51546D"`
		RegistryId int    `json:"registry_id" label:"Registry id" mock:"1"`
	}
)

func NewPublish() *Publish {
	return &Publish{
		request:  &PublishRequest{},
		response: &PublishResponse{},
	}
}

func (o *Publish) Run(ctx context.Context, i iris.Context) (res interface{}) {
	var err error

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

	// Message send progress.
	if err = o.Send(ctx); err != nil {
		return response.With.ErrorCode(
			fmt.Errorf("message publish failed"),
			app.CodeAdapterError,
		)
	}
	return response.With.Data(o.response)
}

func (o *Publish) Send(ctx context.Context) error {
	log.Infofc(ctx, "logic call producer manager: topic=%s, tag=%s, filter=%s, hash=%s", o.registry.TopicName, o.registry.TopicTag, o.registry.FilterTag, o.response.Hash)

	var (
		c = trace.Child(ctx)
		p = base.Pool.AcquirePayload().SetContext(c)
	)

	p.Hash = o.response.Hash
	p.Offset = 0
	p.RegistryId = o.registry.Id
	p.TopicName = o.registry.TopicName
	p.TopicTag = o.registry.TopicTag
	p.FilterTag = o.registry.FilterTag
	p.MessageBody = o.request.MessageBody

	return md.Boot.Producer().Publish(p)
}

func (o *PublishRequest) Validate() error {
	if s, ok := o.Message.(string); ok {
		o.MessageBody = strings.TrimSpace(s)
		return nil
	}

	buf, _ := json.Marshal(o.Message)
	o.MessageBody = string(buf)
	return nil
}
