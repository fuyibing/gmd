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
// date: 2023-03-09

package md

import (
	"github.com/fuyibing/gmd/v8/md/adapters/rocketmq"
	"github.com/fuyibing/gmd/v8/md/base"
	"github.com/fuyibing/gmd/v8/md/conditions/el_condition"
	"github.com/fuyibing/gmd/v8/md/dispatchers/http_json"
	"github.com/fuyibing/gmd/v8/md/results/json_errno"
	"strings"
)

type (
	builtinConsumer string
	builtinProducer string
	builtinRemoter  string
)

func (o builtinConsumer) New() base.ConsumerConstructor {
	switch strings.ToLower(string(o)) {
	case "rocketmq":
		return rocketmq.NewConsumer
	}
	return nil
}

func (o builtinProducer) New() base.ProducerConstructor {
	switch strings.ToLower(string(o)) {
	case "rocketmq":
		return rocketmq.NewProducer
	}
	return nil
}

func (o builtinRemoter) New() base.RemoterConstructor {
	switch strings.ToLower(string(o)) {
	case "rocketmq":
		return rocketmq.NewRemoter
	}
	return nil
}

var (
	builtinConditions = map[string]base.ConditionConstructor{
		"el": el_condition.New,
	}

	builtinDispatchers = map[string]base.DispatcherConstructor{
		"http_json": http_json.New,
	}

	builtinResults = map[string]base.ResultConstructor{
		"json_errno": json_errno.New,
	}
)
