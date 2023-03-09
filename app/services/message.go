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

package services

import (
	"github.com/fuyibing/db/v5"
	"github.com/fuyibing/gmd/v8/app/models"
	"time"
	"xorm.io/xorm"
)

type (
	MessageService struct {
		db.Service
	}
)

func NewMessageService(ss ...*xorm.Session) *MessageService {
	o := &MessageService{}
	o.Use(ss...)
	o.UseConnection(Connection)
	return o
}

// AddFailed
// 添加失败消费.
func (o *MessageService) AddFailed(req *models.Message) (*models.Message, error) {
	req.Status = models.StatusFailed
	return o.add(req)
}

// AddSucceed
// 添加成功消费.
func (o *MessageService) AddSucceed(req *models.Message) (*models.Message, error) {
	req.Status = models.StatusSucceed
	return o.add(req)
}

// GetById
// 读取一条消息.
func (o *MessageService) GetById(id int64) (*models.Message, error) {
	bean := &models.Message{}
	if exists, err := o.Slave().Where("id = ?", id).Get(bean); err != nil || !exists {
		return nil, err
	}
	return bean, nil
}

// GetByMessageId
// 读取一条消息.
func (o *MessageService) GetByMessageId(taskId int, messageId string) (*models.Message, error) {
	bean := &models.Message{}
	if exists, err := o.Slave().Where("task_id = ? AND message_id = ?", taskId, messageId).Get(bean); err != nil || !exists {
		return nil, err
	}
	return bean, nil
}

// SetStatusAsFailed
// 设为消费失败.
func (o *MessageService) SetStatusAsFailed(id int64, duration time.Duration, responseBody string) (int64, error) {
	return o.Master().Cols("status", "duration", "response_body").
		Incr("retry").
		Where("id = ?", id).
		Update(&models.Message{Status: models.StatusFailed, Duration: duration.Seconds(), ResponseBody: responseBody})
}

// SetStatusAsSucceed
// 设为消费成功.
func (o *MessageService) SetStatusAsSucceed(id int64, duration time.Duration, responseBody string) (int64, error) {
	return o.Master().Cols("status", "duration", "response_body").
		Incr("retry").
		Where("id = ?", id).
		Update(&models.Message{Status: models.StatusSucceed, Duration: duration.Seconds(), ResponseBody: responseBody})
}

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *MessageService) add(req *models.Message) (*models.Message, error) {
	var (
		now  = models.Now()
		bean = &models.Message{
			Status:           req.Status,
			Duration:         req.Duration,
			TaskId:           req.TaskId,
			Dequeue:          req.Dequeue,
			Retry:            1,
			PayloadMessageId: req.PayloadMessageId,
			MessageTime:      req.MessageTime,
			MessageId:        req.MessageId,
			MessageBody:      req.MessageBody,
			ResponseBody:     req.ResponseBody,
			GmtCreated:       now,
			GmtUpdated:       now,
		}
	)
	if _, err := o.Master().Insert(bean); err != nil || bean.Id == 0 {
		return nil, err
	}
	return bean, nil
}
