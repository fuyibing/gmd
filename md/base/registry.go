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
// date: 2023-03-07

package base

import (
	"fmt"
	"github.com/fuyibing/gmd/v8/app/models"
	"strings"
)

type (
	// Registry
	// 注册关系.
	//
	// 此配置用于控制生产/发布消息, 只有注册过的主题/标签关系才允许发布到MQ服务器,
	// 反之禁止发布.
	//
	// 不同的MQ服务器, 有大小写敏感差异, 注册关系统一以大写格式输出.
	Registry struct {
		Id int

		// 主题名称.
		// 例如: FINANCE
		TopicName string

		// 主题标签.
		// 例如: CREATED, PAID
		TopicTag string

		// 过滤标签.
		//
		// 在 AliyunMNS 服务中, TopicTag 有最长 16 个字符限制, 执行订阅时使用此
		// 字段替代.
		FilterTag string
	}
)

func (o *Registry) init(m *models.Registry) *Registry {
	o.Id = m.Id
	o.TopicName = strings.ToUpper(m.TopicName)
	o.TopicTag = strings.ToUpper(m.TopicTag)

	// 标签重置.
	if o.FilterTag = strings.ToUpper(m.FilterTag); o.FilterTag == "" {
		o.FilterTag = fmt.Sprintf("T%d", m.Id)
	}

	return o
}
