// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// author: wsfuyibing <websearch@163.com>
// date: 2023-03-08

package base

import (
	"sync"
)

var (
	// Pool
	// 池操作.
	Pool PoolOperation
)

type (
	// PoolOperation
	// 池操作接口.
	PoolOperation interface {
		AcquireMessage() (v *Message)
		AcquireNotification() (v *Notification)
		AcquirePayload() (v *Payload)
		ReleaseMessage(v *Message)
		ReleaseNotification(v *Notification)
		ReleasePayload(v *Payload)
	}

	pool struct {
		messages,
		notifications,
		payloads sync.Pool
	}
)

func (o *pool) AcquireMessage() (v *Message) {
	if g := o.messages.Get(); g != nil {
		v = g.(*Message)
		v.before()
		return
	}

	v = (&Message{}).init()
	v.before()
	return
}

func (o *pool) ReleaseMessage(v *Message) {
	v.after()
	o.messages.Put(v)
}

func (o *pool) AcquireNotification() (v *Notification) {
	if g := o.notifications.Get(); g != nil {
		v = g.(*Notification)
		v.before()
		return
	}

	v = (&Notification{}).init()
	v.before()
	return v
}

func (o *pool) ReleaseNotification(v *Notification) {
	v.after()
	o.notifications.Put(v)
}

func (o *pool) AcquirePayload() (v *Payload) {
	if g := o.payloads.Get(); g != nil {
		v = g.(*Payload)
		v.before()
		return v
	}

	v = (&Payload{}).init()
	v.before()
	return v
}

func (o *pool) ReleasePayload(v *Payload) {
	v.after()
	o.payloads.Put(v)
}

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *pool) init() *pool {
	return o
}
