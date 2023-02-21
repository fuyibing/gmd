// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package base

type (
	// ResultCallable
	// constructor for create ResultManager instance.
	ResultCallable func() ResultManager

	// ResultManager
	// validate dispatcher result from subscription handler
	// is succeed or failed.
	ResultManager interface {
		// Validate
		// verify dispatcher result.
		Validate(code int, body []byte) (err error)
	}
)

// Result enums.

const (
	ResultHttpOk        = "HttpStatusOk"
	ResultJsonErrnoZero = "JsonErrnoZero"
)
