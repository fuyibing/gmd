// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package base

type (
	// Dispatcher
	// defined as dispatcher manager name.
	Dispatcher string

	// DispatcherManager
	// dispatch received message to subscription handler.
	DispatcherManager interface {
	}
)

// Dispatcher enums.

const (
	DispatchHttp Dispatcher = "http"
	DispatchRpc  Dispatcher = "grpc"
)
