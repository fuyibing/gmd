// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package core

import (
	"github.com/fuyibing/gmd/v8/app/md/base"
	"github.com/fuyibing/gmd/v8/app/md/conf"
	"sync"
)

type (
	container struct {
		mu *sync.RWMutex

		consumer base.ConsumerCallable
		producer base.ProducerCallable
		remoting base.RemotingCallable

		conditions  map[base.Condition]base.ContainerManager
		dispatchers map[base.Dispatcher]base.DispatcherManager
		results     map[base.Result]base.ResultManager
	}
)

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

// GetConsumer
// return consumer manager constructor.
func (o *container) GetConsumer() (callable base.ConsumerCallable, exists bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if o.consumer != nil {
		callable = o.consumer
		exists = true
	}
	return
}

// GetProducer
// return producer manager constructor.
func (o *container) GetProducer() (callable base.ProducerCallable, exists bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if o.producer != nil {
		callable = o.producer
		exists = true
	}
	return
}

// GetRemoting
// return remoting manager constructor.
func (o *container) GetRemoting() (callable base.RemotingCallable, exists bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if o.remoting != nil {
		callable = o.remoting
		exists = true
	}
	return
}

// SetAdapter
// configure consumer, producer, remoting with adapter.
func (o *container) SetAdapter(adapter base.Adapter) {
	// Use consumer callable if defined.
	if v, ok := builtInConsumer[adapter]; ok {
		o.SetConsumer(v)
	}

	// Use producer callable if defined.
	if v, ok := builtInProducer[adapter]; ok {
		o.SetProducer(v)
	}

	// Use remoting callable if defined.
	if v, ok := builtInRemoting[adapter]; ok {
		o.SetRemoting(v)
	}
}

// SetConsumer
// configure consumer manager constructor.
func (o *container) SetConsumer(callable base.ConsumerCallable) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.consumer = callable
}

// SetProducer
// configure producer manager constructor.
func (o *container) SetProducer(callable base.ProducerCallable) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.producer = callable
}

// SetRemoting
// configure remoting manager constructor.
func (o *container) SetRemoting(callable base.RemotingCallable) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.remoting = callable
}

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *container) init() *container {
	o.mu = &sync.RWMutex{}

	o.initAdapter()
	o.initCondition()
	o.initDispatcher()
	o.initResult()
	return o
}

func (o *container) initAdapter() {
	o.SetAdapter(conf.Config.GetAdapter())
}

func (o *container) initCondition() {
	o.conditions = map[base.Condition]base.ContainerManager{
		base.ConditionEl: nil,
	}
}

func (o *container) initDispatcher() {
	o.dispatchers = map[base.Dispatcher]base.DispatcherManager{
		base.DispatchHttp: nil,
		base.DispatchRpc:  nil,
	}
}

func (o *container) initResult() {
	o.results = map[base.Result]base.ResultManager{
		base.ResultHttpOk:        nil,
		base.ResultJsonErrnoZero: nil,
	}
}
