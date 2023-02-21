// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package el

import (
	"github.com/fuyibing/gmd/v8/core/base"
)

// Validation
// verify message body.
type Validation struct{}

func New() *Validation {
	return (&Validation{}).init()
}

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *Validation) Validate(_ *base.Message) (ignored bool, err error) { return }

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *Validation) init() *Validation {
	return o
}
