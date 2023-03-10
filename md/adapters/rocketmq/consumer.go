// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// author: wsfuyibing <websearch@163.com>
// date: 2023-03-08

package rocketmq

import (
	"context"
	"fmt"
	rmq "github.com/apache/rocketmq-client-go/v2"
	rmqc "github.com/apache/rocketmq-client-go/v2/consumer"
	rmqp "github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/fuyibing/gmd/v8/app"
	"github.com/fuyibing/gmd/v8/md/base"
	"github.com/fuyibing/log/v5"
	"github.com/fuyibing/util/v8/process"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type (
	// Consumer
	// 消费者.
	Consumer struct {
		sync.RWMutex

		client       rmq.PushConsumer
		delayEnabled bool
		delayMilli   int64
		delayTag     string
		handler      base.ConsumerHandler
		id, parallel int
		key, name    string
		processor    process.Processor
		processing   int32
		suspended    bool
		task         *base.Task
	}
)

func NewConsumer(id, parallel int, key string, handler base.ConsumerHandler) base.ConsumerExecutor {
	return (&Consumer{
		handler: handler,
		id:      id, parallel: parallel,
		key: key,
	}).init()
}

// +---------------------------------------------------------------------------+
// + Interface methods                                                         |
// +---------------------------------------------------------------------------+

func (o *Consumer) Processor() process.Processor { return o.processor }

// +---------------------------------------------------------------------------+
// + Event methods                                                             |
// +---------------------------------------------------------------------------+

func (o *Consumer) onCall(ctx context.Context) (ignored bool) {
	log.Info("<%s.%s> channel listening", o.name, o.key)
	defer log.Info("<%s.%s> channel closed", o.name, o.key)

	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (o *Consumer) onClientBuild(_ context.Context) (ignored bool) {
	var (
		err   error
		group = Agent.GenerateGroupId(o.task.Id)
		node  = fmt.Sprintf("%sT%dP%d", strings.ToUpper(strings.ReplaceAll(uuid.NewString(), "-", "")), o.id, o.parallel)
		opts  = []rmqc.Option{
			rmqc.WithGroupName(group),
			rmqc.WithNsResolver(rmqp.NewPassthroughResolver(app.Config.GetRocketmq().GetServers())),
			rmqc.WithInstance(node),
		}
	)

	// 消费模式.
	if o.task.Broadcasting {
		opts = append(opts, rmqc.WithConsumerModel(rmqc.BroadCasting))
	} else {
		opts = append(opts, rmqc.WithConsumerModel(rmqc.Clustering))
	}

	// 消费选项.
	opts = append(opts,
		rmqc.WithConsumeFromWhere(rmqc.ConsumeFromLastOffset),
		rmqc.WithMaxReconsumeTimes(30),
		rmqc.WithConsumeMessageBatchMaxSize(1),
		rmqc.WithPullBatchSize(1),
		rmqc.WithSuspendCurrentQueueTimeMillis(time.Millisecond),
	)

	// 连接鉴权.
	if k := app.Config.GetRocketmq().GetKey(); k != "" {
		opts = append(opts, rmqc.WithCredentials(rmqp.Credentials{
			AccessKey:     k,
			SecretKey:     app.Config.GetRocketmq().GetSecret(),
			SecurityToken: app.Config.GetRocketmq().GetToken(),
		}))
	}

	// 建立连接.
	if o.client, err = rmq.NewPushConsumer(opts...); err != nil {
		log.Error("<%s.%s> client build error: group=%s, node=%s, %v", o.name, o.key, group, node, err)
		return true
	}

	log.Info("<%s.%s> client built: group=%s, node=%s", o.name, o.key, group, node)
	return
}

func (o *Consumer) onClientShutdown(ctx context.Context) (ignored bool) {
	if atomic.LoadInt32(&o.processing) == 0 {
		if err := o.client.Shutdown(); err != nil {
			log.Error("<%s.%s> client shutdown: %v", o.name, o.key, err)
		} else {
			log.Info("<%s.%s> client shutdown", o.name, o.key)
		}

		o.client = nil
		return
	}

	time.Sleep(time.Millisecond * 100)
	return o.onClientShutdown(ctx)
}

func (o *Consumer) onClientSubscribe(_ context.Context) (ignored bool) {
	var (
		err error
		sel = rmqc.MessageSelector{Type: rmqc.TAG, Expression: o.task.TopicTag}
	)

	// 标签过滤.
	if o.task.DelaySeconds > 0 {
		o.delayEnabled = true
		o.delayMilli = int64(o.task.DelaySeconds) * 1000
		o.delayTag = Agent.GenerateDelayTag(o.id)

		sel.Expression = fmt.Sprintf("%s || %s", o.task.TopicTag, o.delayTag)
	}

	// 创建订阅.
	if err = o.client.Subscribe(Agent.GenerateTopicName(o.task.TopicName), sel, o.pipe); err != nil {
		log.Error("<%s.%s> client subscribe error: %v", o.name, o.key, err)
		return true
	}

	// 启动连接.
	if err = o.client.Start(); err != nil {
		log.Error("<%s.%s> client start error: %v", o.name, o.key, err)
		return true
	}

	log.Info("<%s.%s> client started: topic=%s, tag=%s", o.name, o.key, o.task.TopicName, o.task.TopicTag)
	return
}

func (o *Consumer) onPanic(_ context.Context, v interface{}) {
	log.Fatal("<%s.%s> %v", o.name, o.key, v)
}

func (o *Consumer) onTaskLoad(_ context.Context) (ignored bool) {
	var exists bool

	// 无效任务.
	// 订阅任务停用或已被删除.
	if o.task, exists = base.Memory.GetTask(o.id); !exists {
		log.Error("<%s.%s> task not found: task-id=%d", o.name, o.key, o.id)

		o.processor.UnbindWhenStopped(true)
		return true
	}

	// 并行下降.
	// 订阅任务的最大消费者数量限制.
	if o.parallel >= o.task.Parallels {
		log.Error("<%s.%s> task parallel limited: current=%d, maximum=%d", o.name, o.key, o.parallel, o.task.Parallels)

		o.processor.UnbindWhenStopped(true)
		return true
	}

	log.Info("<%s.%s> task loaded: task-id=%d, task-title=%s", o.name, o.key, o.id, o.task.Title)
	return
}

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *Consumer) delaySecondsToLevel(seconds int) (level int) {
	if seconds < 5 {
		level = 1
	} else if seconds >= 5 && seconds < 10 {
		level = 2
	} else if seconds >= 10 && seconds < 30 {
		level = 3
	} else if seconds >= 30 && seconds < 60 {
		level = 4
	} else if seconds >= 60 && seconds < 120 {
		level = 5
	} else if seconds >= 120 && seconds < 180 {
		level = 6
	} else if seconds >= 180 && seconds < 240 {
		level = 7
	} else if seconds >= 240 && seconds < 300 {
		level = 8
	} else if seconds >= 300 && seconds < 360 {
		level = 9
	} else if seconds >= 360 && seconds < 360 {
		level = 10
	} else if seconds >= 420 && seconds < 360 {
		level = 11
	} else if seconds >= 480 && seconds < 360 {
		level = 12
	} else if seconds >= 540 && seconds < 360 {
		level = 13
	} else if seconds >= 600 && seconds < 1200 {
		level = 14
	} else if seconds >= 1200 && seconds < 1800 {
		level = 15
	} else if seconds >= 1800 && seconds < 3600 {
		level = 16
	} else if seconds >= 3600 && seconds < 7200 {
		level = 17
	} else if seconds >= 7200 {
		level = 18
	}
	return
}

func (o *Consumer) init() *Consumer {
	o.name = "rocketmq.consumer"
	o.processor = process.New(o.key).
		Before(o.onTaskLoad, o.onClientBuild).
		Callback(o.onClientSubscribe, o.onCall, o.onClientShutdown).
		Panic(o.onPanic)
	return o
}

func (o *Consumer) pipe(ctx context.Context, list ...*rmqp.MessageExt) (res rmqc.ConsumeResult, err error) {
	o.pipeSuspend()
	defer o.pipeResume()

	res = rmqc.ConsumeRetryLater

	if len(list) == 1 {
		if o.delayEnabled {
			if !o.pipeDelay(ctx, list[0]) {
				res = rmqc.ConsumeSuccess
			}
		} else {
			if !o.pipeHandle(list[0]) {
				res = rmqc.ConsumeSuccess
			}
		}
	}

	return
}

func (o *Consumer) pipeDelay(ctx context.Context, ext *rmqp.MessageExt) (retry bool) {
	// 消息时间.
	// 生产者发布消息时间.
	bornMilli := ext.BornTimestamp
	if s := ext.GetProperty(PropertyBornTime); s != "" {
		if n, ne := strconv.ParseInt(s, 10, 64); ne == nil {
			bornMilli = n
		}
	}

	// 消息时差.
	bornDiff := time.Now().Sub(time.UnixMilli(bornMilli)).Milliseconds()

	// 应消费时间.
	if bornDiff >= o.delayMilli {
		return o.pipeHandle(ext)
	}

	// 主题消息ID.
	payloadMessageId := ext.MsgId
	if s := ext.GetProperty(PropertyPayloadMessageId); s != "" {
		payloadMessageId = s
	}

	// 重发次数.
	publishCount := 1
	if s := ext.GetProperty(PropertyPublishCount); s != "" {
		if n, ne := strconv.ParseInt(s, 10, 32); ne == nil {
			publishCount = int(n) + 1
		}
	}

	// 延时级别.
	delaySeconds := int((o.delayMilli - bornDiff) / 1000)
	delayLevel := o.delaySecondsToLevel(delaySeconds)

	// 发布消息.
	msg := &rmqp.Message{Topic: Agent.GenerateTopicName(o.task.TopicName), Body: ext.Body}
	msg.WithDelayTimeLevel(delayLevel)
	msg.WithTag(o.delayTag)

	// 绑定属性.
	msg.WithProperty(PropertyBornTime, fmt.Sprintf("%v", bornMilli))
	msg.WithProperty(PropertyPayloadMessageId, payloadMessageId)
	msg.WithProperty(PropertyPublishCount, fmt.Sprintf("%v", publishCount))

	// 重新发布.
	messageId, err := internalProducer.send(ctx, msg)
	if err != nil {
		log.Error("<%s.%s> message republish error: delay-seconds=%d, delay-level=%d, origin-message-id=%s, error=%v", o.name, o.key, delaySeconds, delayLevel, ext.MsgId, err)
		return true
	}

	log.Info("<%s.%s> message republish succeed: delay-seconds=%d, delay-level=%d, origin-message-id=%s, target-message-id=%s", o.name, o.key, delaySeconds, delayLevel, ext.MsgId, messageId)
	return
}

func (o *Consumer) pipeHandle(ext *rmqp.MessageExt) (retry bool) {
	var (
		span = log.NewSpan("message.received")
		msg  = base.Pool.AcquireMessage().SetContext(span.Context())
	)

	defer span.End()

	// 基础字段.
	msg.Dequeue = int(ext.ReconsumeTimes) + 1
	msg.MessageBody = string(ext.Body)
	msg.MessageId = ext.MsgId
	msg.MessageTime = ext.BornTimestamp
	msg.PayloadMessageId = ext.MsgId

	// 消息时间.
	if s := ext.GetProperty(PropertyBornTime); s != "" {
		if n, ne := strconv.ParseInt(s, 10, 64); ne == nil {
			msg.MessageTime = n
		}
	}

	// 主题消息ID.
	if s := ext.GetProperty(PropertyPayloadMessageId); s != "" {
		msg.PayloadMessageId = s
	}

	// 投递过程.
	span.Kv().
		Add("message.received.adapter", o.name).
		Add("message.received.message.id", msg.MessageId).
		Add("message.received.message.time", msg.MessageTime).
		Add("message.received.task.id", o.id).
		Add("message.received.task.parallel", o.parallel)
	span.Logger().Info("message received: topic-name=%s, topic-tag=%s, message-id=%s, message-delay-ms=%d",
		ext.Topic, ext.GetTags(),
		ext.MsgId, span.StartTime().Sub(time.UnixMilli(msg.MessageTime)).Milliseconds())
	retry, _ = o.handler(o.task, msg)
	return
}

func (o *Consumer) pipeResume() {
	o.Lock()
	defer o.Unlock()

	if n := atomic.AddInt32(&o.processing, -1); n < o.task.Concurrency {
		if o.suspended {
			log.Info("<%s.%s> client resume", o.name, o.key)
			o.suspended = false
			o.client.Resume()
		}
	}
}

func (o *Consumer) pipeSuspend() {
	o.Lock()
	defer o.Unlock()

	if n := atomic.AddInt32(&o.processing, 1); n >= o.task.Concurrency {
		if o.suspended {
			return
		}

		log.Info("<%s.%s> client suspend", o.name, o.key)
		o.suspended = true
		o.client.Suspend()
	}
}
