// author: wsfuyibing <websearch@163.com>
// date: 2023-02-02

package conf

import (
	"time"
)

const (
	EventSleepDuration = time.Millisecond * 10
)

type Adapter string

const (
	Aliyunmns Adapter = "aliyunmns"
	Rabbitmq  Adapter = "rabbitmq"
	Rocketmq  Adapter = "rocketmq"
)
