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
	"sync/atomic"

	rmq "github.com/apache/rocketmq-client-go/v2"
	rmqm "github.com/apache/rocketmq-client-go/v2/primitive"
	rmqp "github.com/apache/rocketmq-client-go/v2/producer"

	"github.com/fuyibing/gmd/v8/app"
	"github.com/fuyibing/gmd/v8/md/base"
	"github.com/fuyibing/log/v5"
	"github.com/fuyibing/util/v8/process"
	"github.com/google/uuid"
	"strings"
	"time"
)

var (
	internalProducer *Producer
)

type (
	// Producer
	// 生产者.
	Producer struct {
		client     rmq.Producer
		name       string
		processor  process.Processor
		processing int32
	}
)

func NewProducer() (producer base.ProducerExecutor) {
	o := (&Producer{}).init()

	if internalProducer == nil {
		internalProducer = o
	}

	return o
}

// +---------------------------------------------------------------------------+
// + Interface methods                                                         |
// +---------------------------------------------------------------------------+

func (o *Producer) Processor() process.Processor {
	return o.processor
}

func (o *Producer) Publish(payload *base.Payload) (messageId string, err error) {
	var msg = &rmqm.Message{
		Topic: Agent.GenerateTopicName(payload.TopicName),
		Body:  []byte(payload.MessageBody),
	}

	if payload.TopicTag != "" {
		msg.WithTag(payload.TopicTag)
	}

	return o.send(payload.GetContext(), msg)
}

// +---------------------------------------------------------------------------+
// + Event methods                                                             |
// +---------------------------------------------------------------------------+

func (o *Producer) onClientBuild(_ context.Context) (ignored bool) {
	var (
		err  error
		node = strings.ToUpper(strings.ReplaceAll(uuid.NewString(), "-", ""))

		// 基础选项.
		opts = []rmqp.Option{
			rmqp.WithNsResolver(rmqm.NewPassthroughResolver(app.Config.GetRocketmq().GetServers())),
			rmqp.WithRetry(app.Config.GetProducer().GetMaxRetry()),
			rmqp.WithSendMsgTimeout(time.Duration(app.Config.GetProducer().GetTimeout()) * time.Second),
			rmqp.WithInstanceName(node),
			rmqp.WithDefaultTopicQueueNums(app.Config.GetRocketmq().GetTopicQueueNums()),
		}
	)

	// 连接鉴权.
	if app.Config.GetRocketmq().GetKey() != "" || app.Config.GetRocketmq().GetSecret() != "" || app.Config.GetRocketmq().GetToken() != "" {
		opts = append(opts, rmqp.WithCredentials(rmqm.Credentials{
			AccessKey:     app.Config.GetRocketmq().GetKey(),
			SecretKey:     app.Config.GetRocketmq().GetSecret(),
			SecurityToken: app.Config.GetRocketmq().GetToken(),
		}))
	}

	// 创建/Client.
	if o.client, err = rmq.NewProducer(opts...); err != nil {
		log.Error("<%s> client built: %v", o.name, err)
		return true
	}

	// 启动/Client.
	if err = o.client.Start(); err != nil {
		log.Error("<%s> client start: %v", o.name, err)
		return true
	}

	// 启动完成.
	log.Info("<%s> client built: node=%s", o.name, node)
	return
}

func (o *Producer) onClientShutdown(ctx context.Context) (ignored bool) {
	// 任务完成.
	// 应发布的消息全部处理完成后, 关闭连接.
	if atomic.LoadInt32(&o.processing) == 0 {
		if err := o.client.Shutdown(); err != nil {
			log.Error("<%s> client shutdown: %v", o.name, err)
		} else {
			log.Info("<%s> client shutdown", o.name)
		}
		return
	}

	// 稍后检测.
	time.Sleep(time.Millisecond * 100)
	return o.onClientShutdown(ctx)
}

func (o *Producer) onListen(ctx context.Context) (ignored bool) {
	log.Info("<%s> channel listening", o.name)
	defer log.Info("<%s> channel closed", o.name)

	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (o *Producer) onPanic(_ context.Context, v interface{}) {
	log.Fatal("<%s> %v", o.name, v)
}

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *Producer) init() *Producer {
	o.name = "rocketmq.producer"
	o.processor = process.New(o.name).
		After(o.onClientShutdown).
		Before(o.onClientBuild).
		Callback(o.onListen).
		Panic(o.onPanic)
	return o
}

func (o *Producer) send(ctx context.Context, msg *rmqm.Message) (messageId string, err error) {
	var res *rmqm.SendResult
	if res, err = o.client.SendSync(ctx, msg); err == nil {
		messageId = res.MsgID
	}
	return
}
