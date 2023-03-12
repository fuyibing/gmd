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

const (
	defaultConsumerReloadFrequency = 180
	defaultConsumerRetryFrequency  = 60
	defaultConsumerRetryLimit      = 30
)

type (
	ConfigurationConsumer interface {
		GetConsumer() ConsumerConfiguration
	}

	ConsumerConfiguration interface {
		GetReloadFrequency() int
		GetRetryFrequency() int
		GetRetryLimit() int
		GetSaveFailed() bool
		GetSaveSucceed() bool
	}

	consumerConfiguration struct {
		// Read registry and task
		// from database then update into memory.
		//
		// Default: 180 (Second)
		ReloadFrequency int `yaml:"reload-frequency" json:"reload_frequency"`

		// How many messages are read from the database when retrying
		// again.
		//
		// Default: 30
		RetryLimit int `yaml:"retry-limit" json:"retry_limit"`

		// How often to check whether the message in the database needs
		// to be retried.
		//
		// Default: 60 (Second)
		RetryFrequency int `yaml:"retry-frequency" json:"retry_frequency"`

		// Whether to save the failed message to the database.
		SaveFailed *bool `yaml:"save-failed" json:"save_failed"`

		// Whether to save the succeed message to the database.
		SaveSucceed *bool `yaml:"save-succeed" json:"save_succeed"`
	}
)

// +---------------------------------------------------------------------------+
// + Interface methods                                                         |
// +---------------------------------------------------------------------------+

func (o *configuration) GetConsumer() ConsumerConfiguration { return o.Consumer }

func (o *consumerConfiguration) GetReloadFrequency() int { return o.ReloadFrequency }
func (o *consumerConfiguration) GetRetryFrequency() int  { return o.RetryFrequency }
func (o *consumerConfiguration) GetRetryLimit() int      { return o.RetryLimit }
func (o *consumerConfiguration) GetSaveFailed() bool     { return *o.SaveFailed }
func (o *consumerConfiguration) GetSaveSucceed() bool    { return *o.SaveSucceed }

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *consumerConfiguration) initDefaults() *consumerConfiguration {
	yes := true

	if o.ReloadFrequency == 0 {
		o.ReloadFrequency = defaultConsumerReloadFrequency
	}

	if o.RetryFrequency == 0 {
		o.RetryFrequency = defaultConsumerRetryFrequency
	}
	if o.RetryLimit == 0 {
		o.RetryLimit = defaultConsumerRetryLimit
	}

	if o.SaveFailed == nil {
		o.SaveFailed = &yes
	}
	if o.SaveSucceed == nil {
		o.SaveSucceed = &yes
	}

	return o
}
