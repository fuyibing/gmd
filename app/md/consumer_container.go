// author: wsfuyibing <websearch@163.com>
// date: 2023-02-11

package md

import (
	"context"
	"fmt"
	"github.com/fuyibing/gmd/app/md/adapters"
	"github.com/fuyibing/gmd/app/md/base"
	"github.com/fuyibing/gmd/app/md/conf"
	"github.com/fuyibing/log/v8"
	"sync"
)

type (
	ConsumerContainer interface {
		// Flush
		// start adapter if not started. Stop adapter if task
		// disabled or deleted or parallels changed down.
		Flush(ctx context.Context)

		// IsEmpty
		// return container has adapter or not.
		//
		// Return true means no adapter in container, otherwise false
		// returned.
		IsEmpty() bool

		// Reload
		// call memory manager then flush adapters.
		Reload(ctx context.Context) error

		// Worker
		// return consumer worker interface.
		Worker() ConsumerWorker
	}

	container struct {
		adapters map[string]adapters.ConsumerAdapter
		mu       *sync.RWMutex
		updated  map[int]int64
		worker   ConsumerWorker
	}
)

// /////////////////////////////////////////////////////////////
// Interface methods.
// /////////////////////////////////////////////////////////////

func (o *container) Flush(ctx context.Context)              { o.doFlush(ctx) }
func (o *container) IsEmpty() bool                          { return o.isEmpty() }
func (o *container) Reload(ctx context.Context) (err error) { return o.doReload(ctx) }
func (o *container) Worker() ConsumerWorker                 { return o.worker }

// /////////////////////////////////////////////////////////////
// Action methods.
// /////////////////////////////////////////////////////////////

func (o *container) doClean(mapper map[string]int) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	for key, adapter := range o.adapters {
		if _, ok := mapper[key]; ok {
			continue
		}

		log.Debugf("consumer container: remove %s adapter, name=%s", conf.Config.Adapter, adapter.Processor().Name())
		adapter.Processor().Stop()
	}
}

func (o *container) doFlush(ctx context.Context) {
	var (
		mapper = make(map[string]int)
	)

	log.Debugf("consumer container: flush %s adapters", conf.Config.Adapter)

	// Iterate
	// tasks in memory.
	for _, task := range base.Memory.GetTasks() {
		restart := o.doUpdated(task)
		if ks := o.doTask(ctx, task, restart); len(ks) > 0 {
			for _, k := range ks {
				mapper[k] = task.Id
			}
		}
	}

	// Call
	// clean method.
	o.doClean(mapper)
}

func (o *container) doReload(ctx context.Context) (err error) {
	if err = base.Memory.Reload(); err == nil {
		o.doFlush(ctx)
	}
	return
}

func (o *container) doTask(ctx context.Context, task *base.Task, restart bool) (keys []string) {
	keys = make([]string, 0)
	for parallel := 0; parallel < task.Parallels; parallel++ {
		if key := o.doTaskParallel(ctx, task, parallel, restart); key != "" {
			keys = append(keys, key)
		}
	}
	return
}

func (o *container) doTaskParallel(ctx context.Context, task *base.Task, parallel int, restart bool) string {
	o.mu.Lock()
	defer o.mu.Unlock()

	var (
		adapter adapters.ConsumerAdapter
		err     error
		key     = fmt.Sprintf("%d:%d", task.Id, parallel)
		ok      bool
	)

	// Return key
	// if adapter mapped. Call restart method if task properties
	// changed.
	if adapter, ok = o.adapters[key]; ok {
		if restart {
			adapter.Processor().Restart()
		}
		return key
	}

	// Return error
	// if build adapter failed.
	if adapter, err = adapters.NewConsumer(conf.Config.Adapter, task.Id, parallel); err != nil {
		log.Errorf("consumer container: build %s adapter, error=%v", conf.Config.Adapter, err)
		return ""
	}

	// Bind
	// dispatcher callback on adapter. Called when queue message
	// received.
	adapter.Dispatcher(o.worker.Do)

	// Call
	// processor start method.
	go o.doTaskProcessor(ctx, adapter, key)

	// Set
	// adapter mapping.
	o.adapters[key] = adapter
	return key
}

func (o *container) doTaskProcessor(ctx context.Context, adapter adapters.ConsumerAdapter, key string) {
	// Remove
	// adapter processor mapping when end.
	defer func() {
		o.mu.Lock()
		defer o.mu.Unlock()
		delete(o.adapters, key)
	}()

	// Start
	// adapter processor.
	_ = adapter.Processor().Start(ctx)
}

func (o *container) doUpdated(task *base.Task) (restart bool) {
	o.mu.Lock()
	defer o.mu.Unlock()

	// Return false
	// if adapter started and task properties not changed.
	if u, ok := o.updated[task.Id]; ok && u == task.Updated {
		return false
	}

	// Return true
	// and save last updated timestamp of task.
	o.updated[task.Id] = task.Updated
	return true
}

func (o *container) isEmpty() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return len(o.adapters) == 0
}

// /////////////////////////////////////////////////////////////
// Access methods.
// /////////////////////////////////////////////////////////////

func (o *container) init() *container {
	o.adapters = make(map[string]adapters.ConsumerAdapter)
	o.mu = &sync.RWMutex{}
	o.updated = make(map[int]int64)
	o.worker = (&worker{}).init()
	return o
}
