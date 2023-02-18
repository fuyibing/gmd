// author: wsfuyibing <websearch@163.com>
// date: 2023-01-17

// Package controllers
// MVC Controller with route prefix /.
package controllers

import (
	"encoding/json"
	"github.com/fuyibing/gmd/app/logics"
	"github.com/fuyibing/gmd/app/logics/index"
	"github.com/fuyibing/util/v8/web/response"
	"github.com/kataras/iris/v12"
)

// Controller
// Default.
type Controller struct{}

// Get
// Home.
//
// @Ignore()
// @Response(app/logics/index.HomeResponse)
func (o *Controller) Get(i iris.Context) interface{} {
	return logics.New(i, index.NewHome().Run)
}

// GetPing
// Health check.
//
// @Response(app/logics/index.PingResponse)
func (o *Controller) GetPing(i iris.Context) interface{} {
	return logics.New(i, index.NewPing().Run)
}

// PostConsume
// Example consume.
//
// @Ignore(true)
// todo : debug for self consumed.
func (o *Controller) PostConsume(i iris.Context) interface{} {
	data := make(map[string]interface{})

	if b, be := json.Marshal(i.Request().Header); be == nil {
		data["header"] = string(b)
	}

	if b, be := i.GetBody(); be == nil {
		data["payload"] = string(b)
	}

	return response.With.Data(data)
}
