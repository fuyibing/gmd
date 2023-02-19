// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package rocketmq

import (
	"github.com/fuyibing/util/v8/process"
)

type (
	// Producer
	// for aliyunmns adapter.
	Producer struct {
		processor process.Processor
	}
)

func NewProducer() *Producer {
	return (&Producer{}).init()
}

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *Producer) Processor() process.Processor { return o.processor }

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *Producer) init() *Producer {
	o.processor = process.New("rocketmq-producer")
	return o
}
