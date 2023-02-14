// author: wsfuyibing <websearch@163.com>
// date: 2023-02-09

package base

import (
	"sync"
)

var (
	// Pool
	// instance of pool manager.
	Pool PoolManager
)

type (
	// PoolManager
	// interface of pool manager.
	PoolManager interface {
		// AcquireMessage
		// acquire message instance from pool.
		AcquireMessage() *Message

		// AcquireNotification
		// acquire notification instance from pool.
		AcquireNotification() *Notification

		// AcquirePayload
		// acquire payload instance from pool.
		AcquirePayload() *Payload

		// ReleaseMessage
		// release message instance to pool.
		ReleaseMessage(x *Message)

		// ReleaseNotification
		// release notification instance to pool.
		ReleaseNotification(x *Notification)

		// ReleasePayload
		// release payload instance to pool.
		ReleasePayload(x *Payload)
	}

	pool struct {
		messages, notifications, payloads *sync.Pool
	}
)

func (o *pool) AcquireMessage() *Message {
	x := o.messages.Get().(*Message)
	x.before()
	return x
}

func (o *pool) AcquireNotification() *Notification {
	x := o.notifications.Get().(*Notification)
	x.before()
	return x
}

func (o *pool) AcquirePayload() *Payload {
	x := o.payloads.Get().(*Payload)
	x.before()
	return x
}

func (o *pool) ReleaseMessage(x *Message) {
	x.after()
	o.messages.Put(x)
}

func (o *pool) ReleaseNotification(x *Notification) {
	x.after()
	o.notifications.Put(x)
}

func (o *pool) ReleasePayload(x *Payload) {
	x.after()
	o.payloads.Put(x)
}

// /////////////////////////////////////////////////////////////
// Access methods.
// /////////////////////////////////////////////////////////////

func (o *pool) init() *pool {
	o.messages = &sync.Pool{New: func() interface{} { return (&Message{}).init() }}
	o.notifications = &sync.Pool{New: func() interface{} { return (&Notification{}).init() }}
	o.payloads = &sync.Pool{New: func() interface{} { return (&Payload{}).init() }}
	return o
}
