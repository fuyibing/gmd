// author: wsfuyibing <websearch@163.com>
// date: 2023-02-01

// Package main
// Message dispatcher application.
package main

import (
	"context"
	"github.com/fuyibing/console/v3"
	"github.com/fuyibing/console/v3/managers"
	"github.com/fuyibing/gdoc/adapters/markdown/i18n"
	"github.com/fuyibing/gmd/app"
	"github.com/fuyibing/gmd/app/controllers"
	"github.com/fuyibing/gmd/app/md"
	"github.com/fuyibing/gmd/app/md/conf"
	"github.com/fuyibing/gmd/app/middlewares"
	"github.com/fuyibing/log/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/pprof"
	"github.com/kataras/iris/v12/mvc"
	"time"
)

var (
	ce error
	cm managers.Manager
)

type (
	// Bootstrap
	//
	// struct for bootstrap definitions.
	Bootstrap struct {
		c      managers.Command
		cancel context.CancelFunc
		ctx    context.Context
		fw     *iris.Application
	}
)

// DoBefore
//
// called when application start fired. Start message
// dispatcher boot manager in coroutine.
func (o *Bootstrap) DoBefore(_ *iris.Application) {
	go func() {
		if err := md.Boot.Processor().Start(o.ctx); err != nil {
			log.Errorf("%v", err)
		}
	}()
}

// DoInterrupt
//
// called when SIGTERM/SIGINT signal received. Block coroutine
// until message dispatcher boot manager stopped.
func (o *Bootstrap) DoInterrupt() {
	// Cancel context
	// if it is running.
	if o.ctx != nil && o.ctx.Err() == nil {
		o.cancel()
	}

	// Block process
	// until message dispatcher boot stopped.
	for {
		if md.Boot.Processor().Stopped() {
			break
		}

		time.Sleep(conf.EventSleepDuration)
	}
}

// Initialize
//
// called when registered into console manager.
func (o *Bootstrap) Initialize() *Bootstrap {
	o.InitConsoleCommand()
	o.InitFramework()
	return o
}

// InitConsoleCommand
//
// called in initialize method.
func (o *Bootstrap) InitConsoleCommand() {
	o.c = managers.NewCommand("start").
		SetDescription("Start gmd service").
		SetHandler(o.Run)
}

// InitFramework
//
// called in initialize method.
func (o *Bootstrap) InitFramework() {
	// Register callback
	// which called when application quit signal fired.
	//
	// Example: Ctrl+C
	// Example: SIGTERM, SIGINT
	iris.RegisterOnInterrupt(o.DoInterrupt)

	// Build
	// basic fields and framework settings.
	o.fw = iris.New()
	o.fw.Logger().SetLevel("disable")
	o.fw.Configure(iris.WithConfiguration(iris.Configuration{
		DisableBodyConsumptionOnUnmarshal: true,
		DisableStartupLog:                 true,
		EnableOptimizations:               true,
		TimeFormat:                        "2006-01-02 15:04:05",
	}))

	// Initialize framework extensions.
	o.InitFrameworkMiddlewares()
	o.InitFrameworkStatusCode()
	o.InitFrameworkProfile()
	o.InitFrameworkControllers()
}

// InitFrameworkControllers
//
// called in initialize method. It register mvc controller
// in iris framework.
func (o *Bootstrap) InitFrameworkControllers() {
	for k, c := range controllers.ControllerRegistration {
		func(key string, controller interface{}) {
			mvc.Configure(o.fw.Party(key), func(application *mvc.Application) {
				application.Handle(controller)
			})
		}(k, c)
	}
}

// InitFrameworkMiddlewares
//
// called in initialize method, It register middlewares on
// each request.
func (o *Bootstrap) InitFrameworkMiddlewares() {
	o.fw.UseGlobal(middlewares.Tracer, middlewares.Panic)
}

// InitFrameworkProfile
//
// called in initialize method. It register debug profile
// routes.
func (o *Bootstrap) InitFrameworkProfile() {
	p := pprof.New()
	o.fw.Any("/debug/pprof", p)
	o.fw.Any("/debug/pprof/{action:path}", p)
}

// InitFrameworkStatusCode
//
// called in initialize method. It register error http status
// callback handler.
func (o *Bootstrap) InitFrameworkStatusCode() {
	o.fw.OnAnyErrorCode(middlewares.ErrCode)
}

// Run
//
// called by console manager.
func (o *Bootstrap) Run(_ managers.Manager, _ managers.Arguments) error {
	// Context definitions
	// for message dispatcher boot manager.
	o.ctx, o.cancel = context.WithCancel(context.Background())

	// Clean called
	// when main process ended.
	defer func() {
		// Catch panic.
		if r := recover(); r != nil {
			log.Panicf("%v", r)
		}

		// Cancel context
		// if running.
		if o.ctx.Err() == nil {
			o.cancel()
		}

		// Unset context.
		o.cancel = nil
		o.ctx = nil
	}()

	// Run iris framework.
	log.Infof("server begin: pid=%d, name=%s, host=%s, port=%v", app.Config.Pid, app.Config.Name, app.Config.Host, app.Config.Port)
	defer log.Infof("server finish")
	_ = o.fw.Run(iris.Addr(app.Config.Addr), o.DoBefore)
	return nil
}

// /////////////////////////////////////////////////////////////
// Package methods
// /////////////////////////////////////////////////////////////

func init() {
	// Initialize console
	// then add bootstrap command into manager.
	if cm, ce = console.Default(); ce == nil {
		ce = cm.AddCommand((&Bootstrap{}).Initialize().c)

		// Document language as nil.
		i18n.SetLang(nil)
	}
}

func main() {
	// Run console
	// with terminal arguments.
	if ce == nil {
		ce = cm.RunTerminal()
	}

	// Logger console error
	// if executed failed.
	if ce != nil {
		log.Errorf("console error: %v", ce)
	}
}
