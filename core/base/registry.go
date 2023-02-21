// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package base

import (
	"fmt"
	"github.com/fuyibing/gmd/v8/app/models"
	"strings"
)

// Registry
// config topic name and tag relation.
type Registry struct {
	Id        int
	TopicName string
	TopicTag  string
	FilterTag string
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
