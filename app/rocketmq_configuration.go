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
// date: 2023-03-09

package app

var (
	defaultRocketmqPrefix         = "X-"
	defaultRocketmqTopicQueueNums = 15
)

type (
	ConfigurationAdapterRocketmq interface {
		GetRocketmq() RocketmqConfiguration
	}

	RocketmqConfiguration interface {
		GetKey() string
		GetPrefix() string
		GetSecret() string
		GetServers() []string
		GetToken() string
		GetTopicQueueNums() int
	}

	rocketmqConfiguration struct {
		Servers        []string `yaml:"servers" json:"servers"`
		Prefix         *string  `yaml:"prefix" json:"prefix"`
		TopicQueueNums int      `yaml:"topic-queue-nums" json:"topic_queue_nums"`

		Key    string `yaml:"key" json:"key"`
		Secret string `yaml:"secret" json:"secret"`
		Token  string `yaml:"token" json:"token"`
	}
)

func (o *configuration) GetRocketmq() RocketmqConfiguration { return o.AdapterRocketmq }

func (o *rocketmqConfiguration) GetPrefix() string      { return *o.Prefix }
func (o *rocketmqConfiguration) GetServers() []string   { return o.Servers }
func (o *rocketmqConfiguration) GetKey() string         { return o.Key }
func (o *rocketmqConfiguration) GetSecret() string      { return o.Secret }
func (o *rocketmqConfiguration) GetToken() string       { return o.Token }
func (o *rocketmqConfiguration) GetTopicQueueNums() int { return o.TopicQueueNums }

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *rocketmqConfiguration) initDefaults() {
	if o.Prefix == nil {
		o.Prefix = &defaultRocketmqPrefix
	}
	if o.TopicQueueNums == 0 {
		o.TopicQueueNums = defaultRocketmqTopicQueueNums
	}
}
