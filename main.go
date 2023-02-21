// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package main

import (
	"context"
	"github.com/fuyibing/gmd/v8/core/managers"
	"github.com/fuyibing/log/v8"
	"time"
)

func init() {
}

func main() {
	defer log.Client.Close()

	ctx, canceler := context.WithCancel(context.Background())
	t := time.Now().Format("15:04:05")

	go func() {
		time.Sleep(time.Second * 3)
		log.Warnf("---------------- [restart=%s] ----------------", t)
		managers.Boot.Restart()

		time.Sleep(time.Second * 5)
		log.Warnf("---------------- [cancelled=%s] ----------------", t)
		canceler()
	}()

	log.Warnf("---------------- [ready=%s] ----------------", t)
	_ = managers.Boot.Start(ctx)
}
