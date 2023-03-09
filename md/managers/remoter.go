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
	"github.com/fuyibing/log/v5"
	"github.com/fuyibing/util/v8/process"
)

type (
	// RemoterManager
	// 服务端管理器.
	RemoterManager interface {
		// Processor
		// 类进程.
		Processor() process.Processor
	}

	remoter struct {
		name      string
		processor process.Processor
	}
)

func (o *remoter) Processor() process.Processor { return o.processor }

// +---------------------------------------------------------------------------+
// + Event methods                                                             |
// +---------------------------------------------------------------------------+

func (o *remoter) onAfter(_ context.Context) (ignored bool) {
	return
}

func (o *remoter) onBefore(_ context.Context) (ignored bool) {
	return
}

func (o *remoter) onCall(ctx context.Context) (ignored bool) {
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (o *remoter) onPanic(_ context.Context, v interface{}) {
	log.Fatal("<%s> %v", o.name, v)
}

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *remoter) init() *remoter {
	o.name = "remoter.manager"
	o.processor = process.New(o.name).
		After(o.onAfter).
		Before(o.onBefore).
		Callback(o.onCall).
		Panic(o.onPanic)

	return o
}
