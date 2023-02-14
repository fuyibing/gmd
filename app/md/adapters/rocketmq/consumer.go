// author: wsfuyibing <websearch@163.com>
// date: 2023-02-03

package rocketmq

import (
	"context"
	"fmt"
	"github.com/fuyibing/gmd/app/md/base"
	"github.com/fuyibing/log/v3"
	"github.com/fuyibing/util/v2/process"
)

type (
	// Consumer
	// 消息消费者.
	Consumer struct {
		name      string
		processor process.Processor

		id, parallel int
		task         *base.Task
		dispatcher   func(task *base.Task, message *base.Message) (retry bool)
	}
)

func NewConsumer(id, parallel int) *Consumer {
	return (&Consumer{id: id, parallel: parallel}).init()
}

func (o *Consumer) Dispatcher(dispatcher func(*base.Task, *base.Message) bool) {
	o.dispatcher = dispatcher
}

func (o *Consumer) Processor() process.Processor {
	return o.processor
}

// /////////////////////////////////////////////////////////////
// 生产者构造
// /////////////////////////////////////////////////////////////

func (o *Consumer) init() *Consumer {
	o.name = fmt.Sprintf("rocketmq-consumer-%d-%d", o.id, o.parallel)
	o.processor = process.New(o.name).
		After(o.onAfter).
		Before(o.onBefore).
		Callback(o.onCallerBefore, o.onCaller, o.onCallerAfter).
		Panic(o.onPanic)

	return o
}

// /////////////////////////////////////////////////////////////
// 生产者事件
// /////////////////////////////////////////////////////////////

func (o *Consumer) onAfter(ctx context.Context) (ignored bool) {
	log.Infof("[%s] consumer finish", o.name)
	return
}

func (o *Consumer) onBefore(ctx context.Context) (ignored bool) {
	log.Infof("[%s] consumer begin", o.name)
	return
}

func (o *Consumer) onCaller(ctx context.Context) (ignored bool) {
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (o *Consumer) onCallerAfter(ctx context.Context) (ignored bool) { return }

func (o *Consumer) onCallerBefore(ctx context.Context) (ignored bool) { return }

func (o *Consumer) onPanic(ctx context.Context, v interface{}) {
	log.Panicfc(ctx, "[%s] %v", o.name, v)
}
