// author: wsfuyibing <websearch@163.com>
// date: 2023-02-17

package models

type (
	// Task
	//
	// subscription task.
	Task struct {
		Id     int    `xorm:"id pk autoincr"`
		Status int    `xorm:"status"`
		Title  string `xorm:"title"`
		Remark string `xorm:"remark"`

		// Parallels
		// maximum consumers per node.
		//
		// Default: 1
		Parallels int `xorm:"parallels"`

		// Concurrency
		// maximum consuming process per consumer.
		//
		// Default: 10
		Concurrency int32 `xorm:"concurrency"`

		// MaxRetry
		// maximum consume times if delivered failed.
		//
		// If the first consumption fails, wait for 1 minute and
		// try again. If the fifth consumption fails, wait for 5
		// minutes and try again.
		//
		// Default: 3
		MaxRetry int `xorm:"max_retry"`

		// DelaySeconds
		// maximum delay seconds.
		//
		// After the message is produced, it will not be consumed
		// immediately. You need to wait for the specified time
		// before the first consumption.
		//
		// Default: 0 (not delay).
		DelaySeconds int `xorm:"delay_seconds"`

		// Broadcasting
		// switch.
		//
		// When enabled, each consumer will receive a message, that is,
		// the same message will be consumed multiple times.
		//
		// Accept: Rabbitmq, Rocketmq.
		// NotAccept: Aliyunmns.
		Broadcasting int `xorm:"broadcasting"`

		RegistryId int `xorm:"registry_id"`

		Handler             string `xorm:"handler"`
		HandlerTimeout      int    `xorm:"handler_timeout"`
		HandlerMethod       string `xorm:"handler_method"`
		HandlerCondition    string `xorm:"handler_condition"`
		HandlerResponseType int    `xorm:"handler_response_type"`
		HandlerIgnoreCodes  string `xorm:"handler_ignore_codes"`

		Failed             string `xorm:"failed"`
		FailedTimeout      int    `xorm:"failed_timeout"`
		FailedMethod       string `xorm:"failed_method"`
		FailedCondition    string `xorm:"failed_condition"`
		FailedResponseType int    `xorm:"failed_response_type"`
		FailedIgnoreCodes  string `xorm:"failed_ignore_codes"`

		Succeed             string `xorm:"succeed"`
		SucceedTimeout      int    `xorm:"succeed_timeout"`
		SucceedMethod       string `xorm:"succeed_method"`
		SucceedCondition    string `xorm:"succeed_condition"`
		SucceedResponseType int    `xorm:"succeed_response_type"`
		SucceedIgnoreCodes  string `xorm:"succeed_ignore_codes"`

		GmtCreated Timeline `xorm:"gmt_created"`
		GmtUpdated Timeline `xorm:"gmt_updated"`
	}
)

func (o *Task) IsEnabled() bool {
	return o.Status == StatusEnabled
}
