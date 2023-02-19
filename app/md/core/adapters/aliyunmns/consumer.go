// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package aliyunmns

import (
	"github.com/fuyibing/util/v8/process"
)

type (
	// Consumer
	// for aliyunmns adapter.
	Consumer struct {
		id, parallel int
		processor    process.Processor
	}
)

func NewConsumer(id, parallel int) *Consumer {
	return &Consumer{id: id, parallel: parallel}
}

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *Consumer) Processor() process.Processor { return o.processor }

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *Consumer) init() *Consumer {
	o.processor = process.New("aliyunmns-consumer")
	return o
}
