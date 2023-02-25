// author: wsfuyibing <websearch@163.com>
// date: 2023-02-21

package app

const (
	rocketmqDefaultMaxRetry         = 30
	rocketmqDefaultPrefix           = "X-"
	rocketmqDefaultServer           = "127.0.0.1:9092"
	rocketmqDefaultConsumeSuspendMs = 100
)

type (
	RocketmqConfiguration interface {
		GetConsumeMaxRetry() int32
		GetConsumePosition() int
		GetConsumeSuspendMs() int
		GetKey() string
		GetPrefix() string
		GetSecret() string
		GetServers() []string
		GetToken() string
	}

	rocketmqConfiguration struct {
		Servers []string `yaml:"servers"`
		Prefix  string   `yaml:"prefix"`

		// Consuming
		// point on consumer booting.
		//
		// 0: LastOffset
		// 1: FirstOffset
		// 2: Timestamp
		ConsumePosition int `yaml:"consume-position"`

		ConsumeMaxRetry int32 `yaml:"consume-max-retry"`

		ConsumeSuspendMs int `yaml:"consume-suspend-ms"`

		Key    string `yaml:"key"`
		Secret string `yaml:"secret"`
		Token  string `yaml:"token"`
	}
)

// /////////////////////////////////////////////////////////////
// Rocketmq initialize
// /////////////////////////////////////////////////////////////

func (o *rocketmqConfiguration) initDefaults() {
	if o.ConsumeMaxRetry == 0 {
		o.ConsumeMaxRetry = rocketmqDefaultMaxRetry
	}
	if o.Prefix == "" {
		o.Prefix = rocketmqDefaultPrefix
	}
	if len(o.Servers) == 0 {
		o.Servers = []string{rocketmqDefaultServer}
	}
	if o.ConsumeSuspendMs == 0 {
		o.ConsumeSuspendMs = rocketmqDefaultConsumeSuspendMs
	}
}

// /////////////////////////////////////////////////////////////
// Rocketmq getter
// /////////////////////////////////////////////////////////////

func (o *rocketmqConfiguration) GetConsumeMaxRetry() int32 { return o.ConsumeMaxRetry }
func (o *rocketmqConfiguration) GetConsumePosition() int   { return o.ConsumePosition }
func (o *rocketmqConfiguration) GetConsumeSuspendMs() int  { return o.ConsumeSuspendMs }
func (o *rocketmqConfiguration) GetPrefix() string         { return o.Prefix }
func (o *rocketmqConfiguration) GetServers() []string      { return o.Servers }

func (o *rocketmqConfiguration) GetKey() string    { return o.Key }
func (o *rocketmqConfiguration) GetSecret() string { return o.Secret }
func (o *rocketmqConfiguration) GetToken() string  { return o.Token }

// /////////////////////////////////////////////////////////////
// Rocketmq setter
// /////////////////////////////////////////////////////////////

func SetRocketmqConsumeMaxRetry(n int32) Option {
	return func(c *configuration) { c.AdapterRocketmq.ConsumeMaxRetry = n }
}

func SetRocketmqConsumePosition(n int) Option {
	return func(c *configuration) { c.AdapterRocketmq.ConsumePosition = n }
}

func SetRocketmqConsumeSuspendMs(n int) Option {
	return func(c *configuration) { c.AdapterRocketmq.ConsumeSuspendMs = n }
}

func SetRocketmqServers(s ...string) Option {
	return func(c *configuration) { c.AdapterRocketmq.Servers = s }
}

func SetRocketmqPrefix(s string) Option {
	return func(c *configuration) { c.AdapterRocketmq.Prefix = s }
}

func SetRocketmqKey(s string) Option {
	return func(c *configuration) { c.AdapterRocketmq.Key = s }
}

func SetRocketmqSecret(s string) Option {
	return func(c *configuration) { c.AdapterRocketmq.Secret = s }
}

func SetRocketmqToken(s string) Option {
	return func(c *configuration) { c.AdapterRocketmq.Token = s }
}
