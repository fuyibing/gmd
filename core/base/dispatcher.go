// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package base

type (
	// DispatcherCallable
	// 消费分发构造函数.
	DispatcherCallable func() DispatcherManager

	// DispatcherManager
	// 消费分发管理器.
	DispatcherManager interface {
		// Dispatch
		// 分发过程.
		Dispatch(task *Task, subscriber *Subscriber, message *Message) (err error)
	}
)

// 分发类型枚举.

const (
	DispatchHttpGet      = "HTTP_GET"
	DispatchHttpPostForm = "HTTP_POST_FORM"
	DispatchHttpPostJson = "HTTP_POST_JSON"
	DispatchRpc          = "RPC"
	DispatchTcp          = "TCP"
	DispatchWebsocket    = "WSS"
)
