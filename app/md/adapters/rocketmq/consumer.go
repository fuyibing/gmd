// author: wsfuyibing <websearch@163.com>
// date: 2023-02-03

package rocketmq

import (
	"context"
	"fmt"
	sdk "github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/fuyibing/gmd/app/md/base"
	"github.com/fuyibing/gmd/app/md/conf"
	"github.com/fuyibing/log/v3"
	"github.com/fuyibing/util/v2/process"
	"github.com/google/uuid"
	"strings"
	"sync"
	"time"
)

type (
	Consumer struct {
		mu        *sync.RWMutex
		name      string
		processor process.Processor

		client        sdk.PushConsumer
		clientSuspend bool
		dispatcher    func(task *base.Task, message *base.Message) (retry bool)
		id, parallel  int
		received      *Received
		task          *base.Task
	}
)

func NewConsumer(id, parallel int) *Consumer {
	return (&Consumer{
		id:       id,
		parallel: parallel,
	}).init()
}

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *Consumer) Dispatcher(x func(*base.Task, *base.Message) bool) { o.dispatcher = x }
func (o *Consumer) Processor() process.Processor                      { return o.processor }

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *Consumer) doClientResume() {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.clientSuspend {
		o.clientSuspend = false

		// Call in coroutine.
		go func() {
			if o.client != nil {
				if log.Config.DebugOn() {
					log.Debugf("%s: client resume", o.name)
				}
				o.client.Resume()
			}
		}()
	}
}

func (o *Consumer) doClientSuspend() {
	o.mu.Lock()
	defer o.mu.Unlock()

	if !o.clientSuspend {
		o.clientSuspend = true

		// Call in coroutine.
		go func() {
			if o.client != nil {
				if log.Config.DebugOn() {
					log.Debugf("%s: client suspend", o.name)
				}
				o.client.Suspend()
			}
		}()
	}
}

func (o *Consumer) doInterceptor(ctx context.Context, req, reply interface{}, next primitive.Invoker) error {
	return next(ctx, req, reply)
}

// /////////////////////////////////////////////////////////////
// Construct method
// /////////////////////////////////////////////////////////////

func (o *Consumer) init() *Consumer {
	o.mu = &sync.RWMutex{}
	o.name = fmt.Sprintf("rocketmq-consumer-%d-%d", o.id, o.parallel)
	o.processor = process.New(o.name).After(
		o.onAfter,
	).Before(
		o.onBefore,
	).Callback(
		o.onCallTaskCheck,
		o.onCallTaskSync,
		o.onCallClientBuild,
		o.onCallClientSubscribe,
		o.onCallChannel,
		o.onCallClientWaitIdle,
		o.onCallClientDestroy,
	).Panic(
		o.onPanic,
	)

	return o
}

// /////////////////////////////////////////////////////////////
// Processor events
// /////////////////////////////////////////////////////////////

func (o *Consumer) onAfter(_ context.Context) (ignored bool) {
	log.Debugf("%s: processor stopped", o.name)
	return
}

func (o *Consumer) onBefore(_ context.Context) (ignored bool) {
	log.Debugf("%s: start processor", o.name)
	return
}

func (o *Consumer) onCallChannel(ctx context.Context) (ignored bool) {
	log.Debugf("%s: listen channel signal", o.name)

	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (o *Consumer) onCallClientBuild(_ context.Context) (ignored bool) {
	var (
		err  error
		node = fmt.Sprintf("%sT%dP%d", strings.ToUpper(strings.ReplaceAll(uuid.NewString(), "-", "")), o.id, o.parallel)
		opts = []consumer.Option{
			consumer.WithGroupName(Agent.GenGroupName(o.task.Id)),
			consumer.WithNsResolver(primitive.NewPassthroughResolver(conf.Config.Account.Rocketmq.Servers)),
			consumer.WithInstance(node),
		}
	)

	// Set consume mode.
	if o.task.Broadcasting {
		opts = append(opts, consumer.WithConsumerModel(consumer.BroadCasting))
	} else {
		opts = append(opts, consumer.WithConsumerModel(consumer.Clustering))
	}

	// Extension
	opts = append(opts,
		consumer.WithInterceptor(o.doInterceptor),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromLastOffset),
		consumer.WithMaxReconsumeTimes(DefaultReconsumeTimes),
		consumer.WithConsumeMessageBatchMaxSize(1),
		consumer.WithPullBatchSize(1),
		consumer.WithSuspendCurrentQueueTimeMillis(DefaultConsumeSuspendDuration),
	)

	// Append authentication.
	if conf.Config.Account.Rocketmq.Key != "" || conf.Config.Account.Rocketmq.Secret != "" || conf.Config.Account.Rocketmq.Token != "" {
		opts = append(opts, consumer.WithCredentials(primitive.Credentials{
			AccessKey:     conf.Config.Account.Rocketmq.Key,
			SecretKey:     conf.Config.Account.Rocketmq.Secret,
			SecurityToken: conf.Config.Account.Rocketmq.Token,
		}))
	}

	// Return if client built failed.
	if o.client, err = sdk.NewPushConsumer(opts...); err != nil {
		log.Errorf("%s: client built failed, error=%v", o.name, err)
		return true
	}

	// Client create succeed.
	log.Infof("%s: client built, node=%s", o.name, node)
	return
}

func (o *Consumer) onCallClientDestroy(_ context.Context) (ignored bool) {
	// Shutdown client.
	if err := o.client.Shutdown(); err != nil {
		log.Errorf("%s: client shutdown failed: error=%v", o.name, err)
	}

	// Reset and next
	// event callee.
	o.client = nil
	o.received = nil
	log.Infof("%s: client destroy", o.name)
	return
}

func (o *Consumer) onCallClientWaitIdle(ctx context.Context) (ignored bool) {
	// Recall after specified milliseconds
	// if message consume not completed.
	if !o.received.IsIdle() {
		time.Sleep(conf.EventSleepDuration)
		return o.onCallClientWaitIdle(ctx)
	}

	// Next event callee.
	log.Infof("%s: client idle", o.name)
	return
}

func (o *Consumer) onCallClientSubscribe(_ context.Context) (ignored bool) {
	// Build received manager.
	o.received = &Received{
		client:     o.client,
		dispatcher: o.dispatcher,
		name:       o.name,
		selector:   consumer.MessageSelector{Type: consumer.TAG, Expression: o.task.TopicTag},
		task:       o.task,
		topic:      Agent.GenTopicName(o.task.TopicName),

		callbackResume:  o.doClientResume,
		callbackSuspend: o.doClientSuspend,
	}

	// Delay subscription.
	if o.task.DelaySeconds > 0 {
		o.received.delayer = true
		o.received.delayerMilliSeconds = int64(o.task.DelaySeconds * 1000)
		o.received.delayerTag = fmt.Sprintf("%s%d-%s", DefaultDelayTagPrefix, o.task.Id, o.task.TopicTag)
		o.received.selector.Expression = fmt.Sprintf("%s || %s", o.task.TopicTag, o.received.delayerTag)
	}

	// Subscribe topic.
	if err := o.client.Subscribe(o.received.topic, o.received.selector, o.received.Consume); err != nil {
		log.Errorf("%s: client subscribe failed, topic=%s, tag=%s, error=%v", o.name, o.received.topic, o.received.selector.Expression, err)
		return true
	} else {
		log.Infof("%s: client subscribed, topic=%s, tag=%v", o.name, o.received.topic, o.received.selector.Expression)
	}

	// Start client.
	if err := o.client.Start(); err != nil {
		log.Errorf("%s: client start failed, error=%v", o.name, err)
		return true
	}

	// Return succeed.
	log.Infof("%s: client started", o.name)
	return
}

func (o *Consumer) onCallTaskCheck(_ context.Context) (ignored bool) {
	// Return true
	// if subscription task not found in memory.
	if o.task = base.Memory.GetTask(o.id); o.task == nil {
		log.Errorf("%s: subscription task not found", o.name)
		return true
	}

	// Return true
	// if adapter parallel is greater or equal to task parallels.
	if o.parallel >= o.task.Parallels {
		log.Errorf("%s: subscription task parallels limited", o.name)
		return true
	}

	// Next
	// event callee.
	log.Infof("%s: subscription task loaded, topic=%v, tag=%v, filter=%v, delay-seconds=%d, title=%s", o.name, o.task.TopicName, o.task.TopicTag, o.task.FilterTag, o.task.DelaySeconds, o.task.Title)
	return false
}

func (o *Consumer) onCallTaskSync(_ context.Context) bool {
	// Next
	// event callee.
	return false
}

func (o *Consumer) onPanic(ctx context.Context, v interface{}) {
	log.Panicfc(ctx, "%s: %v", o.name, v)
}
