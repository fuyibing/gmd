// author: wsfuyibing <websearch@163.com>
// date: 2023-01-17

// Package task
// MVC Controller with route prefix /task.
package task

import (
	"github.com/fuyibing/gmd/app/logics"
	"github.com/fuyibing/gmd/app/logics/task"
	"github.com/kataras/iris/v12"
)

type (
	// Controller
	// Task.
	//
	// @RoutePrefix(/task)
	Controller struct{}
)

// PostAdd
// Add new task.
//
// @Request(app/logics/task.AddRequest)
// @Response(app/logics/task.AddResponse)
func (o *Controller) PostAdd(i iris.Context) interface{} {
	return logics.New(i, task.NewAdd().Run)
}

// PostDel
// Delete task.
func (o *Controller) PostDel(i iris.Context) interface{} {
	return logics.New(i)
}

// PostDisable
// Disable task.
//
// @Request(app/logics/task.EditStatus)
// @Response(app/logics/task.EditResponse)
func (o *Controller) PostDisable(i iris.Context) interface{} {
	return logics.New(i, task.NewEditDisable().Run)
}

// PostEdit
// Edit task basic fields.
//
// @Request(app/logics/task.EditRequest)
// @Response(app/logics/task.EditResponse)
func (o *Controller) PostEdit(i iris.Context) interface{} {
	return logics.New(i, task.NewEdit().Run)
}

// PostEditFailed
// Edit task failed notification.
//
// When message consumption fails, forward the last delivery result
// to the failed callback.
//
// @Request(app/logics/task.EditSubscriber)
// @Response(app/logics/task.EditResponse)
func (o *Controller) PostEditFailed(i iris.Context) interface{} {
	return logics.New(i, task.NewEditFailed().Run)
}

// PostEditHandler
// Edit task subscriber.
//
// When the consumer receives the message, it will be delivered to the
// specified callback.
//
// @Request(app/logics/task.EditSubscriber)
// @Response(app/logics/task.EditResponse)
func (o *Controller) PostEditHandler(i iris.Context) interface{} {
	return logics.New(i, task.NewEditHandler().Run)
}

// PostEditSucceed
// Edit task succeed notification.
//
// When the message consumption is successful, forward the delivery
// result to the successful callback.
//
// @Request(app/logics/task.EditSubscriber)
// @Response(app/logics/task.EditResponse)
func (o *Controller) PostEditSucceed(i iris.Context) interface{} {
	return logics.New(i, task.NewEditSucceed().Run)
}

// PostEnable
// Enable task.
//
// @Request(app/logics/task.EditStatus)
// @Response(app/logics/task.EditResponse)
func (o *Controller) PostEnable(i iris.Context) interface{} {
	return logics.New(i, task.NewEditEnable().Run)
}

// PostRemoteBuild
// Build task remote relations on mq server.
//
// @Request(app/logics/task.RemoteBuildRequest)
// @Response(app/logics/task.RemoteBuildResponse)
func (o *Controller) PostRemoteBuild(i iris.Context) interface{} {
	return logics.New(i, task.NewRemoteBuild().Run)
}

// PostRemoteDestroy
// Destroy task remote relations of mq server.
//
// @Request(app/logics/task.RemoteDestroyRequest)
// @Response(app/logics/task.RemoteDestroyResponse)
func (o *Controller) PostRemoteDestroy(i iris.Context) interface{} {
	return logics.New(i, task.NewRemoteDestroy().Run)
}
