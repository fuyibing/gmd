// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package base

type (
	// Condition
	// defined as condition manager name.
	Condition string

	// ConditionManager
	// validate received message should dispatcher to subscription
	// handler or not.
	ConditionManager interface {
	}
)

// Condition enums.

const (
	ConditionEl Condition = "EL"
)
