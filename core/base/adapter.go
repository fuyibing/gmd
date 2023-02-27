// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package base

import (
	"github.com/fuyibing/util/v8/process"
)

type (
	// Adapter
	// 适配器名称.
	Adapter string

	// ConsumerCallable
	// MQ中间件的消费者构造函数.
	ConsumerCallable func(id, parallel int, name string, process ConsumerProcess) ConsumerManager

	// ConsumerManager
	// MQ中间件的消费者管理器.
	ConsumerManager interface {
		// Processor
		// 获取执行器.
		Processor() process.Processor
	}

	// ConsumerProcess
	// 消息执行器.
	//
	// 当从MQ队列(AliyunMNS, RocketMQ等)收到消息时, 通过此执行器处理消费过程.
	//
	// 若返回的 ignored 值为 true, 表示消息不满足投递条件(条件校验)被忽略了, 反之
	// 需分发给订阅方. 若返回的 err 值非空, 则表示条件校验出错.
	ConsumerProcess func(task *Task, message *Message) (ignored bool, err error)

	// ProducerCallable
	// MQ中间件的生产者构造函数.
	ProducerCallable func() ProducerManager

	// ProducerManager
	// MQ中间件的生产者管理器.
	ProducerManager interface {
		// Processor
		// 获取执行器.
		Processor() process.Processor
	}

	// RemotingCallable
	// MQ中间件的服务商构造函数.
	RemotingCallable func() RemotingManager

	// RemotingManager
	// MQ中间件的服务端管理器.
	RemotingManager interface {
		// Processor
		// 获取执行器.
		Processor() process.Processor
	}
)

// 适配器枚举.

const (
	AliyunMns Adapter = "aliyunmns"
	RocketMq  Adapter = "rocketmq"
)

// String
// 适配器名称.
func (a Adapter) String() string { return string(a) }
