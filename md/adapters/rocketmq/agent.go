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

var (
	Agent AgentManager
)

type (
	AgentManager interface {
		GenerateTopicName(name string) string
	}

	agent struct {
	}
)

func (o *agent) GenerateTopicName(name string) string {
	return fmt.Sprintf("%s%s", app.Config.GetRocketmq().GetPrefix(), name)
}

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *agent) init() *agent {
	rlog.SetLogLevel("error")
	return o
}

func init() { new(sync.Once).Do(func() { Agent = (&agent{}).init() }) }
