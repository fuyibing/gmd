// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package tcp

import (
	"github.com/fuyibing/gmd/v8/core/base"
)

// Dispatcher
// 消息分发.
type Dispatcher struct{}

func New() *Dispatcher { return (&Dispatcher{}).init() }

func (o *Dispatcher) Dispatch(_ *base.Task, _ *base.Subscriber, _ *base.Message) (err error) {
	return
}

func (o *Dispatcher) init() *Dispatcher {
	return o
}
