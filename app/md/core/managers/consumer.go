// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package managers

import (
	"context"
	"github.com/fuyibing/gmd/v8/app/md/base"
	"github.com/fuyibing/log/v8"
	"github.com/fuyibing/util/v8/process"
)

type (
	ConsumerManager struct {
		container base.ContainerManager
		processor process.Processor
	}
)

func NewConsumerManager(container base.ContainerManager) *ConsumerManager {
	return (&ConsumerManager{
		container: container,
	}).init()
}

// /////////////////////////////////////////////////////////////
// Exported methods
// /////////////////////////////////////////////////////////////

func (o *ConsumerManager) Processor() process.Processor { return o.processor }

// /////////////////////////////////////////////////////////////
// Processor events
// /////////////////////////////////////////////////////////////

// /////////////////////////////////////////////////////////////
// Processor events
// /////////////////////////////////////////////////////////////

func (o *ConsumerManager) OnAfter(_ context.Context) (ignored bool) {
	// log.Infof("%s: processor stopped", o.processor.Name())
	return
}

func (o *ConsumerManager) OnBefore(_ context.Context) (ignored bool) {
	// log.Infof("%s: start processor", o.processor.Name())
	return
}

func (o *ConsumerManager) OnListen(ctx context.Context) (ignored bool) {
	for {
		select {
		case <-ctx.Done():
			// log.Infof("%s: %v", o.processor.Name(), ctx.Err())
			return
		}
	}
}

func (o *ConsumerManager) OnPanic(ctx context.Context, v interface{}) {
	log.Panicfc(ctx, "%s: %v", o.processor.Name(), v)
}

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *ConsumerManager) init() *ConsumerManager {
	o.processor = process.New("consumer-manager").After(
		o.OnAfter,
	).Before(
		o.OnBefore,
	).Callback(o.OnListen).Panic(o.OnPanic)

	return o
}
