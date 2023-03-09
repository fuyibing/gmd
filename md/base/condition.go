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
	"context"
)

type (
	// ConditionConstructor
	// 条件构造器.
	ConditionConstructor func(str string) ConditionExecutor

	// ConditionExecutor
	// 条件执行器.
	ConditionExecutor interface {
		// Name
		// 执行器名称.
		Name() string

		// Validate
		// 校验条件.
		Validate(ctx context.Context, task, source *Task, message *Message) (ignored bool, err error)
	}
)
