// author: wsfuyibing <websearch@163.com>
// date: 2023-02-20

package base

type (
	// ContainerManager
	// 容器管理器.
	ContainerManager interface {
		// GetConsumer
		// 获取MQ中间件的消费者构造函数.
		GetConsumer() (callable ConsumerCallable)

		// SetConsumer
		// 设置MQ中间件的消费者构造函数.
		SetConsumer(callable ConsumerCallable)

		// GetProducer
		// 获取MQ中间件的生产者构造函数.
		GetProducer() (callable ProducerCallable)

		// SetProducer
		// 设置MQ中间件的生产者构造函数.
		SetProducer(callable ProducerCallable)

		// GetRemoting
		// 获取MQ中间件的服务端构造函数.
		GetRemoting() (callable RemotingCallable)

		// SetRemoting
		// 设置MQ中间件的服务端构造函数.
		SetRemoting(callable RemotingCallable)
	}
)
