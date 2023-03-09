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
	// Remoter
	// 服务端.
	Remoter struct {
		name       string
		processor  process.Processor
		processing int32
	}
)

func NewRemoter() (remoter base.RemoterExecutor) {
	return (&Remoter{}).init()
}

// +---------------------------------------------------------------------------+
// + Interface methods                                                         |
// +---------------------------------------------------------------------------+

func (o *Remoter) Processor() process.Processor {
	return o.processor
}

// +---------------------------------------------------------------------------+
// + Event methods                                                             |
// +---------------------------------------------------------------------------+

func (o *Remoter) onAfter(_ context.Context) (ignored bool) {
	return
}

func (o *Remoter) onBefore(_ context.Context) (ignored bool) {
	return
}

func (o *Remoter) onListen(ctx context.Context) (ignored bool) {
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (o *Remoter) onPanic(_ context.Context, v interface{}) {
	log.Fatal("<%s> %v", o.name, v)
}

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *Remoter) init() *Remoter {
	o.name = "rocketmq.remoter"
	o.processor = process.New(o.name).
		After(o.onAfter).
		Before(o.onBefore).
		Callback(o.onListen).
		Panic(o.onPanic)
	return o
}
