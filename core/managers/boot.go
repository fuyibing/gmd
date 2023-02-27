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
		// Restart
		// 重启Processor.
		Restart()

		// Start
		// 启动Processor.
		Start(ctx context.Context) error

		// Stop
		// 退出Processor.
		Stop()

		// Stopped
		// 退出状态.
		Stopped() bool
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

func (o *boot) Restart()                        { o.processor.Restart() }
func (o *boot) Start(ctx context.Context) error { return o.processor.Start(ctx) }
func (o *boot) Stop()                           { o.processor.Stop() }
func (o *boot) Stopped() bool                   { return o.processor.Stopped() }

func (o *boot) OnAfter(_ context.Context) (ignored bool) {
	return
}

func (o *boot) OnBefore(_ context.Context) (ignored bool) {
	return
}

func (o *boot) OnBeforeMemory(ctx context.Context) (ignored bool) {
	var (
		trace = log.NewTraceFromContext(ctx, "memory")
		span  = trace.NewSpan("boot.memory.update")
	)

	defer span.End()

	if err := base.Memory.Reload(span.GetContext()); err != nil {
		span.Logger().Error("%v", err)
		return true
	}

	span.Logger().Info("succeed")
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
	if sp, spe := log.Span(ctx); spe {
		sp.Logger().Fatal("processor {%s} fatal: %v", o.name, v)
	} else {
		log.Fatal("processor {%s} fatal: %v", o.name, v)
	}
}

func (o *boot) init() *boot {
	o.consumer = (&Consumer{}).init()
	o.producer = (&Producer{}).init()
	o.remoting = (&Remoting{}).init()
	o.retry = (&Retry{}).init()

	o.name = "boot-manager"
	o.processor = process.New(o.name).After(
		o.OnAfter,
	).Before(
		o.OnBefore,
		o.OnBeforeMemory,
	).Callback(
		o.OnCall,
	).Panic(o.OnPanic)

	o.processor.Add(
		o.consumer.processor,
		// o.producer.processor,
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
