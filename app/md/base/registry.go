// author: wsfuyibing <websearch@163.com>
// date: 2023-01-17

package base

import (
	"fmt"
	"github.com/fuyibing/gmd/app/models"
	"strings"
)

// Registry
// memory registry relation.
type Registry struct {
	Id                             int
	FilterTag, TopicTag, TopicName string
}

// /////////////////////////////////////////////////////////////
// Access methods.
// /////////////////////////////////////////////////////////////

func (o *Registry) init(m *models.Registry) *Registry {
	// Basic fields.
	o.Id = m.Id
	o.TopicName = strings.ToUpper(m.TopicName)
	o.TopicTag = strings.ToUpper(m.TopicTag)

	// Execution field.
	if o.FilterTag = strings.ToUpper(m.FilterTag); o.FilterTag == "" {
		o.FilterTag = fmt.Sprintf("T%d", m.Id)
	}

	return o
}
