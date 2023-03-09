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

const (
	DispatcherHeaderDequeueCount = "X-Gmd-Dequeue-Count"
	DispatcherHeaderMessageId    = "X-Gmd-Message-Id"
	DispatcherHeaderMessageTime  = "X-Gmd-Message-Time"
	DispatcherHeaderSoftware     = "X-Gmd-Software"
	DispatcherHeaderTopicTag     = "X-Gmd-Topic-Tag"
	DispatcherHeaderTopicName    = "X-Gmd-Topic-Name"
)

type (
	// DispatcherConstructor
	// 分发构造器.
	DispatcherConstructor func(addr, method string, timeout int) DispatcherExecutor

	// DispatcherExecutor
	// 分发执行器.
	DispatcherExecutor interface {
		// Dispatch
		// 分发过程.
		Dispatch(task, source *Task, message *Message) (body []byte, err error)

		// Name
		// 执行器名称.
		Name() string
	}
)
