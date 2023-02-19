// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package conf

import (
	"github.com/fuyibing/gmd/v8/app/md/base"
)

type Option func(c *configuration)

func SetAdapter(a base.Adapter) Option {
	return func(c *configuration) {
		c.Adapter = a
	}
}
