// author: wsfuyibing <websearch@163.com>
// date: 2023-02-20

package app

type (
	// Option
	// 配置选项接口.
	Option func(c *configuration)
)

func SetAdapter(adapter string) Option {
	return func(c *configuration) {
		c.Adapter = adapter
	}
}

func SetMemoryReloadSeconds(n int) Option {
	return func(c *configuration) { c.MemoryReloadSeconds = n }
}
