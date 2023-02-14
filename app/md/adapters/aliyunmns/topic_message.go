// author: wsfuyibing <websearch@163.com>
// date: 2023-02-07

package aliyunmns

import (
	"sync"
)

var (
	TopicMessagePool *sync.Pool
)

type (
	TopicMessage struct {
		Message   string `json:"Message"`
		MessageId string `json:"MessageId"`
		TopicName string `json:"TopicName"`
	}
)

func (o *TopicMessage) Release() {
	o.Message = ""
	o.MessageId = ""
	o.TopicName = ""
	TopicMessagePool.Put(o)
}
