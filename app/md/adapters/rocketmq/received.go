// author: wsfuyibing <websearch@163.com>
// date: 2023-02-15

package rocketmq

import (
	"context"
	"fmt"
	sdk "github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/fuyibing/gmd/app/md/base"
	"github.com/fuyibing/log/v8"
	"strconv"
	"sync/atomic"
	"time"
)

type (
	Received struct {
		client              sdk.PushConsumer
		consuming           int32
		delayer             bool
		delayerMilliSeconds int64
		delayerTag          string
		dispatcher          func(task *base.Task, message *base.Message) (retry bool)
		name                string
		selector            consumer.MessageSelector
		task                *base.Task
		topic               string

		callbackResume, callbackSuspend func()
	}
)

// Consume
// main process for received message.
func (o *Received) Consume(ctx context.Context, ms ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	// Return error
	// if context cancelled.
	if ctx != nil && ctx.Err() != nil {
		return consumer.ConsumeRetryLater,
			ctx.Err()
	}

	// Return error
	// if received message count is not allowed.
	if n := len(ms); n != 1 {
		return consumer.ConsumeRetryLater,
			fmt.Errorf("received message count not one")
	}

	// Check message.
	return o.doCheck(ctx, ms[0])
}

// IsIdle
// return all delivering messages are completed or not.
func (o *Received) IsIdle() bool {
	return atomic.LoadInt32(&o.consuming) == 0
}

// /////////////////////////////////////////////////////////////
// Consume methods
// /////////////////////////////////////////////////////////////

func (o *Received) doCheck(ctx context.Context, m *primitive.MessageExt) (consumer.ConsumeResult, error) {
	// Consume immediately
	// if subscription task not configure delay time.
	if !o.delayer {
		return o.doConsume(m)
	}

	// Time diff.
	var (
		bornTime    = m.BornTimestamp
		currTime    = time.Now().UnixMilli()
		diffSeconds int
	)

	// Read message born time
	// from property.
	if s := m.GetProperty(DefaultDelayMessageTime); s != "" {
		if n, ne := strconv.ParseInt(s, 0, 64); ne == nil && n > 0 {
			bornTime = n
		}
	}

	// Consume immediately
	// if the current time minus the message time is greater than delay time.
	if diffSeconds = int((o.delayerMilliSeconds + bornTime - currTime) / 1000); diffSeconds <= 0 {
		return o.doConsume(m)
	}

	// Publish delay message.
	level := o.delaySecondsToLevel(diffSeconds)
	return o.doPublish(ctx, m, bornTime, diffSeconds, level)
}

func (o *Received) doConsume(m *primitive.MessageExt) (consumer.ConsumeResult, error) {
	// Increment consuming count, Call suspend
	// if concurrency is greater than configuration.
	if concurrency := atomic.AddInt32(&o.consuming, 1); concurrency >= o.task.Concurrency {
		if o.callbackSuspend != nil {
			o.callbackSuspend()
		}
	}

	// Decrement consuming count, call resume
	// if concurrency is greater than configuration.
	defer func() {
		if concurrency := atomic.AddInt32(&o.consuming, -1); concurrency < o.task.Concurrency {
			if o.callbackResume != nil {
				o.callbackResume()
			}
		}
	}()

	// Prepare
	// for message consume process.
	var (
		ctx            = log.NewContext()
		msg            *base.Message
		topicMessageId = m.MsgId
		bornTime       = m.BornTimestamp
		currTime       = time.Now().UnixMilli()
		diffTime       int64
		consumeTimes   = int(m.ReconsumeTimes + 1)
		information    = ""
	)

	// Parse topic message id.
	if s := m.GetProperty(DefaultTopicMessageId); s != "" {
		topicMessageId = s
	}

	// Parse born time.
	if s := m.GetProperty(DefaultDelayMessageTime); s != "" {
		if n, ne := strconv.ParseInt(s, 0, 64); ne == nil && n > 0 {
			bornTime = n
		}
	}

	diffTime = currTime - bornTime

	// Info: message + topic + delay time.
	information += fmt.Sprintf("message received, Message {Dequeue=%d, MessageId=%s, MessageTag=%s}",
		consumeTimes, m.MsgId, m.GetTags(),
	) + fmt.Sprintf(", Topic {Name=%s, MessageId=%s}",
		m.Topic, topicMessageId,
	) + fmt.Sprintf(", DelayTime {Expected=%v, Really=%v}",
		o.delayerMilliSeconds, diffTime,
	)

	if m.Queue != nil {
		information += fmt.Sprintf(", Queue {Queue=%d, broker=%s}",
			m.Queue.QueueId, m.Queue.BrokerName,
		)
	}

	// Info: logger trace.
	log.Infofc(ctx, "%s: %s", o.name, information)

	// Acquire message instance.
	msg = base.Pool.AcquireMessage().SetContext(ctx)
	msg.Dequeue = consumeTimes
	msg.MessageId = m.MsgId
	msg.MessageTime = bornTime
	msg.MessageBody = string(m.Body)
	msg.PayloadMessageId = topicMessageId

	// Keyword trace.
	if k := m.GetKeys(); k != "" {
		msg.Keyword = k
	}

	// Call dispatcher.
	if retry := o.dispatcher(o.task, msg); retry {
		log.Infofc(ctx, "%s: consume later", o.name)
		return consumer.ConsumeRetryLater, nil
	}

	// Return succeed response.
	log.Infofc(ctx, "%s: consume succeed", o.name)
	return consumer.ConsumeSuccess, nil
}

func (o *Received) doPublish(ctx context.Context, m *primitive.MessageExt, bt int64, seconds, level int) (consumer.ConsumeResult, error) {
	var (
		bornTime       = fmt.Sprintf("%v", bt)
		err            error
		information    = ""
		messageId      string
		publishCount   = "1"
		topicMessageId = m.MsgId
	)

	// Parse topic message id.
	if s := m.GetProperty(DefaultTopicMessageId); s != "" {
		topicMessageId = s
	}

	// Parse publish count.
	if s := m.GetProperty(DefaultDelayPublishCount); s != "" {
		if n, ne := strconv.ParseInt(s, 0, 32); ne == nil && n > 0 {
			publishCount = fmt.Sprintf("%v", n+1)
		}
	}

	information += fmt.Sprintf("message publish delay, Message {Dequeue=%d, MessageId=%s, MessageTag=%s}",
		m.ReconsumeTimes+1, m.MsgId, m.GetTags(),
	) + fmt.Sprintf(", Topic {Name=%s, MessageId=%s}",
		m.Topic, topicMessageId,
	) + fmt.Sprintf(", Target {Topic=%s, MessageTag=%s}",
		o.topic, o.delayerTag,
	) + fmt.Sprintf(", DelayTime {Publish=%v, Seconds=%d, Level=%d}",
		publishCount, seconds, level,
	)

	// Generate message param.
	x := &primitive.Message{Topic: o.topic, Body: m.Body}
	x.WithProperty(DefaultDelayPublishCount, publishCount)
	x.WithProperty(DefaultDelayMessageTime, bornTime)
	x.WithProperty(DefaultTopicMessageId, topicMessageId)
	x.WithTag(o.delayerTag)
	x.WithDelayTimeLevel(level)

	// Copy key.
	if k := m.GetKeys(); k != "" {
		x.WithKeys([]string{k})
	}

	// Delay message publish failed.
	if messageId, err = defaultProducer.doSend(ctx, x); err != nil {
		log.Errorf("%s: %s, error=%v", o.name, information, err)
		return consumer.ConsumeRetryLater, nil
	}

	// Delay message published.
	log.Infof("%s: %s, Target {MessageId=%s}", o.name, information, messageId)
	return consumer.ConsumeSuccess, nil
}

func (o *Received) delaySecondsToLevel(s int) (l int) {
	if s < 5 {
		l = 1
	} else if s >= 5 && s < 10 {
		l = 2
	} else if s >= 10 && s < 30 {
		l = 3
	} else if s >= 30 && s < 60 {
		l = 4
	} else if s >= 60 && s < 120 {
		l = 5
	} else if s >= 120 && s < 180 {
		l = 6
	} else if s >= 180 && s < 240 {
		l = 7
	} else if s >= 240 && s < 300 {
		l = 8
	} else if s >= 300 && s < 360 {
		l = 9
	} else if s >= 360 && s < 360 {
		l = 10
	} else if s >= 420 && s < 360 {
		l = 11
	} else if s >= 480 && s < 360 {
		l = 12
	} else if s >= 540 && s < 360 {
		l = 13
	} else if s >= 600 && s < 1200 {
		l = 14
	} else if s >= 1200 && s < 1800 {
		l = 15
	} else if s >= 1800 && s < 3600 {
		l = 16
	} else if s >= 3600 && s < 7200 {
		l = 17
	} else if s >= 7200 {
		l = 18
	}
	return
}
