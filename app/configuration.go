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
		GetAdapter() string
		GetHost() string
		GetMemoryReloadSeconds() int
		GetPort() int
		GetRocketmq() RocketmqConfiguration
		GetStartedTime() time.Time
	}

	configuration struct {
		Adapter             string                 `yaml:"adapter"`
		AdapterRocketmq     *rocketmqConfiguration `yaml:"adapter-rocketmq"`
		Host                string                 `yaml:"host"`
		MemoryReloadSeconds int                    `yaml:"memory-reload-seconds"`
		Port                int                    `yaml:"port"`
		StartedTime         time.Time              `yaml:"-"`
	}
)

func (o *configuration) GetAdapter() string                 { return o.Adapter }
func (o *configuration) GetHost() string                    { return "0.0.0.0" }
func (o *configuration) GetMemoryReloadSeconds() int        { return o.MemoryReloadSeconds }
func (o *configuration) GetPort() int                       { return 8101 }
func (o *configuration) GetRocketmq() RocketmqConfiguration { return o.AdapterRocketmq }
func (o *configuration) GetStartedTime() time.Time          { return o.StartedTime }

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

// /////////////////////////////////////////////////////////////
// Configuration setter
// /////////////////////////////////////////////////////////////

func SetAdapter(adapter string) Option {
	return func(c *configuration) {
		c.Adapter = adapter
	}
}

func SetMemoryReloadSeconds(n int) Option {
	return func(c *configuration) { c.MemoryReloadSeconds = n }
}
