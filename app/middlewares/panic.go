// author: wsfuyibing <websearch@163.com>
// date: 2023-01-16

package middlewares

import (
	"context"
	"github.com/fuyibing/log/v8"
	"github.com/fuyibing/log/v8/conf"
	"github.com/kataras/iris/v12"
	"net/http"
)

// Panic
//
// catch request panic.
func Panic(i iris.Context) {
	defer func() {
		var (
			c   context.Context
			err = recover()
		)

		if err == nil || i.IsStopped() {
			return
		}

		defer i.StopExecution()

		if t := i.Values().Get(conf.OpenTracingKey); t != nil {
			c = t.(context.Context)
		}

		log.Panicfc(c, "%v", err)

		ErrSend(i, http.StatusInternalServerError, err)
	}()

	i.Next()
}
