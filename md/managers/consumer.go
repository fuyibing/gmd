// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// author: wsfuyibing <websearch@163.com>
// date: 2023-03-08

package managers

import (
	"context"
	"fmt"
	"github.com/fuyibing/gmd/v8/app"
	"github.com/fuyibing/gmd/v8/md/base"
	"github.com/fuyibing/log/v5"
	"github.com/fuyibing/log/v5/tracers"
	"github.com/fuyibing/util/v8/process"
	"sync"
	"time"
)

type (
	// ConsumerManager
	// 消息者管理.
	ConsumerManager interface {
		// Processor
		// 类进程.
		Processor() process.Processor

		// Reload
		// 重新加载.
		Reload()
	}

	consumer struct {
		sync.RWMutex

		adapterConstructor base.ConsumerConstructor
		adapterKeys        map[string]bool
		adapterUpdated     map[int]int64
		executor           ConsumeExecutor
		name               string
		processor          process.Processor
		re                 chan bool
	}
)

// +---------------------------------------------------------------------------+
// + Interface methods                                                         |
// +---------------------------------------------------------------------------+

func (o *consumer) Processor() process.Processor {
	return o.processor
}

func (o *consumer) Reload() {
	if o.re != nil && o.processor.Healthy() {
		o.re <- true
	}
}

// +---------------------------------------------------------------------------+
// + Event methods                                                             |
// +---------------------------------------------------------------------------+

func (o *consumer) onAdapterCheck(_ context.Context) (ignored bool) {
	if constructor := base.Container.GetConsumer(); constructor != nil {
		o.adapterConstructor = constructor
		return
	}

	log.Error("<%s> adapter not injected into container", o.name)
	return true
}

func (o *consumer) onAdapterFlush(ctx context.Context) (ignored bool) {
	go o.flush(ctx)
	return
}

func (o *consumer) onAfter(ctx context.Context) (ignored bool) {
	if o.executor.Idle() {
		return
	}

	time.Sleep(time.Millisecond * 100)
	return o.onAfter(ctx)
}

func (o *consumer) onListen(ctx context.Context) (ignored bool) {
	rf := time.Duration(app.Config.GetConsumer().GetReloadFrequency()) * time.Second
	ti := time.NewTicker(rf)

	o.re = make(chan bool)

	for {
		select {
		case <-o.re:
			go o.loader()
		case <-ti.C:
			go o.loader()
		case <-ctx.Done():
			return
		}
	}
}

func (o *consumer) onPanic(_ context.Context, v interface{}) {
	log.Fatal("<%s> %v", o.name, v)
}

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *consumer) checkParallel(span tracers.Span, task *base.Task, parallel int, restart bool) (key string) {
	key = fmt.Sprintf("%d.%d", task.Id, parallel)

	// 已经启动.
	if p, exists := o.processor.Get(key); exists {
		if restart {
			p.Restart()
		}

		span.Logger().Info("adapter started already: id=%d, parallel=%d, restart=%v", task.Id, parallel, restart)
		return
	}

	// 首次启动.
	p := o.adapterConstructor(task.Id, parallel, key, o.executor.Do)
	o.processor.Add(p.Processor())

	span.Logger().Info("adapter first start: id=%d, parallel=%d", task.Id, parallel)
	_ = o.processor.StartChild(key)
	return
}

func (o *consumer) checkUpdated(span tracers.Span, task *base.Task) (restart bool) {
	o.Lock()
	defer o.Unlock()

	span.Logger().Info("task loaded: id=%d, updated=%s", task.Id, time.Unix(task.Updated, 0).Format("2006-01-02 15:04:05"))

	// 不需重启.
	if n, ok := o.adapterUpdated[task.Id]; ok && n == task.Updated {
		return false
	}

	// 需要重启.
	o.adapterUpdated[task.Id] = task.Updated
	return true
}

func (o *consumer) flush(ctx context.Context) {
	var (
		keys = make(map[string]bool)
		span = log.NewSpanFromContext(ctx, "consumer.memory.flush")
	)

	defer span.End()

	// 遍历任务.
	for _, task := range base.Memory.GetTasks() {
		restart := o.checkUpdated(span, task)
		// 遍历并行.
		for parallel := 0; parallel < task.Parallels; parallel++ {
			key := o.checkParallel(span, task, parallel, restart)
			keys[key] = true
		}
	}

	// 移除任务.
	for key, _ := range func() map[string]bool {
		o.RLock()
		defer o.RUnlock()
		return o.adapterKeys
	}() {
		if p, exists := o.processor.Get(key); exists {
			if _, ok := keys[key]; ok {
				continue
			}

			span.Logger().Info("adapter remove: key=%s", key)
			p.UnbindWhenStopped(true).Stop()
		}
	}

	// 更新结果.
	o.Lock()
	o.adapterKeys = keys
	o.Unlock()
}

func (o *consumer) init() *consumer {
	o.adapterKeys = make(map[string]bool)
	o.adapterUpdated = make(map[int]int64)

	o.executor = (&executor{}).init()
	o.name = "consumer.manager"
	o.processor = process.New(o.name).
		After(o.onAfter).
		Before(o.onAdapterCheck).
		Callback(o.onAdapterFlush, o.onListen).
		Panic(o.onPanic)
	return o
}

func (o *consumer) loader() {
	span := log.NewSpan("consumer.memory.reload")
	defer span.End()

	// 更新内存.
	if base.Memory.Reload(span.Context()) != nil {
		return
	}

	// 应用内存.
	// 刷新基于内存数据启动的消费者适配器.
	o.flush(span.Context())
}
