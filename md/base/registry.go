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
	Registry struct {
		Id        int
		TopicName string
		TopicTag  string
		FilterTag string
	}
)

func (o *Registry) init(m *models.Registry) *Registry {
	o.Id = m.Id
	o.TopicName = strings.ToUpper(m.TopicName)
	o.TopicTag = strings.ToUpper(m.TopicTag)

	if o.FilterTag = strings.ToUpper(m.FilterTag); o.FilterTag == "" {
		o.FilterTag = fmt.Sprintf("T%d", m.Id)
	}
	return o
}
