// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package md

import (
	"github.com/fuyibing/gmd/v8/app/md/conf"
	"github.com/fuyibing/gmd/v8/app/md/core"
	"sync"
)

var (
	Boot   core.BootManager
	Config conf.Configuration
)

func init() {
	new(sync.Once).Do(func() {
		Boot = core.Boot
		Config = conf.Config
	})
}
