// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package base

type (
	// ResultCallable
	// 结果校验构造函数.
	ResultCallable func() ResultManager

	// ResultManager
	// 结果校验管理器.
	ResultManager interface {
		// Validate
		// 校验结果.
		Validate(code int, body []byte) (err error)
	}
)

// 结果类型枚举.

const (
	ResultHttpOk        = "HttpStatusOk"
	ResultJsonErrnoZero = "JsonErrnoZero"
)
