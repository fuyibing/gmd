// author: wsfuyibing <websearch@163.com>
// date: 2023-01-17

// Package controllers
// MVC Controller with route prefix /.
package controllers

import (
	"github.com/fuyibing/gmd/app/logics"
	"github.com/fuyibing/gmd/app/logics/index"
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
