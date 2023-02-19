// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package core

import (
	"context"
	"github.com/fuyibing/gmd/v8/app/md/base"
	"github.com/fuyibing/gmd/v8/app/md/core/managers"
	"github.com/fuyibing/log/v8"
	"github.com/fuyibing/util/v8/process"
)

type (
	// BootManager
	// core process manager.
	BootManager interface {
		// Container
		// return container manager.
		Container() base.ContainerManager

		// Prepare
		// boot manager instance.
		Prepare() error

		// Restart
		// boot processor.
		Restart()

		// Start
		// boot processor.
		Start(ctx context.Context) error

		// Stop
		// boot processor.
		Stop()
	}

	boot struct {
		container base.ContainerManager
		processor process.Processor

		consumer *managers.ConsumerManager
		producer *managers.ProducerManager
		remoting *managers.RemotingManager
		retry    *managers.RetryManager
	}
)

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *boot) Container() base.ContainerManager { return o.container }
func (o *boot) Prepare() error                   { return o.prepare() }
func (o *boot) Restart()                         { o.processor.Restart() }
func (o *boot) Start(ctx context.Context) error  { return o.processor.Start(ctx) }
func (o *boot) Stop()                            { o.processor.Stop() }

// /////////////////////////////////////////////////////////////
// Processor events
// /////////////////////////////////////////////////////////////

func (o *boot) OnAfter(_ context.Context) (ignored bool) {
	// log.Infof("%s: processor stopped", o.processor.Name())
	return
}

func (o *boot) OnBefore(_ context.Context) (ignored bool) {
	// log.Infof("%s: start processor", o.processor.Name())
	return
}

func (o *boot) OnListen(ctx context.Context) (ignored bool) {
	for {
		select {
		case <-ctx.Done():
			// log.Infof("%s: %v", o.processor.Name(), ctx.Err())
			return
		}
	}
}

func (o *boot) OnPanic(ctx context.Context, v interface{}) {
	log.Panicfc(ctx, "%s: %v", o.processor.Name(), v)
}

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *boot) init() *boot {
	o.container = (&container{}).init()

	o.consumer = managers.NewConsumerManager(o.container)
	o.producer = managers.NewProducerManager(o.container)
	o.remoting = managers.NewRemotingManager(o.container)
	o.retry = managers.NewRetryManager(o.container)

	// Create process
	// then register event callback.
	o.processor = process.New("boot-manager").Add(
		o.consumer.Processor(),
		o.producer.Processor(),
		o.remoting.Processor(),
		o.retry.Processor(),
	).After(
		o.OnAfter,
	).Before(
		o.OnBefore,
	).Callback(o.OnListen).Panic(o.OnPanic)

	return o
}

func (o *boot) prepare() error {
	// Consumer callable required.
	if _, ok := o.container.GetConsumer(); !ok {
		return base.ErrConsumerCallableNotConfigured
	}

	// Producer callable required.
	if _, ok := o.container.GetProducer(); !ok {
		return base.ErrProducerCallableNotConfigured
	}

	// Remoting callable required.
	if _, ok := o.container.GetRemoting(); !ok {
		return base.ErrRemotingCallableNotConfigured
	}

	// Ready already.
	log.Infof("configured adapter: %s", Config.GetAdapter())
	return nil
}
