// author: wsfuyibing <websearch@163.com>
// date: 2023-02-22

package rocketmq

import (
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/fuyibing/gmd/v8/app"
	"sync"
)

var (
	Agent AgentManager
)

type (
	AgentManager interface {
		GenDelayTag(tag string, id int) string
		GenTopicName(name string) string
	}

	agent struct {
	}
)

func (o *agent) GenDelayTag(tag string, id int) string {
	return fmt.Sprintf("%s-%d", tag, id)
}

func (o *agent) GenTopicName(name string) string {
	return fmt.Sprintf("%s%s", app.Config.GetRocketmq().GetPrefix(), name)
}

func init() {
	new(sync.Once).Do(func() {
		rlog.SetLogLevel("error")

		Agent = &agent{}
	})
}
