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
	RetryManager struct {
		container base.ContainerManager
		processor process.Processor
	}
)

func NewRetryManager(container base.ContainerManager) *RetryManager {
	return (&RetryManager{container: container}).init()
}

// /////////////////////////////////////////////////////////////
// Exported methods
// /////////////////////////////////////////////////////////////

func (o *RetryManager) Processor() process.Processor { return o.processor }

// /////////////////////////////////////////////////////////////
// Processor events
// /////////////////////////////////////////////////////////////

func (o *RetryManager) OnAfter(_ context.Context) (ignored bool) {
	// log.Infof("%s: processor stopped", o.processor.Name())
	return
}

func (o *RetryManager) OnBefore(_ context.Context) (ignored bool) {
	// log.Infof("%s: start processor", o.processor.Name())
	return
}

func (o *RetryManager) OnListen(ctx context.Context) (ignored bool) {
	for {
		select {
		case <-ctx.Done():
			// log.Infof("%s: %v", o.processor.Name(), ctx.Err())
			return
		}
	}
}

func (o *RetryManager) OnPanic(ctx context.Context, v interface{}) {
	log.Panicfc(ctx, "%s: %v", o.processor.Name(), v)
}

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *RetryManager) init() *RetryManager {
	o.processor = process.New("retry-manager").After(
		o.OnAfter,
	).Before(
		o.OnBefore,
	).Callback(o.OnListen).Panic(o.OnPanic)

	return o
}
