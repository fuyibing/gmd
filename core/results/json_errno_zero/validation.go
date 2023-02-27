// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package json_errno_zero

// Validation
// 结果校验.
type Validation struct{}

func New() *Validation { return (&Validation{}).init() }

func (o *Validation) Validate(_ int, _ []byte) (err error) { return }

func (o *Validation) init() *Validation {
	return o
}
