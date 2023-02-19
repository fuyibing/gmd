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
	RemotingManager struct {
		adapter   base.RemotingManager
		container base.ContainerManager
		processor process.Processor
	}
)

func NewRemotingManager(container base.ContainerManager) *RemotingManager {
	return (&RemotingManager{container: container}).init()
}

// /////////////////////////////////////////////////////////////
// Exported methods
// /////////////////////////////////////////////////////////////

func (o *RemotingManager) Processor() process.Processor { return o.processor }

// /////////////////////////////////////////////////////////////
// Processor events
// /////////////////////////////////////////////////////////////

func (o *RemotingManager) OnAfter(_ context.Context) (ignored bool) {
	// log.Infof("%s: processor stopped", o.processor.Name())
	return
}

func (o *RemotingManager) OnBefore(_ context.Context) (ignored bool) {
	// log.Infof("%s: start processor", o.processor.Name())
	return
}

func (o *RemotingManager) OnListen(ctx context.Context) (ignored bool) {
	for {
		select {
		case <-ctx.Done():
			// log.Infof("%s: %v", o.processor.Name(), ctx.Err())
			return
		}
	}
}

func (o *RemotingManager) OnPanic(ctx context.Context, v interface{}) {
	log.Panicfc(ctx, "%s: %v", o.processor.Name(), v)
}

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *RemotingManager) init() *RemotingManager {
	// Create process
	// then register event callback.
	o.processor = process.New("remoting-manager").After(
		o.OnAfter,
	).Before(
		o.OnBefore,
	).Callback(o.OnListen).Panic(o.OnPanic)

	// Create adapter as child, then add child to this processor.
	if call, exists := o.container.GetRemoting(); exists {
		o.adapter = call()
		o.processor.Add(o.adapter.Processor())
	}

	return o
}
