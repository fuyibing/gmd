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
	// 消费者构造/适配MQ中间件.
	ConsumerCallable func(id, parallel int, name string, process ConsumerProcess) ConsumerManager

	// ConsumerManager
	// 消费者管理器/适配MQ中间件.
	ConsumerManager interface {
		// Processor
		// 类进程.
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
	// 生产者构造/适配MQ中间件.
	ProducerCallable func() ProducerManager

	// ProducerManager
	// 生产者管理器/适配MQ中间件.
	ProducerManager interface {
		// Processor
		// 类进程.
		Processor() process.Processor

		// Publish
		// 发布消息.
		Publish(v *Payload) (mid string, err error)
	}

	// RemotingCallable
	// 服务端构造/适配MQ中间件.
	RemotingCallable func() RemotingManager

	// RemotingManager
	// 服务端管理器/适配MQ中间件.
	RemotingManager interface {
		// Processor
		// 类进程.
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
