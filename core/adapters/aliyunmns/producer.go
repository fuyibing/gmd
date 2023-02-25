// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package aliyunmns

import (
	"context"
	"github.com/fuyibing/util/v8/process"
)

type (
	// Producer
	// for aliyunmns adapter.
	Producer struct {
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

// /////////////////////////////////////////////////////////////
// Processor events
// /////////////////////////////////////////////////////////////

func (o *Producer) OnAfter(_ context.Context) (ignored bool) {
	// log.Infof("{%s} stopped", o.name)
	return
}

func (o *Producer) OnBefore(_ context.Context) (ignored bool) {
	// log.Infof("{%s} start", o.name)
	return
}

func (o *Producer) OnCall(_ context.Context) (ignored bool) {
	// log.Infof("{%s} listen channel signal", o.name)
	return
}

func (o *Producer) OnCallListen(ctx context.Context) (ignored bool) {
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (o *Producer) OnPanic(ctx context.Context, v interface{}) {
	// log.Panicfc(ctx, "processor {%s} fatal: %v", o.name, v)
}

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *Producer) init() *Producer {
	o.name = "aliyunmns-producer"
	o.processor = process.New(o.name).After(
		o.OnAfter,
	).Before(
		o.OnBefore,
	).Callback(
		o.OnCall,
		o.OnCallListen,
	).Panic(o.OnPanic)

	return o
}
