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
	"github.com/fuyibing/gmd/v8/md/base"
	"github.com/fuyibing/log/v5"
	"github.com/fuyibing/util/v8/process"
	"sync"
)

var (
	// Boot
	// 根管理器.
	Boot BootManager
)

type (
	// BootManager
	// 根管理器接口.
	BootManager interface {
		// Consumer
		// 消费者管理器.
		Consumer() ConsumerManager

		// Processor
		// 类进程.
		Processor() process.Processor

		// Producer
		// 生产者管理器.
		Producer() ProducerManager

		// Remoter
		// 服务端管理器.
		Remoter() RemoterManager

		// Retry
		// 重试管理器.
		Retry() RetryManager

		// Start
		// 启动管理器.
		Start(ctx context.Context) error

		// Stop
		// 退出管理器.
		Stop()
	}

	boot struct {
		consumer  ConsumerManager
		name      string
		processor process.Processor
		producer  ProducerManager
		retry     RetryManager
		remoter   RemoterManager
	}
)

// +---------------------------------------------------------------------------+
// + Interface methods                                                         |
// +---------------------------------------------------------------------------+

func (o *boot) Consumer() ConsumerManager       { return o.consumer }
func (o *boot) Processor() process.Processor    { return o.processor }
func (o *boot) Producer() ProducerManager       { return o.producer }
func (o *boot) Remoter() RemoterManager         { return o.remoter }
func (o *boot) Retry() RetryManager             { return o.retry }
func (o *boot) Start(ctx context.Context) error { return o.processor.Start(ctx) }
func (o *boot) Stop()                           { o.processor.Stop() }

// +---------------------------------------------------------------------------+
// + Event methods                                                             |
// +---------------------------------------------------------------------------+

func (o *boot) onBefore(ctx context.Context) (ignored bool) {
	span := log.NewSpanFromContext(ctx, "boot.memory.init")
	defer span.End()
	err := base.Memory.Reload(span.Context())
	return err != nil
}

func (o *boot) onCall(ctx context.Context) (ignored bool) {
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (o *boot) onPanic(_ context.Context, v interface{}) {
	log.Fatal("<%s> %v", o.name, v)
}

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *boot) init() *boot {
	o.consumer = (&consumer{}).init()
	o.producer = (&producer{}).init()
	o.remoter = (&remoter{}).init()
	o.retry = (&retry{}).init()

	o.name = "boot.manager"
	o.processor = process.New(o.name).Add(
		o.consumer.Processor(),
		o.producer.Processor(),
		o.remoter.Processor(),
		o.retry.Processor(),
	).Before(o.onBefore).Callback(o.onCall).Panic(
		o.onPanic,
	)

	return o
}

func init() { new(sync.Once).Do(func() { Boot = (&boot{}).init() }) }
