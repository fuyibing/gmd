// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package main

import (
	"context"
	"github.com/fuyibing/gmd/v8/core/managers"
	"github.com/fuyibing/log/v5"
	"github.com/fuyibing/log/v5/exporters/logger_term"
	"github.com/fuyibing/log/v5/exporters/tracer_jaeger"
	"time"
)

func init() {
	log.Manager.LoggerManager().SetExporter(logger_term.NewExporter())
	log.Manager.TracerManager().SetExporter(tracer_jaeger.NewExporter())
	go log.Manager.Start(context.Background())
}

func main() {
	var (
		cancel context.CancelFunc
		ctx    context.Context
		err    error
	)

	defer func() {
		if err != nil {
			log.Error("main error: %v", err)
		}

		if ctx != nil && ctx.Err() == nil {
			cancel()
		}

		managers.Boot.Stop()

		log.Manager.Stop()
	}()

	go func() {
		time.Sleep(time.Second * 2)
		if err == nil {

			log.Warn("------ call restart")
			managers.Boot.Restart()

			time.Sleep(time.Second * 5)
			log.Warn("------ call canceller")
			cancel()
		}
	}()

	ctx, cancel = context.WithCancel(context.Background())
	err = managers.Boot.Start(ctx)
}
