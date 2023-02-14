// author: wsfuyibing <websearch@163.com>
// date: 2023-02-02

package conf

type (
	// AccountConfig
	// configurations for adapters connection account.
	AccountConfig struct {
		Aliyunmns *AccountAliyunmnsConfig `yaml:"aliyunmns" json:"aliyunmns"`
		Rabbitmq  *AccountRabbitmqConfig  `yaml:"rabbitmq" json:"rabbitmq"`
		Rocketmq  *AccountRocketmqConfig  `yaml:"rocketmq" json:"rocketmq"`
	}
)

func (o *AccountConfig) init() *AccountConfig {
	return o
}

func (o *AccountConfig) initDefaults() {
	if o.Aliyunmns == nil {
		o.Aliyunmns = (&AccountAliyunmnsConfig{}).init()
	}
	o.Aliyunmns.initDefaults()

	if o.Rabbitmq == nil {
		o.Rabbitmq = (&AccountRabbitmqConfig{}).init()
	}
	o.Rabbitmq.initDefaults()

	if o.Rocketmq == nil {
		o.Rocketmq = (&AccountRocketmqConfig{}).init()
	}
	o.Rocketmq.initDefaults()
}
