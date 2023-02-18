// author: wsfuyibing <websearch@163.com>
// date: 2023-02-08

package md

import (
	"context"
	"github.com/fuyibing/gmd/app/md/conf"
	"github.com/fuyibing/log/v8"
	"github.com/fuyibing/util/v8/process"
	"time"
)

type (
	ConsumerManager interface {
		// Container
		// return consumer container interface.
		Container() ConsumerContainer

		// Processor
		// return consumer processor interface.
		Processor() process.Processor

		// Reload
		// send consumer reload process.
		Reload()
	}

	consumer struct {
		container ConsumerContainer
		processor process.Processor
		reload    chan bool
	}
)

// /////////////////////////////////////////////////////////////
// Interface methods.
// /////////////////////////////////////////////////////////////

func (o *consumer) Container() ConsumerContainer { return o.container }
func (o *consumer) Processor() process.Processor { return o.processor }
func (o *consumer) Reload()                      { o.sendReload() }

// /////////////////////////////////////////////////////////////
// Event methods.
// /////////////////////////////////////////////////////////////

// OnAfter
// called when processor stopped.
func (o *consumer) OnAfter(_ context.Context) (ignored bool) {
	log.Debugf("consumer manager: processor stopped")
	return
}

// OnBefore
// called when processor start.
func (o *consumer) OnBefore(_ context.Context) (ignored bool) {
	log.Debugf("consumer manager: start processor")
	return
}

// OnCallChannel
// listen channel signal.
func (o *consumer) OnCallChannel(ctx context.Context) (ignored bool) {
	log.Debugf("consumer manager: listen channel signal")

	// Create ticker
	// for pop coroutine concurrency.
	re := time.NewTicker(time.Duration(conf.Config.Consumer.ReloadSeconds) * time.Second)
	o.reload = make(chan bool)

	// Stop and unset
	// ticker when end.
	defer func() {
		close(o.reload)
		o.reload = nil

		re.Stop()
		re = nil
	}()

	// Select
	// channel messages.
	for {
		select {
		case <-re.C:
			go func() { _ = o.container.Reload(ctx) }()
		case <-o.reload:
			go func() { _ = o.container.Reload(ctx) }()
		case <-ctx.Done():
			return
		}
	}
}

// OnCallAdapterStart
// call consumer container and flush adapters.
func (o *consumer) OnCallAdapterStart(ctx context.Context) (ignored bool) {
	log.Debugf("consumer manager: call consumer container and flush adapters")
	o.container.Flush(ctx)
	return
}

// OnCallAdapterStopped
// wait until consumer container is empty.
func (o *consumer) OnCallAdapterStopped(ctx context.Context) (ignored bool) {
	// Recall after specified milliseconds
	// if consumer container is not empty.
	if !o.container.IsEmpty() {
		time.Sleep(conf.EventSleepDuration)
		return o.OnCallAdapterStopped(ctx)
	}

	// Next
	// event callee.
	log.Debugf("consumer manager: consumer container cleaned")
	return
}

// OnCallWorkerIdle
// wait until consumer worker is idle.
func (o *consumer) OnCallWorkerIdle(ctx context.Context) (ignored bool) {
	// Recall after specified milliseconds
	// if consumer worker is not idle.
	if !o.container.Worker().IsIdle() {
		time.Sleep(conf.EventSleepDuration)
		return o.OnCallWorkerIdle(ctx)
	}

	// Next
	// event callee.
	log.Debugf("consumer manager: consumer worker completed")
	return
}

// OnPanic
// called with panic at runtime.
func (o *consumer) OnPanic(ctx context.Context, v interface{}) {
	log.Panicfc(ctx, "consumer manager: %v", v)
}

func (o *consumer) sendReload() {
	if o.reload != nil {
		o.reload <- true
	}
}

// /////////////////////////////////////////////////////////////
// Access methods.
// /////////////////////////////////////////////////////////////

func (o *consumer) init() *consumer {
	// Prepare consumer container.
	o.container = (&container{}).init()

	// Register consumer processor event callbacks.
	o.processor = process.New("consumer manager").After(
		o.OnAfter,
	).Before(
		o.OnBefore,
	).Callback(
		o.OnCallAdapterStart,
		o.OnCallChannel,
		o.OnCallAdapterStopped,
		o.OnCallWorkerIdle,
	).Panic(o.OnPanic)

	return o
}
