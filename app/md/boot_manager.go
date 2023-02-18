// author: wsfuyibing <websearch@163.com>
// date: 2023-02-08

package md

import (
	"context"
	"github.com/fuyibing/gmd/app/md/base"
	"github.com/fuyibing/gmd/app/md/conf"
	"github.com/fuyibing/log/v8"
	"github.com/fuyibing/util/v8/process"
	"time"
)

var Boot BootManager

type (
	BootManager interface {
		// Consumer
		// return consumer manager interface.
		//
		//   x := md.Boot.Consumer()
		//   x.Container().Flush(ctx)
		Consumer() ConsumerManager

		// Processor
		// return boot processor interface.
		//
		//   x := md.Boot.Processor()
		//   x.Start(ctx)
		Processor() process.Processor

		// Producer
		// return producer manager interface.
		//
		//   x := md.Boot.Producer()
		//   x.Publish(payload)
		Producer() ProducerManager

		// Retry
		// return retry manager interface.
		//
		//   x := md.Boot.Retry()
		//   x.Publish()
		Retry() RetryManager

		// Remoter
		// return remoter manager interface.
		//
		//   x := md.Boot.Remoter()
		//   x.Build(task)
		Remoter() RemoterManager
	}

	boot struct {
		consumer ConsumerManager
		producer ProducerManager
		retry    RetryManager
		remoter  RemoterManager

		children  []process.Processor
		processor process.Processor
	}
)

// /////////////////////////////////////////////////////////////
// Interface methods.
// /////////////////////////////////////////////////////////////

func (o *boot) Consumer() ConsumerManager    { return o.consumer }
func (o *boot) Processor() process.Processor { return o.processor }
func (o *boot) Producer() ProducerManager    { return o.producer }
func (o *boot) Retry() RetryManager          { return o.retry }
func (o *boot) Remoter() RemoterManager      { return o.remoter }

// /////////////////////////////////////////////////////////////
// Event methods.
// /////////////////////////////////////////////////////////////

// OnAfter
// called when processor stopped.
func (o *boot) OnAfter(_ context.Context) (ignored bool) {
	log.Debugf("boot manager: processor stopped")
	return
}

// OnBefore
// called when processor start.
func (o *boot) OnBefore(_ context.Context) (ignored bool) {
	log.Debugf("boot manager: start processor")
	return
}

// OnBeforeLoadMemory
// call memory manager reload.
func (o *boot) OnBeforeLoadMemory(_ context.Context) (ignored bool) {
	// Load memory for the first time.
	// The consumer and producer manager needs this data.
	log.Debugf("boot manager: load memory for the first time")
	if err := base.Memory.Reload(); err != nil {
		log.Errorf("boot manager: memory load failed, error=%v", err)
		return true
	}

	// Next
	// event callee.
	return
}

// OnCallChannel
// listen channel signal.
func (o *boot) OnCallChannel(ctx context.Context) (ignored bool) {
	log.Debugf("boot manager: listen channel signal")

	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

// OnCallChildStart
// start child processors in coroutine.
func (o *boot) OnCallChildStart(ctx context.Context) (ignored bool) {
	log.Debugf("boot manager: start %d child processors in coroutine", len(o.children))

	// Start
	// child processors in coroutine.
	for _, x := range o.children {
		go func(c context.Context, p process.Processor) {
			_ = p.Start(c)
		}(ctx, x)
	}

	// Next
	// event callee.
	return
}

// OnCallChildStopped
// wait until child processors stopped.
func (o *boot) OnCallChildStopped(ctx context.Context) (ignored bool) {
	// Recall
	// after specified millisecond if any child processor not stopped.
	for _, x := range o.children {
		if !x.Stopped() {
			time.Sleep(conf.EventSleepDuration)
			return o.OnCallChildStopped(ctx)
		}
	}

	// Next
	// event callee.
	log.Debugf("boot manager: %d child processors stopped", len(o.children))
	return
}

// OnPanic
// called with panic at runtime.
func (o *boot) OnPanic(ctx context.Context, v interface{}) {
	log.Panicfc(ctx, "boot manager: %v", v)
}

// /////////////////////////////////////////////////////////////
// Constructor methods.
// /////////////////////////////////////////////////////////////

func (o *boot) init() *boot {
	// Prepare child managers.
	o.consumer = (&consumer{}).init()
	o.producer = (&producer{}).init()
	o.retry = (&retry{}).init()
	o.remoter = (&remoter{}).init()

	// Initialize child processors.
	o.children = []process.Processor{
		o.consumer.Processor(),
		o.producer.Processor(),
		o.retry.Processor(),
		o.remoter.Processor(),
	}

	// Register boot processor event callbacks.
	o.processor = process.New("boot manager").After(
		o.OnAfter,
	).Before(
		o.OnBefore,
		o.OnBeforeLoadMemory,
	).Callback(
		o.OnCallChildStart,
		o.OnCallChannel,
		o.OnCallChildStopped,
	).Panic(o.OnPanic)

	return o
}
