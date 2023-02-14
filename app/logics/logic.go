// author: wsfuyibing <websearch@163.com>
// date: 2023-01-18

// Package logics
// Application logical.
package logics

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fuyibing/gmd/app"
	"github.com/fuyibing/log/v3"
	"github.com/fuyibing/log/v3/trace"
	"github.com/fuyibing/util/v2/web/response"
	"github.com/kataras/iris/v12"
	"net/http"
	"regexp"
	"time"
)

type (
	// Logic
	//
	// interface for main process.
	Logic func(c context.Context, i iris.Context) interface{}
)

// New
// create and return logic interface .
//
//   logics.New(i)
//   logics.New(i, example.Fn)
//   logics.New(i, example.NewExample().Run)
func New(i iris.Context, logics ...Logic) (res interface{}) {
	var (
		ctx, cx context.Context
		ok      bool
		t       = time.Now()
	)

	// Create open tracing context
	// based on middleware definitions.
	if x := i.Values().Get(trace.OpenTracingKey); x != nil {
		if cx, ok = x.(context.Context); ok {
			ctx = cx
		}
	}

	// Build default context
	// if request headers not specified.
	if ctx == nil {
		ctx = trace.FromRequest(i.Request())
	}

	// Called
	// before logic result returned.
	defer func() {
		// Override response result
		// if panic occurred in logic.
		if r := recover(); r != nil {
			log.Panicfc(ctx, "logic panic: %v", r)
			res = response.With.ErrorCode(fmt.Errorf("%v", r), app.CodeInternalError)
		}

		// Logger response result.
		dur := time.Now().Sub(t).Seconds()
		buf, _ := json.Marshal(res)
		log.Infofc(ctx, "logic finish: duration=%.06f, response: %s", dur, buf)

		// Next
		// middleware caller.
		i.Next()
	}()

	// Prepare standard logger info
	// before logic process.
	func() {
		text := fmt.Sprintf("logic begin: %s %s %s", i.Request().Proto, i.Request().Method, i.Request().RequestURI)

		// Append headers.
		if b, be := json.Marshal(i.Request().Header); be == nil {
			text += fmt.Sprintf(", headers: %s", b)
		}

		// Append request body.
		switch i.Request().Method {
		case http.MethodPost, http.MethodPut:
			if b, be := i.GetBody(); be == nil {
				if s := fmt.Sprintf("%s", b); s != "" {
					text += fmt.Sprintf(", body: %s",
						regexp.MustCompile(`\n\s*`).ReplaceAllString(s, ""),
					)
				}
			}
		}

		// Logger request params.
		log.Infofc(ctx, text)
	}()

	// Return
	// registered logic response.
	if len(logics) > 0 {
		res = logics[0](ctx, i)
		return
	}

	// Return
	// unknown logic error.
	res = response.With.ErrorCode(
		fmt.Errorf("logic callable not specified"),
		http.StatusBadRequest,
	)
	return
}
