// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package managers

import (
	"context"
	"github.com/fuyibing/gmd/v8/core/base"
	"github.com/fuyibing/util/v8/process"
)

type (
	// Remoting
	// 服务端管理器.
	Remoting struct {
		adapter   base.RemotingManager
		callable  base.RemotingCallable
		name      string
		processor process.Processor
	}
)

func (o *Remoting) OnAfter(_ context.Context) (ignored bool) {
	return
}

func (o *Remoting) OnBefore(_ context.Context) (ignored bool) {
	return
}

func (o *Remoting) OnBeforeSubprocess(_ context.Context) (ignored bool) {
	// Return true
	// if constructor not injected into container.
	if o.callable = Container.GetRemoting(); o.callable == nil {
		// log.Errorf("remoting constructor for {%s} adapter not injected", app.Config.GetAdapter())
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

func (o *Remoting) OnCall(_ context.Context) (ignored bool) {
	return
}

func (o *Remoting) OnCallListen(ctx context.Context) (ignored bool) {
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (o *Remoting) OnPanic(ctx context.Context, v interface{}) {
	// log.Panicfc(ctx, "processor {%s} fatal: %v", o.name, v)
}

func (o *Remoting) init() *Remoting {
	o.name = "remoting-manager"
	o.processor = process.New(o.name).After(
		o.OnAfter,
	).Before(
		o.OnBefore,
		o.OnBeforeSubprocess,
	).Callback(
		o.OnCall,
		o.OnCallListen,
	).Panic(o.OnPanic)

	return o
}
