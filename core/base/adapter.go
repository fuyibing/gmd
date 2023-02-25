// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package base

import (
	"github.com/fuyibing/util/v8/process"
)

type (
	// Adapter
	// defined for mq middlewares name.
	Adapter string

	// ConsumerProcess
	// message consume process.
	//
	// Call this handler when message received from mq.
	//
	// If the ignored value is true, it means that the message is ignored
	// because the consumption condition is not met.
	ConsumerProcess func(task *Task, message *Message) (ignored bool, err error)

	// ConsumerCallable
	// constructor for create ConsumerManager instance.
	ConsumerCallable func(id, parallel int, name string, process ConsumerProcess) ConsumerManager

	// ConsumerManager
	// receive message from queue of mq middleware server.
	ConsumerManager interface {
		Processor() process.Processor
	}

	// ProducerCallable
	// constructor for create ProducerManager instance.
	ProducerCallable func() ProducerManager

	// ProducerManager
	// publish message to mq middleware server.
	ProducerManager interface {
		Processor() process.Processor
	}

	// RemotingCallable
	// constructor for create RemotingManager instance.
	RemotingCallable func() RemotingManager

	// RemotingManager
	// manager remote with mq middleware server.
	RemotingManager interface {
		Processor() process.Processor
	}
)

// Adapter enums.

const (
	AliyunMns Adapter = "aliyunmns"
	RocketMq  Adapter = "rocketmq"
)

func (adapter Adapter) String() {}