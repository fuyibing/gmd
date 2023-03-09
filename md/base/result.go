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
	// ResultConstructor
	// 结果构造器.
	ResultConstructor func(ignoreCodes string) ResultExecutor

	// ResultExecutor
	// 结果执行器.
	ResultExecutor interface {
		// Name
		// 执行器名称.
		Name() string

		// Validate
		// 校验结果.
		Validate(ctx context.Context, task, source *Task, body []byte) (code int, err error)
	}
)