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
	"github.com/fuyibing/gmd/v8/app"
	"github.com/fuyibing/log/v5"
	"github.com/fuyibing/util/v8/process"
	"sync"
	"time"
)

type (
	// RetryManager
	// 重试管理器.
	RetryManager interface {
		// Processor
		// 类进程.
		Processor() process.Processor
	}

	retry struct {
		sync.RWMutex

		name      string
		processor process.Processor

		consuming,
		publishing bool
	}
)

func (o *retry) Processor() process.Processor { return o.processor }

// +---------------------------------------------------------------------------+
// + Event methods                                                             |
// +---------------------------------------------------------------------------+

// consumeMessage
// 消费消息.
func (o *retry) consumeMessage() {
}

// publishPayload
// 发布消息.
func (o *retry) publishPayload() {
}

func (o *retry) onCall(ctx context.Context) (ignored bool) {
	dc := time.Duration(app.Config.GetConsumer().GetRetryFrequency()) * time.Second
	tc := time.NewTimer(dc)

	dp := time.Duration(app.Config.GetProducer().GetRetryFrequency()) * time.Second
	tp := time.NewTicker(dp)

	for {
		select {
		case <-tc.C:
			go o.consumeMessage()
		case <-tp.C:
			go o.publishPayload()
		case <-ctx.Done():
			return
		}
	}
}

func (o *retry) onPanic(_ context.Context, v interface{}) {
	log.Fatal("<%s> %v", o.name, v)
}

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *retry) init() *retry {
	o.name = "retry.manager"
	o.processor = process.New(o.name).
		Callback(o.onCall).
		Panic(o.onPanic)

	return o
}
