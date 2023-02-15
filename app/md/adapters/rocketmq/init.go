// author: wsfuyibing <websearch@163.com>
// date: 2023-02-12

// Package rocketmq
// Message queue adapter on RocketMQ.
package rocketmq

import (
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"sync"
	"time"
)

const (
	DefaultConsumerGroupName      = "GID"
	DefaultProducerGroupName      = "GmdProducers"
	DefaultConsumeSuspendDuration = time.Millisecond * 10
	DefaultReconsumeTimes         = 5

	DefaultDelayTagPrefix    = "GMD-DELAY-"
	DefaultDelayMessageTime  = "GMD_DELAY_MESSAGE_TIME"
	DefaultDelayPublishCount = "GMD_DELAY_PUBLISH_COUNT"
	DefaultTopicMessageId    = "GMD_TOPIC_MESSAGE_ID"
)

func init() {
	new(sync.Once).Do(func() {
		Agent = (&agent{}).init()

		// rlog.SetLogLevel("debug")
		// rlog.SetLogLevel("error")
		rlog.SetLogLevel("fatal")
	})
}
