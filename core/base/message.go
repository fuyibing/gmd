// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package base

// Message
// received from mq middleware server.
type Message struct {
}

func (o *Message) Release() { Pool.ReleaseMessage(o) }

// /////////////////////////////////////////////////////////////
// Access methods.
// /////////////////////////////////////////////////////////////

func (o *Message) after() {
}

func (o *Message) before() {
}

func (o *Message) init() *Message {
	return o
}

func (o *Message) save() {

}
