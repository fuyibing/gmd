// author: wsfuyibing <websearch@163.com>
// date: 2023-02-20

package managers

import (
	"github.com/fuyibing/gmd/v8/app"
	"github.com/fuyibing/gmd/v8/core/base"
	"sync"
)

type (
	// 适配器容器.
	container struct {
		ac base.ConsumerCallable
		ap base.ProducerCallable
		ar base.RemotingCallable
		cs map[string]base.ConditionCallable
		ds map[string]base.DispatcherCallable
		rs map[string]base.ResultCallable
		mu *sync.RWMutex
	}
)

// GetConsumer
// return consumer manager constructor.
func (o *container) GetConsumer() base.ConsumerCallable { return o.ac }

// SetConsumer
// configure consumer manager constructor, singleton instance.
func (o *container) SetConsumer(callable base.ConsumerCallable) { o.ac = callable }

// GetProducer
// return producer manager constructor.
func (o *container) GetProducer() base.ProducerCallable { return o.ap }

// SetProducer
// configure producer manager constructor, singleton instance.
func (o *container) SetProducer(callable base.ProducerCallable) { o.ap = callable }

// GetRemoting
// return remoting manager constructor.
func (o *container) GetRemoting() base.RemotingCallable { return o.ar }

// SetRemoting
// configure remoting manager constructor, singleton instance.
func (o *container) SetRemoting(callable base.RemotingCallable) { o.ar = callable }

func (o *container) init() *container {
	o.mu = &sync.RWMutex{}
	o.cs = make(map[string]base.ConditionCallable)
	o.ds = make(map[string]base.DispatcherCallable)
	o.rs = make(map[string]base.ResultCallable)

	o.initAdapter()
	o.initConditions()
	o.initDispatchers()
	o.initResults()
	return o
}

func (o *container) initAdapter() {
	adapter := base.Adapter(app.Config.GetAdapter())

	if v, ok := builtInConsumer[adapter]; ok {
		o.SetConsumer(v)
	}

	if v, ok := builtInProducer[adapter]; ok {
		o.SetProducer(v)
	}

	if v, ok := builtInRemoting[adapter]; ok {
		o.SetRemoting(v)
	}
}

func (o *container) initConditions()  { o.cs = buildInConditions }
func (o *container) initDispatchers() { o.ds = buildInDispatchers }
func (o *container) initResults()     { o.rs = buildInResults }
