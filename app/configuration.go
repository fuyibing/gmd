// author: wsfuyibing <websearch@163.com>
// date: 2023-02-20

package app

import (
	"gopkg.in/yaml.v3"
	"os"
)

type configuration struct {
	Adapter             string `yaml:"adapter"`
	MemoryReloadSeconds int    `yaml:"memory-reload-seconds"`
}

func (o *configuration) init() *configuration {
	o.initYaml()
	o.initDefaults()
	return o
}

func (o *configuration) initDefaults() {
	if o.MemoryReloadSeconds == 0 {
		o.MemoryReloadSeconds = DefaultMemoryReloadSeconds
	}
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
