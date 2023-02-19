// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package conf

import (
	"github.com/fuyibing/gmd/v8/app/md/base"
)

type (
	Configuration interface {
		GetAdapter() base.Adapter
		Set(options ...Option)
	}
)

func (o *configuration) GetAdapter() base.Adapter {
	return o.Adapter
}

func (o *configuration) Set(options ...Option) {
	for _, option := range options {
		option(o)
	}
}
