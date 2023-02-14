// author: wsfuyibing <websearch@163.com>
// date: 2023-01-16

package controllers

import (
	"github.com/fuyibing/gmd/app/controllers/task"
	"github.com/fuyibing/gmd/app/controllers/topic"
	"sync"
)

var ControllerRegistration map[string]interface{}

func init() {
	new(sync.Once).Do(func() {
		ControllerRegistration = map[string]interface{}{
			"/":      &Controller{},
			"task":   &task.Controller{},
			"/topic": &topic.Controller{},
		}
	})
}
