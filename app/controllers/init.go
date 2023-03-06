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
// date: 2023-02-27

package controllers

import (
	"github.com/fuyibing/gmd/v8/app/controllers/task"
	"github.com/fuyibing/gmd/v8/app/controllers/topic"
	"sync"
)

var (
	Containers map[string]interface{}
)

func init() {
	new(sync.Once).Do(func() {
		Containers = map[string]interface{}{
			"/":      &Controller{},
			"/task":  &task.Controller{},
			"/topic": &topic.Controller{},
		}
	})
}
