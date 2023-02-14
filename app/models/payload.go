// author: wsfuyibing <websearch@163.com>
// date: 2023-01-17

package models

type (
	// Payload
	//
	// producer message struct.
	Payload struct {
		Id       int64   `xorm:"id pk autoincr"`
		Status   int     `xorm:"status"`
		Duration float64 `xorm:"duration"`
		Retry    int     `xorm:"retry"`

		MessageTaskId    int    `xorm:"message_task_id"`
		MessageMessageId string `xorm:"message_message_id"`

		Hash         string `xorm:"hash"`
		Offset       int    `xorm:"offset"`
		RegistryId   int    `xorm:"registry_id"`
		MessageId    string `xorm:"message_id"`
		MessageBody  string `xorm:"message_body"`
		ResponseBody string `xorm:"response_body"`

		GmtCreated Timeline `xorm:"gmt_created"`
		GmtUpdated Timeline `xorm:"gmt_updated"`
	}
)
