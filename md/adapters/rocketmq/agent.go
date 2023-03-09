// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// author: wsfuyibing <websearch@163.com>
// date: 2023-03-09

package rocketmq

import (
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/fuyibing/gmd/v8/app"
	"sync"
)

const (
	PropertyBornTime         = "_GMD_BORN_TIME_"
	PropertyPayloadMessageId = "_GMD_PAYLOAD_MESSAGE_ID_"
	PropertyPublishCount     = "_GMD_PUBLISH_COUNT_"
)

var (
	// Agent
	// 代理实例.
	Agent AgentManager
)

type (
	// AgentManager
	// 代理管理器.
	AgentManager interface {
		// GenerateDelayTag
		// 延时订阅标签.
		//
		// 在 Rocketmq 中, 延时订阅除了自有标签外, 额外增加一组随机标签, 当收到
		// 的消息未到消费时间时, 重要发一条对应的额外标签消息, 直到应消费时间达到
		// 阈值.
		//
		// 入参: 101
		// 出参: X-DELAY-101
		GenerateDelayTag(id int) string

		// GenerateGroupId
		// 消费者分组.
		//
		// 出参: 101
		// 出参: X-GROUP-101
		GenerateGroupId(id int) string

		// GenerateTopicName
		// 主题名称.
		//
		// 入参: Topic
		// 出参: X-Topic
		GenerateTopicName(name string) string
	}

	agent struct{}
)

func (o *agent) GenerateDelayTag(id int) string {
	return fmt.Sprintf("%sDELAY-%d",
		app.Config.GetRocketmq().GetPrefix(),
		id,
	)
}

func (o *agent) GenerateGroupId(id int) string {
	return fmt.Sprintf("%sGROUP-%d",
		app.Config.GetRocketmq().GetPrefix(),
		id,
	)
}

func (o *agent) GenerateTopicName(name string) string {
	return fmt.Sprintf("%s%s",
		app.Config.GetRocketmq().GetPrefix(),
		name,
	)
}

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *agent) init() *agent {
	rlog.SetLogLevel("error")
	return o
}

func init() { new(sync.Once).Do(func() { Agent = (&agent{}).init() }) }
