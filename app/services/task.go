// author: wsfuyibing <websearch@163.com>
// date: 2023-01-17

package services

import (
	"github.com/fuyibing/db/v5"
	"github.com/fuyibing/gmd/v8/app/models"
	"xorm.io/xorm"
)

type (
	TaskService struct {
		db.Service
	}
)

func NewTaskService(ss ...*xorm.Session) *TaskService {
	o := &TaskService{}
	o.Use(ss...)
	o.UseConnection(models.ConnectionName)
	return o
}

func (o *TaskService) Add(req *models.Task) (*models.Task, error) {
	var (
		now  = models.NewTimeline()
		bean = &models.Task{
			Status:       models.StatusDisabled,
			Title:        req.Title,
			Remark:       req.Remark,
			DelaySeconds: req.DelaySeconds,
			RegistryId:   req.RegistryId,
			Handler:      req.Handler,
			GmtCreated:   now,
			GmtUpdated:   now,
		}
		err error
	)
	if _, err = o.Master().Insert(bean); err != nil {
		return nil, err
	}
	return bean, nil
}

func (o *TaskService) GetByHandler(id int, handler string) (*models.Task, error) {
	var (
		bean   = &models.Task{}
		err    error
		exists bool
	)
	if exists, err = o.Slave().
		Where("registry_id = ? AND handler = ?", id, handler).
		Get(bean); err != nil || !exists {
		return nil, err
	}
	return bean, nil
}

func (o *TaskService) GetById(id int) (*models.Task, error) {
	var (
		bean   = &models.Task{}
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

func (o *TaskService) ListEnables() (list []*models.Task, err error) {
	list = make([]*models.Task, 0)
	err = o.Slave().Where("status = ?", models.StatusEnabled).Find(&list)
	return
}

func (o *TaskService) SetBasicFields(req *models.Task) (int64, error) {
	return o.Master().Cols(
		"title",
		"remark",
		"parallels",
		"concurrency",
		"max_retry",
		"delay_seconds",
		"broadcasting",
	).Where("id = ?", req.Id).Update(&models.Task{
		Title:        req.Title,
		Remark:       req.Remark,
		Parallels:    req.Parallels,
		Concurrency:  req.Concurrency,
		MaxRetry:     req.MaxRetry,
		DelaySeconds: req.DelaySeconds,
		Broadcasting: req.Broadcasting,
	})
}

func (o *TaskService) SetStatusAsDisabled(id int) (int64, error) {
	return o.Master().Cols("status").Where("id = ?", id).Update(&models.Task{
		Status: models.StatusDisabled,
	})
}

func (o *TaskService) SetStatusAsEnabled(id int) (int64, error) {
	return o.Master().Cols("status").Where("id = ?", id).Update(&models.Task{
		Status: models.StatusEnabled,
	})
}

func (o *TaskService) SetSubscriberForFailed(req *models.Task) (int64, error) {
	return o.Master().Cols(
		"failed",
		"failed_timeout",
		"failed_method",
		"failed_condition",
		"failed_response_type",
		"failed_ignore_codes",
	).Where("id = ?", req.Id).Update(&models.Task{
		Failed:             req.Failed,
		FailedCondition:    req.FailedCondition,
		FailedTimeout:      req.FailedTimeout,
		FailedMethod:       req.FailedMethod,
		FailedResponseType: req.FailedResponseType,
		FailedIgnoreCodes:  req.FailedIgnoreCodes,
	})
}

func (o *TaskService) SetSubscriberForHandler(req *models.Task) (int64, error) {
	return o.Master().Cols(
		"handler",
		"handler_timeout",
		"handler_method",
		"handler_condition",
		"handler_response_type",
		"handler_ignore_codes",
	).Where("id = ?", req.Id).Update(&models.Task{
		Handler:             req.Handler,
		HandlerCondition:    req.HandlerCondition,
		HandlerTimeout:      req.HandlerTimeout,
		HandlerMethod:       req.HandlerMethod,
		HandlerResponseType: req.HandlerResponseType,
		HandlerIgnoreCodes:  req.HandlerIgnoreCodes,
	})
}

func (o *TaskService) SetSubscriberForSucceed(req *models.Task) (int64, error) {
	return o.Master().Cols(
		"succeed",
		"succeed_timeout",
		"succeed_method",
		"succeed_condition",
		"succeed_response_type",
		"succeed_ignore_codes",
	).Where("id = ?", req.Id).Update(&models.Task{
		Succeed:             req.Succeed,
		SucceedCondition:    req.SucceedCondition,
		SucceedTimeout:      req.SucceedTimeout,
		SucceedMethod:       req.SucceedMethod,
		SucceedResponseType: req.SucceedResponseType,
		SucceedIgnoreCodes:  req.SucceedIgnoreCodes,
	})
}
