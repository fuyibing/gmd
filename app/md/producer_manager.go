// author: wsfuyibing <websearch@163.com>
// date: 2023-02-08

package md

import (
	"context"
	"fmt"
	"github.com/fuyibing/gmd/app/md/adapters"
	"github.com/fuyibing/gmd/app/md/base"
	"github.com/fuyibing/gmd/app/md/conf"
	"github.com/fuyibing/log/v8"
	"github.com/fuyibing/util/v8/process"
	"sync"
	"sync/atomic"
	"time"
)

type (
	ProducerManager interface {
		// Bucket
		// return producer bucket interface.
		Bucket() ProducerBucket

		// Processor
		// return producer processor interface.
		Processor() process.Processor

		// Publish
		// send payloads to channel.
		//
		// Return immediately do not wait process completed.
		Publish(ps ...*base.Payload) (err error)

		// PublishDirect
		// send payload directly in sync coroutine.
		//
		// Return when process completed.
		PublishDirect(p *base.Payload) (err error)
	}

	producer struct {
		adapter               adapters.ProducerAdapter
		bucket                ProducerBucket
		ch                    chan *base.Payload
		processor             process.Processor
		publishing, releasing int32
	}
)

// /////////////////////////////////////////////////////////////
// Interface methods.
// /////////////////////////////////////////////////////////////

func (o *producer) Bucket() ProducerBucket                    { return o.bucket }
func (o *producer) Processor() process.Processor              { return o.processor }
func (o *producer) Publish(ps ...*base.Payload) (err error)   { return o.doChannel(ps...) }
func (o *producer) PublishDirect(p *base.Payload) (err error) { return o.doSend(p) }

// /////////////////////////////////////////////////////////////
// Event methods.
// /////////////////////////////////////////////////////////////

// OnAfter
// called when processor stopped.
func (o *producer) OnAfter(_ context.Context) (ignored bool) {
	log.Debugf("producer manager: processor stopped")
	return
}

// OnAfterClean
// clean bucket.
func (o *producer) OnAfterClean(ctx context.Context) (ignored bool) {
	// Pop
	// 30 messages.
	if s, n := o.bucket.Popn(30); n > 0 {
		e := fmt.Errorf("clean bucket")
		w := &sync.WaitGroup{}
		log.Debugf("producer manager: producer bucket clean %d payloads", n)

		// Release
		// with parallel.
		for _, x := range s {
			w.Add(1)
			go func(p *base.Payload) {
				defer w.Done()
				p.SetError(e)
				o.doRelease(p)
			}(x)
		}
		w.Wait()

		// Recall
		// until bucket is empty.
		return o.OnAfterClean(ctx)
	}

	// Next
	// event callee.
	log.Debugf("producer manager: producer bucket clean finish")
	return
}

// OnAfterIdle
// clean bucket.
func (o *producer) OnAfterIdle(ctx context.Context) (ignored bool) {
	// Recall
	// if any process not completed.
	if atomic.LoadInt32(&o.releasing) > 0 || atomic.LoadInt32(&o.publishing) > 0 {
		time.Sleep(conf.EventSleepDuration)
		return o.OnAfterIdle(ctx)
	}

	// Next
	// event callee.
	return
}

// OnBefore
// called when processor start.
func (o *producer) OnBefore(_ context.Context) (ignored bool) {
	log.Debugf("producer manager: start processor")
	return
}

// OnCallAdapterBuild
// build producer adapter.
func (o *producer) OnCallAdapterBuild(_ context.Context) (ignored bool) {
	var err error

	// Return error
	// if create adapter failed.
	if o.adapter, err = adapters.NewProducer(conf.Config.Adapter); err != nil {
		log.Errorf("producer manager: build %s adapter, error=%v", conf.Config.Adapter, err)
		return true
	}

	// Next
	// event callee.
	log.Debugf("producer manager: build %s adapter", conf.Config.Adapter)
	return
}

// OnCallAdapterDestroy
// wait until producer adapter stopped.
func (o *producer) OnCallAdapterDestroy(ctx context.Context) (ignored bool) {
	// Next event caller
	// if adapter processor stopped.
	if o.adapter.Processor().Stopped() {
		o.adapter = nil
		log.Debugf("producer manager: destroy %s adapter", conf.Config.Adapter)
		return
	}

	// Recall
	// wait for a while.
	time.Sleep(conf.EventSleepDuration)
	return o.OnCallAdapterDestroy(ctx)
}

// OnCallAdapterStart
// start remoter adapter in coroutine.
func (o *producer) OnCallAdapterStart(ctx context.Context) (ignored bool) {
	go func() { _ = o.adapter.Processor().Start(ctx) }()
	return
}

// OnCallChannel
// listen channel signal.
func (o *producer) OnCallChannel(ctx context.Context) (ignored bool) {
	log.Debugf("producer manager: listen channel signal")

	// Prepare
	// channel and ticker.
	o.ch = make(chan *base.Payload)
	pop := time.NewTicker(time.Second * 3)

	// Called
	// when end.
	defer func() {
		// Close channel.
		close(o.ch)
		o.ch = nil

		// Stop ticker and unset.
		pop.Stop()
		pop = nil
	}()

	// Receive
	// channel message.
	for {
		select {
		case <-pop.C:
			o.rePop(ctx)
		case p := <-o.ch:
			go o.doPush(ctx, p)
		case <-ctx.Done():
			return
		}
	}
}

// OnPanic
// called with panic at runtime.
func (o *producer) OnPanic(ctx context.Context, v interface{}) {
	log.Panicfc(ctx, "producer manager: %v", v)
}

// /////////////////////////////////////////////////////////////
// Action methods.
// /////////////////////////////////////////////////////////////

func (o *producer) doChannel(ps ...*base.Payload) (err error) {
	// Return error
	// if processor is not healthy.
	if !o.processor.Healthy() {
		err = fmt.Errorf("processor stopping or restarting")
		return
	}

	// Return error
	// if bucket is full.
	if o.bucket.IsFull() {
		err = fmt.Errorf("bucket is full")
		return
	}

	// Send payloads
	// to channel.
	for _, x := range ps {
		func(ch chan *base.Payload, p *base.Payload) {
			// Release
			// if panic occurred.
			defer func() {
				if r := recover(); r != nil {
					go o.doRelease(p.SetError(fmt.Errorf("%v", r)))
				}
			}()

			// Send channel.
			ch <- p
		}(o.ch, x)
	}
	return
}

func (o *producer) doPop(ctx context.Context) {
	var (
		payload    *base.Payload
		publishing int32
	)

	// Return
	// if context cancelled.
	if ctx == nil || ctx.Err() != nil {
		return
	}

	// Execute
	// when end.
	defer func() {
		atomic.AddInt32(&o.publishing, -1)

		// Release
		// if payload received.
		if payload != nil {
			go o.doRelease(payload)

			// Recall pop.
			o.doPop(ctx)
		}
	}()

	// Return
	// if publishing coroutines count is greater than maximum
	// configured.
	if publishing = atomic.AddInt32(&o.publishing, 1); publishing > conf.Config.Producer.Concurrency {
		return
	}

	// Return
	// if bucket is empty.
	if payload = o.bucket.Pop(); payload == nil {
		return
	}

	// Send process.
	if log.Config.DebugOn() {
		log.Debugfc(payload.GetContext(), "producer manager: pop payload from bucket")
	}
	_ = o.doSend(payload)
}

func (o *producer) doPush(ctx context.Context, payload *base.Payload) {
	// Push
	// into bucket.
	if err := o.bucket.Push(payload); err != nil {
		go o.doRelease(payload.SetError(err))
	} else {
		if log.Config.DebugOn() {
			log.Debugfc(payload.GetContext(), "producer manager: push payload into bucket")
		}
	}

	// Pop immediately.
	o.doPop(ctx)
}

func (o *producer) doRelease(p *base.Payload) {
	atomic.AddInt32(&o.releasing, 1)
	defer atomic.AddInt32(&o.releasing, -1)
	p.Release()
}

func (o *producer) doSend(p *base.Payload) (err error) {
	var (
		messageId string
		t         = time.Now()
	)

	// Called
	// when end.
	defer func() {
		d := time.Now().Sub(t).Seconds()

		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			log.Panicfc(p.GetContext(), "%v", err)
		}

		// Update
		// publish result.
		p.SetError(err).SetMessageId(messageId).SetDuration(d)

		// Logger result.
		if err != nil {
			log.Errorfc(p.GetContext(), "%s adapter: duration=%03f, error=%v", conf.Config.Adapter, d, err)
		} else {
			log.Infofc(p.GetContext(), "%s adapter: duration=%03f, message-id=%s", conf.Config.Adapter, d, messageId)
		}
	}()

	// Publish process.
	log.Infofc(p.GetContext(), "producer manager: call %s adapter and publish", conf.Config.Adapter)
	messageId, err = o.adapter.Publish(p)
	return
}

func (o *producer) rePop(ctx context.Context) {
	var (
		backlog int
		idle    int
	)

	// return
	// if no item in bucket.
	if backlog = o.bucket.Length(); backlog == 0 {
		return
	}

	// Return
	// if idle coroutines is zero.
	if idle = int(conf.Config.Producer.Concurrency - atomic.LoadInt32(&o.publishing)); idle <= 0 {
		return
	}

	// Reset
	// maximum idle coroutines.
	if idle > backlog {
		idle = backlog
	}

	// Call pop
	// in coroutines.
	log.Warnf("producer manager: redo pop goroutine, max-concurrency=%d, idle=%d, backlog=%d", conf.Config.Producer.Concurrency, idle, backlog)
	for i := 0; i < idle; i++ {
		go o.doPop(ctx)
	}
}

// /////////////////////////////////////////////////////////////
// Access methods.
// /////////////////////////////////////////////////////////////

func (o *producer) init() *producer {
	// Prepare producer bucket.
	o.bucket = (&bucket{size: conf.Config.Producer.BucketSize}).init()

	// Register producer processor event callbacks.
	o.processor = process.New("producer manager").After(
		o.OnAfterClean,
		o.OnAfterIdle,
		o.OnAfter,
	).Before(
		o.OnBefore,
	).Callback(
		o.OnCallAdapterBuild,
		o.OnCallAdapterStart,
		o.OnCallChannel,
		o.OnCallAdapterDestroy,
	).Panic(o.OnPanic)

	return o
}
