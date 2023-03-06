// author: wsfuyibing <websearch@163.com>
// date: 2023-02-20

package base

type (
	// ContainerManager
	// 容器管理器.
	ContainerManager interface {
		// GetConsumer
		// 消费者构造.
		GetConsumer() (callable ConsumerCallable)

		// SetConsumer
		// 设置消费者构造.
		SetConsumer(callable ConsumerCallable)

		// GetProducer
		// 生产者构造.
		GetProducer() (callable ProducerCallable)

		// SetProducer
		// 设置生产者构造.
		SetProducer(callable ProducerCallable)

		// GetRemoting
		// 服务端构造.
		GetRemoting() (callable RemotingCallable)

		// SetRemoting
		// 设置服务端构造.
		SetRemoting(callable RemotingCallable)
	}
)
