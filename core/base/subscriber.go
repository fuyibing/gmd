// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package base

type (
	Subscriber struct {
		// Subscriber handler.
		//
		// - http://example.com/route/path?key=value
		// - 172.16.10.110:8088
		Handler string

		// Condition  ConditionManager
		// Dispatcher DispatcherManager
		// Result     ResultManager
	}
)
