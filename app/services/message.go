// author: wsfuyibing <websearch@163.com>
// date: 2023-01-17

package services

import (
	"github.com/fuyibing/db/v8"
	"github.com/fuyibing/gmd/app/models"
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
	o.UseConnection(models.ConnectionName)
	return o
}

func (o *MessageService) AddFailed(r *models.Message) (*models.Message, error) {
	r.Status = models.StatusFailed
	return o.add(r)
}

func (o *MessageService) AddIgnored(r *models.Message) (*models.Message, error) {
	r.Status = models.StatusIgnored
	return o.add(r)
}

func (o *MessageService) AddSucceed(r *models.Message) (*models.Message, error) {
	r.Status = models.StatusSucceed
	return o.add(r)
}

func (o *MessageService) GetById(id int64) (*models.Message, error) {
	var (
		bean   = &models.Message{}
		err    error
		exists bool
	)
	if exists, err = o.Slave().
		Where("id = ?", id).
		Get(bean); err != nil || !exists {
		return nil, err
	}
	return bean, nil
}

func (o *MessageService) GetByMessageId(taskId int, messageId string) (*models.Message, error) {
	var (
		bean   = &models.Message{}
		err    error
		exists bool
	)
	if exists, err = o.Slave().
		Where("task_id = ? AND message_id = ?", taskId, messageId).
		Get(bean); err != nil || !exists {
		return nil, err
	}
	return bean, nil
}

func (o *MessageService) ListWaiting(limit int) (list []*models.Message, err error) {
	list = make([]*models.Message, 0)
	err = o.Slave().Where("status = ?", models.StatusWaiting).Limit(limit).Find(&list)
	return
}

func (o *MessageService) SetStatusAsFailed(id int64, duration float64, responseBody string) (int64, error) {
	return o.Master().Cols(
		"status",
		"duration",
		"response_body",
	).Incr("retry", 1).Where("id = ?", id).Update(&models.Message{
		Status:       models.StatusFailed,
		Duration:     duration,
		ResponseBody: responseBody,
	})
}

func (o *MessageService) SetStatusAsIgnored(id int64) (int64, error) {
	return o.Master().Cols(
		"status",
	).Incr("retry", 1).Where("id = ?", id).Update(&models.Message{
		Status: models.StatusIgnored,
	})
}

func (o *MessageService) SetStatusAsProcessing(id int64) (int64, error) {
	return o.Master().Cols(
		"status",
	).Where(
		"id = ? AND status = ?",
		id,
		models.StatusWaiting,
	).Update(&models.Message{
		Status: models.StatusProcessing,
	})
}

func (o *MessageService) SetStatusAsSucceed(id int64, duration float64, responseBody string) (int64, error) {
	return o.Master().Cols(
		"status",
		"duration",
		"response_body",
	).Incr("retry", 1).Where("id = ?", id).Update(&models.Message{
		Status:       models.StatusSucceed,
		Duration:     duration,
		ResponseBody: responseBody,
	})
}

// /////////////////////////////////////////////////////////////
// Access methods.
// /////////////////////////////////////////////////////////////

func (o *MessageService) add(req *models.Message) (*models.Message, error) {
	var (
		now  = models.NewTimeline()
		bean = &models.Message{
			Status:           req.Status,
			Duration:         req.Duration,
			Retry:            1,
			PayloadMessageId: req.PayloadMessageId,
			TaskId:           req.TaskId,
			MessageDequeue:   req.MessageDequeue,
			MessageTime:      req.MessageTime,
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
