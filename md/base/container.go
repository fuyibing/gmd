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
	"sync"
)

var (
	Container ContainerOperation
)

type (
	ContainerOperation interface {
		GetConsumer() (constructor ConsumerConstructor)
		GetProducer() (constructor ProducerConstructor)
		GetRemoter() (constructor RemoterConstructor)
		GetCondition(key string) (constructor ConditionConstructor, exists bool)
		GetDispatcher(key string) (constructor DispatcherConstructor, exists bool)
		GetResult(key string) (constructor ResultConstructor, exists bool)
		RegisterConsumer(constructor ConsumerConstructor)
		RegisterProducer(constructor ProducerConstructor)
		RegisterRemoter(v RemoterConstructor)
		RegisterCondition(key string, constructor ConditionConstructor)
		RegisterDispatcher(key string, constructor DispatcherConstructor)
		RegisterResult(key string, constructor ResultConstructor)
	}

	container struct {
		sync.RWMutex

		cc ConsumerConstructor
		pc ProducerConstructor
		rc RemoterConstructor

		cs map[string]ConditionConstructor
		ds map[string]DispatcherConstructor
		rs map[string]ResultConstructor
	}
)

// +---------------------------------------------------------------------------+
// + Interface methods                                                         |
// +---------------------------------------------------------------------------+

func (o *container) GetConsumer() ConsumerConstructor {
	o.RLock()
	defer o.RUnlock()
	return o.cc
}

func (o *container) GetProducer() ProducerConstructor {
	o.RLock()
	defer o.RUnlock()
	return o.pc
}

func (o *container) GetRemoter() RemoterConstructor {
	o.RLock()
	defer o.RUnlock()
	return o.rc
}

func (o *container) GetCondition(k string) (v ConditionConstructor, exists bool) {
	o.RLock()
	defer o.RUnlock()
	v, exists = o.cs[k]
	return
}

func (o *container) GetDispatcher(k string) (v DispatcherConstructor, exists bool) {
	o.RLock()
	defer o.RUnlock()
	v, exists = o.ds[k]
	return
}

func (o *container) GetResult(k string) (v ResultConstructor, exists bool) {
	o.RLock()
	defer o.RUnlock()
	v, exists = o.rs[k]
	return
}

func (o *container) RegisterConsumer(v ConsumerConstructor) {
	o.Lock()
	defer o.Unlock()
	o.cc = v
}

func (o *container) RegisterProducer(v ProducerConstructor) {
	o.Lock()
	defer o.Unlock()
	o.pc = v
}

func (o *container) RegisterRemoter(v RemoterConstructor) {
	o.Lock()
	defer o.Unlock()
	o.rc = v
}

func (o *container) RegisterCondition(k string, v ConditionConstructor) {
	o.Lock()
	defer o.Unlock()
	o.cs[k] = v
}

func (o *container) RegisterDispatcher(k string, v DispatcherConstructor) {
	o.Lock()
	defer o.Unlock()
	o.ds[k] = v
}

func (o *container) RegisterResult(k string, v ResultConstructor) {
	o.Lock()
	defer o.Unlock()
	o.rs[k] = v
}

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *container) init() *container {
	o.cs = make(map[string]ConditionConstructor)
	o.ds = make(map[string]DispatcherConstructor)
	o.rs = make(map[string]ResultConstructor)
	return o
}
