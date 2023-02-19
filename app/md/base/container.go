// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package base

// ContainerManager
// registered manager operations.
type ContainerManager interface {
	// GetConsumer
	// return consumer manager constructor.
	GetConsumer() (manager ConsumerCallable, exists bool)

	// GetProducer
	// return producer manager constructor.
	GetProducer() (manager ProducerCallable, exists bool)

	// GetRemoting
	// return remoting manager constructor.
	GetRemoting() (manager RemotingCallable, exists bool)

	// SetAdapter
	// configure consumer, producer, remoting with adapter.
	SetAdapter(adapter Adapter)

	// SetConsumer
	// configure consumer manager constructor, singleton instance.
	SetConsumer(manager ConsumerCallable)

	// SetProducer
	// configure producer manager constructor, singleton instance.
	SetProducer(manager ProducerCallable)

	// SetRemoting
	// configure remoting manager constructor, singleton instance.
	SetRemoting(manager RemotingCallable)
}
