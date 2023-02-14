// author: wsfuyibing <websearch@163.com>
// date: 2023-02-03

package aliyunmns

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fuyibing/gmd/app/md/base"
	"github.com/fuyibing/gmd/app/md/conf"
	"github.com/fuyibing/log/v3"
	"github.com/fuyibing/log/v3/trace"
	"github.com/fuyibing/util/v2/process"
	"sync/atomic"
	"time"

	mns "github.com/aliyun/aliyun-mns-go-sdk"
)

type (
	// Consumer
	// struct for aliyun mns consumer.
	Consumer struct {
		name      string
		processor process.Processor

		dispatcher   func(task *base.Task, message *base.Message) (retry bool)
		id, parallel int
		task         *base.Task

		cli           mns.AliMNSQueue
		cliErr        chan error
		cliRes        chan mns.MessageReceiveResponse
		cliProcessing int32
	}
)

func NewConsumer(id, parallel int) *Consumer {
	return (&Consumer{
		id:       id,
		parallel: parallel,
	}).init()
}

func (o *Consumer) Dispatcher(v func(*base.Task, *base.Message) bool) { o.dispatcher = v }
func (o *Consumer) Processor() process.Processor                      { return o.processor }

// /////////////////////////////////////////////////////////////
// Action methods
// /////////////////////////////////////////////////////////////

func (o *Consumer) doReceivedError(err error) {
	if err != nil && !RegexMessageNotExist.MatchString(err.Error()) {
		log.Errorf("%s: receive message failed, error=%v", o.name, err)
	}
}

func (o *Consumer) doReceivedMessage(res mns.MessageReceiveResponse) {
	atomic.AddInt32(&o.cliProcessing, 1)
	defer atomic.AddInt32(&o.cliProcessing, -1)

	// Prepare
	// for message consume process.
	var (
		ctx = trace.New()
		msg *base.Message
	)

	log.Infofc(ctx, "%s: message received, dequeue=%d, message-id=%v", o.name, res.DequeueCount, res.MessageId)

	msg = base.Pool.AcquireMessage().SetContext(ctx)
	msg.Dequeue = int(res.DequeueCount)
	msg.MessageId = res.MessageId
	msg.MessageTime = res.EnqueueTime

	// Check topic message.
	if ok, mi, mb := o.parseTopicMessage(res.MessageBody); ok {
		msg.MessageBody = mb
		msg.PayloadMessageId = mi
	} else {
		msg.MessageBody = res.MessageBody
	}

	// Call dispatcher.
	if retry := o.dispatcher(o.task, msg); retry {
		o.sendRetry(ctx, res.ReceiptHandle, res.DequeueCount)
	} else {
		o.sendDelete(ctx, res.ReceiptHandle)
	}
}

func (o *Consumer) doReceiver(ctx context.Context) {
	// Return
	// if context cancelled.
	if ctx == nil || ctx.Err() != nil {
		return
	}

	// Recall
	// if channel signal creating.
	if o.cliRes == nil {
		time.Sleep(time.Millisecond)
		o.doReceiver(ctx)
		return
	}

	// Recall
	// if concurrency is greater than configured.
	if processing := atomic.LoadInt32(&o.cliProcessing); processing >= o.task.Concurrency {
		time.Sleep(time.Millisecond * 50)
		o.doReceiver(ctx)
		return
	}

	// Polling message
	// from aliyunmns queue.
	func(cli mns.AliMNSQueue, cliRes chan mns.MessageReceiveResponse, cliErr chan error) {
		defer func() { recover() }()
		cli.ReceiveMessage(cliRes, cliErr, conf.Config.Consumer.PollingWaitSeconds)
	}(o.cli, o.cliRes, o.cliErr)

	// Recall receiver.
	o.doReceiver(ctx)
}

func (o *Consumer) parseTopicMessage(str string) (yes bool, messageId, messageBody string) {
	v := TopicMessagePool.Get().(*TopicMessage)
	defer v.Release()
	if err := json.Unmarshal([]byte(str), v); err == nil && v.MessageId != "" {
		return true, v.MessageId, v.Message
	}
	return false, "", ""
}

func (o *Consumer) sendDelete(ctx context.Context, key string) {
	if err := o.cli.DeleteMessage(key); err != nil {
		log.Warnfc(ctx, "%s: delete message, error=%v", o.name, err)
	} else {
		log.Infofc(ctx, "%s: delete message", o.name)
	}
}

func (o *Consumer) sendRetry(ctx context.Context, key string, minutes int64) {
	if _, err := o.cli.ChangeMessageVisibility(key, minutes*60); err != nil {
		log.Warnfc(ctx, "%s: change queue message visibility time, minutes=%d, error=%v", o.name, minutes, err)
	} else {
		log.Infofc(ctx, "%s: change queue message visibility time, minutes=%d", o.name, minutes)
	}
}

// /////////////////////////////////////////////////////////////
// Event methods
// /////////////////////////////////////////////////////////////

// onAfter
// called when processor stopped.
func (o *Consumer) onAfter(_ context.Context) (ignored bool) {
	log.Debugf("%s: processor stopped", o.name)
	return
}

// onBefore
// called when processor start.
func (o *Consumer) onBefore(_ context.Context) (ignored bool) {
	log.Debugf("%s: start processor", o.name)
	return
}

// onCallAfter
// unset clint and task instance.
func (o *Consumer) onCallAfter(_ context.Context) (ignored bool) {
	if o.cli != nil {
		o.cli = nil
	}
	if o.task != nil {
		o.task = nil
	}
	return
}

// onCallChannel
// listen channel signal.
func (o *Consumer) onCallChannel(ctx context.Context) (ignored bool) {
	log.Debugf("%s: listen channel signal", o.name)

	// Create channel.
	o.cliErr = make(chan error)
	o.cliRes = make(chan mns.MessageReceiveResponse)

	// Close and unset channel
	// when end.
	defer func() {
		close(o.cliErr)
		o.cliErr = nil

		close(o.cliRes)
		o.cliRes = nil
	}()

	// Wait and listen
	// channel signal.
	for {
		select {
		case err := <-o.cliErr:
			go o.doReceivedError(err)
		case res := <-o.cliRes:
			go o.doReceivedMessage(res)
		case <-ctx.Done():
			return
		}
	}
}

// onCallClient
// create client and receive message in coroutine.
func (o *Consumer) onCallClient(ctx context.Context) bool {
	log.Debugf("%s: create aliyunmns consumer client", o.name)
	o.cli = Agent.GetQueueClient(o.task.Id)
	go o.doReceiver(ctx)
	return false
}

// onCallTaskCheck
// check subscription task.
func (o *Consumer) onCallTaskCheck(_ context.Context) bool {
	// Return true
	// if subscription task not found in memory.
	if o.task = base.Memory.GetTask(o.id); o.task == nil {
		log.Errorf("%s: subscription task not found", o.name)
		return true
	}

	// Return true
	// if adapter parallel is greater or equal to task parallels.
	if o.parallel >= o.task.Parallels {
		log.Errorf("%s: consumer parallels limited", o.name)
		return true
	}

	// Next
	// event callee.
	log.Debugf("%s: subscription task loaded, topic=%v, tag=%v, filter=%v, title=%s", o.name, o.task.TopicName, o.task.TopicTag, o.task.FilterTag, o.task.Title)
	return false
}

// onCallTaskSync
// sync subscribed relations to aliyunmns server.
func (o *Consumer) onCallTaskSync(ctx context.Context) bool {
	// Return true
	// if error occurred.
	if o.parallel == 0 {
		if err := Agent.Build(ctx, o.task); err != nil {
			log.Errorf("%s: sync subscribed relation to remote failed, error=%v", o.name, err)
			return true
		}
	}

	// Next
	// event callee.
	return false
}

// onCallWaiting
// recall after specified milliseconds until consume completed.
func (o *Consumer) onCallWaiting(ctx context.Context) bool {
	// Recall
	// if consume process not completed.
	if atomic.LoadInt32(&o.cliProcessing) > 0 {
		time.Sleep(conf.EventSleepDuration)
		return o.onCallWaiting(ctx)
	}

	// Next
	// event callee.
	return false
}

// onPanic
// called with panic at runtime.
func (o *Consumer) onPanic(ctx context.Context, v interface{}) {
	log.Panicfc(ctx, "%s: %v", o.name, v)
}

// /////////////////////////////////////////////////////////////
// Construct method
// /////////////////////////////////////////////////////////////

func (o *Consumer) init() *Consumer {
	o.name = fmt.Sprintf("aliyunmns-consumer-%d-%d", o.id, o.parallel)
	o.processor = process.New(o.name).After(
		o.onAfter,
	).Before(
		o.onBefore,
	).Callback(
		o.onCallTaskCheck,
		o.onCallTaskSync,
		o.onCallClient,
		o.onCallChannel,
		o.onCallWaiting,
		o.onCallAfter,
	).Panic(o.onPanic)

	return o
}
