// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package rocketmq

import (
	"context"
	"fmt"
	rmq "github.com/apache/rocketmq-client-go/v2"
	rmqc "github.com/apache/rocketmq-client-go/v2/consumer"
	rmqp "github.com/apache/rocketmq-client-go/v2/primitive"

	"github.com/fuyibing/gmd/v8/app"
	"github.com/fuyibing/gmd/v8/core/base"
	"github.com/fuyibing/util/v8/process"
	"github.com/google/uuid"
	"strings"
	"time"
)

type (
	// Consumer
	// for rocketmq adapter.
	Consumer struct {
		consume      base.ConsumerProcess
		id, parallel int
		name         string
		processor    process.Processor

		client       rmq.PushConsumer
		delayEnabled bool
		delayTag     string
		task         *base.Task
	}
)

func NewConsumer(id, parallel int, name string, consume base.ConsumerProcess) *Consumer {
	return (&Consumer{
		id: id, parallel: parallel,
		name: name, consume: consume,
	}).init()
}

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *Consumer) Processor() process.Processor { return o.processor }

// /////////////////////////////////////////////////////////////
// Internal methods
// /////////////////////////////////////////////////////////////

func (o *Consumer) DoPipe(ctx context.Context, ext ...*rmqp.MessageExt) (r rmqc.ConsumeResult, err error) {
	// log.Errorf("{%s} message received", o.name)
	return
}

func (o *Consumer) DoResume() {}

func (o *Consumer) DoSuspend() {}

// /////////////////////////////////////////////////////////////
// Processor events
// /////////////////////////////////////////////////////////////

func (o *Consumer) OnAfter(_ context.Context) (ignored bool) {
	// log.Infof("{%s} stopped", o.name)
	return
}

func (o *Consumer) OnAfterClientDestroy(_ context.Context) (ignored bool) {
	if o.client != nil {
		if err := o.client.Shutdown(); err != nil {
			// log.Errorf("{%s} rmq client shutdown error: %v", o.name, err)
		}

		o.client = nil
		// log.Infof("{%s} rmq client destroyed", o.name)
	}

	return
}

func (o *Consumer) OnBefore(_ context.Context) (ignored bool) {
	// log.Infof("{%s} start", o.name)
	return
}

func (o *Consumer) OnCall(_ context.Context) (ignored bool) {
	// log.Infof("{%s} listen channel signal", o.name)
	return
}

func (o *Consumer) OnCallChannel(ctx context.Context) (ignored bool) {
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (o *Consumer) OnCallClientBuild(_ context.Context) (ignored bool) {
	var (
		err   error
		group = fmt.Sprintf("%sGID-%d", app.Config.GetRocketmq().GetPrefix(), o.task.Id)
		node  = fmt.Sprintf("%sT%dP%d", strings.ToUpper(strings.ReplaceAll(uuid.NewString(), "-", "")), o.id, o.parallel)
		opts  = []rmqc.Option{
			rmqc.WithGroupName(group),
			rmqc.WithNsResolver(rmqp.NewPassthroughResolver(app.Config.GetRocketmq().GetServers())),
			rmqc.WithInstance(node),
		}
	)

	// Set consume mode.
	if o.task.Broadcasting {
		opts = append(opts, rmqc.WithConsumerModel(rmqc.BroadCasting))
	} else {
		opts = append(opts, rmqc.WithConsumerModel(rmqc.Clustering))
	}

	// Advanced options.
	opts = append(opts,
		rmqc.WithConsumeFromWhere(rmqc.ConsumeFromWhere(app.Config.GetRocketmq().GetConsumePosition())),
		rmqc.WithConsumeFromWhere(rmqc.ConsumeFromLastOffset),
		rmqc.WithMaxReconsumeTimes(app.Config.GetRocketmq().GetConsumeMaxRetry()),
		rmqc.WithConsumeMessageBatchMaxSize(1),
		rmqc.WithPullBatchSize(1),
		rmqc.WithSuspendCurrentQueueTimeMillis(time.Duration(app.Config.GetRocketmq().GetConsumeSuspendMs())*time.Millisecond),
	)

	// Authentication options.
	if k := app.Config.GetRocketmq().GetKey(); k != "" {
		opts = append(opts, rmqc.WithCredentials(rmqp.Credentials{
			AccessKey:     k,
			SecretKey:     app.Config.GetRocketmq().GetSecret(),
			SecurityToken: app.Config.GetRocketmq().GetToken(),
		}))
	}

	// Build consumer client.
	if o.client, err = rmq.NewPushConsumer(opts...); err != nil {
		// log.Errorf("{%s} rmq client build error: group=%s, node=%s, %v", o.name, group, node, err)
		return true
	}

	// log.Infof("{%s} rmq client built: group=%s, node=%s", o.name, group, node)
	return
}

func (o *Consumer) OnCallClientSubscribe(_ context.Context) (ignored bool) {
	var (
		err error
		sel = rmqc.MessageSelector{Type: rmqc.TAG, Expression: o.task.TopicTag}
	)

	// Append delay tag.
	if o.task.DelaySeconds > 0 {
		o.delayEnabled = true
		o.delayTag = Agent.GenDelayTag(o.task.TopicTag, o.task.Id)

		sel.Expression = fmt.Sprintf("%s || %s", o.task.TopicTag, o.delayTag)
	} else {
		o.delayEnabled = false
		o.delayTag = ""
	}

	// Subscribe topic message.
	if err = o.client.Subscribe(Agent.GenTopicName(o.task.TopicName), sel, o.DoPipe); err != nil {
		// log.Errorf("{%s} rmq client subscribe error: %v", o.name, err)
		return true
	}

	// Start rocketmq consumer client.
	if err = o.client.Start(); err != nil {
		// log.Errorf("{%s} rmq client start error: %v", o.name, err)
		return true
	}

	// log.Infof("{%s} rmq client started: topic=%s, tag=%s", o.name, o.task.TopicName, o.task.TopicTag)
	return
}

func (o *Consumer) OnCallTaskLoad(_ context.Context) (ignored bool) {
	// Return true
	// if task not in memory.
	if o.task = base.Memory.GetTask(o.id); o.task == nil {
		// log.Errorf("{%s} task not found: task=%d", o.name, o.id)
		return true
	}

	// Return true if current parallel number is greater than
	// maximum configuration.
	if o.parallel >= o.task.Parallels {
		// log.Errorf("{%s} task parallel number down: task=%d, current=%d, maximum=%d",
		// 	o.name,
		// 	o.id,
		// 	o.parallel+1,
		// 	o.task.Parallels,
		// )
		return true
	}

	// Task verified.
	// log.Infof("{%s} load task: id=%d", o.name, o.task.Id)
	return
}

func (o *Consumer) OnPanic(ctx context.Context, v interface{}) {
	// log.Panicfc(ctx, "processor {%s} fatal: %v", o.name, v)
}

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *Consumer) init() *Consumer {
	o.processor = process.New(o.name).After(
		o.OnAfterClientDestroy,
		o.OnAfter,
	).Before(
		o.OnBefore,
	).Callback(
		o.OnCall,
		o.OnCallTaskLoad,
		o.OnCallClientBuild,
		o.OnCallClientSubscribe,
		o.OnCallChannel,
	).Panic(o.OnPanic)

	return o
}
