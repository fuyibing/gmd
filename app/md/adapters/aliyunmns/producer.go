// author: wsfuyibing <websearch@163.com>
// date: 2023-02-03

package aliyunmns

import (
	"context"
	mns "github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/fuyibing/gmd/app/md/base"
	"github.com/fuyibing/log/v8"
	"github.com/fuyibing/util/v8/process"
	"sync"
	"sync/atomic"
	"time"
)

type (
	// Producer
	// struct for aliyun mns producer.
	Producer struct {
		mu         *sync.RWMutex
		name       string
		processing int32
		processor  process.Processor
	}
)

func NewProducer() *Producer {
	o := (&Producer{}).init()
	return o
}

func (o *Producer) Processor() process.Processor { return o.processor }

// Publish
// send topic message to aliyunmns.
func (o *Producer) Publish(payload *base.Payload) (string, error) {
	atomic.AddInt32(&o.processing, 1)
	defer atomic.AddInt32(&o.processing, -1)

	res, err := Agent.GetTopicClient(payload.TopicName).PublishMessage(mns.MessagePublishRequest{
		MessageBody: payload.MessageBody, MessageTag: payload.FilterTag,
	})

	if err != nil {
		return "", err
	}

	return res.MessageId, nil
}

// /////////////////////////////////////////////////////////////
// Event methods
// /////////////////////////////////////////////////////////////

func (o *Producer) onAfter(ctx context.Context) (ignored bool) {
	if atomic.LoadInt32(&o.processing) > 0 {
		time.Sleep(time.Millisecond * 10)
		return o.onAfter(ctx)
	}

	log.Infof("%s: processor stopped", o.name)
	return
}

func (o *Producer) onBefore(_ context.Context) (ignored bool) {
	log.Infof("%s: start processor", o.name)
	return
}

func (o *Producer) onCaller(ctx context.Context) (ignored bool) {
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (o *Producer) onCallerAfter(_ context.Context) (ignored bool) {
	return
}

func (o *Producer) onCallerBefore(_ context.Context) (ignored bool) {
	return
}

func (o *Producer) onPanic(ctx context.Context, v interface{}) {
	log.Panicfc(ctx, "%s: %v", o.name, v)
}

// /////////////////////////////////////////////////////////////
// Construct method
// /////////////////////////////////////////////////////////////

func (o *Producer) init() *Producer {
	o.name = "aliyunmns-producer"
	o.processor = process.New(o.name).After(
		o.onAfter,
	).Before(
		o.onBefore,
	).Callback(
		o.onCallerBefore,
		o.onCaller,
		o.onCallerAfter,
	).Panic(o.onPanic)

	return o
}
