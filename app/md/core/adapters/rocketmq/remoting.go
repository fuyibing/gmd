// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package rocketmq

import (
	"github.com/fuyibing/util/v8/process"
)

type (
	// Remoting
	// for aliyunmns adapter.
	Remoting struct {
		processor process.Processor
	}
)

func NewRemoting() *Remoting {
	return (&Remoting{}).init()
}

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *Remoting) Processor() process.Processor { return o.processor }

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *Remoting) init() *Remoting {
	o.processor = process.New("rocketmq-remoting")
	return o
}
