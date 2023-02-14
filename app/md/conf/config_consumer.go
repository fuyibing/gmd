// author: wsfuyibing <websearch@163.com>
// date: 2023-02-02

package conf

type (
	// ConsumerConfig
	// configurations for consumer manager.
	ConsumerConfig struct {
		Concurrency int32 `yaml:"concurrency" json:"concurrency"`

		DispatchTimeout int `yaml:"dispatch-timeout" json:"dispatch-timeout"`

		// StoreDispatchFailed
		// 投递失败消息是否存储.
		StoreDispatchFailed *bool `yaml:"store-dispatch-failed" json:"store-dispatch-failed"`

		// StoreDispatchIgnored
		// 被忽略消息是否存储.
		StoreDispatchIgnored *bool `yaml:"store-dispatch-ignored" json:"store-dispatch-ignored"`

		// StoreDispatchSucceed
		// 投递成功消息是否存储.
		StoreDispatchSucceed *bool `yaml:"store-dispatch-succeed" json:"store-dispatch-succeed"`

		// MaxRetry
		// 最大重试次数.
		//
		// 当收到消息并向订阅方投递时, 如果投递失败允许最大重试次数.
		//
		// 默认: 3
		MaxRetry int `yaml:"retry" json:"retry"`

		// Parallels
		// 最大并行数.
		//
		// 控制订阅任务在单个节点(单个服务器)上, 同时开启几个消费者. 当
		// 数据表 schema.task 表中的 parallels 字段值为0时, 此控制生
		// 效.
		//
		// 默认: 1
		Parallels int `yaml:"parallels" json:"parallels"`

		// PollingWaitSeconds
		// 轮询阻塞时长.
		//
		// 当以轮询模式收取消息时, 本次轮询最大超时时长.
		//
		// 默认: 30
		PollingWaitSeconds int64 `yaml:"polling-wait-seconds" json:"polling-wait-seconds"`

		// ReloadSeconds
		// 重载时长.
		//
		// 定时从DB中同步订阅任务, 并启动/停止消费者.
		//
		// 默认: 180
		ReloadSeconds int `yaml:"reload-seconds" json:"reload-seconds"`
	}
)

func (o *ConsumerConfig) init() *ConsumerConfig {
	return o
}

func (o *ConsumerConfig) initDefaults() {
	var (
		bt = true
	)

	if o.Concurrency == 0 {
		o.Concurrency = 10
	}

	if o.DispatchTimeout == 0 {
		o.DispatchTimeout = 10
	}

	if o.StoreDispatchFailed == nil {
		o.StoreDispatchFailed = &bt
	}
	if o.StoreDispatchIgnored == nil {
		o.StoreDispatchIgnored = &bt
	}
	if o.StoreDispatchSucceed == nil {
		o.StoreDispatchSucceed = &bt
	}

	if o.MaxRetry == 0 {
		o.MaxRetry = 3
	}

	if o.Parallels < 1 {
		o.Parallels = 1
	}

	if o.PollingWaitSeconds == 0 {
		o.PollingWaitSeconds = 30
	}

	if o.ReloadSeconds == 0 {
		o.ReloadSeconds = 180
	}
}
