// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package aliyunmns

import (
	"context"
	"github.com/fuyibing/log/v8"
	"github.com/fuyibing/util/v8/process"
)

type (
	// Remoting
	// for aliyunmns adapter.
	Remoting struct {
		name      string
		processor process.Processor
	}
)

func NewRemoting() *Remoting {
	return (&Remoting{}).init()
}

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *Remoting) Processor() process.Processor { return o.processor }

// /////////////////////////////////////////////////////////////
// Processor events
// /////////////////////////////////////////////////////////////

func (o *Remoting) OnAfter(_ context.Context) (ignored bool) {
	log.Infof("{%s} stopped", o.name)
	return
}

func (o *Remoting) OnBefore(_ context.Context) (ignored bool) {
	log.Infof("{%s} start", o.name)
	return
}

func (o *Remoting) OnCall(_ context.Context) (ignored bool) {
	log.Infof("{%s} listen channel signal", o.name)
	return
}

func (o *Remoting) OnCallChannel(ctx context.Context) (ignored bool) {
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (o *Remoting) OnPanic(ctx context.Context, v interface{}) {
	log.Panicfc(ctx, "processor {%s} fatal: %v", o.name, v)
}

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *Remoting) init() *Remoting {
	o.name = "aliyunmns-remoting"
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
