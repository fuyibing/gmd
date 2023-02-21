// author: wsfuyibing <websearch@163.com>
// date: 2023-02-09

// Package base
// Secondary level of core library.
package base

import (
	"sync"
)

func init() {
	new(sync.Once).Do(func() {
		Memory = (&memory{}).init()
		Pool = (&pool{}).init()
	})
}
