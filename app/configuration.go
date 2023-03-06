// author: wsfuyibing <websearch@163.com>
// date: 2023-02-20

package app

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type (
	// Configuration
	// 应用配置接口.
	Configuration interface {
		GetHost() string
		GetPort() int
		GetName() string
		GetVersion() string

		GetAdapter() string
		GetMemoryReloadSeconds() int
		GetRocketmq() RocketmqConfiguration
		GetStartedTime() time.Time

		ConfigurationProducer
	}

	configuration struct {
		Name    string `yaml:"name"`
		Host    string `yaml:"host"`
		Port    int    `yaml:"port"`
		Version string `yaml:"version"`

		Adapter         string                 `yaml:"adapter"`
		AdapterRocketmq *rocketmqConfiguration `yaml:"adapter-rocketmq"`

		Producer *producerConfiguration `yaml:"producer"`

		MemoryReloadSeconds int       `yaml:"memory-reload-seconds"`
		StartedTime         time.Time `yaml:"-"`
	}
)

func (o *configuration) GetName() string    { return o.Name }
func (o *configuration) GetHost() string    { return o.Host }
func (o *configuration) GetPort() int       { return o.Port }
func (o *configuration) GetVersion() string { return o.Version }

func (o *configuration) GetStartedTime() time.Time { return o.StartedTime }

func (o *configuration) GetAdapter() string                 { return o.Adapter }
func (o *configuration) GetRocketmq() RocketmqConfiguration { return o.AdapterRocketmq }

func (o *configuration) GetMemoryReloadSeconds() int { return o.MemoryReloadSeconds }

// /////////////////////////////////////////////////////////////
// Configuration initialize
// /////////////////////////////////////////////////////////////

func (o *configuration) init() *configuration {
	o.initYaml()
	o.initDefaults()
	o.initExtensions()

	o.StartedTime = time.Now()
	return o
}

func (o *configuration) initDefaults() {
	if o.MemoryReloadSeconds == 0 {
		o.MemoryReloadSeconds = DefaultMemoryReloadSeconds
	}
}

func (o *configuration) initExtensions() {
	if o.AdapterRocketmq == nil {
		o.AdapterRocketmq = &rocketmqConfiguration{}
	}
	o.AdapterRocketmq.initDefaults()

	// 消费者.

	// 生产者.
	if o.Producer == nil {
		o.Producer = &producerConfiguration{}
	}
	o.Producer.initDefaults()
}

func (o *configuration) initYaml() {
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
