// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package base

type (
	// DispatcherCallable
	// constructor for create DispatcherManager instance.
	DispatcherCallable func() DispatcherManager

	// DispatcherManager
	// dispatch received message to subscription handler.
	DispatcherManager interface {
		Dispatch(task *Task, subscriber *Subscriber, message *Message) (err error)
	}
)

// Dispatcher enums.

const (
	DispatchHttpGet      = "HTTP_GET"
	DispatchHttpPostForm = "HTTP_POST_FORM"
	DispatchHttpPostJson = "HTTP_POST_JSON"
	DispatchRpc          = "RPC"
	DispatchTcp          = "TCP"
	DispatchWebsocket    = "WSS"
)
