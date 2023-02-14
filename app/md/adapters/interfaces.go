// author: wsfuyibing <websearch@163.com>
// date: 2023-02-03

// Package adapters
// Message queue adapters.
package adapters

import (
	"context"
	"github.com/fuyibing/gmd/app/md/base"
	"github.com/fuyibing/util/v2/process"
)

type (
	// ConsumerAdapter
	// interface of consumer adapter.
	ConsumerAdapter interface {
		// Dispatcher
		// bind consume callback for consumer adapter.
		Dispatcher(dispatcher func(task *base.Task, message *base.Message) (retry bool))

		// Processor
		// return processor instance of consumer adapter.
		Processor() process.Processor
	}

	// ProducerAdapter
	// interface of producer adapter.
	ProducerAdapter interface {
		// Processor
		// return processor instance of producer adapter.
		Processor() process.Processor

		// Publish
		// send message for producer adapter.
		Publish(p *base.Payload) (messageId string, err error)
	}

	// RemoterAdapter
	// interface of remoter adapter.
	RemoterAdapter interface {
		Build(ctx context.Context, task *base.Task) (err error)
		BuildById(ctx context.Context, id int) (err error)
		Destroy(ctx context.Context, task *base.Task) (err error)
		DestroyById(ctx context.Context, id int) (err error)

		// Processor
		// return processor instance of producer adapter.
		Processor() process.Processor
	}
)
