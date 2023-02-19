// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package aliyunmns

import (
	"context"
	"github.com/fuyibing/log/v8"
	"github.com/fuyibing/util/v8/process"
)

type (
	// Producer
	// for aliyunmns adapter.
	Producer struct {
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
	log.Infof("%s: processor stopped", o.processor.Name())
	return
}

func (o *Producer) OnBefore(_ context.Context) (ignored bool) {
	log.Infof("%s: start processor", o.processor.Name())
	return
}

func (o *Producer) OnListen(ctx context.Context) (ignored bool) {
	for {
		select {
		case <-ctx.Done():
			log.Debugf("%s: %v", o.processor.Name(), ctx.Err())
			return
		}
	}
}

func (o *Producer) OnPanic(ctx context.Context, v interface{}) {
	log.Panicfc(ctx, "%s: %v", o.processor.Name(), v)
}

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *Producer) init() *Producer {
	o.processor = process.New("aliyunmns-producer").After(
		o.OnAfter,
	).Before(
		o.OnBefore,
	).Callback(o.OnListen).Panic(o.OnPanic)

	return o
}
