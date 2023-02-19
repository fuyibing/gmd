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
	log.Infof("%s: processor stopped", o.processor.Name())
	return
}

func (o *Remoting) OnBefore(_ context.Context) (ignored bool) {
	log.Infof("%s: start processor", o.processor.Name())
	return
}

func (o *Remoting) OnListen(ctx context.Context) (ignored bool) {
	for {
		select {
		case <-ctx.Done():
			log.Debugf("%s: %v", o.processor.Name(), ctx.Err())
			return
		}
	}
}

func (o *Remoting) OnPanic(ctx context.Context, v interface{}) {
	log.Panicfc(ctx, "%s: %v", o.processor.Name(), v)
}

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *Remoting) init() *Remoting {
	o.processor = process.New("aliyunmns-remoting").After(
		o.OnAfter,
	).Before(
		o.OnBefore,
	).Callback(o.OnListen).Panic(o.OnPanic)

	return o
}
