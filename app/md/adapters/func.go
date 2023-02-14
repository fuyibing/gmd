// author: wsfuyibing <websearch@163.com>
// date: 2023-02-03

package adapters

import (
	"fmt"
	"github.com/fuyibing/gmd/app/md/adapters/aliyunmns"
	"github.com/fuyibing/gmd/app/md/adapters/rabbitmq"
	"github.com/fuyibing/gmd/app/md/adapters/rocketmq"
	"github.com/fuyibing/gmd/app/md/conf"
)

// NewConsumer
// create and return consumer adapter instance.
func NewConsumer(a conf.Adapter, id, parallel int) (adapter ConsumerAdapter, err error) {
	switch a {
	case conf.Aliyunmns:
		adapter = aliyunmns.NewConsumer(id, parallel)
	case conf.Rabbitmq:
		adapter = rabbitmq.NewConsumer(id, parallel)
	case conf.Rocketmq:
		adapter = rocketmq.NewConsumer(id, parallel)
	default:
		err = fmt.Errorf("undefined")
	}
	return
}

// NewProducer
// create and return producer adapter instance.
func NewProducer(a conf.Adapter) (adapter ProducerAdapter, err error) {
	switch a {
	case conf.Aliyunmns:
		adapter = aliyunmns.NewProducer()
	case conf.Rabbitmq:
		adapter = rabbitmq.NewProducer()
	case conf.Rocketmq:
		adapter = rocketmq.NewProducer()
	default:
		err = fmt.Errorf("undefined")
	}
	return
}

// NewRemoter
// create and return remoter adapter instance.
func NewRemoter(a conf.Adapter) (adapter RemoterAdapter, err error) {
	switch a {
	case conf.Aliyunmns:
		adapter = aliyunmns.NewRemoter()
	case conf.Rabbitmq:
		adapter = rabbitmq.NewRemoter()
	case conf.Rocketmq:
		adapter = rocketmq.NewRemoter()
	default:
		err = fmt.Errorf("undefined")
	}
	return
}
