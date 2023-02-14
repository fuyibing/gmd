// author: wsfuyibing <websearch@163.com>
// date: 2023-01-17

package models

type (
	// Message
	//
	// received message struct from queue.
	Message struct {
		Id       int64   `xorm:"id pk autoincr"`
		Status   int     `xorm:"status"`
		Duration float64 `xorm:"duration"`
		Retry    int     `xorm:"retry"`

		TaskId           int    `xorm:"task_id"`
		PayloadMessageId string `xorm:"payload_message_id"`

		MessageDequeue int    `xorm:"message_dequeue"`
		MessageTime    int64  `xorm:"message_time"`
		MessageId      string `xorm:"message_id"`
		MessageBody    string `xorm:"message_body"`
		ResponseBody   string `xorm:"response_body"`

		GmtCreated Timeline `xorm:"gmt_created"`
		GmtUpdated Timeline `xorm:"gmt_updated"`
	}
)
