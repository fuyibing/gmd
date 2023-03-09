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
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

var (
	Config Configuration
)

const (
	DefaultConfigurationName    = "gmd"
	DefaultConfigurationVersion = "1.0"
	DefaultConfigurationPort    = 8101
)

type (
	// Configuration
	// 全局配置.
	Configuration interface {
		GetHost() string
		GetName() string
		GetPort() int
		GetSoftware() string
		GetStartTime() time.Time
		GetVersion() string

		GetAdapter() string
		ConfigurationAdapterRocketmq

		ConfigurationConsumer
		ConfigurationProducer
	}

	configuration struct {
		Name    string `yaml:"-" json:"name"`
		Host    string `yaml:"host" json:"host"`
		Port    int    `yaml:"port" json:"port"`
		Version string `yaml:"-" json:"version"`

		// +-------------------------------------------------------------------+
		// + Advanced configurations                                           |
		// +-------------------------------------------------------------------+

		Consumer *consumerConfiguration `yaml:"consumer" json:"consumer"`
		Producer *producerConfiguration `yaml:"producer" json:"producer"`

		// +-------------------------------------------------------------------+
		// + Adapter (AliyunMNS, Rocketmq, RabbitMQ)                           |
		// +-------------------------------------------------------------------+

		Adapter         string                 `yaml:"adapter" json:"adapter"`
		AdapterRocketmq *rocketmqConfiguration `yaml:"adapter-rocketmq" json:"adapter_rocketmq"`

		// +-------------------------------------------------------------------+
		// + Internal                                                          |
		// +-------------------------------------------------------------------+

		software  string
		startTime time.Time
	}
)

func (o *configuration) GetAdapter() string      { return o.Adapter }
func (o *configuration) GetHost() string         { return o.Host }
func (o *configuration) GetName() string         { return o.Name }
func (o *configuration) GetPort() int            { return o.Port }
func (o *configuration) GetSoftware() string     { return o.software }
func (o *configuration) GetStartTime() time.Time { return o.startTime }
func (o *configuration) GetVersion() string      { return o.Version }

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *configuration) init() *configuration {
	o.scan()
	o.initDefaults()

	o.initAdapterRocketmq()

	o.initConsumer()
	o.initProducer()
	return o
}

func (o *configuration) initAdapterRocketmq() {
	if o.AdapterRocketmq == nil {
		o.AdapterRocketmq = &rocketmqConfiguration{}
	}
	o.AdapterRocketmq.initDefaults()
}

func (o *configuration) initDefaults() {
	o.Name = DefaultConfigurationName
	o.Port = DefaultConfigurationPort
	o.Version = DefaultConfigurationVersion

	o.software = fmt.Sprintf("%v/%v", o.Name, o.Version)
	o.startTime = time.Now()
}

func (o *configuration) initConsumer() {
	if o.Consumer == nil {
		o.Consumer = &consumerConfiguration{}
	}
	o.Consumer.initDefaults()
}

func (o *configuration) initProducer() {
	if o.Producer == nil {
		o.Producer = &producerConfiguration{}
	}
	o.Producer.initDefaults()
}

func (o *configuration) scan() {
	for _, path := range []string{
		"config/app.yaml", "../config/app.yaml",
		"tmp/app.yaml", "../tmp/app.yaml",
	} {
		if buf, err := os.ReadFile(path); err == nil {
			if yaml.Unmarshal(buf, o) == nil {
				break
			}
		}
	}
}
