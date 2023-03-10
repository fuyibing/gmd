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

package base

type (
	// Subscriber
	// 订阅者接口.
	Subscriber interface {
		GetCondition() ConditionExecutor
		GetDispatcher() DispatcherExecutor
		GetResult() ResultExecutor
		HasCondition() bool
		HasDispatcher() bool
		HasResult() bool
		SetCondition(kind, filter string) Subscriber
		SetDispatcher(kind, addr, method string, timeout int) Subscriber
		SetResult(kind, ignores string) Subscriber
	}

	// 订阅者.
	subscriber struct {
		ce ConditionExecutor
		ch bool
		de DispatcherExecutor
		dh bool
		re ResultExecutor
		rh bool
	}
)

// NewSubscriber
// 创建订阅者.
func NewSubscriber() Subscriber { return (&subscriber{}).init() }

// +---------------------------------------------------------------------------+
// + Interface methods                                                         |
// +---------------------------------------------------------------------------+

func (o *subscriber) GetCondition() ConditionExecutor   { return o.ce }
func (o *subscriber) GetDispatcher() DispatcherExecutor { return o.de }
func (o *subscriber) GetResult() ResultExecutor         { return o.re }

func (o *subscriber) HasCondition() bool  { return o.ch }
func (o *subscriber) HasDispatcher() bool { return o.dh }
func (o *subscriber) HasResult() bool     { return o.rh }

func (o *subscriber) SetCondition(kind, filter string) Subscriber {
	if caller, exists := Container.GetCondition(kind); exists {
		o.ce = caller(filter)
		o.ch = true
	}
	return o
}

func (o *subscriber) SetDispatcher(kind, addr, method string, timeout int) Subscriber {
	if caller, exists := Container.GetDispatcher(kind); exists {
		o.de = caller(addr, method, timeout)
		o.dh = true
	}
	return o
}

func (o *subscriber) SetResult(kind, ignoreCodes string) Subscriber {
	if caller, exists := Container.GetResult(kind); exists {
		o.re = caller(ignoreCodes)
		o.rh = true
	}
	return o
}

// +---------------------------------------------------------------------------+
// + Constructor and access methods                                            |
// +---------------------------------------------------------------------------+

func (o *subscriber) init() *subscriber { return o }
