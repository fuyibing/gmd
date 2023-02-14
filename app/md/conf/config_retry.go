// author: wsfuyibing <websearch@163.com>
// date: 2023-02-02

package conf

type (
	// RetryConfig
	// configurations for retry manager.
	RetryConfig struct {
		MessageCount   int
		MessageSeconds int

		PayloadCount   int
		PayloadSeconds int
	}
)

func (o *RetryConfig) init() *RetryConfig {
	return o
}

func (o *RetryConfig) initDefaults() {
	if o.MessageCount == 0 {
		o.MessageCount = 10
	}
	if o.MessageSeconds == 0 {
		o.MessageSeconds = 60
	}
	if o.PayloadCount == 0 {
		o.PayloadCount = 10
	}
	if o.PayloadSeconds == 0 {
		o.PayloadSeconds = 60
	}
}
