// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package rocketmq

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/fuyibing/gmd/v8/core/base"
	"github.com/fuyibing/log/v5"
	"github.com/fuyibing/util/v8/process"
)

type (
	// Producer
	// for rocketmq adapter.
	Producer struct {
		client    rocketmq.Producer
		name      string
		processor process.Processor
	}
)

func NewProducer() *Producer {
	return (&Producer{}).init()
}

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *Producer) Processor() process.Processor { return o.processor }

func (o *Producer) Publish(v *base.Payload) (messageId string, err error) {
	var (
		msg = &primitive.Message{
			Topic: Agent.GenTopicName(v.TopicName),
			Body:  []byte(v.MessageBody),
		}
		res *primitive.SendResult
	)
	// 绑定标签.
	if v.TopicTag != "" {
		msg.WithTag(v.TopicTag)
	}
	// 发布过程.
	if res, err = o.client.SendSync(v.GetContext(), msg); err == nil {
		messageId = res.MsgID
	}
	return
}

// /////////////////////////////////////////////////////////////
// Event methods
// /////////////////////////////////////////////////////////////

func (o *Producer) onAfter(_ context.Context) (ignored bool) {
	return
}

func (o *Producer) onBefore(_ context.Context) (ignored bool) {
	return
}

func (o *Producer) onCall(_ context.Context) (ignored bool) {
	return
}

func (o *Producer) onCallListen(ctx context.Context) (ignored bool) {
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (o *Producer) onPanic(ctx context.Context, v interface{}) {
	if spa, exists := log.Span(ctx); exists {
		spa.Logger().Fatal("<%s> %v", o.name, v)
	} else {
		log.Fatal("<%s> %v", o.name, v)
	}
}

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *Producer) init() *Producer {
	o.name = "rocketmq-producer"
	o.processor = process.New(o.name).After(o.onAfter).
		Before(o.onBefore).
		Callback(o.onCall, o.onCallListen).
		Panic(o.onPanic)

	return o
}
