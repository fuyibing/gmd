// author: wsfuyibing <websearch@163.com>
// date: 2023-02-09

package base

import (
	"encoding/json"
	"fmt"
)

type (
	// Notification
	// struct for notification produce.
	Notification struct {
		MessageBody string `json:"__gmd__message_body_"`
		TaskId      int    `json:"__gmd__task_id_"`
	}
)

// Parse
// unmarshal json string into instance.
func (o *Notification) Parse(s string) error {
	if err := json.Unmarshal([]byte(s), o); err != nil {
		return fmt.Errorf("invalid notification message body")
	}
	return nil
}

// Release
// instance to pool.
func (o *Notification) Release() {
	Pool.ReleaseNotification(o)
}

// /////////////////////////////////////////////////////////////
// Access methods.
// /////////////////////////////////////////////////////////////

func (o *Notification) after()              { o.MessageBody = ""; o.TaskId = 0 }
func (o *Notification) before()             {}
func (o *Notification) init() *Notification { return o }
