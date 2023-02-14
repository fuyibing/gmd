// author: wsfuyibing <websearch@163.com>
// date: 2023-01-17

package services

import (
	"github.com/fuyibing/db/v3"
	"github.com/fuyibing/gmd/app/models"
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
	o.UseConnection(models.ConnectionName)
	return o
}

func (o *PayloadService) AddFailed(req *models.Payload) (*models.Payload, error) {
	req.Status = models.StatusFailed
	return o.add(req)
}

func (o *PayloadService) AddIgnored(req *models.Payload) (*models.Payload, error) {
	req.Status = models.StatusIgnored
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
	var (
		bean   = &models.Payload{}
		err    error
		exists bool
	)
	if exists, err = o.Slave().
		Where("hash = ? AND offset = ?", hash, offset).
		Get(bean); err != nil || !exists {
		return nil, err
	}
	return bean, nil
}

func (o *PayloadService) GetById(id int64) (*models.Payload, error) {
	var (
		bean   = &models.Payload{}
		err    error
		exists bool
	)
	if exists, err = o.Slave().Where("id = ?", id).Get(bean); err != nil || !exists {
		return nil, err
	}
	return bean, nil
}

func (o *PayloadService) ListWaiting(limit int) (list []*models.Payload, err error) {
	list = make([]*models.Payload, 0)
	err = o.Slave().Where("status = ?", models.StatusWaiting).Limit(limit).Find(&list)
	return
}

func (o *PayloadService) SetStatusAsFailed(id int64, duration float64, responseBody string) (int64, error) {
	return o.Master().Cols(
		"status",
		"duration",
		"response_body",
	).Incr("retry", 1).Where("id = ?", id).Update(&models.Payload{
		Status:       models.StatusFailed,
		Duration:     duration,
		ResponseBody: responseBody,
	})
}

func (o *PayloadService) SetStatusAsIgnored(id int64) (int64, error) {
	return o.Master().Cols(
		"status",
	).Incr("retry", 1).Where("id = ?", id).Update(&models.Payload{
		Status: models.StatusIgnored,
	})
}

func (o *PayloadService) SetStatusAsProcessing(id int64) (int64, error) {
	return o.Master().Cols(
		"status",
	).Where(
		"id = ? AND status = ?",
		id,
		models.StatusWaiting,
	).Update(&models.Payload{
		Status: models.StatusProcessing,
	})
}

func (o *PayloadService) SetStatusAsPublished(id int64) (int64, error) {
	return o.Master().Cols(
		"status",
	).Where("id = ?", id).Update(&models.Payload{
		Status: models.StatusSucceed,
	})
}

func (o *PayloadService) SetStatusAsSucceed(id int64, duration float64, messageId string) (int64, error) {
	return o.Master().Cols(
		"status",
		"duration",
		"message_id",
		"response_body",
	).Incr("retry", 1).Where("id = ?", id).Update(&models.Payload{
		Status:       models.StatusSucceed,
		Duration:     duration,
		MessageId:    messageId,
		ResponseBody: "",
	})
}

func (o *PayloadService) SetStatusAsWaiting(id int64, duration float64, responseBody string) (int64, error) {
	return o.Master().Cols(
		"status",
		"duration",
		"response_body",
	).Incr("retry", 1).Where("id = ?", id).Update(&models.Payload{
		Status:       models.StatusWaiting,
		Duration:     duration,
		ResponseBody: responseBody,
	})
}

// /////////////////////////////////////////////////////////////
// Access methods.
// /////////////////////////////////////////////////////////////

func (o *PayloadService) add(req *models.Payload) (*models.Payload, error) {
	var (
		now  = models.NewTimeline()
		bean = &models.Payload{
			Status:           req.Status,
			Duration:         req.Duration,
			Retry:            1,
			MessageTaskId:    req.MessageTaskId,
			MessageMessageId: req.MessageMessageId,
			Hash:             req.Hash,
			Offset:           req.Offset,
			RegistryId:       req.RegistryId,
			MessageId:        req.MessageId,
			MessageBody:      req.MessageBody,
			ResponseBody:     req.ResponseBody,
			GmtCreated:       now,
			GmtUpdated:       now,
		}
	)
	_, err := o.Master().Insert(bean)
	return bean, err
}
