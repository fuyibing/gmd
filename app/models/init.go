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

import (
	"time"
)

const (
	DefaultDatetimeLayout  = "2006-01-02 15:04:05"
	DefaultTaskParallels   = 1
	DefaultTaskConcurrency = 10
	DefaultTaskMaxRetry    = 3
)

const (
	StatusDisabled = 0
	StatusEnabled  = 1

	StatusSucceed    = 1
	StatusFailed     = 2
	StatusWaiting    = 3
	StatusProcessing = 4
)

var (
	nilDatetime = time.Unix(0, 0)
)

type (
	Datetime string
)

func Now() Datetime {
	return Datetime(time.Now().Format(DefaultDatetimeLayout))
}

func (o Datetime) Time() time.Time {
	if t, te := time.Parse(DefaultDatetimeLayout, string(o)); te == nil {
		return t
	}
	return nilDatetime
}
