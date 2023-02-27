// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package base

// Message
// 从MQ队列收到的消息.
type Message struct {
}

func (o *Message) Release() { Pool.ReleaseMessage(o) }

func (o *Message) after() {}

func (o *Message) before() {}

func (o *Message) init() *Message {
	return o
}

func (o *Message) save() {
}
