// author: wsfuyibing <websearch@163.com>
// date: 2023-02-08

package md

import (
	"context"
	"github.com/fuyibing/db/v3"
	"github.com/fuyibing/gmd/app/md/base"
	"github.com/fuyibing/gmd/app/md/conf"
	"github.com/fuyibing/gmd/app/models"
	"github.com/fuyibing/gmd/app/services"
	"github.com/fuyibing/log/v3"
	"github.com/fuyibing/log/v3/trace"
	"github.com/fuyibing/util/v2/process"
	"sync"
	"time"
)

type (
	RetryKind int
)

const (
	_ RetryKind = iota

	RetryKindMessage
	RetryKindPayload
)

type (
	// RetryManager
	// interface of retry manager.
	RetryManager interface {
		// Message
		// read waiting messages in database then call consume.
		//
		// - Waiting messages:
		//   SELECT * FROM `message` WHERE `status` = 3 LIMIT 10
		//
		// - Call retry:
		//   x := md.Boot.Retry()
		//   x.Message()
		Message()

		// Payload
		// read waiting payloads in database then call publish.
		//
		// - Waiting payloads:
		//   SELECT * FROM `payload` WHERE `status` = 3 LIMIT 10
		//
		// - Call:
		//   x := md.Boot.Retry()
		//   x.Payload()
		Payload()

		// Processor
		// return retry processor interface.
		//
		//   x := md.Boot.Retry().Processor()
		//   x.Start(ctx)
		Processor() process.Processor
	}

	retry struct {
		chm, chp chan bool
		dom, dop bool
		tkm, tkp *time.Ticker

		mu        *sync.RWMutex
		processor process.Processor
	}
)

// /////////////////////////////////////////////////////////////
// Interface methods.
// /////////////////////////////////////////////////////////////

func (o *retry) Message()                     { o.chanMessage() }
func (o *retry) Payload()                     { o.chanPayload() }
func (o *retry) Processor() process.Processor { return o.processor }

// /////////////////////////////////////////////////////////////
// Event methods.
// /////////////////////////////////////////////////////////////

// OnAfter
// called when processor stopped.
func (o *retry) OnAfter(_ context.Context) (ignored bool) {
	log.Debugf("retry manager: processor stopped")
	return
}

// OnBefore
// called when processor start.
func (o *retry) OnBefore(_ context.Context) (ignored bool) {
	log.Debugf("retry manager: start processor")
	return
}

// OnCallChannel
// listen channel signal.
func (o *retry) OnCallChannel(ctx context.Context) (ignored bool) {
	log.Debugf("retry manager: listen channel signal")

	// Create
	// channel and ticker.
	o.chm = make(chan bool)
	o.chp = make(chan bool)
	o.tkm = time.NewTicker(time.Duration(conf.Config.Retry.MessageSeconds) * time.Second)
	o.tkp = time.NewTicker(time.Duration(conf.Config.Retry.PayloadSeconds) * time.Second)

	// Unset
	// channel and ticker.
	defer func() {
		// Close and unset
		// message channel.
		close(o.chm)
		o.chm = nil

		// Close and unset
		// payload channel.
		close(o.chp)
		o.chp = nil

		// Stop and unset
		// message ticker.
		o.tkm.Stop()
		o.tkm = nil

		// Stop and unset
		// payload ticker.
		o.tkp.Stop()
		o.tkp = nil
	}()

	// Range
	// channel message.
	for {
		select {
		case <-o.chm:
			go o.CallMessage()
		case <-o.chp:
			go o.CallPayload()
		case <-o.tkm.C:
			go o.CallMessage()
		case <-o.tkp.C:
			go o.CallPayload()
		case <-ctx.Done():
			return
		}
	}
}

// OnPanic
// called with panic at runtime.
func (o *retry) OnPanic(ctx context.Context, v interface{}) {
	log.Panicfc(ctx, "retry manager: %v", v)
}

// /////////////////////////////////////////////////////////////
// Actions methods.
// /////////////////////////////////////////////////////////////

func (o *retry) CallMessage() {
	// Return
	// if process is running.
	if o.lockExists(RetryKindMessage) {
		return
	}

	// Lock
	// when begin.
	o.lockSet(RetryKindMessage)

	// Unlock
	// when end.
	redo := false
	defer func() {
		o.lockUnset(RetryKindMessage)

		// Recall
		// if not empty.
		if redo {
			o.CallMessage()
		}
	}()

	// Wait
	// message retry process.
	redo = o.SendMessages() > 0
}

func (o *retry) CallPayload() {
	// Return
	// if process is running.
	if o.lockExists(RetryKindPayload) {
		return
	}

	// Lock
	// when begin.
	o.lockSet(RetryKindPayload)

	// Unlock
	// when end.
	redo := false
	defer func() {
		o.lockUnset(RetryKindPayload)

		// Recall
		// if not empty.
		if redo {
			o.CallPayload()
		}
	}()

	// Wait
	// message retry process.
	redo = o.SendPayloads() > 0
}

func (o *retry) SendMessage(ctx context.Context, bean *models.Message, index int) {
	var (
		affects int64
		err     error
		message *base.Message
		task    *base.Task
		sess    = db.Connector.GetMasterWithContext(ctx)
		service = services.NewMessageService(sess)
	)

	// Called
	// when end.
	defer func() {
		_ = sess.Close()
	}()

	// Return error
	// if change status as processing failed.
	if affects, err = service.SetStatusAsProcessing(bean.Id); err != nil {
		log.Errorfc(ctx,
			"retry-manager: change message status as processing failed, bean-id=%d, index=%d, error=%v",
			bean.Id,
			index,
			err,
		)
		return
	}

	// Return
	// if returned affects is zero. It executed by other process or coroutine.
	if affects == 0 {
		log.Infofc(ctx,
			"retry-manager: message processing by other process or coroutine, bean-id=%d, index=%d",
			bean.Id,
			index,
		)
		return
	}

	// Return
	// if task not found.
	if task = base.Memory.GetTask(bean.TaskId); task == nil {
		log.Errorfc(ctx,
			"retry-manager: task belongs to not found, bean-id=%d, index=%d, task-id=%d",
			bean.Id,
			index,
			bean.TaskId,
		)
		return
	}

	// Prepare payload.
	message = base.Pool.AcquireMessage().SetContext(ctx)
	message.Dequeue = bean.MessageDequeue
	message.MessageBody = bean.MessageBody
	message.MessageId = bean.MessageId
	message.MessageTime = bean.MessageTime
	message.PayloadMessageId = bean.PayloadMessageId
	message.TaskId = bean.TaskId

	_ = Boot.Consumer().Container().Worker().Do(task, message)
}

func (o *retry) SendMessages() (count int) {
	var (
		ctx  context.Context
		err  error
		list []*models.Message
		wg   *sync.WaitGroup
	)

	// Return
	// if list waiting messages failed.
	if list, err = services.NewMessageService().ListWaiting(conf.Config.Retry.MessageCount); err != nil {
		log.Errorf("retry manager: list waiting message failed, error=%v", err)
		return
	}

	// Return
	// if message not found.
	if count = len(list); count == 0 {
		return
	}

	// Publish with parallel mode.
	ctx = trace.New()
	log.Infofc(ctx, "retry manager: waiting messages loaded, count=%d", count)

	wg = &sync.WaitGroup{}
	for i0, b0 := range list {
		wg.Add(1)
		c0 := trace.Child(ctx)
		go func(c1 context.Context, b1 *models.Message, i1 int) {
			defer wg.Done()
			o.SendMessage(c1, b1, i1)
		}(c0, b0, i0)
	}
	wg.Wait()

	return
}

func (o *retry) SendPayload(ctx context.Context, bean *models.Payload, index int) {
	var (
		affects  int64
		err      error
		payload  *base.Payload
		registry *base.Registry
		sess     = db.Connector.GetMasterWithContext(ctx)
		service  = services.NewPayloadService(sess)
	)

	// Called
	// when end.
	defer func() {
		_ = sess.Close()

		// Release payload.
		if payload != nil {
			payload.Release()
		}
	}()

	// Return error
	// if change status as processing failed.
	if affects, err = service.SetStatusAsProcessing(bean.Id); err != nil {
		log.Errorfc(ctx,
			"retry-manager: change payload status as processing failed, bean-id=%d, index=%d, error=%v",
			bean.Id,
			index,
			err,
		)
		return
	}

	// Return
	// if returned affects is zero. It executed by other process or coroutine.
	if affects == 0 {
		log.Infofc(ctx,
			"retry-manager: payload processing by other process or coroutine, bean-id=%d, index=%d",
			bean.Id,
			index,
		)
		return
	}

	// Return
	// if payload published successful.
	if bean.MessageId != "" {
		_, _ = service.SetStatusAsPublished(bean.Id)

		log.Infofc(ctx,
			"retry-manager: payload publish succeed by other process or coroutine, bean-id=%d, index=%d, message-id=%s",
			bean.Id,
			index,
			bean.MessageId,
		)
		return
	}

	// Return
	// if registry not found.
	if registry = base.Memory.GetRegistry(bean.RegistryId); registry == nil {
		log.Errorfc(ctx,
			"retry-manager: registry belongs to not found, bean-id=%d, index=%d, registry=%d",
			bean.Id,
			index,
			bean.RegistryId,
		)
		return
	}

	// Prepare payload.
	payload = base.Pool.AcquirePayload().SetContext(ctx)
	payload.FilterTag = registry.FilterTag
	payload.Hash = bean.Hash
	payload.MessageBody = bean.MessageBody
	payload.Offset = bean.Offset
	payload.RegistryId = bean.RegistryId
	payload.TopicName = registry.TopicName
	payload.TopicTag = registry.TopicTag
	_ = Boot.Producer().PublishDirect(payload)
}

func (o *retry) SendPayloads() (count int) {
	var (
		ctx  context.Context
		err  error
		list []*models.Payload
		wg   *sync.WaitGroup
	)

	// Return
	// if list waiting payloads failed.
	if list, err = services.NewPayloadService().ListWaiting(conf.Config.Retry.PayloadCount); err != nil {
		log.Errorf("retry manager: list waiting payload failed, error=%v", err)
		return
	}

	// Return
	// if payload not found.
	if count = len(list); count == 0 {
		return
	}

	// Publish with parallel mode.
	ctx = trace.New()
	log.Infofc(ctx, "retry manager: load waiting payloads, count=%d", count)

	wg = &sync.WaitGroup{}
	for i0, b0 := range list {
		wg.Add(1)
		c0 := trace.Child(ctx)
		go func(c1 context.Context, b1 *models.Payload, i1 int) {
			defer wg.Done()
			o.SendPayload(c1, b1, i1)
		}(c0, b0, i0)
	}
	wg.Wait()

	return
}

// /////////////////////////////////////////////////////////////
// Channel send
// /////////////////////////////////////////////////////////////

func (o *retry) chanMessage() {
	if o.processor.Healthy() && o.chm != nil {
		o.chm <- true
	}
}

func (o *retry) chanPayload() {
	if o.processor.Healthy() && o.chp != nil {
		o.chp <- true
	}
}

// /////////////////////////////////////////////////////////////
// Locker operations
// /////////////////////////////////////////////////////////////

func (o *retry) lockExists(kind RetryKind) bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	switch kind {
	case RetryKindMessage:
		return o.dom
	case RetryKindPayload:
		return o.dop
	}
	return false
}

func (o *retry) lockSet(kind RetryKind) {
	o.mu.Lock()
	defer o.mu.Unlock()
	switch kind {
	case RetryKindMessage:
		o.dom = true
	case RetryKindPayload:
		o.dop = true
	}
}

func (o *retry) lockUnset(kind RetryKind) {
	o.mu.Lock()
	defer o.mu.Unlock()
	switch kind {
	case RetryKindMessage:
		o.dom = false
	case RetryKindPayload:
		o.dop = false
	}
}

// /////////////////////////////////////////////////////////////
// Constructor methods.
// /////////////////////////////////////////////////////////////

func (o *retry) init() *retry {
	o.mu = &sync.RWMutex{}

	// Create
	// processor instance.
	o.processor = process.New("retry manager").After(
		o.OnAfter,
	).Before(
		o.OnBefore,
	).Callback(
		o.OnCallChannel,
	).Panic(o.OnPanic)

	return o
}
