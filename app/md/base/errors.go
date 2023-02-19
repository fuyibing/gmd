// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package base

import (
	"fmt"
)

var (
	ErrUnknown = fmt.Errorf("unknown error")

	ErrConsumerCallableNotConfigured = fmt.Errorf("consumer callable not configured")
	ErrProducerCallableNotConfigured = fmt.Errorf("producer callable not configured")
	ErrRemotingCallableNotConfigured = fmt.Errorf("remoting callable not configured")
)
