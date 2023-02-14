// author: wsfuyibing <websearch@163.com>
// date: 2023-01-16

package middlewares

import (
	"github.com/fuyibing/log/v3/trace"
	"github.com/kataras/iris/v12"
)

// Tracer
//
// initialize open tracing.
func Tracer(i iris.Context) {
	i.Values().Set(trace.OpenTracingKey, trace.FromRequest(i.Request()))
	i.Next()
}
