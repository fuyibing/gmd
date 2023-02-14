// author: wsfuyibing <websearch@163.com>
// date: 2023-01-16

package middlewares

import (
	"fmt"
	"github.com/fuyibing/util/v2/web/response"
	"github.com/kataras/iris/v12"
	"net/http"
)

// ErrCode
//
// handler error http status code.
func ErrCode(i iris.Context) {
	ErrSend(i, i.GetStatusCode(), nil)
}

func ErrSend(i iris.Context, c int, v interface{}) {
	defer i.StatusCode(http.StatusOK)

	s := fmt.Sprintf("%s %s %s", i.Request().Proto, i.Request().Method, i.Request().RequestURI)

	if v != nil {
		s += fmt.Sprintf(", %v", v)
	} else {
		s += fmt.Sprintf(", %v", http.StatusText(c))
	}

	_, _ = i.JSON(response.With.ErrorCode(fmt.Errorf(s), c))
}
