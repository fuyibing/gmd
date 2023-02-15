// author: wsfuyibing <websearch@163.com>
// date: 2023-02-14

package rocketmq

import (
	"fmt"
	"github.com/fuyibing/gmd/app/md/conf"
)

var (
	Agent AgentManager
)

type (
	AgentManager interface {
		// GenGroupName
		// return subscription group name.
		//
		//   return "X-GID-1"
		GenGroupName(id int) string

		// GenTopicName
		// return topic name.
		//
		//   return "X-TOPIC"
		GenTopicName(name string) string
	}

	agent struct{}
)

func (o *agent) GenDelayTag(tag string) string {
	return ""
}

func (o *agent) GenGroupName(id int) string {
	return fmt.Sprintf("%s%s-%d",
		conf.Config.Account.Rocketmq.Prefix,
		DefaultConsumerGroupName,
		id,
	)
}

func (o *agent) GenTopicName(name string) string {
	return fmt.Sprintf("%s%s",
		conf.Config.Account.Rocketmq.Prefix,
		name,
	)
}

// func (o *agent) GenRetryTopicName(id int) string {
// 	return fmt.Sprintf("%%RETRY%%%s%s-%d",
// 		conf.Config.Account.Rocketmq.Prefix,
// 		GROUP_NAME,
// 		id,
// 	)
// }

func (o *agent) init() *agent {
	return o
}
