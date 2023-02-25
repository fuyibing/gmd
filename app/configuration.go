// author: wsfuyibing <websearch@163.com>
// date: 2023-02-20

package app

import (
	"gopkg.in/yaml.v3"
	"os"
)

type (
	// Configuration
	// exported methods, like getter.
	Configuration interface {
		GetAdapter() string
		GetMemoryReloadSeconds() int
		GetRocketmq() RocketmqConfiguration
	}

	configuration struct {
		Adapter             string                 `yaml:"adapter"`
		AdapterRocketmq     *rocketmqConfiguration `yaml:"adapter-rocketmq"`
		MemoryReloadSeconds int                    `yaml:"memory-reload-seconds"`
	}
)

// /////////////////////////////////////////////////////////////
// Configuration initialize
// /////////////////////////////////////////////////////////////

func (o *configuration) init() *configuration {
	o.initYaml()
	o.initDefaults()
	o.initExtensions()
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
// Configuration getter
// /////////////////////////////////////////////////////////////

func (o *configuration) GetAdapter() string                 { return o.Adapter }
func (o *configuration) GetMemoryReloadSeconds() int        { return o.MemoryReloadSeconds }
func (o *configuration) GetRocketmq() RocketmqConfiguration { return o.AdapterRocketmq }

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
