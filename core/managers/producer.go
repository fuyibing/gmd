// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package managers

import (
	"context"
	"fmt"
	"github.com/fuyibing/gmd/v8/app"
	"github.com/fuyibing/gmd/v8/core/base"
	"github.com/fuyibing/log/v5"
	"github.com/fuyibing/util/v8/process"
	"sync/atomic"
	"time"
)

type (
	// Producer
	// 生产者管理器.
	Producer struct {
		adapter    base.ProducerManager
		bucket     Bucket
		callable   base.ProducerCallable
		name       string
		processor  process.Processor
		processing int32
	}
)

// /////////////////////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////////////////////

func (o *Producer) Publish(v *base.Payload) (err error) {
	// 消息入桶.
	if _, err = o.bucket.Add(v); err != nil {
		v.SetError(err)
		v.Release()
		return
	}

	// 立即取出.
	go o.pop()
	return
}

// /////////////////////////////////////////////////////////////////////////////
// Event methods
// /////////////////////////////////////////////////////////////////////////////

func (o *Producer) onBeforeSubprocess(_ context.Context) (ignored bool) {
	if o.callable = Container.GetProducer(); o.callable == nil {
		log.Error("producer constructor for {%s} adapter not injected", app.Config.GetAdapter())
		return true
	}

	if o.adapter == nil {
		o.adapter = o.callable()
		o.processor.Add(o.adapter.Processor())
	}

	return
}

func (o *Producer) onCallListen(ctx context.Context) (ignored bool) {
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (o *Producer) onPanic(ctx context.Context, v interface{}) {
	if spa, exists := log.Span(ctx); exists {
		spa.Logger().Fatal("<%s> %v", o.name, v)
	} else {
		log.Fatal("<%s> %v", o.name, v)
	}
}

// /////////////////////////////////////////////////////////////////////////////
// Access and constructor
// /////////////////////////////////////////////////////////////////////////////

func (o *Producer) init() *Producer {
	o.bucket = NewBucket(BucketCapacity)

	o.name = "producer-manager"
	o.processor = process.New(o.name).
		Before(o.onBeforeSubprocess).
		Callback(o.onCallListen).
		Panic(o.onPanic)

	return o
}

func (o *Producer) pop() {
	// 并发限流.
	if processing := atomic.AddInt32(&o.processing, 1); processing > app.Config.GetProducer().GetConcurrency() {
		atomic.AddInt32(&o.processing, -1)
		return
	}

	// 准备处理.
	var (
		err    error
		exists bool
		mid    string
		v      *base.Payload
	)

	// 空数据桶.
	if v, exists = o.bucket.Pop(); !exists {
		atomic.AddInt32(&o.processing, -1)
		return
	}

	var (
		span = log.NewSpanFromContext(v.GetContext(), "producer.async")
		t    = time.Now()
	)

	span.Kv().Add("producer.payload.hash", v.Hash).
		Add("producer.payload.offset", v.Offset).
		Add("producer.publish.adapter", app.Config.GetAdapter())

	defer func() {
		// 发布异常.
		if r := recover(); r != nil {
			span.Logger().Fatal("producer.payload.fatal: %v", r)

			if err == nil {
				err = fmt.Errorf("%v", r)
			}
		}

		// 结束链路.
		span.End()

		// 释放消息.
		v.SetDuration(t.Sub(time.Now())).SetError(err).SetMessageId(mid)
		v.Release()

		// 恢复并发.
		atomic.AddInt32(&o.processing, -1)
		o.pop()
	}()

	// 发布过程.
	if mid, err = o.adapter.Publish(v); err != nil {
		span.Kv().Add("producer.publish.error", err)
	} else {
		span.Kv().Add("producer.publish.message.id", mid)
	}
}
