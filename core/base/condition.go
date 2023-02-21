// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package base

type (
	// ConditionCallable
	// constructor for create ConditionManager instance.
	ConditionCallable func() ConditionManager

	// ConditionManager
	// validate received message should dispatcher to subscription
	// handler or not.
	ConditionManager interface {
		// Validate
		// verify whether the message body matches the configuration.
		//
		// If the return value is true, it means that it does not match
		// the specified configuration, and message needs to be ignored,
		// otherwise needs to deliver to subscriber.
		//
		// If an error is returned, the message body verification error.
		Validate(message *Message) (ignored bool, err error)
	}
)

// Condition enums.

const (
	ConditionEl = "EL"
)
