// author: wsfuyibing <websearch@163.com>
// date: 2023-02-03

package aliyunmns

import (
	"context"
	"github.com/fuyibing/gmd/app/md/base"
	"github.com/fuyibing/log/v3"
	"github.com/fuyibing/util/v2/process"
	"sync"
)

type (
	// Remoter
	// struct for aliyun mns remoter.
	Remoter struct {
		mu        *sync.RWMutex
		name      string
		processor process.Processor
	}
)

func NewRemoter() *Remoter {
	o := (&Remoter{}).init()
	return o
}

func (o *Remoter) Build(ctx context.Context, task *base.Task) (err error) {
	return Agent.Build(ctx, task)
}

func (o *Remoter) BuildById(ctx context.Context, id int) (err error) {
	var task *base.Task
	if task, err = base.Memory.GetTaskFromBean(ctx, id); err != nil {
		return
	}
	return o.Build(ctx, task)
}

func (o *Remoter) Destroy(ctx context.Context, task *base.Task) (err error) {
	return Agent.Destroy(ctx, task)
}

func (o *Remoter) DestroyById(ctx context.Context, id int) (err error) {
	var task *base.Task
	if task, err = base.Memory.GetTaskFromBean(ctx, id); err != nil {
		return
	}
	return o.Destroy(ctx, task)
}

func (o *Remoter) Processor() process.Processor { return o.processor }

// /////////////////////////////////////////////////////////////
// Event methods
// /////////////////////////////////////////////////////////////

func (o *Remoter) onAfter(ctx context.Context) (ignored bool) {
	log.Infof("%s: processor stopped", o.name)
	return
}

func (o *Remoter) onBefore(_ context.Context) (ignored bool) {
	log.Infof("%s: start processor", o.name)
	return
}

func (o *Remoter) onCall(ctx context.Context) (ignored bool) {
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (o *Remoter) onCallAfter(_ context.Context) (ignored bool) {
	return
}

func (o *Remoter) onCallBefore(_ context.Context) (ignored bool) {
	return
}

func (o *Remoter) onPanic(ctx context.Context, v interface{}) {
	log.Panicfc(ctx, "%s: %v", o.name, v)
}

// /////////////////////////////////////////////////////////////
// Construct method
// /////////////////////////////////////////////////////////////

func (o *Remoter) init() *Remoter {
	// Create
	// processor instance.
	o.name = "aliyunmns-remoter"
	o.processor = process.New(o.name).After(
		o.onAfter,
	).Before(
		o.onBefore,
	).Callback(
		o.onCallBefore,
		o.onCall,
		o.onCallAfter,
	).Panic(o.onPanic)

	return o
}
