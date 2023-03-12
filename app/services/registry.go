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

package services

import (
	"github.com/fuyibing/db/v5"
	"github.com/fuyibing/gmd/v8/app/models"
	"xorm.io/xorm"
)

type (
	RegistryService struct {
		db.Service
	}
)

func NewRegistryService(ss ...*xorm.Session) *RegistryService {
	o := &RegistryService{}
	o.Use(ss...)
	o.UseConnection(Connection)
	return o
}

func (o *RegistryService) ListAll() (list []*models.Registry, err error) {
	list = make([]*models.Registry, 0)
	err = o.Slave().Find(&list)
	return
}
