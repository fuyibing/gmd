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
	"context"
	"fmt"
	"github.com/fuyibing/db/v5"
	"github.com/fuyibing/gmd/v8/app/models"
	"github.com/fuyibing/gmd/v8/app/services"
	"strings"
	"sync"
	"xorm.io/xorm"
)

var (
	Memory MemoryOperation
)

type (
	MemoryOperation interface {
		GetRegistry(id int) (registry *Registry, exists bool)
		GetRegistryByNames(topicName, topicTag string) (registry *Registry, exists bool)
		GetTask(id int) (task *Task, exists bool)
		GetTasks() (task map[int]*Task)
		Reload(ctx context.Context) (err error)
	}

	memory struct {
		sync.RWMutex

		registryKeys   map[string]int
		registryMapper map[int]*Registry
		taskMapper     map[int]*Task
	}
)

// +---------------------------------------------------------------------------+
// + Interface methods                                                         |
// +---------------------------------------------------------------------------+

func (o *memory) GetRegistry(id int) (registry *Registry, exists bool) {
	o.RLock()
	defer o.RUnlock()

	registry, exists = o.registryMapper[id]
	return
}

func (o *memory) GetRegistryByNames(topicName, topicTag string) (registry *Registry, exists bool) {
	o.RLock()
	defer o.RUnlock()

	k := o.key(topicName, topicTag)

	if id, ok := o.registryKeys[k]; ok {
		registry, exists = o.registryMapper[id]
	}
	return
}

func (o *memory) GetTask(id int) (task *Task, exists bool) {
	o.RLock()
	defer o.RUnlock()

	task, exists = o.taskMapper[id]
	return
}

func (o *memory) GetTasks() (task map[int]*Task) {
	o.RLock()
	defer o.RUnlock()

	return o.taskMapper
}

func (o *memory) Reload(ctx context.Context) (err error) {
	sess := db.Connector.GetSlaveWithContext(ctx)
	for _, call := range []func(*xorm.Session) error{
		o.loadRegistry,
		o.loadTask,
	} {
		if err = call(sess); err != nil {
			return
		}
	}
	return
}

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *memory) init() *memory {
	o.registryKeys = make(map[string]int, 0)
	o.registryMapper = make(map[int]*Registry, 0)
	o.taskMapper = make(map[int]*Task, 0)
	return o
}

func (o *memory) key(topicName, topicTag string) string {
	return strings.ToUpper(fmt.Sprintf("%s:%s", topicName, topicTag))
}

func (o *memory) loadRegistry(sess *xorm.Session) (err error) {
	var (
		list   []*models.Registry
		keys   = make(map[string]int)
		mapper = make(map[int]*Registry)
	)
	// 读取列表.
	if list, err = services.NewRegistryService(sess).ListAll(); err != nil {
		return
	}
	// 遍历列表.
	for _, bean := range list {
		k := o.key(bean.TopicName, bean.TopicTag)
		keys[k] = bean.Id
		mapper[bean.Id] = (&Registry{}).init(bean)
	}
	// 更新内存.
	o.Lock()
	o.registryKeys = keys
	o.registryMapper = mapper
	o.Unlock()
	return
}

func (o *memory) loadTask(sess *xorm.Session) (err error) {
	var (
		list   []*models.Task
		mapper = make(map[int]*Task)
	)
	// 读取列表.
	if list, err = services.NewTaskService(sess).ListEnables(); err != nil {
		return
	}
	// 遍历列表.
	for _, bean := range list {
		x := &Task{}
		if registry, exists := func(id int) (registry *Registry, exists bool) {
			o.RLock()
			defer o.RUnlock()
			registry, exists = o.registryMapper[id]
			return
		}(bean.RegistryId); exists {
			x.bind(registry)
		}
		mapper[bean.Id] = x.init(bean)
	}
	// 更新内存.
	o.Lock()
	o.taskMapper = mapper
	o.Unlock()
	return
}
