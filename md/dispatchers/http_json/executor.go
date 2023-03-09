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

package http_json

import (
	"context"
	"fmt"
	"github.com/fuyibing/gmd/v8/app"
	"github.com/fuyibing/gmd/v8/md/base"
	"github.com/fuyibing/log/v5"
	"github.com/valyala/fasthttp"
	"strings"
	"time"
)

type Executor struct {
	addr, method string
	contentType  string
	name         string
	timeout      time.Duration
}

func New(addr, method string, timeout int) base.DispatcherExecutor {
	return (&Executor{
		addr:   addr,
		method: strings.ToUpper(method),
	}).init(timeout)
}

// +---------------------------------------------------------------------------+
// + Interface methods                                                         |
// +---------------------------------------------------------------------------+

func (o *Executor) Dispatch(ctx context.Context, task, _ *base.Task, message *base.Message) (body []byte, err error) {
	var (
		span     = log.NewSpanFromContext(ctx, "dispatch.http")
		request  = fasthttp.AcquireRequest()
		response = fasthttp.AcquireResponse()
	)

	// 结束投递.
	defer func() {
		// 捕获异常.
		if r := recover(); r != nil {
			span.Logger().Fatal("dispatch fatal: %v", r)

			if err == nil {
				err = fmt.Errorf("%v", err)
			}
		}

		// 释放入池.
		fasthttp.ReleaseRequest(request)
		fasthttp.ReleaseResponse(response)

		// 结束跨度.
		if err != nil {
			span.Logger().Error("dispatch error: %v", err)
		} else {
			span.Logger().Info("dispatch completed: %s", body)
		}

		span.End()
	}()

	// 投递参数.
	//
	// - 投递地址
	// - 请求方式
	// - 请求参数
	request.Header.SetRequestURI(o.addr)
	request.Header.SetMethod(o.method)
	request.SetBodyRaw([]byte(message.DispatcherBody))

	// 请求头.
	//
	//   {
	//       "Content-Type": "application/json",
	//
	//       "X-Gmd-Tag": "CREATED",
	//       "X-Gmd-Topic": "FINANCE",
	//
	//       "X-Gmd-Dequeue": "1",
	//       "X-Gmd-Software": "gmd/1.0",
	//       "X-Gmd-Time": "1234567890123",
	//       "X-Gmd-Message": "Topic"
	//   }
	request.Header.SetContentType(o.contentType)
	request.Header.Add(base.DispatcherHeaderDequeueCount, fmt.Sprintf("%d", message.Dequeue))
	request.Header.Add(base.DispatcherHeaderMessageId, message.MessageId)
	request.Header.Add(base.DispatcherHeaderMessageTime, fmt.Sprintf("%d", message.MessageTime))
	request.Header.Add(base.DispatcherHeaderSoftware, app.Config.GetSoftware())
	request.Header.Add(base.DispatcherHeaderTopicTag, task.TopicTag)
	request.Header.Add(base.DispatcherHeaderTopicName, task.TopicName)

	// 请求过程.
	span.Logger().Info("dispatch on: %s %s %v", request.Header.Protocol(), o.method, o.addr)
	err = fasthttp.DoTimeout(request, response, o.timeout)
	body = response.Body()
	return
}

func (o *Executor) Name() string { return o.name }

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *Executor) send() {}

func (o *Executor) init(timeout int) *Executor {
	o.contentType = "application/json"
	o.name = "dispatcher:http"

	if timeout > 0 {
		o.timeout = time.Duration(timeout) * time.Second
	} else {
		o.timeout = time.Duration(10) * time.Second
	}

	return o
}