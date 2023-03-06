// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package managers

import (
	"context"
	"github.com/fuyibing/log/v5"
	"github.com/fuyibing/util/v8/process"
)

type (
	// Retry
	// 重试管理器.
	Retry struct {
		name      string
		processor process.Processor
	}
)

func (o *Retry) OnAfter(_ context.Context) (ignored bool) {
	return
}

func (o *Retry) OnBefore(_ context.Context) (ignored bool) {
	return
}

func (o *Retry) OnCall(ctx context.Context) (ignored bool) {
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (o *Retry) OnPanic(ctx context.Context, v interface{}) {
	if spa, exists := log.Span(ctx); exists {
		spa.Logger().Fatal("<%s> %v", o.name, v)
	} else {
		log.Fatal("<%s> %v", o.name, v)
	}
}

func (o *Retry) init() *Retry {
	o.name = "retry-manager"
	o.processor = process.New(o.name).After(
		o.OnAfter,
	).Before(
		o.OnBefore,
	).Callback(
		o.OnCall,
	).Panic(o.OnPanic)

	return o
}
