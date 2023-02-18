// author: wsfuyibing <websearch@163.com>
// date: 2023-01-19

package index

import (
	"context"
	"github.com/fuyibing/gmd/app"
	"github.com/fuyibing/util/v8/web/response"
	"github.com/kataras/iris/v12"
	"time"
)

type (
	Home struct {
		response *HomeResponse
	}

	HomeResponse struct {
		Time time.Time `json:"time" label:"Started time"`
	}
)

func NewHome() *Home {
	return &Home{
		response: &HomeResponse{
			Time: app.Config.StartTime,
		},
	}
}

func (o *Home) Run(_ context.Context, _ iris.Context) interface{} {
	return response.With.Data(o.response)
}
