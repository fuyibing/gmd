// author: wsfuyibing <websearch@163.com>
// date: 2023-02-02

package conf

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"os"
)

var (
	// Config
	// instance of configuration.
	Config *Configuration
)

type (
	Configuration struct {
		// Adapter name.
		// Accept: aliyunmns, rabbitmq, rocketmq
		Adapter Adapter `yaml:"adapter" json:"adapter"`

		// Account
		// for connection adapter.
		Account *AccountConfig `yaml:"account" json:"account"`

		Consumer *ConsumerConfig `yaml:"consumer" json:"consumer"`
		Producer *ProducerConfig `yaml:"producer" json:"producer"`
		Retry    *RetryConfig    `yaml:"retry" json:"retry"`
	}
)

func (o *Configuration) LoadJson(name string) error {
	buf, err := os.ReadFile(name)

	if err != nil {
		return err
	}

	if err = json.Unmarshal(buf, o); err == nil {
		o.initDefaults()
	}

	return err
}

func (o *Configuration) LoadYaml(name string) error {
	buf, err := os.ReadFile(name)

	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(buf, o); err == nil {
		o.initDefaults()
	}

	return err
}

func (o *Configuration) init() *Configuration {
	var err error
	for _, f := range []string{"./tmp/md.yaml", "../tmp/md.yaml", "./config/md.yaml", "../config/md.yaml"} {
		if err = o.LoadYaml(f); err == nil {
			break
		}
	}
	if err != nil {
		o.initDefaults()
	}

	return o
}

func (o *Configuration) initDefaults() {
	if o.Account == nil {
		o.Account = (&AccountConfig{}).init()
	}
	o.Account.initDefaults()

	if o.Consumer == nil {
		o.Consumer = (&ConsumerConfig{}).init()
	}
	o.Consumer.initDefaults()

	if o.Producer == nil {
		o.Producer = (&ProducerConfig{}).init()
	}
	o.Producer.initDefaults()

	if o.Retry == nil {
		o.Retry = (&RetryConfig{}).init()
	}
	o.Retry.initDefaults()
}
