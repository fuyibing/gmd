// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package conf

import (
	"github.com/fuyibing/gmd/v8/app/md/base"
	"gopkg.in/yaml.v3"
	"os"
)

type (
	configuration struct {
		Adapter base.Adapter `yaml:"adapter"`
	}
)

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *configuration) init() *configuration {
	o.initYaml()
	return o
}

func (o *configuration) initDefaults() {}

func (o *configuration) initYaml() {
	for _, path := range []string{
		"config/md.yaml", "../config/md.yaml",
		"tmp/md.yaml", "../tmp/md.yaml",
	} {
		if buf, err := os.ReadFile(path); err == nil {
			if yaml.Unmarshal(buf, o) == nil {
				break
			}
		}
	}
}
