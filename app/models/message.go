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
// date: 2023-03-07

package models

type (
	// Message
	// 队列消息.
	Message struct {
		Id       int64   `xorm:"id pk autoincr"`
		Status   int     `xorm:"status"`
		Duration float64 `xorm:"duration"`

		TaskId  int `xorm:"task_id"`
		Dequeue int `xorm:"dequeue"`
		Retry   int `xorm:"retry"`

		PayloadMessageId string `xorm:"payload_message_id"`
		MessageTime      int64  `xorm:"message_time"`
		MessageId        string `xorm:"message_id"`
		MessageBody      string `xorm:"message_body"`
		ResponseBody     string `xorm:"response_body"`

		GmtCreated Datetime `xorm:"gmt_created"`
		GmtUpdated Datetime `xorm:"gmt_updated"`
	}
)
