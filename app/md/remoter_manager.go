// author: wsfuyibing <websearch@163.com>
// date: 2023-02-08

package md

import (
	"context"
	"github.com/fuyibing/gmd/app/md/adapters"
	"github.com/fuyibing/gmd/app/md/conf"
	"github.com/fuyibing/log/v8"
	"github.com/fuyibing/util/v8/process"
	"time"
)

type (
	RemoterManager interface {
		// Adapter
		// return remoter adapter interface.
		Adapter() adapters.RemoterAdapter

		// Processor
		// return remoter processor interface.
		//
		//   x := md.Boot.Remoter().Processor()
		//   x.Start(ctx)
		Processor() process.Processor
	}

	remoter struct {
		adapter   adapters.RemoterAdapter
		processor process.Processor
	}
)

// /////////////////////////////////////////////////////////////
// Interface methods.
// /////////////////////////////////////////////////////////////

func (o *remoter) Adapter() adapters.RemoterAdapter { return o.adapter }
func (o *remoter) Processor() process.Processor     { return o.processor }

// /////////////////////////////////////////////////////////////
// Event methods.
// /////////////////////////////////////////////////////////////

// OnAfter
// called when processor stopped.
func (o *remoter) OnAfter(_ context.Context) (ignored bool) {
	log.Debugf("remoter manager: processor stopped")
	return
}

// OnBefore
// called when processor start.
func (o *remoter) OnBefore(_ context.Context) (ignored bool) {
	log.Debugf("remoter manager: start processor")
	return
}

// OnCallAdapterBuild
// build remoter adapter.
func (o *remoter) OnCallAdapterBuild(_ context.Context) (ignored bool) {
	var err error

	// Return error
	// if create adapter failed.
	if o.adapter, err = adapters.NewRemoter(conf.Config.Adapter); err != nil {
		log.Errorf("remoter manager: build %s adapter, error=%v", conf.Config.Adapter, err)
		return true
	}

	// Next
	// event callee.
	log.Debugf("remoter manager: build %s adapter", conf.Config.Adapter)
	return
}

// OnCallAdapterDestroy
// wait until remoter adapter stopped.
func (o *remoter) OnCallAdapterDestroy(ctx context.Context) (ignored bool) {
	// Next event caller
	// if adapter processor stopped.
	if o.adapter.Processor().Stopped() {
		o.adapter = nil
		log.Debugf("remoter manager: destroy %s adapter", conf.Config.Adapter)
		return
	}

	// Recall
	// wait for a while.
	time.Sleep(conf.EventSleepDuration)
	return o.OnCallAdapterDestroy(ctx)
}

// OnCallAdapterStart
// start remoter adapter in coroutine.
func (o *remoter) OnCallAdapterStart(ctx context.Context) (ignored bool) {
	go func() { _ = o.adapter.Processor().Start(ctx) }()
	return
}

// OnCallChannel
// listen channel signal.
func (o *remoter) OnCallChannel(ctx context.Context) (ignored bool) {
	log.Debugf("remoter manager: listen channel signal")

	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

// OnPanic
// called with panic at runtime.
func (o *remoter) OnPanic(ctx context.Context, v interface{}) {
	log.Panicfc(ctx, "remoter manager: %v", v)
}

// /////////////////////////////////////////////////////////////
// Constructor methods.
// /////////////////////////////////////////////////////////////

func (o *remoter) init() *remoter {
	// Register remoter processor event callbacks.
	o.processor = process.New("remoter manager").After(
		o.OnAfter,
	).Before(
		o.OnBefore,
	).Callback(
		o.OnCallAdapterBuild,
		o.OnCallAdapterStart,
		o.OnCallChannel,
		o.OnCallAdapterDestroy,
	).Panic(o.OnPanic)

	return o
}
