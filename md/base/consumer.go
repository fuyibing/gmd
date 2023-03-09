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
// date: 2023-03-07

package base

import (
	"github.com/fuyibing/util/v8/process"
)

type (
	// ConsumerConstructor
	// 消费者构造器.
	ConsumerConstructor func(id, parallel int, name string, handler ConsumerHandler) ConsumerExecutor

	// ConsumerExecutor
	// 消费者执行器.
	ConsumerExecutor interface {
		// Processor
		// 类进程.
		Processor() process.Processor
	}

	// ConsumerHandler
	// 消费过程.
	ConsumerHandler func(task *Task, message *Message) (retry bool, err error)
)
