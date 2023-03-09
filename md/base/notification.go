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
	"encoding/json"
	"fmt"
)

type (
	// Notification
	// 通知消息.
	Notification struct {
		MessageBody string `json:"notification_message_body"`
		MessageId   string `json:"notification_message_id"`
		TaskId      int    `json:"notification_task_id"`
	}
)

// Decoder
// 解码通知消息.
func (o *Notification) Decoder(message *Message) (*Task, error) {
	// 解码消息.
	if err := json.Unmarshal([]byte(message.MessageBody), o); err != nil || o.TaskId == 0 {
		err = fmt.Errorf("illegal message format")
		return nil, err
	}

	// 加载任务.
	if source, exists := Memory.GetTask(o.TaskId); exists {
		return source, nil
	}

	// 无效任务.
	return nil, fmt.Errorf("subscription task disabled or deleted")
}

func (o *Notification) Release() { Pool.ReleaseNotification(o) }

func (o *Notification) String() (str string) {
	if buf, err := json.Marshal(o); err == nil {
		str = string(buf)
	}
	return
}

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *Notification) after() {}

func (o *Notification) before() {}

func (o *Notification) init() *Notification {
	return o
}
