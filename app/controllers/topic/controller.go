// author: wsfuyibing <websearch@163.com>
// date: 2023-01-17

// Package topic
// MVC Controller with route prefix /topic.
package topic

import (
	"github.com/fuyibing/gmd/app/logics"
	"github.com/fuyibing/gmd/app/logics/topic"
	"github.com/kataras/iris/v12"
)

type (
	// Controller
	// Topic.
	//
	// @RoutePrefix(/topic)
	Controller struct{}
)

// PostBatch
// Publish multiple.
//
// Each request can publish multiple messages, up to 100.
// Asynchronous mode.
//
// @Request(app/logics/topic.BatchRequest)
// @Response(app/logics/topic.BatchResponse)
func (o *Controller) PostBatch(i iris.Context) interface{} {
	return logics.New(i, topic.NewBatch().Run)
}

// PostPublish
// Publish one.
//
// Only 1 message can be published per request.
// Asynchronous mode.
//
// @Request(app/logics/topic.PublishRequest)
// @Response(app/logics/topic.PublishResponse)
func (o *Controller) PostPublish(i iris.Context) interface{} {
	return logics.New(i, topic.NewPublish().Run)
}
