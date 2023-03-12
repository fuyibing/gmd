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

package controllers

import (
	"github.com/fuyibing/gmd/v8/app/logics"
	"github.com/fuyibing/gmd/v8/app/logics/index"
	"github.com/kataras/iris/v12"
)

// Controller
// Default.
type Controller struct{}

// Get
// Home page.
//
// @Response(app/logics/index.GetHomeResponse)
func (o *Controller) Get(i iris.Context) interface{} {
	return logics.New(i, index.NewGetHome)
}

// GetPing
// Health check.
//
// @Response(app/logics/index.GetPingResponse)
func (o *Controller) GetPing(i iris.Context) interface{} {
	return logics.New(i, index.NewGetPing)
}
