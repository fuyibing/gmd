// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package el

import (
	"github.com/fuyibing/gmd/v8/core/base"
)

// Validation
// 条件校验.
type Validation struct{}

func New() *Validation { return (&Validation{}).init() }

func (o *Validation) Validate(_ *base.Message) (ignored bool, err error) { return }

func (o *Validation) init() *Validation {
	return o
}
