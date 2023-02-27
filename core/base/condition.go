// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package base

type (
	// ConditionCallable
	// 条件管理器构造函数.
	ConditionCallable func() ConditionManager

	// ConditionManager
	// 条件管理器接口.
	ConditionManager interface {
		// Validate
		// 条件校验.
		//
		// 如果返回的 ignored 值为 true 表示消息不满足条件需要被忽略, 反之需分
		// 发给订阅方, 如果返回的 err 值非 nil, 表示条件校验格式错误.
		Validate(message *Message) (ignored bool, err error)
	}
)

// 条件枚举.

const (
	ConditionEl = "EL"
)
