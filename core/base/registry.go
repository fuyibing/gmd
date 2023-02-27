// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package base

import (
	"fmt"
	"github.com/fuyibing/gmd/v8/app/models"
	"strings"
)

// Registry
// 注册组合.
type Registry struct {
	Id        int
	TopicName string
	TopicTag  string
	FilterTag string
}

func (o *Registry) init(m *models.Registry) *Registry {
	o.Id = m.Id
	o.TopicName = strings.ToUpper(m.TopicName)
	o.TopicTag = strings.ToUpper(m.TopicTag)

	if o.FilterTag = strings.ToUpper(m.FilterTag); o.FilterTag == "" {
		o.FilterTag = fmt.Sprintf("T%d", m.Id)
	}

	return o
}
