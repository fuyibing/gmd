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
	Container base.ContainerOperation
	Manager   managers.BootManager
)

func init() {
	new(sync.Once).Do(func() {
		Manager = managers.Boot
		Container = base.Container

		// +-------------------------------------------------------------------+
		// + Adapter: rocketmq/rabbitmq/aliyunmns                              |
		// +   1. consumer                                                     |
		// +   2. producer                                                     |
		// +   3. remoter                                                      |
		// +-------------------------------------------------------------------+
		// + Register follow constructors when package initialized             |
		// +-------------------------------------------------------------------+

		adapter := app.Config.GetAdapter()

		// Register consumer constructor use adapter.
		if v := builtinConsumer(adapter).New(); v != nil {
			Container.RegisterConsumer(v)
		}

		// Register producer constructor use adapter.
		if v := builtinProducer(adapter).New(); v != nil {
			Container.RegisterProducer(v)
		}

		// Register remoter constructor use adapter.
		if v := builtinRemoter(adapter).New(); v != nil {
			Container.RegisterRemoter(v)
		}

		// +-------------------------------------------------------------------+
		// + Utilities and Dependency:                                         |
		// +   condition                                                       |
		// +   dispatcher                                                      |
		// +   result                                                          |
		// +-------------------------------------------------------------------+

		// Register builtin conditions.
		for k, v := range builtinConditions {
			Container.RegisterCondition(k, v)
		}

		// Register builtin dispatchers.
		for k, v := range builtinDispatchers {
			Container.RegisterDispatcher(k, v)
		}

		// Register builtin result parser.
		for k, v := range builtinResults {
			Container.RegisterResult(k, v)
		}
	})
}
