// author: wsfuyibing <websearch@163.com>
// date: 2023-02-02

// Package conf
// Top level of core library configurations.
package conf

import (
	"sync"
)

func init() {
	new(sync.Once).Do(func() {
		Config = (&Configuration{}).init()
	})
}
