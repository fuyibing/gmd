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
	"fmt"
	"github.com/fuyibing/db/v5"
	"github.com/fuyibing/gmd/v8/app"
	"github.com/fuyibing/gmd/v8/app/models"
	"github.com/fuyibing/gmd/v8/app/services"
	"github.com/fuyibing/log/v5"
	"time"
)

type (
	// Message
	// 消息结构.
	Message struct {
		ctx          context.Context
		dur          time.Duration
		err          error
		responseBody string

		TaskId int

		Dequeue     int
		MessageBody string
		MessageId   string
		MessageTime int64

		// 主题消息ID.
		PayloadMessageId string

		// 原始正文.
		//
		// 字段 MessageBody 为自来队列的原始消息, 本字段 DispatcherBody 为投递
		// 消息时使用的正文.
		//
		// 1. 原始消息.
		// 2. 通知消息.
		DispatcherBody string
	}
)

func (o *Message) GetContext() context.Context { return o.ctx }
func (o *Message) GetError() error             { return o.err }
func (o *Message) GetResponseBody() string     { return o.responseBody }

func (o *Message) Release() { Pool.ReleaseMessage(o) }

func (o *Message) SetContext(ctx context.Context) *Message { o.ctx = ctx; return o }
func (o *Message) SetDuration(dur time.Duration) *Message  { o.dur = dur; return o }
func (o *Message) SetError(err error) *Message             { o.err = err; return o }
func (o *Message) SetResponseBody(str string) *Message     { o.responseBody = str; return o }

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *Message) after() {
	// 保存消息.
	if o.err != nil {
		if app.Config.GetConsumer().GetSaveFailed() {
			o.save()
		}
	} else {
		if app.Config.GetConsumer().GetSaveSucceed() {
			o.save()
		}
	}

	// 清理数据.
	o.ctx = nil
	o.dur = 0
	o.err = nil
	o.responseBody = ""

	o.TaskId = 0
	o.Dequeue = 0
	o.MessageBody = ""
	o.MessageId = ""
	o.MessageTime = 0
	o.PayloadMessageId = ""
	o.DispatcherBody = ""
}

func (o *Message) before() {}

func (o *Message) init() *Message { return o }

func (o *Message) save() {
	var (
		affects, beanId int64
		bean            *models.Message
		err             error
		span            = log.NewSpanFromContext(o.ctx, "message.save")
		sess            = db.Connector.GetMasterWithContext(span.Context())
		service         = services.NewMessageService(sess)
	)

	span.Kv().Add("message.save.task.id", o.TaskId).
		Add("message.save.message.id", o.MessageId)

	// 结束保存.
	defer func() {
		// 保存异常.
		if r := recover(); r != nil {
			span.Logger().Fatal("message save fatal: %v", r)

			if err == nil {
				err = fmt.Errorf("%v", r)
			}
		}

		if err != nil {
			span.Logger().Error("message save error: bean-id=%d, affects=%d, error=%v", beanId, affects, err)
		} else {
			span.Logger().Info("message save succeed: bean-id=%d, affects=%d", beanId, affects)
		}

		span.End()
	}()

	// 读取出错.
	if bean, err = service.GetByMessageId(o.TaskId, o.MessageId); err != nil {
		return
	}

	// 首次保存.
	if bean == nil {
		req := &models.Message{
			Duration:         o.dur.Seconds(),
			TaskId:           o.TaskId,
			Dequeue:          o.Dequeue,
			PayloadMessageId: o.PayloadMessageId,
			MessageTime:      o.MessageTime,
			MessageId:        o.MessageId,
			MessageBody:      o.MessageBody,
			ResponseBody:     o.responseBody,
		}

		// 保存消息.
		if o.err != nil {
			bean, err = service.AddFailed(req)
		} else {
			bean, err = service.AddSucceed(req)
		}

		if bean != nil {
			affects = 1
			beanId = bean.Id
		}
		return
	}

	// 更新状态.
	if o.err != nil {
		affects, err = service.SetStatusAsFailed(bean.Id, o.dur, o.responseBody)
	} else {
		affects, err = service.SetStatusAsSucceed(bean.Id, o.dur, o.responseBody)
	}
	beanId = bean.Id
}
