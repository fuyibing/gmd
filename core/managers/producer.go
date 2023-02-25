// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package managers

import (
	"context"
	"github.com/fuyibing/gmd/v8/core/base"
	"github.com/fuyibing/util/v8/process"
)

type Producer struct {
	adapter   base.ProducerManager
	callable  base.ProducerCallable
	name      string
	processor process.Processor
}

// /////////////////////////////////////////////////////////////
// Exported methods
// /////////////////////////////////////////////////////////////

// /////////////////////////////////////////////////////////////
// Processor events
// /////////////////////////////////////////////////////////////

func (o *Producer) OnAfter(_ context.Context) (ignored bool) {
	return
}

func (o *Producer) OnBefore(_ context.Context) (ignored bool) {
	return
}

func (o *Producer) OnBeforeSubprocess(_ context.Context) (ignored bool) {
	// Return true
	// if constructor not injected into container.
	if o.callable = Container.GetProducer(); o.callable == nil {
		// log.Errorf("producer constructor for {%s} adapter not injected", app.Config.GetAdapter())
		return true
	}

	// Dynamically
	// configure adapter instance and name.
	if o.adapter == nil {
		o.adapter = o.callable()
		o.processor.Add(o.adapter.Processor())
	}

	return
}

func (o *Producer) OnCall(_ context.Context) (ignored bool) {
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
	o.name = "producer-manager"
	o.processor = process.New(o.name).After(
		o.OnAfter,
	).Before(
		o.OnBefore,
		o.OnBeforeSubprocess,
	).Callback(
		o.OnCall,
		o.OnCallListen,
	).Panic(
		o.OnPanic,
	)

	return o
}
