// author: wsfuyibing <websearch@163.com>
// date: 2023-02-03

// Package aliyunmns
// Message queue adapter on AliyunMNS.
package aliyunmns

import (
	"sync"
)

func init() {
	new(sync.Once).Do(func() {
		Agent = (&agent{}).init()
		TopicMessagePool = &sync.Pool{New: func() interface{} { return &TopicMessage{} }}
	})
}
