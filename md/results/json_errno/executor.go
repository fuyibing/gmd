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

// Package json_errno
// 校验JSON数据.
//
// 消息投递结果返回 JSON 字符串, 且 errno 字段的值为 0 表示消费成功, 或者其值在忽
// 略列表中.
//
// 消费成功:
//   {
//       "data": {
//           "key": "value"
//       },
//       "data_type": "OBJECT",
//       "errno": 0,
//       "error": ""
//   }
//
// 消费失败:
//   {
//       "data": {
//           "key": "value"
//       },
//       "data_type": "OBJECT",
//       "errno": 1001,
//       "error": "错误原因"
//   }
package json_errno

import (
	"encoding/json"
	"fmt"
	"github.com/fuyibing/gmd/v8/md/base"
	"github.com/fuyibing/log/v5"
	"strconv"
	"strings"
)

const (
	ErrnoZero = "0"
)

type (
	// Executor
	// 执行器.
	Executor struct {
		// 忽略列表.
		// 当 JSON 字段的 errno 值为 0 或在此范围内, 表示消费成功.
		ignoreCodes []string

		// 执行器名称.
		name string
	}

	result struct {
		Errno interface{} `json:"errno"`
		Error interface{} `json:"error"`
	}
)

func New(ic string) base.ResultExecutor {
	return (&Executor{
		ignoreCodes: make([]string, 0),
	}).init(ic)
}

// +---------------------------------------------------------------------------+
// + Interface methods                                                         |
// +---------------------------------------------------------------------------+

func (o *Executor) Name() string {
	return o.name
}

// Validate
// 校验结果.
func (o *Executor) Validate(_, _ *base.Task, message *base.Message, body []byte) (code int, err error) {
	var (
		errno string
		ptr   = &result{}
		span  = log.NewSpanFromContext(message.GetContext(), "message.result.json.errno")
	)

	// 完成校验.
	defer func() {
		// 捕获异常.
		if r := recover(); r != nil {
			span.Logger().Fatal("result fatal: %v", r)

			if err == nil {
				err = fmt.Errorf("%v", r)
			}
		}

		// 转错误码.
		if errno != "" {
			if n, ne := strconv.ParseInt(errno, 10, 32); ne == nil {
				code = int(n)
			}
		}

		// 记录结果.
		if err != nil {
			span.Logger().Error("result error: code=%v, error=%v", code, err)
		} else {
			span.Logger().Info("result succeed: code=%v", code)
		}

		span.End()
	}()

	// 转成JSON.
	if err = json.Unmarshal(body, ptr); err != nil {
		err = fmt.Errorf("illegal result format: %v", err)
		return
	}

	// 错误编码.
	if errno = fmt.Sprintf("%v", ptr.Errno); errno == "" {
		err = fmt.Errorf("result code can not be blank")
		return
	}

	// 默认编码.
	if errno == ErrnoZero {
		return
	}

	// 忽略编码.
	if errno != ErrnoZero {
		for _, c := range o.ignoreCodes {
			if c == errno {
				span.Logger().Info("result code ignored: code=%s", c)
				return
			}
		}
	}

	// 校验出错.
	err = fmt.Errorf("%v", ptr.Error)
	return
}

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *Executor) init(ic string) *Executor {
	o.name = "result.json.errno"

	// 计算忽略码.
	for _, s := range strings.Split(ic, ",") {
		if s = strings.TrimSpace(s); s != "" {
			o.ignoreCodes = append(o.ignoreCodes, s)
		}
	}

	return o
}
