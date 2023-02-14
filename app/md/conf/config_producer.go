// author: wsfuyibing <websearch@163.com>
// date: 2023-02-02

package conf

import (
	"strings"
)

type (
	// ProducerConfig
	// configurations for producer manager.
	ProducerConfig struct {
		// BucketSize
		// 数据桶尺寸.
		//
		// 每个节点允许最多存在多少条消息同步处于发送中, 此控制可以较好地
		// 解决瞬间流量问题.
		//
		// 默认: 30000
		BucketSize int `yaml:"bucket-size" json:"bucket-size"`

		// Concurrency
		// 生产者最大并发.
		//
		// 最大允许多少条消息正在向MQ服务器发送, 此控制以避免瞬间消息过多
		// 致MQ服务器死掉.
		//
		// 默认: 100
		Concurrency int32 `yaml:"concurrency" json:"concurrency"`

		// MaxRetry
		// 最大重试次数.
		//
		// 当发布消息MQ服务器时, 如果发布失败允许最大重试次数.
		//
		// 默认: 5
		MaxRetry int `yaml:"max-retry" json:"max-retry"`

		NotificationTagFailed  string `yaml:"notification-tag-failed" json:"notification-tag-failed"`
		NotificationTagSucceed string `yaml:"notification-tag-succeed" json:"notification-tag-succeed"`
		NotificationTopic      string `yaml:"notification-topic" json:"notification-topic"`

		StorePublishFailed  *bool `yaml:"store-publish-failed" json:"store-publish-failed"`
		StorePublishIgnored *bool `yaml:"store-publish-ignored" json:"store-publish-ignored"`
		StorePublishSucceed *bool `yaml:"store-publish-succeed" json:"store-publish-succeed"`
	}
)

func (o *ProducerConfig) init() *ProducerConfig {
	return o
}

func (o *ProducerConfig) initDefaults() {
	var (
		bt = true
	)

	if o.BucketSize == 0 {
		o.BucketSize = 30000
	}

	if o.Concurrency == 0 {
		o.Concurrency = 100
	}

	if o.MaxRetry == 0 {
		o.MaxRetry = 5
	}

	if o.StorePublishFailed == nil {
		o.StorePublishFailed = &bt
	}
	if o.StorePublishIgnored == nil {
		o.StorePublishIgnored = &bt
	}
	if o.StorePublishSucceed == nil {
		o.StorePublishSucceed = &bt
	}

	if o.NotificationTagFailed == "" {
		o.NotificationTagFailed = "FAILED"
	} else {
		o.NotificationTagFailed = strings.ToUpper(o.NotificationTagFailed)
	}

	if o.NotificationTagSucceed == "" {
		o.NotificationTagSucceed = "SUCCEED"
	} else {
		o.NotificationTagSucceed = strings.ToUpper(o.NotificationTagSucceed)
	}

	if o.NotificationTopic == "" {
		o.NotificationTopic = "NOTIFICATION"
	} else {
		o.NotificationTopic = strings.ToUpper(o.NotificationTopic)
	}
}
