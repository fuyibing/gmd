// author: wsfuyibing <websearch@163.com>
// date: 2023-02-20

package base

type (
	// ContainerManager
	// registered manager operations.
	ContainerManager interface {
		// GetConsumer
		// return consumer manager constructor.
		GetConsumer() (callable ConsumerCallable)

		// SetConsumer
		// configure consumer manager constructor, singleton instance.
		SetConsumer(callable ConsumerCallable)

		// GetProducer
		// return producer manager constructor.
		GetProducer() (callable ProducerCallable)

		// SetProducer
		// configure producer manager constructor, singleton instance.
		SetProducer(callable ProducerCallable)

		// GetRemoting
		// return remoting manager constructor.
		GetRemoting() (callable RemotingCallable)

		// SetRemoting
		// configure remoting manager constructor, singleton instance.
		SetRemoting(callable RemotingCallable)
	}
)
