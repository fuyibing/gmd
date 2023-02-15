// author: wsfuyibing <websearch@163.com>
// date: 2023-02-03

package rocketmq

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/fuyibing/gmd/app/md/base"
	"github.com/fuyibing/gmd/app/md/conf"
	"github.com/fuyibing/log/v3"
	"github.com/fuyibing/util/v2/process"
	"github.com/google/uuid"
	"strings"
	"time"

	sdk "github.com/apache/rocketmq-client-go/v2"
)

var (
	defaultProducer *Producer
)

type (
	Producer struct {
		name      string
		processor process.Processor

		client sdk.Producer
	}
)

func NewProducer() *Producer {
	o := (&Producer{}).init()
	if defaultProducer == nil {
		defaultProducer = o
	}
	return o
}

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *Producer) Processor() process.Processor            { return o.processor }
func (o *Producer) Publish(p *base.Payload) (string, error) { return o.doPublish(p) }

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *Producer) doPublish(p *base.Payload) (string, error) {
	if !o.processor.Healthy() {
		return "", fmt.Errorf("producer is starting or restarting")
	}
	return o.doSend(p.GetContext(), (&primitive.Message{
		Topic: Agent.GenTopicName(p.TopicName),
		Body:  []byte(p.MessageBody),
	}).WithTag(p.TopicTag))
}

func (o *Producer) doSend(ctx context.Context, m *primitive.Message) (string, error) {
	var (
		err error
		res *primitive.SendResult
	)
	if res, err = o.client.SendSync(ctx, m); err != nil {
		return "", err
	}
	return res.MsgID, nil
}

// /////////////////////////////////////////////////////////////
// Constructor method
// /////////////////////////////////////////////////////////////

func (o *Producer) init() *Producer {
	o.name = "rocketmq-producer"
	o.processor = process.New(o.name).After(
		o.onAfter,
	).Before(
		o.onBefore,
	).Callback(
		o.onCallClientBuild,
		o.onCallClientStart,
		o.onCallChannel,
		o.onCallClientDestroy,
	).Panic(
		o.onPanic,
	)

	return o
}

// /////////////////////////////////////////////////////////////
// Processor events
// /////////////////////////////////////////////////////////////

func (o *Producer) onAfter(_ context.Context) (ignored bool) {
	log.Debugf("%s processor stopped", o.name)
	return
}

func (o *Producer) onBefore(_ context.Context) (ignored bool) {
	log.Debugf("%s start processor", o.name)
	return
}

func (o *Producer) onCallChannel(ctx context.Context) (ignored bool) {
	log.Debugf("%s: listen channel signal", o.name)

	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (o *Producer) onCallClientBuild(_ context.Context) (ignored bool) {
	var (
		err  error
		node = strings.ToUpper(strings.ReplaceAll(uuid.NewString(), "-", ""))
		op   = []producer.Option{
			producer.WithNsResolver(primitive.NewPassthroughResolver(conf.Config.Account.Rocketmq.Servers)),
			producer.WithDefaultTopicQueueNums(conf.Config.Account.Rocketmq.QueueCount),
			producer.WithRetry(conf.Config.Producer.MaxRetry),
			producer.WithSendMsgTimeout(time.Duration(conf.Config.Producer.PublishTimeout) * time.Second),
			producer.WithInstanceName(node),
			producer.WithGroupName(DefaultProducerGroupName),
		}
	)

	// Extension
	if conf.Config.Account.Rocketmq.TopicTemplate != "" {
		op = append(op, producer.WithCreateTopicKey(
			conf.Config.Account.Rocketmq.TopicTemplate,
		))
	}

	// Append authentication.
	if conf.Config.Account.Rocketmq.Key != "" || conf.Config.Account.Rocketmq.Secret != "" || conf.Config.Account.Rocketmq.Token != "" {
		op = append(op, producer.WithCredentials(primitive.Credentials{
			AccessKey:     conf.Config.Account.Rocketmq.Key,
			SecretKey:     conf.Config.Account.Rocketmq.Secret,
			SecurityToken: conf.Config.Account.Rocketmq.Token,
		}))
	}

	// Return if client built failed.
	if o.client, err = sdk.NewProducer(op...); err != nil {
		log.Errorf("%s: client build failed, error=%v", o.name)
		return true
	}

	// Client create succeed.
	log.Infof("%s: client built, name=%s, node=%s", o.name, DefaultProducerGroupName, node)
	return
}

func (o *Producer) onCallClientDestroy(_ context.Context) (ignored bool) {
	if err := o.client.Shutdown(); err != nil {
		log.Errorf("%s: client shutdown failed, error=%v", o.name, err)
	}

	o.client = nil
	return
}

func (o *Producer) onCallClientStart(_ context.Context) (ignored bool) {
	if err := o.client.Start(); err != nil {
		log.Errorf("%s: client start failed: error=%v", o.name, err)
		return false
	}

	log.Debugf("%s: client started", o.name)
	return
}

func (o *Producer) onPanic(ctx context.Context, v interface{}) {
	log.Panicfc(ctx, "%s panic: %v", o.name, v)
}
