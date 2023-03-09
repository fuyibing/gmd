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
	"github.com/google/uuid"
	"strings"
	"time"
)

type (
	// Payload
	// 消息结构.
	Payload struct {
		ctx       context.Context
		dur       time.Duration
		err       error
		messageId string

		Hash   string
		Offset int

		FilterTag   string
		MessageBody string
		RegistryId  int
		TopicName   string
		TopicTag    string
	}
)

func (o *Payload) GenHash()                                { o.Hash = strings.ToUpper(strings.ReplaceAll(uuid.NewString(), "-", "")) }
func (o *Payload) GetContext() context.Context             { return o.ctx }
func (o *Payload) Release()                                { Pool.ReleasePayload(o) }
func (o *Payload) SetContext(ctx context.Context) *Payload { o.ctx = ctx; return o }
func (o *Payload) SetDuration(dur time.Duration) *Payload  { o.dur = dur; return o }
func (o *Payload) SetError(err error) *Payload             { o.err = err; return o }
func (o *Payload) SetMessageId(str string) *Payload        { o.messageId = str; return o }

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *Payload) after() {
	// 保存消息.
	if o.err != nil {
		if app.Config.GetProducer().GetSaveFailed() {
			o.save()
		}
	} else {
		if app.Config.GetProducer().GetSaveSucceed() {
			o.save()
		}
	}

	// 清理数据.
	o.ctx = nil
	o.dur = 0
	o.err = nil
	o.messageId = ""

	o.Hash = ""
	o.Offset = 0
	o.TopicName = ""
	o.TopicTag = ""
	o.FilterTag = ""
	o.MessageBody = ""
}

func (o *Payload) before() {}

func (o *Payload) init() *Payload { return o }

func (o *Payload) save() {
	var (
		affects, beanId int64
		bean            *models.Payload
		err             error
		span            = log.NewSpanFromContext(o.ctx, "payload.save")
		sess            = db.Connector.GetMasterWithContext(span.Context())
		service         = services.NewPayloadService(sess)
	)

	span.Kv().Add("payload.save.hash", o.Hash).
		Add("payload.save.offset", o.Offset).
		Add("payload.save.registry.id", o.RegistryId)

	// 结束保存.
	defer func() {
		// 保存异常.
		if r := recover(); r != nil {
			span.Logger().Fatal("payload save fatal: %v", r)

			if err == nil {
				err = fmt.Errorf("%v", r)
			}
		}

		// 记录结果.
		if err != nil {
			span.Logger().Error("payload save error: bean-id=%d, affects=%d, error=%v", beanId, affects, err)
		} else {
			span.Logger().Info("payload save succeed: bean-id=%d, affects=%d", beanId, affects)
		}
		span.End()
	}()

	// 读取出错.
	if bean, err = service.GetByHash(o.Hash, o.Offset); err != nil {
		return
	}

	// 首次保存.
	if bean == nil {
		req := &models.Payload{
			Duration:    o.dur.Seconds(),
			Hash:        o.Hash,
			Offset:      o.Offset,
			RegistryId:  o.RegistryId,
			MessageId:   o.messageId,
			MessageBody: o.MessageBody,
		}

		// 保存消息.
		if o.err != nil {
			req.ResponseBody = o.err.Error()
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
		if (bean.Retry + 1) >= app.Config.GetProducer().GetMaxRetry() {
			affects, err = service.SetStatusAsWaiting(bean.Id, o.dur, o.err.Error())
		} else {
			affects, err = service.SetStatusAsFailed(bean.Id, o.dur, o.err.Error())
		}
	} else {
		affects, err = service.SetStatusAsSucceed(bean.Id, o.dur, o.messageId)
	}
	beanId = bean.Id
}
