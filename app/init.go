// author: wsfuyibing <websearch@163.com>
// date: 2023-02-20

package app

import (
	"sync"
)

var (
	Config Configuration
)

func init() {
	new(sync.Once).Do(func() {
		Config = (&configuration{}).init()
	})
}
