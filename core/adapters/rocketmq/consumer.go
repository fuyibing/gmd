// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package rocketmq

import (
	"context"
	"github.com/fuyibing/gmd/v8/core/base"
	"github.com/fuyibing/log/v8"
	"github.com/fuyibing/util/v8/process"
)

type (
	// Consumer
	// for rocketmq adapter.
	Consumer struct {
		consume      base.ConsumerProcess
		id, parallel int
		name         string
		processor    process.Processor
	}
)

func NewConsumer(id, parallel int, name string, consume base.ConsumerProcess) *Consumer {
	return (&Consumer{
		id: id, parallel: parallel,
		name: name, consume: consume,
	}).init()
}

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *Consumer) Processor() process.Processor { return o.processor }

// /////////////////////////////////////////////////////////////
// Processor events
// /////////////////////////////////////////////////////////////

func (o *Consumer) OnAfter(_ context.Context) (ignored bool) {
	log.Infof("{%s} stopped", o.name)
	return
}

func (o *Consumer) OnBefore(_ context.Context) (ignored bool) {
	log.Infof("{%s} start", o.name)
	return
}

func (o *Consumer) OnCall(_ context.Context) (ignored bool) {
	log.Infof("{%s} listen channel signal", o.name)
	return
}

func (o *Consumer) OnCallChannel(ctx context.Context) (ignored bool) {
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (o *Consumer) OnPanic(ctx context.Context, v interface{}) {
	log.Panicfc(ctx, "processor {%s} fatal: %v", o.name, v)
}

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *Consumer) init() *Consumer {
	o.processor = process.New(o.name).After(
		o.OnAfter,
	).Before(
		o.OnBefore,
	).Callback(
		o.OnCall,
		o.OnCallChannel,
	).Panic(o.OnPanic)

	return o
}
