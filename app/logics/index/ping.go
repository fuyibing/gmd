// author: wsfuyibing <websearch@163.com>
// date: 2023-01-19

package index

import (
	"context"
	"github.com/fuyibing/gmd/app"
	"github.com/fuyibing/gmd/app/models"
	"github.com/fuyibing/util/v2/web/response"
	"github.com/kataras/iris/v12"
	"os"
	"runtime"
)

type (
	Ping struct {
		response *PingResponse
	}

	PingResponse struct {
		Cpu       int     `json:"cpu" label:"CPU core count" mock:"8"`
		Goroutine int     `json:"goroutines" label:"Coroutine counts" mock:"32"`
		Memory    float64 `json:"memory" label:"Used system memory" desc:"Unit: MB" mock:"16.57"`
		Pid       int     `json:"pid" label:"Process ID" mock:"3721"`
		StartTime string  `json:"start_time" label:"Started time" mock:"2022-01-19 14:21:25"`
	}
)

func NewPing() *Ping {
	return &Ping{
		response: &PingResponse{
			Cpu:       runtime.NumCPU(),
			Goroutine: runtime.NumGoroutine(),
			Pid:       os.Getpid(),
			StartTime: app.Config.StartTime.Format(models.GmtTimeLayout),
		},
	}
}

func (o *Ping) Run(_ context.Context, _ iris.Context) interface{} {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)

	o.response.Memory = float64(int((float64(m.Sys)/1024/1024)*100)) / 100
	return response.With.Data(o.response)
}
