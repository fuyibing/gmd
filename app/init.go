// author: wsfuyibing <websearch@163.com>
// date: 2021-08-04

// Package app
// Application core.
package app

import (
	"sync"
)

func init() {
	new(sync.Once).Do(func() {
		Config = (&Configuration{}).init()
	})
}
