// author: wsfuyibing <websearch@163.com>
// date: 2023-02-02

package conf

type (
	AccountRocketmqConfig struct {
		Servers       []string `yaml:"servers" json:"servers"`
		Brokers       []string `yaml:"brokers" json:"brokers"`
		Prefix        string   `yaml:"prefix" json:"prefix"`
		QueueCount    int      `yaml:"queue-count" json:"queue-count"`
		TopicTemplate string   `yaml:"topic-template" json:"topic-template"`

		Key    string
		Secret string
		Token  string
	}
)

func (o *AccountRocketmqConfig) init() *AccountRocketmqConfig {
	return o
}

func (o *AccountRocketmqConfig) initDefaults() {
	if len(o.Servers) == 0 {
		o.Servers = []string{
			"127.0.0.1:9876",
		}
	}
	if len(o.Brokers) == 0 {
		o.Brokers = []string{
			"127.0.0.1:10911",
		}
	}

	if o.QueueCount == 0 {
		o.QueueCount = 6
	}

	if o.TopicTemplate == "" {
		o.TopicTemplate = "TBW102"
	}
}
