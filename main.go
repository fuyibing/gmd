// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package main

import (
	"context"
	"fmt"
	cs "github.com/fuyibing/console/v3"
	cm "github.com/fuyibing/console/v3/managers"
	"github.com/fuyibing/gmd/v8/app"
	"github.com/fuyibing/gmd/v8/app/controllers"
	"github.com/fuyibing/gmd/v8/app/middlewares"
	"github.com/fuyibing/gmd/v8/core/managers"
	"github.com/fuyibing/log/v5"
	"github.com/fuyibing/log/v5/conf"
	"github.com/fuyibing/log/v5/cores"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/pprof"
	"github.com/kataras/iris/v12/mvc"
	"os"
	"time"
)

var (
	ctx context.Context
	err error
	mng cm.Manager
	my  *gmd
)

type (
	// gmd
	// golang message dispatcher.
	gmd struct {
		cancel context.CancelFunc
		cmd    cm.Command
		ctx    context.Context

		framework *iris.Application
	}
)

func (o *gmd) error(text string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, fmt.Sprintf(text, args...))
}

// init
// 构造项目.
func (o *gmd) init() *gmd {
	o.cmd = cm.NewCommand("start").
		SetDescription("start gmd service").
		SetHandler(o.run)
	return o
}

// initFramework
// 初始化IRIS框架.
func (o *gmd) initFramework() {
	o.framework = iris.New()
	o.framework.Logger().SetLevel("disable")
	o.framework.Configure(iris.WithConfiguration(iris.Configuration{
		DisableBodyConsumptionOnUnmarshal: true,
		DisableStartupLog:                 true,
		EnableOptimizations:               true,
		TimeFormat:                        "2006-01-02 15:04:05",
	}))

	o.framework.OnAnyErrorCode(middlewares.ErrCode)
}

// initFrameworkControllers
// 注册MVC控制器.
func (o *gmd) initFrameworkControllers() {
	for k, c := range controllers.Containers {
		func(key string, controller interface{}) {
			mvc.Configure(o.framework.Party(key), func(application *mvc.Application) {
				application.Handle(controller)
			})
		}(k, c)
	}
}

// initFrameworkMiddlewares
// 注册中间件.
func (o *gmd) initFrameworkMiddlewares() {
	o.framework.Use(
		middlewares.Tracer,
		middlewares.Panic,
	)
}

func (o *gmd) initFrameworkProfiles() {
	p := pprof.New()
	o.framework.Any("/debug/pprof", p)
	o.framework.Any("/debug/pprof/{action:path}", p)
}

// run
// 执行项目.
func (o *gmd) run(_ cm.Manager, _ cm.Arguments) error {
	iris.RegisterOnInterrupt(o.runInterrupt)

	o.initFramework()
	o.initFrameworkMiddlewares()
	o.initFrameworkControllers()
	o.initFrameworkProfiles()

	// 启动服务.
	o.runServe()

	// 卸载日志, 阻塞协程直到全部上报完成.
	log.Manager.Stop()
	return nil
}

// 加载内核.
func (o *gmd) runBeforeLoadCore(_ *iris.Application) {
	o.ctx, o.cancel = context.WithCancel(ctx)
	go func(c context.Context) { _ = managers.Boot.Start(c) }(o.ctx)
}

// 加载日志.
func (o *gmd) runBeforeLoadLogger(_ *iris.Application) {
	go func() {
		// 覆盖配置.
		conf.Config.With(
			conf.ServiceName(app.Config.GetName()),
			conf.ServicePort(app.Config.GetPort()),
			conf.ServiceVersion(app.Config.GetVersion()),
		)

		// 更新资源.
		cores.Registry.Update()

		// 启动日志.
		if el := log.Manager.Start(ctx); el != nil {
			_, _ = fmt.Fprintf(os.Stderr, fmt.Sprintf("%v", el))
		}
	}()
}

func (o *gmd) runInterrupt() {
	// 退出信号.
	if o.ctx != nil && o.ctx.Err() == nil {
		o.cancel()
	}

	// 等待完成.
	for {
		if managers.Boot.Stopped() {
			return
		}
		time.Sleep(time.Millisecond * 30)
	}
}

func (o *gmd) runServe() {
	if err = o.framework.Configure(
		o.runBeforeLoadLogger,
		o.runBeforeLoadCore,
	).Run(
		iris.Addr(fmt.Sprintf("%s:%d", app.Config.GetHost(), app.Config.GetPort())),
	); err != nil {
		log.Error("%v", err)
	}
}

// init
// 项目构造.
func init() {
	if mng, err = cs.New(); err == nil {
		my = (&gmd{}).init()
		err = mng.AddCommand(my.cmd)
	}
}

// main
// 项目执行.
func main() {
	if err == nil {
		ctx = context.Background()
		err = mng.RunTerminal()
	}
	if err != nil {
		my.error("%v", err)
	}
}
