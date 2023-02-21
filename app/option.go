// author: wsfuyibing <websearch@163.com>
// date: 2023-02-20

package app

type (
	// Configuration
	// exported methods, like getter.
	Configuration interface {
		GetAdapter() string
		GetMemoryReloadSeconds() int
	}

	// Option
	// for configure easily, like setter.
	Option func(c *configuration)
)

func (o *configuration) GetAdapter() string { return o.Adapter }

func SetAdapter(adapter string) Option {
	return func(c *configuration) {
		c.Adapter = adapter
	}
}

func (o *configuration) GetMemoryReloadSeconds() int { return o.MemoryReloadSeconds }

func SetMemoryReloadSeconds(n int) Option {
	return func(c *configuration) { c.MemoryReloadSeconds = n }
}
