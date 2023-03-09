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

package json_errno

import (
	"context"
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
	Executor struct {
		ignoreCodes []string
		name        string
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

func (o *Executor) Name() string { return o.name }

// Validate
// 校验结果.
func (o *Executor) Validate(ctx context.Context, _, _ *base.Task, body []byte) (code int, err error) {
	var (
		errno string
		ptr   = &result{}
		span  = log.NewSpanFromContext(ctx, "result.json.errno")
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
		err = fmt.Errorf("json format: %v", err)
		return
	}

	// 错误编码.
	if errno = fmt.Sprintf("%v", ptr.Errno); errno == "" {
		err = fmt.Errorf("empty code")
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
				span.Logger().Info("result ignored: code=%s", c)
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
	o.name = "result:http"

	// 计算忽略码.
	for _, s := range strings.Split(ic, ",") {
		if s = strings.TrimSpace(s); s != "" {
			o.ignoreCodes = append(o.ignoreCodes, s)
		}
	}

	return o
}
