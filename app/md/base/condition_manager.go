// author: wsfuyibing <websearch@163.com>
// date: 2023-02-08

package base

type (
	// ConditionManager
	// interface of condition manager.
	ConditionManager interface {
		// Expression
		// return registered expression string.
		//
		// Defined in follow columns.
		//   - task.handler_condition
		//   - task.failed_condition
		//   - task.succeed_condition
		Expression() string

		// MatchJsonString
		// return result after matched json string.
		//
		// Return true if matched with json string, otherwise
		// false return.
		MatchJsonString(s string) (ignored bool, err error)
	}

	condition struct {
		s string
	}
)

// /////////////////////////////////////////////////////////////
// Interface methods.
// /////////////////////////////////////////////////////////////

func (o *condition) Expression() string                     { return o.s }
func (o *condition) MatchJsonString(s string) (bool, error) { return o.matchJsonString(s) }

// /////////////////////////////////////////////////////////////
// Match methods.
// /////////////////////////////////////////////////////////////

func (o *condition) matchJsonString(_ string) (ignored bool, err error) {
	// todo : match json string on condition
	// return true, nil
	return false, nil
}

// /////////////////////////////////////////////////////////////
// Access methods.
// /////////////////////////////////////////////////////////////

func (o *condition) init() *condition {
	return o
}
