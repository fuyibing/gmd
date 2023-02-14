// author: wsfuyibing <websearch@163.com>
// date: 2023-02-03

package rabbitmq

import (
	"context"
	"fmt"
	"github.com/fuyibing/gmd/app/md/base"
	"github.com/fuyibing/log/v3"
	"github.com/fuyibing/util/v2/process"
)

type (
	// Producer
	// 消息生产者.
	Producer struct {
		name      string
		processor process.Processor
	}
)

func NewProducer() *Producer {
	return (&Producer{}).init()
}

func (o *Producer) Processor() process.Processor {
	return o.processor
}

func (o *Producer) Publish(payload *base.Payload) (string, error) {
	return "", fmt.Errorf("not support")
}

// /////////////////////////////////////////////////////////////
// 生产者构造
// /////////////////////////////////////////////////////////////

func (o *Producer) init() *Producer {
	o.name = "rabbitmq-producer"
	o.processor = process.New(o.name).
		After(o.onAfter).
		Before(o.onBefore).
		Callback(o.onCallerBefore, o.onCaller, o.onCallerAfter).
		Panic(o.onPanic)

	return o
}

// /////////////////////////////////////////////////////////////
// 生产者事件
// /////////////////////////////////////////////////////////////

func (o *Producer) onAfter(ctx context.Context) (ignored bool) {
	log.Infof("[%s] producer finish", o.name)
	return
}

func (o *Producer) onBefore(ctx context.Context) (ignored bool) {
	log.Infof("[%s] producer begin", o.name)
	return
}

func (o *Producer) onCaller(ctx context.Context) (ignored bool) {
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (o *Producer) onCallerAfter(ctx context.Context) (ignored bool) { return }

func (o *Producer) onCallerBefore(ctx context.Context) (ignored bool) { return }

func (o *Producer) onPanic(ctx context.Context, v interface{}) {
	log.Panicfc(ctx, "[%s] %v", o.name, v)
}
