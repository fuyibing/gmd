// author: wsfuyibing <websearch@163.com>
// date: 2023-02-08

// Package md
// Core library for mq dispatcher.
package md

import (
	"sync"
)

func init() {
	new(sync.Once).Do(func() {
		Boot = (&boot{}).init()
	})
}
