// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package main

import (
	"context"
	"github.com/fuyibing/gmd/v8/app/md"
	"github.com/fuyibing/log/v8"
	"time"
)

func init() {
}

func main() {
	var (
		cancel context.CancelFunc
		ctx    context.Context
		err    error
	)

	defer func() {
		if err != nil {
			log.Errorf("main error: %v", err)
		}

		if ctx != nil && ctx.Err() == nil {
			cancel()
		}

		log.Client.Close()
	}()

	if err = md.Boot.Prepare(); err != nil {
		log.Errorf("boot: %v", err)
		return
	}

	go func() {
		time.Sleep(time.Second * 2)
		if err == nil {

			log.Warnf("------ call restart")
			md.Boot.Restart()

			time.Sleep(time.Second * 5)
			log.Warnf("------ call canceller")
			cancel()
		}
	}()

	ctx, cancel = context.WithCancel(context.Background())
	err = md.Boot.Start(ctx)
}
