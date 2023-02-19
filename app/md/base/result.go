// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package base

type (
	// Result
	// defined for result validation name.
	Result string

	// ResultManager
	// validate dispatcher result from subscription handler
	// is succeed or failed.
	ResultManager interface {
	}
)

// Result enums.

const (
	ResultHttpOk        Result = "HttpStatusOk"
	ResultJsonErrnoZero Result = "JsonErrnoZero"
)
