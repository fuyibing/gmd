// author: wsfuyibing <websearch@163.com>
// date: 2023-02-20

package managers

import (
	"context"
	"github.com/fuyibing/gmd/v8/core/base"
	"github.com/fuyibing/log/v5"
	"github.com/fuyibing/util/v8/process"
	"sync"
)

type (
	// BootManager
	// 入口管理器.
	BootManager interface {
		Restart()
		Start(ctx context.Context) error
		Stop()
		Stopped() bool

		Producer() *Producer
	}

	boot struct {
		consumer *Consumer
		producer *Producer
		remoting *Remoting
		retry    *Retry

		name      string
		processor process.Processor
	}
)

// /////////////////////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////////////////////

func (o *boot) Restart()                        { o.processor.Restart() }
func (o *boot) Start(ctx context.Context) error { return o.processor.Start(ctx) }
func (o *boot) Stop()                           { o.processor.Stop() }
func (o *boot) Stopped() bool                   { return o.processor.Stopped() }

func (o *boot) Producer() *Producer { return o.producer }

// /////////////////////////////////////////////////////////////////////////////
// Event methods
// /////////////////////////////////////////////////////////////////////////////

func (o *boot) OnBeforeMemory(ctx context.Context) (ignored bool) {
	var (
		span = log.NewSpanFromContext(ctx, "memory")
	)

	defer span.End()

	if err := base.Memory.Reload(span.Context()); err != nil {
		span.Logger().Error("%v", err)
		return true
	}
	return
}

func (o *boot) OnCall(ctx context.Context) (ignored bool) {
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (o *boot) OnPanic(ctx context.Context, v interface{}) {
	if spa, exists := log.Span(ctx); exists {
		spa.Logger().Fatal("<%s> %v", o.name, v)
	} else {
		log.Fatal("<%s> %v", o.name, v)
	}
}

// /////////////////////////////////////////////////////////////////////////////
// Access and constructor
// /////////////////////////////////////////////////////////////////////////////

func (o *boot) init() *boot {
	o.consumer = (&Consumer{}).init()
	o.producer = (&Producer{}).init()
	o.remoting = (&Remoting{}).init()
	o.retry = (&Retry{}).init()

	o.name = "boot-manager"
	o.processor = process.New(o.name).
		Before(o.OnBeforeMemory).
		Callback(o.OnCall).
		Panic(o.OnPanic)

	o.processor.Add(
		// o.consumer.processor,
		o.producer.processor,
		// o.remoting.processor,
		// o.retry.processor,
	)

	return o
}

var (
	Boot      BootManager
	Container base.ContainerManager
)

func init() {
	new(sync.Once).Do(func() {
		Container = (&container{}).init()
		Boot = (&boot{}).init()
	})
}
