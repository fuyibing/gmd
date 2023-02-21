// author: wsfuyibing <websearch@163.com>
// date: 2023-01-17

package services

import (
	"fmt"
	"github.com/fuyibing/db/v8"
	"github.com/fuyibing/gmd/v8/app/models"
	"xorm.io/xorm"
)

type (
	RegistryService struct {
		db.Service
	}
)

func NewRegistryService(ss ...*xorm.Session) *RegistryService {
	o := &RegistryService{}
	o.Use(ss...)
	o.UseConnection(models.ConnectionName)
	return o
}

func (o *RegistryService) AddByNames(topic, tag string) (*models.Registry, error) {
	var (
		now  = models.NewTimeline()
		bean = &models.Registry{
			TopicName:  topic,
			TopicTag:   tag,
			GmtCreated: now,
			GmtUpdated: now,
		}
		err error
	)
	if _, err = o.Master().Insert(bean); err != nil {
		return nil, err
	}
	return bean, nil
}

func (o *RegistryService) GetById(id int) (*models.Registry, error) {
	var (
		bean   = &models.Registry{}
		err    error
		exists bool
	)
	if exists, err = o.Slave().Where("id = ?", id).Get(bean); err != nil || !exists {
		return nil, err
	}
	return bean, nil
}

func (o *RegistryService) GetByNames(topic, tag string) (*models.Registry, error) {
	var (
		bean   = &models.Registry{}
		err    error
		exists bool
	)
	if exists, err = o.Slave().
		Where("topic_name = ? AND topic_tag = ?", topic, tag).
		Get(bean); err != nil || !exists {
		return nil, err
	}
	return bean, nil
}

func (o *RegistryService) ListAll() (list []*models.Registry, err error) {
	list = make([]*models.Registry, 0)
	err = o.Slave().Find(&list)
	return
}

func (o *RegistryService) SetFilterTag(id int) (int64, error) {
	return o.Master().Cols(
		"filter_tag",
	).Where(
		"id = ? AND (filter_tag = '' OR filter_tag IS NULL)",
		id,
	).Update(&models.Registry{
		FilterTag: fmt.Sprintf("T%d", id),
	})
}
