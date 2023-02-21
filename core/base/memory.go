// author: wsfuyibing <websearch@163.com>
// date: 2023-01-17

package base

import (
	"context"
	"fmt"
	"github.com/fuyibing/db/v8"
	"github.com/fuyibing/gmd/v8/app/models"
	"github.com/fuyibing/gmd/v8/app/services"
	"strings"
	"sync"
)

var (
	// Memory
	// instance of memory manager.
	Memory MemoryManager
)

type (
	// MemoryManager
	// interface of memory manager.
	MemoryManager interface {
		// GetRegistries
		// return registry map in memory.
		GetRegistries() map[int]*Registry

		// GetRegistry
		// return registry instance in memory by id.
		GetRegistry(id int) *Registry

		// GetRegistryByName
		// return registry instance in memory by topic name
		// and tag.
		GetRegistryByName(topic, tag string) *Registry

		// GetTask
		// return task instance in memory by id.
		GetTask(id int) *Task

		// GetTaskFromBean
		// return task instance use task id.
		GetTaskFromBean(ctx context.Context, id int) (task *Task, err error)

		// GetTasks
		// return task instance map in memory.
		GetTasks() map[int]*Task

		// Reload
		// read from database into memory.
		Reload(c context.Context) error
	}

	memory struct {
		mu             *sync.RWMutex
		registryKey    map[string]int
		registryMapper map[int]*Registry
		taskMapper     map[int]*Task
	}
)

func (o *memory) GetRegistries() map[int]*Registry {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.registryMapper
}

func (o *memory) GetRegistry(id int) *Registry {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if v, ok := o.registryMapper[id]; ok {
		return v
	}
	return nil
}

func (o *memory) GetRegistryByName(topic, tag string) *Registry {
	key := o.key(topic, tag)
	o.mu.RLock()
	defer o.mu.RUnlock()
	if id, exists := o.registryKey[key]; exists {
		if v, ok := o.registryMapper[id]; ok {
			return v
		}
	}
	return nil
}

func (o *memory) GetTask(id int) *Task {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if v, ok := o.taskMapper[id]; ok {
		return v
	}
	return nil
}

func (o *memory) GetTaskFromBean(ctx context.Context, id int) (task *Task, err error) {
	var (
		bean *models.Task
		sess = db.Connector.GetSlaveWithContext(ctx)
	)

	if bean, err = services.NewTaskService(sess).GetById(id); err != nil || bean == nil {
		return
	}
	if r := o.GetRegistry(bean.RegistryId); r != nil {
		task = (&Task{}).bind(r).init(bean)
	}
	return
}

func (o *memory) GetTasks() map[int]*Task {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.taskMapper
}

func (o *memory) Reload(c context.Context) (err error) {
	for _, call := range []func(context.Context) error{o.loadRegistry, o.loadTask} {
		if err = call(c); err != nil {
			return
		}
	}
	return
}

// /////////////////////////////////////////////////////////////
// Access methods.
// /////////////////////////////////////////////////////////////

func (o *memory) init() *memory {
	o.mu = &sync.RWMutex{}
	o.registryKey = make(map[string]int)
	o.registryMapper = make(map[int]*Registry)
	o.taskMapper = make(map[int]*Task)
	return o
}

func (o *memory) key(topic, tag string) string {
	return strings.ToUpper(fmt.Sprintf("%s:%s", topic, tag))
}

func (o *memory) loadRegistry(c context.Context) (err error) {
	var list []*models.Registry

	// List all registry
	// from database.
	if list, err = services.NewRegistryService(db.Connector.GetSlaveWithContext(c)).ListAll(); err != nil {
		return
	}

	// Iterate registries
	// and generate to buffer.
	keys := make(map[string]int)
	mapper := make(map[int]*Registry)
	for _, bean := range list {
		key := o.key(bean.TopicName, bean.TopicTag)
		keys[key] = bean.Id
		mapper[bean.Id] = (&Registry{}).init(bean)
	}

	// Set buffer
	// as memory.
	o.mu.Lock()
	o.registryKey = keys
	o.registryMapper = mapper
	o.mu.Unlock()
	return
}

func (o *memory) loadTask(c context.Context) (err error) {
	var list []*models.Task

	// List enable tasks
	// from database.
	if list, err = services.NewTaskService(db.Connector.GetSlaveWithContext(c)).ListEnables(); err != nil {
		return
	}

	// Iterate tasks
	// and generate to buffer.
	tasks := make(map[int]*Task)
	for _, bean := range list {
		if r := o.GetRegistry(bean.RegistryId); r != nil {
			tasks[bean.Id] = (&Task{}).init(bean).bind(r)
		}
	}

	// Set buffer
	// as memory.
	o.mu.Lock()
	o.taskMapper = tasks
	o.mu.Unlock()
	return
}
