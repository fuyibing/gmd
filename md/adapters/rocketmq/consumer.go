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

package rocketmq

import (
	"context"
	"github.com/fuyibing/gmd/v8/md/base"
	"github.com/fuyibing/log/v5"
	"github.com/fuyibing/util/v8/process"
)

type (
	// Consumer
	// 消费者.
	Consumer struct {
		handler      base.ConsumerHandler
		id, parallel int

		key, name string
		processor process.Processor
	}
)

func NewConsumer(id, parallel int, key string, handler base.ConsumerHandler) base.ConsumerExecutor {
	return (&Consumer{
		handler: handler,
		id:      id, parallel: parallel,
		key: key,
	}).init()
}

// +---------------------------------------------------------------------------+
// + Interface methods                                                         |
// +---------------------------------------------------------------------------+

func (o *Consumer) Processor() process.Processor { return o.processor }

// +---------------------------------------------------------------------------+
// + Event methods                                                             |
// +---------------------------------------------------------------------------+

func (o *Consumer) onAfter(_ context.Context) (ignored bool) {
	return
}

func (o *Consumer) onBefore(_ context.Context) (ignored bool) {
	return
}

func (o *Consumer) onCall(ctx context.Context) (ignored bool) {
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (o *Consumer) onPanic(_ context.Context, v interface{}) {
	log.Fatal("<%s.%s> %v", o.name, o.key, v)
}

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *Consumer) init() *Consumer {
	o.name = "rocketmq.consumer"
	o.processor = process.New(o.key).
		After(o.onAfter).
		Before(o.onBefore).
		Callback(o.onCall).
		Panic(o.onPanic)
	return o
}
