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
	Task struct {
		Id         int `xorm:"id pk autoincr"`
		RegistryId int `xorm:"registry_id"`
		Status     int `xorm:"status"`

		// +-------------------------------------------------------------------+
		// + Subscription attributes                                           |
		// +-------------------------------------------------------------------+

		Broadcasting int    `yaml:"broadcasting"`
		Parallels    int    `yaml:"parallels"`
		Concurrency  int32  `yaml:"concurrency"`
		MaxRetry     int    `yaml:"max_retry"`
		DelaySeconds int    `yaml:"delay_seconds"`
		Title        string `xorm:"title"`
		Remark       string `xorm:"remark"`

		// +-------------------------------------------------------------------+
		// + Normal subscriber                                                 |
		// +-------------------------------------------------------------------+

		HandlerConditionFilter   string `xorm:"handler_condition_filter"`
		HandlerConditionKind     string `xorm:"handler_condition_kind"`
		HandlerDispatcherAddr    string `xorm:"handler_dispatcher_addr"`
		HandlerDispatcherKind    string `xorm:"handler_dispatcher_kind"`
		HandlerDispatcherMethod  string `xorm:"handler_dispatcher_method"`
		HandlerDispatcherTimeout int    `xorm:"handler_dispatcher_timeout"`
		HandlerResultIgnoreCodes string `xorm:"handler_result_ignore_codes"`
		HandlerResultKind        string `xorm:"handler_result_kind"`

		// +-------------------------------------------------------------------+
		// + Dispatch failed notification                                      |
		// +-------------------------------------------------------------------+

		FailedConditionFilter   string `xorm:"failed_condition_filter"`
		FailedConditionKind     string `xorm:"failed_condition_kind"`
		FailedDispatcherAddr    string `xorm:"failed_dispatcher_addr"`
		FailedDispatcherKind    string `xorm:"failed_dispatcher_kind"`
		FailedDispatcherMethod  string `xorm:"failed_dispatcher_method"`
		FailedDispatcherTimeout int    `xorm:"failed_dispatcher_timeout"`
		FailedResultIgnoreCodes string `xorm:"failed_result_ignore_codes"`
		FailedResultKind        string `xorm:"failed_result_kind"`

		// +-------------------------------------------------------------------+
		// + Dispatch succeed notification                                     |
		// +-------------------------------------------------------------------+

		SucceedConditionFilter   string `xorm:"succeed_condition_filter"`
		SucceedConditionKind     string `xorm:"succeed_condition_kind"`
		SucceedDispatcherAddr    string `xorm:"succeed_dispatcher_addr"`
		SucceedDispatcherKind    string `xorm:"succeed_dispatcher_kind"`
		SucceedDispatcherMethod  string `xorm:"succeed_dispatcher_method"`
		SucceedDispatcherTimeout int    `xorm:"succeed_dispatcher_timeout"`
		SucceedResultIgnoreCodes string `xorm:"succeed_result_ignore_codes"`
		SucceedResultKind        string `xorm:"succeed_result_kind"`

		GmtCreated Datetime `xorm:"gmt_created"`
		GmtUpdated Datetime `xorm:"gmt_updated"`
	}
)
