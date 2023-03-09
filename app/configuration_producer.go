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

package app

import (
	"strings"
)

const (
	defaultProducerBucketBatch       = 100
	defaultProducerBucketCapacity    = 30000
	defaultProducerBucketConcurrency = 50
	defaultProducerBucketFrequency   = 200

	defaultProducerMaxRetry = 30

	defaultProducerNotifyTagFailed  = "FAILED"
	defaultProducerNotifyTagSucceed = "SUCCEED"
	defaultProducerNotifyTopic      = "NOTIFICATION"
)

type (
	ConfigurationProducer interface {
		GetProducer() ProducerConfiguration
	}

	ProducerConfiguration interface {
		GetBucketBatch() int
		GetBucketCapacity() int
		GetBucketConcurrency() int32
		GetBucketFrequency() int
		GetMaxRetry() int
		GetNotifyTagFailed() string
		GetNotifyTagSucceed() string
		GetNotifyTopic() string
		GetSaveFailed() bool
		GetSaveSucceed() bool
		GetTimeout() int
	}

	producerConfiguration struct {
		BucketCapacity    int   `yaml:"bucket-capacity" json:"bucket_capacity"`
		BucketConcurrency int32 `yaml:"bucket-concurrency" json:"bucket_concurrency"`
		BucketBatch       int   `yaml:"bucket-batch" json:"bucket_batch"`
		BucketFrequency   int   `yaml:"bucket-frequency" json:"bucket_frequency"`

		MaxRetry    int   `yaml:"max-retry" json:"max_retry"`
		SaveFailed  *bool `yaml:"save-failed" json:"save_failed"`
		SaveSucceed *bool `yaml:"save-succeed" json:"save_succeed"`
		Timeout     int   `yaml:"timeout" json:"timeout"`

		NotifyTagFailed  string `yaml:"notify-tag-failed" json:"notify_tag_failed"`
		NotifyTagSucceed string `yaml:"notify-tag-succeed" json:"notify_tag_succeed"`
		NotifyTopic      string `yaml:"notify-topic" json:"notify_topic"`
	}
)

// +---------------------------------------------------------------------------+
// + Interface methods                                                         |
// +---------------------------------------------------------------------------+

func (o *configuration) GetProducer() ProducerConfiguration { return o.Producer }

func (o *producerConfiguration) GetBucketBatch() int         { return o.BucketBatch }
func (o *producerConfiguration) GetBucketCapacity() int      { return o.BucketCapacity }
func (o *producerConfiguration) GetBucketConcurrency() int32 { return o.BucketConcurrency }
func (o *producerConfiguration) GetBucketFrequency() int     { return o.BucketFrequency }
func (o *producerConfiguration) GetMaxRetry() int            { return o.MaxRetry }
func (o *producerConfiguration) GetNotifyTagFailed() string  { return o.NotifyTagFailed }
func (o *producerConfiguration) GetNotifyTagSucceed() string { return o.NotifyTagSucceed }
func (o *producerConfiguration) GetNotifyTopic() string      { return o.NotifyTopic }
func (o *producerConfiguration) GetSaveFailed() bool         { return *o.SaveFailed }
func (o *producerConfiguration) GetSaveSucceed() bool        { return *o.SaveSucceed }
func (o *producerConfiguration) GetTimeout() int             { return o.Timeout }

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *producerConfiguration) initDefaults() *producerConfiguration {
	if o.BucketBatch == 0 {
		o.BucketBatch = defaultProducerBucketBatch
	}
	if o.BucketCapacity == 0 {
		o.BucketCapacity = defaultProducerBucketCapacity
	}
	if o.BucketConcurrency == 0 {
		o.BucketConcurrency = defaultProducerBucketConcurrency
	}
	if o.BucketFrequency == 0 {
		o.BucketFrequency = defaultProducerBucketFrequency
	}

	yes := true

	if o.SaveFailed == nil {
		o.SaveFailed = &yes
	}
	if o.SaveSucceed == nil {
		o.SaveSucceed = &yes
	}

	if o.MaxRetry == 0 {
		o.MaxRetry = defaultProducerMaxRetry
	}

	if o.NotifyTagFailed == "" {
		o.NotifyTagFailed = defaultProducerNotifyTagFailed
	} else {
		o.NotifyTagFailed = strings.ToUpper(o.NotifyTagFailed)
	}

	if o.NotifyTagSucceed == "" {
		o.NotifyTagSucceed = defaultProducerNotifyTagSucceed
	} else {
		o.NotifyTagSucceed = strings.ToUpper(o.NotifyTagSucceed)
	}

	if o.NotifyTopic == "" {
		o.NotifyTopic = defaultProducerNotifyTopic
	} else {
		o.NotifyTopic = strings.ToUpper(o.NotifyTopic)
	}

	return o
}
