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

package managers

import (
	"context"
	"fmt"
	"github.com/fuyibing/gmd/v8/app"
	"github.com/fuyibing/gmd/v8/md/base"
	"github.com/fuyibing/log/v5"
	"github.com/fuyibing/util/v8/process"
	"sync/atomic"
	"time"
)

type (
	ProducerManager interface {
		Processor() process.Processor
		Publish(payload *base.Payload) error
		PublishSync(payload *base.Payload) error
	}

	producer struct {
		bucket                                         Bucket
		concurrency, processing, releasing, truncating int32
		executor                                       base.ProducerExecutor
		name                                           string
		processor                                      process.Processor
	}
)

func (o *producer) Processor() process.Processor {
	return o.processor
}

func (o *producer) Publish(payload *base.Payload) (err error) {
	var (
		span  = log.NewSpanFromContext(payload.GetContext(), "payload.push.into.bucket")
		total int
	)

	span.Kv().
		Add("payload.push.hash", payload.Hash).
		Add("payload.push.offset", payload.Offset)

	defer span.End()

	if total, err = o.bucket.Add(payload); err != nil {
		span.Logger().Error("payload push into bucket: %v", err)
		o.release(payload)
		return err
	}

	span.Logger().Info("payload push into bucket: total=%d", total)
	go o.pop()
	return
}

func (o *producer) PublishSync(payload *base.Payload) (err error) {
	err = o.send(payload)
	o.release(payload)
	return
}

// +---------------------------------------------------------------------------+
// + Event methods                                                             |
// +---------------------------------------------------------------------------+

func (o *producer) onAfter(ctx context.Context) (ignored bool) {
	if atomic.LoadInt32(&o.concurrency) == 0 && atomic.LoadInt32(&o.processing) == 0 && atomic.LoadInt32(&o.releasing) == 0 && atomic.LoadInt32(&o.truncating) == 0 {
		return
	}

	time.Sleep(time.Millisecond * 100)
	return o.onAfter(ctx)
}

func (o *producer) onAdapterBound(_ context.Context) (ignored bool) {
	if _, exists := o.processor.Get(o.executor.Processor().Name()); exists {
		return
	}

	o.processor.Add(o.executor.Processor())
	return
}

func (o *producer) onAdapterCheck(_ context.Context) (ignored bool) {
	if fn := base.Container.GetProducer(); fn != nil {
		if ex := fn(); ex != nil {
			o.executor = ex
			return
		}
	}

	log.Error("<%s> adapter not injected into container", o.name)
	return true
}

func (o *producer) onBucketClean(_ context.Context) (ignored bool) {
	if count := o.bucket.Count(); count > 0 {
		if max := int(app.Config.GetProducer().GetBucketConcurrency()); count > max {
			count = max
		}
		for i := 0; i < count; i++ {
			go o.truncate()
		}
	}
	return
}

func (o *producer) onCall(ctx context.Context) (ignored bool) {
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (o *producer) onPanic(_ context.Context, v interface{}) {
	log.Fatal("<%s> %v", o.name, v)
}

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *producer) init() *producer {
	o.bucket = NewBucket(app.Config.GetProducer().GetBucketCapacity())

	o.name = "producer.manager"
	o.processor = process.New(o.name).
		After(o.onBucketClean, o.onAfter).
		Before(o.onAdapterCheck, o.onAdapterBound).
		Callback(o.onCall).
		Panic(o.onPanic)

	return o
}

// 取出消息.
// 从数据桶中取出一条消息并发布.
func (o *producer) pop() {
	// 并发限流.
	if concurrency := atomic.AddInt32(&o.concurrency, 1); concurrency > app.Config.GetProducer().GetBucketConcurrency() {
		atomic.AddInt32(&o.concurrency, -1)
		return
	}

	var (
		exists  bool
		payload *base.Payload
	)

	// 取出消息.
	// 若未取到消息(空数据桶), 取消此协程, 最低系统资源.
	if payload, exists = o.bucket.Pop(); !exists {
		atomic.AddInt32(&o.concurrency, -1)
		return
	}

	// 发布消息.
	_ = o.send(payload)
	atomic.AddInt32(&o.concurrency, -1)
	o.pop()
}

// 异步释放.
func (o *producer) release(payload *base.Payload) {
	atomic.AddInt32(&o.releasing, 1)
	go func() {
		defer atomic.AddInt32(&o.releasing, -1)
		payload.Release()
	}()
}

// 发布过程.
// 同步阻塞, 结束后释放消息实例.
func (o *producer) send(payload *base.Payload) (err error) {
	// 内部计数.
	//
	// - 开始时/+1
	// - 结束时/-1
	atomic.AddInt32(&o.processing, 1)
	defer atomic.AddInt32(&o.processing, -1)

	var (
		messageId string
		span      = log.NewSpanFromContext(payload.GetContext(), "payload.publish")
	)

	defer span.End()

	span.Kv().
		Add("payload.publish.adapter", o.executor.Processor().Name()).
		Add("payload.publish.hash", payload.Hash).
		Add("payload.publish.offset", payload.Offset)

	// 结束发布.
	defer func() {
		// 发布异常.
		if r := recover(); r != nil {
			span.Logger().Fatal("payload publish: %v", r)

			// 重置错误.
			if err == nil {
				err = fmt.Errorf("%v", r)
			}
		}

		// 记录结果.
		if err != nil {
			span.Logger().Error("payload publish error: %v", err)
		} else {
			span.Logger().Info("payload publish: message-id=%s", messageId)
		}

		// 释放消息.
		payload.SetDuration(time.Now().Sub(span.StartTime())).SetError(err).SetMessageId(messageId)
		o.release(payload)
	}()

	// 禁止发布.
	// 适配器处于中间状态: 启动中/重启中/退出中.
	if !o.executor.Processor().Healthy() {
		err = fmt.Errorf("adapter not healthy")
		return
	}

	// 发布过程.
	// 调用具体的适配器(Rocketmq, AliyunMNS等)执行发布过程.
	messageId, err = o.executor.Publish(payload)
	return
}

// 清数据桶.
func (o *producer) truncate() {
	// 并发限流.
	if truncating := atomic.AddInt32(&o.truncating, 1); truncating > app.Config.GetProducer().GetBucketConcurrency() {
		atomic.AddInt32(&o.truncating, -1)
		return
	}

	// 取出消息.
	payload, exists := o.bucket.Pop()

	// 退出协程.
	if !exists {
		atomic.AddInt32(&o.truncating, -1)
		return
	}

	// 释放消息.
	payload.Release()
	atomic.AddInt32(&o.truncating, -1)
	o.truncate()
}
