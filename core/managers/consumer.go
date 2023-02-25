// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package managers

import (
	"context"
	"fmt"
	"github.com/fuyibing/gmd/v8/app"
	"github.com/fuyibing/gmd/v8/core/base"
	"github.com/fuyibing/util/v8/process"
	"sync"
	"time"
)

type Consumer struct {
	callable  base.ConsumerCallable
	name      string
	processor process.Processor

	mu           *sync.RWMutex
	subprocesses map[string]int
	updates      map[int]int64
}

// /////////////////////////////////////////////////////////////
// Exported methods
// /////////////////////////////////////////////////////////////

func (o *Consumer) DoConsume(_ *base.Task, _ *base.Message) (retry bool, err error) {
	return
}

func (o *Consumer) DoMemoryReload() {
	var (
		c   = context.Background() // log.NewContextInfo("{%s} begin load memory", o.name)
		err error
	)

	defer func() {
		if err != nil {
			// log.Errorfc(c, "{%s} memory load error: %v", o.name, err)
		} else {
			// log.Infofc(c, "{%s} memory load finish", o.name)
		}
	}()

	if err = base.Memory.Reload(c); err != nil {
		return
	}

	o.DoSubprocess()
	return
}

func (o *Consumer) DoSubprocess() {
	var mapper = make(map[string]int)

	for _, task := range base.Memory.GetTasks() {
		for k, i := range o.loadParallel(task) {
			mapper[k] = i
		}
	}

	if rm := func(m map[string]int) (r map[string]int) {
		o.mu.Lock()
		defer o.mu.Unlock()

		r = make(map[string]int)
		for k, i := range o.subprocesses {
			if _, ok := m[k]; ok {
				continue
			}
			r[k] = i
		}

		o.subprocesses = m
		return
	}(mapper); len(rm) > 0 {
		o.unloadSubprocess(rm)
	}
}

// /////////////////////////////////////////////////////////////
// Processor events
// /////////////////////////////////////////////////////////////

func (o *Consumer) OnAfter(_ context.Context) (ignored bool) {
	return
}

func (o *Consumer) OnBefore(_ context.Context) (ignored bool) {
	return
}

func (o *Consumer) OnBeforeCallable(_ context.Context) (ignored bool) {
	// Return
	// if not first start.
	if o.callable != nil {
		return
	}

	// Return true
	// if constructor not injected into container.
	if o.callable = Container.GetConsumer(); o.callable == nil {
		// log.Errorf("consumer constructor for {%s} adapter not injected", app.Config.GetAdapter())
		return true
	}

	return
}

func (o *Consumer) OnCallChannel(ctx context.Context) (ignored bool) {
	// Register ticker
	// for memory update frequency.
	tick := time.NewTicker(time.Duration(app.Config.GetMemoryReloadSeconds()) * time.Second)
	// log.Infof("{%s} register ticker: type=memory, seconds=%d", o.name, app.Config.GetMemoryReloadSeconds())

	// Listen
	// channel signal.
	for {
		select {
		case <-tick.C:
			go o.DoMemoryReload()
		case <-ctx.Done():
			tick.Stop()
			return
		}
	}
}

func (o *Consumer) OnCallSubprocessLoad(_ context.Context) (ignored bool) {
	o.DoSubprocess()
	return
}

func (o *Consumer) OnPanic(ctx context.Context, v interface{}) {
	// log.Panicfc(ctx, "processor {%s} fatal: %v", o.name, v)
}

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *Consumer) init() *Consumer {
	o.mu = &sync.RWMutex{}
	o.subprocesses = make(map[string]int)
	o.updates = make(map[int]int64)

	o.name = "consumer-manager"
	o.processor = process.New(o.name).After(
		o.OnAfter,
	).Before(
		o.OnBefore,
		o.OnBeforeCallable,
	).Callback(
		o.OnCallSubprocessLoad,
		o.OnCallChannel,
	).Panic(o.OnPanic)

	return o
}

func (o *Consumer) loadParallel(task *base.Task) (mapper map[string]int) {
	mapper = make(map[string]int)

	// Check restart status.
	//
	// Return false if task updated timestamp equal to previous,
	// otherwise false returned.
	restart := func() bool {
		o.mu.Lock()
		defer o.mu.Unlock()

		// Equal to previous.
		if u, ok := o.updates[task.Id]; ok && task.Updated == u {
			return false
		}

		// First load or task property changed.
		o.updates[task.Id] = task.Updated
		return true
	}()

	// Range by parallels.
	for parallel := 0; parallel < task.Parallels; parallel++ {
		k := o.loadSubprocess(task, parallel, restart)
		mapper[k] = task.Id
	}
	return
}

func (o *Consumer) loadSubprocess(task *base.Task, parallel int, restart bool) (name string) {
	name = fmt.Sprintf("%s-consumer-%d-%d", app.Config.GetAdapter(), task.Id, parallel)

	// Subprocess exists.
	if sp, exists := o.processor.Get(name); exists {
		// Skip if task property never changed.
		if !restart {
			return
		}

		// Start
		// if subprocess never started.
		if sp.Stopped() {
			_ = o.processor.StartChild(name)
			return
		}

		// Restart if subprocess is healthy.
		sp.Restart()
		return
	}

	// Create and start subprocess.
	sp := o.callable(task.Id, parallel, name, o.DoConsume)
	o.processor.Add(sp.Processor().UnbindWhenStopped(true))
	_ = o.processor.StartChild(name)
	return
}

func (o *Consumer) unloadSubprocess(mapper map[string]int) {
	for name := range mapper {
		if sp, exists := o.processor.Get(name); exists {
			sp.Stop()
		}
	}
}
