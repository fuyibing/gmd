// author: wsfuyibing <websearch@163.com>
// date: 2023-02-09

// Package dispatchers
// Top level of core library for dispatcher message interface.
package dispatchers

import (
	"sync"
)

var (
	// Pool
	// instance of pool manager.
	Pool PoolManager
)

type (
	// PoolManager
	// interface of pool manager.
	PoolManager interface {
		// AcquireHttp
		// acquire http dispatcher instance from pool.
		AcquireHttp() *HttpDispatcher

		// ReleaseHttp
		// release http dispatcher instance into pool.
		ReleaseHttp(x *HttpDispatcher)
	}

	pool struct {
		httpDispatchers *sync.Pool
	}
)

func (o *pool) AcquireHttp() *HttpDispatcher {
	x := o.httpDispatchers.Get().(*HttpDispatcher)
	x.before()
	return x
}

func (o *pool) ReleaseHttp(x *HttpDispatcher) {
	x.after()
	o.httpDispatchers.Put(x)
}

// /////////////////////////////////////////////////////////////
// Access methods.
// /////////////////////////////////////////////////////////////

func (o *pool) init() *pool {
	return o
}

func init() {
	new(sync.Once).Do(func() {
		Pool = (&pool{
			httpDispatchers: &sync.Pool{
				New: func() interface{} {
					return (&HttpDispatcher{}).init()
				},
			},
		}).init()
	})
}
