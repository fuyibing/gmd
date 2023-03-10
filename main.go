// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// author: wsfuyibing <websearch@163.com>
// date: 2023-03-07

package main

import (
	"context"
	"fmt"
	cc "github.com/fuyibing/console/v3"
	cm "github.com/fuyibing/console/v3/managers"
	"github.com/fuyibing/gmd/v8/app"
	"github.com/fuyibing/gmd/v8/app/controllers"
	"github.com/fuyibing/gmd/v8/app/middlewares"
	"github.com/fuyibing/gmd/v8/md"
	"github.com/fuyibing/log/v5"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/pprof"
	"github.com/kataras/iris/v12/mvc"
	"time"
)

var (
	App *GmdApp
)

type (
	// GmdApp
	// 消息分发.
	GmdApp struct {
		Cancel         context.CancelFunc
		ConsoleCommand cm.Command
		ConsoleManager cm.Manager
		Ctx            context.Context
		Err            error
		Framework      *iris.Application
	}
)

// +---------------------------------------------------------------------------+
// + Common methods                                                            |
// +---------------------------------------------------------------------------+

func (o *GmdApp) OnInterrupt() {
	var (
		max = 300
		ms  = 100 * time.Millisecond
	)

	// 待等结束.
	// 连续 30 秒内, 每隔 100 毫秒检查 MD 是否结束.
	md.Manager.Stop()
	for i := 0; i < max; i++ {
		if md.Manager.Processor().Stopped() {
			break
		}
		time.Sleep(ms)
	}
}

func (o *GmdApp) OnStart(_ *iris.Application) {
	go func() {
		if err := md.Manager.Start(o.Ctx); err != nil {
			log.Error("gmd error: %v", err)
		}
	}()
}

func (o *GmdApp) Start() {
	// 处理管理.
	// - 开始时启动日志.
	// - 结束前关闭日志.
	log.Manager.Start(o.Ctx)
	defer log.Manager.Stop()

	// 启动服务.
	if err := o.Framework.Configure(o.OnStart).Run(
		iris.Addr(fmt.Sprintf("%s:%d", app.Config.GetHost(), app.Config.GetPort())),
	); err != nil {
		log.Error("server closed: %v", err)
	}
}

func (o *GmdApp) Stop() { o.Cancel() }

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *GmdApp) cmdRegister() {
	o.ConsoleCommand = cm.NewCommand("start").
		SetDescription("start application").
		SetHandler(o.gmdRunner)

	if o.ConsoleManager, o.Err = cc.Latest(); o.Err == nil {
		o.Err = o.ConsoleManager.AddCommand(o.ConsoleCommand)
	}
}

func (o *GmdApp) gmdRunner(_ cm.Manager, _ cm.Arguments) error {
	o.Start()
	return nil
}

func (o *GmdApp) init() *GmdApp {
	iris.RegisterOnInterrupt(o.OnInterrupt)

	o.Ctx, o.Cancel = context.WithCancel(context.Background())
	o.Framework = iris.New()
	o.initFramework()
	o.initControllers()
	o.initDebugProfile()
	o.initMiddlewares()

	o.cmdRegister()
	return o
}

func (o *GmdApp) initControllers() {
	for k, c := range controllers.Registry {
		func(path string, controller interface{}) {
			mvc.Configure(o.Framework.Party(path), func(application *mvc.Application) {
				application.Handle(controller)
			})
		}(k, c)
	}
}

func (o *GmdApp) initDebugProfile() {
	p := pprof.New()
	o.Framework.Any("/debug/pprof", p)
	o.Framework.Any("/debug/pprof/{action:path}", p)
}

func (o *GmdApp) initFramework() {
	o.Framework.Configure(iris.WithConfiguration(iris.Configuration{
		DisableBodyConsumptionOnUnmarshal: true,
		DisableStartupLog:                 true,
		EnableOptimizations:               true,
		TimeFormat:                        "2006-01-02 15:04:05",
	})).Logger().SetLevel("disable")
}

func (o *GmdApp) initMiddlewares() {
	// 中间件.
	for _, middleware := range middlewares.Registry {
		o.Framework.UseGlobal(middleware)
	}

	// 错误码.
	o.Framework.OnAnyErrorCode(middlewares.ErrCode)
}

func init() { App = (&GmdApp{}).init() }

func main() {
	App.ConsoleManager.RunTerminal()
	// App.Start()
}
