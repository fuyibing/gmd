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

package md

import (
	"github.com/fuyibing/gmd/v8/app"
	"github.com/fuyibing/gmd/v8/md/base"
	"github.com/fuyibing/gmd/v8/md/managers"
	"sync"
)

var (
	// Container
	// 容器实例.
	Container base.ContainerOperation

	// Manager
	// 管理器实例.
	Manager managers.BootManager
)

func init() {
	new(sync.Once).Do(func() {
		Manager = managers.Boot
		Container = base.Container

		// +-------------------------------------------------------------------+
		// + Adapter for: consumer & producer & remoter                        |
		// +-------------------------------------------------------------------+

		adapter := app.Config.GetAdapter()

		// 1. 消费者.
		if v := builtinConsumer(adapter).New(); v != nil {
			Container.RegisterConsumer(v)
		}

		// 2. 生产者.
		if v := builtinProducer(adapter).New(); v != nil {
			Container.RegisterProducer(v)
		}

		// 3. 服务端.
		if v := builtinRemoter(adapter).New(); v != nil {
			Container.RegisterRemoter(v)
		}

		// +-------------------------------------------------------------------+
		// + Dependency: condition & dispatcher & result                        |
		// +-------------------------------------------------------------------+

		// 1. 条件校验.
		for k, v := range builtinConditions {
			Container.RegisterCondition(k, v)
		}

		// 2. 投递消息.
		for k, v := range builtinDispatchers {
			Container.RegisterDispatcher(k, v)
		}

		// 3. 结果校验.
		for k, v := range builtinResults {
			Container.RegisterResult(k, v)
		}
	})
}
