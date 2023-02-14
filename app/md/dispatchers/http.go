// author: wsfuyibing <websearch@163.com>
// date: 2023-02-09

package dispatchers

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"net/http"
	"time"
)

type (
	// HttpDispatcher
	// struct of http dispatcher.
	HttpDispatcher struct {
		Request  *fasthttp.Request
		Response *fasthttp.Response
	}
)

func (o *HttpDispatcher) Release() {
	Pool.ReleaseHttp(o)
}

// Run
// 执行投递.
func (o *HttpDispatcher) Run(timeout int) (body []byte, err error) {
	// Called
	// when end.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	// Send
	// request process.
	if timeout > 0 {
		err = fasthttp.DoTimeout(o.Request, o.Response, time.Duration(timeout)*time.Second)
	} else {
		err = fasthttp.Do(o.Request, o.Response)
	}

	// Set response body
	// if no error occurred.
	if err == nil {
		body = o.Response.Body()

		// Set error
		// if response status code not matched.
		if code := o.Response.StatusCode(); code != http.StatusOK {
			err = fmt.Errorf("HTTP %d %s", code, http.StatusText(code))
		}
	}

	return
}

// /////////////////////////////////////////////////////////////
// Access methods.
// /////////////////////////////////////////////////////////////

func (o *HttpDispatcher) after() {
	fasthttp.ReleaseRequest(o.Request)
	o.Request = nil

	fasthttp.ReleaseResponse(o.Response)
	o.Response = nil
}

func (o *HttpDispatcher) before() {
	o.Request = fasthttp.AcquireRequest()
	o.Response = fasthttp.AcquireResponse()
}

func (o *HttpDispatcher) init() *HttpDispatcher {
	return o
}
