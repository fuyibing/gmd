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

package el_condition

import (
	"github.com/fuyibing/gmd/v8/md/base"
	"github.com/fuyibing/log/v5"
)

type Executor struct {
	filter string
	name   string
}

func New(filter string) base.ConditionExecutor {
	return (&Executor{
		filter: filter,
	}).init()
}

// +---------------------------------------------------------------------------+
// + Interface methods                                                         |
// +---------------------------------------------------------------------------+

func (o *Executor) Name() string { return o.name }

func (o *Executor) Validate(_ *base.Task, message *base.Message) (ignored bool, err error) {
	var (
		span = log.NewSpanFromContext(message.GetContext(), "message.condition.validate")
	)

	span.Kv().Add("message.condition.filter", o.filter)

	defer span.End()

	span.Logger().Info("condition parse error")
	return
}

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *Executor) init() *Executor {
	o.name = "condition:el"
	return o
}
