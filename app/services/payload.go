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
	PayloadService struct {
		db.Service
	}
)

func NewPayloadService(ss ...*xorm.Session) *PayloadService {
	o := &PayloadService{}
	o.Use(ss...)
	o.UseConnection(Connection)
	return o
}

func (o *PayloadService) AddFailed(req *models.Payload) (*models.Payload, error) {
	req.Status = models.StatusFailed
	return o.add(req)
}

func (o *PayloadService) AddSucceed(req *models.Payload) (*models.Payload, error) {
	req.Status = models.StatusSucceed
	return o.add(req)
}

func (o *PayloadService) AddWaiting(req *models.Payload) (*models.Payload, error) {
	req.Status = models.StatusWaiting
	return o.add(req)
}

func (o *PayloadService) GetByHash(hash string, offset int) (*models.Payload, error) {
	bean := &models.Payload{}
	if exists, err := o.Slave().Where("hash = ? AND offset = ?", hash, offset).Get(bean); err != nil || !exists {
		return nil, err
	}
	return bean, nil
}

func (o *PayloadService) GetById(id int64) (*models.Payload, error) {
	bean := &models.Payload{}
	if exists, err := o.Slave().Where("id = ?", id).Get(bean); err != nil || !exists {
		return nil, err
	}
	return bean, nil
}

func (o *PayloadService) SetStatusAsFailed(id int64, duration time.Duration, responseBody string) (int64, error) {
	return o.Master().Cols("status", "duration", "response_body").
		Incr("retry").
		Where("id = ?", id).
		Update(&models.Payload{Status: models.StatusFailed, Duration: duration.Seconds(), ResponseBody: responseBody})
}

func (o *PayloadService) SetStatusAsSucceed(id int64, duration time.Duration, messageId string) (int64, error) {
	return o.Master().Cols("status", "duration", "message_id", "response_body").
		Incr("retry").
		Where("id = ?", id).
		Update(&models.Payload{
			Status: models.StatusSucceed, Duration: duration.Seconds(),
			MessageId:    messageId,
			ResponseBody: "",
		})
}

func (o *PayloadService) SetStatusAsWaiting(id int64, duration time.Duration, responseBody string) (int64, error) {
	return o.Master().Cols("status", "duration", "response_body").
		Incr("retry").
		Where("id = ?", id).
		Update(&models.Payload{Status: models.StatusFailed, Duration: duration.Seconds(), ResponseBody: responseBody})
}

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *PayloadService) add(req *models.Payload) (*models.Payload, error) {
	var (
		now  = models.Now()
		bean = &models.Payload{
			Status:       req.Status,
			Duration:     req.Duration,
			Retry:        1,
			Hash:         req.Hash,
			Offset:       req.Offset,
			RegistryId:   req.RegistryId,
			MessageId:    req.MessageId,
			MessageBody:  req.MessageBody,
			ResponseBody: req.ResponseBody,
			GmtCreated:   now,
			GmtUpdated:   now,
		}
	)
	if _, err := o.Master().Insert(bean); err != nil || bean.Id == 0 {
		return nil, err
	}
	return bean, nil
}
