// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package http_post_json

import (
	"github.com/fuyibing/gmd/v8/core/base"
)

// Dispatcher
// dispatch message to subscriber.
type Dispatcher struct{}

func New() *Dispatcher {
	return (&Dispatcher{}).init()
}

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *Dispatcher) Dispatch(_ *base.Task, _ *base.Subscriber, _ *base.Message) (err error) {
	return
}

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *Dispatcher) init() *Dispatcher {
	return o
}
