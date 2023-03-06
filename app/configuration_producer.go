// author: wsfuyibing <websearch@163.com>
// date: 2023-02-21

package app

const (
	defaultProducerConcurrency = 100
)

type (
	ConfigurationProducer interface {
		GetProducer() ProducerConfiguration
	}

	// ProducerConfiguration
	// 生产者配置.
	ProducerConfiguration interface {
		GetConcurrency() int32
	}

	producerConfiguration struct {
		Concurrency int32 `yaml:"concurrency"`
	}
)

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *configuration) GetProducer() ProducerConfiguration { return o.Producer }

func (o *producerConfiguration) GetConcurrency() int32 { return o.Concurrency }

// /////////////////////////////////////////////////////////////
// Access and constructor
// /////////////////////////////////////////////////////////////

func (o *producerConfiguration) initDefaults() {
	if o.Concurrency == 0 {
		o.Concurrency = defaultProducerConcurrency
	}
}
