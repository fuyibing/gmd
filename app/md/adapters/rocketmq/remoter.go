// author: wsfuyibing <websearch@163.com>
// date: 2023-02-03

package rocketmq

import (
	"context"
	"fmt"
	"github.com/fuyibing/gmd/app/md/base"
	"github.com/fuyibing/log/v8"
	"github.com/fuyibing/util/v8/process"
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

func (o *Remoter) Build(_ context.Context, _ *base.Task) (err error) { return fmt.Errorf("undefined") }
func (o *Remoter) BuildById(_ context.Context, _ int) (err error)    { return fmt.Errorf("undefined") }
func (o *Remoter) Destroy(_ context.Context, _ *base.Task) (err error) {
	return fmt.Errorf("undefined")
}
func (o *Remoter) DestroyById(_ context.Context, _ int) (err error) { return fmt.Errorf("undefined") }
func (o *Remoter) Processor() process.Processor                     { return o.processor }

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

func (o *Remoter) onCaller(ctx context.Context) (ignored bool) {
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (o *Remoter) onCallerAfter(_ context.Context) (ignored bool) {
	return
}

func (o *Remoter) onCallerBefore(_ context.Context) (ignored bool) {
	return
}

func (o *Remoter) onPanic(ctx context.Context, v interface{}) {
	log.Panicfc(ctx, "%s: %v", o.name, v)
}

// /////////////////////////////////////////////////////////////
// Construct method
// /////////////////////////////////////////////////////////////

func (o *Remoter) init() *Remoter {
	o.name = "rocketmq-remoter"
	o.processor = process.New(o.name).After(
		o.onAfter,
	).Before(
		o.onBefore,
	).Callback(
		o.onCallerBefore,
		o.onCaller,
		o.onCallerAfter,
	).Panic(o.onPanic)

	return o
}
