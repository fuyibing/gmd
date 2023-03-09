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
// date: 2023-03-06

package managers

import (
	"fmt"
	"github.com/fuyibing/gmd/v8/md/base"
	"sync"
)

const (
	BucketBatch       = 100
	BucketCapacity    = 30000
	BucketConcurrency = 50
	BucketFrequency   = 300
)

var (
	ErrBucketIsFull = fmt.Errorf("bucket is full")
)

type (
	// Bucket
	// 数据桶.
	Bucket interface {
		// Add
		// 添加消息.
		Add(v *base.Payload) (total int, err error)

		// Count
		// 桶消息数.
		Count() int

		// IsEmpty
		// 是否为空.
		IsEmpty() bool

		// Pop
		// 取出一条.
		Pop() (v *base.Payload, exists bool)

		// Popn
		// 取出N条.
		Popn(limit int) (vs []*base.Payload, total, count int)
	}

	bucket struct {
		sync.RWMutex

		cached   []*base.Payload
		capacity int
	}
)

func NewBucket(capacity int) Bucket { return (&bucket{capacity: capacity}).init() }

// /////////////////////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////////////////////

func (o *bucket) Add(v *base.Payload) (total int, err error) {
	o.Lock()
	defer o.Unlock()

	if total = len(o.cached) + 1; o.capacity > 0 && total > o.capacity {
		err = ErrBucketIsFull
		return
	}

	o.cached = append(o.cached, v)
	return
}

func (o *bucket) Count() int {
	o.RLock()
	defer o.RUnlock()
	return len(o.cached)
}

func (o *bucket) IsEmpty() bool {
	return o.Count() == 0
}

func (o *bucket) Pop() (v *base.Payload, exists bool) {
	if vs, _, count := o.Popn(1); count > 0 {
		return vs[0], true
	}
	return nil, false
}

func (o *bucket) Popn(limit int) (vs []*base.Payload, total, count int) {
	o.Lock()
	defer o.Unlock()

	if total = len(o.cached); total == 0 {
		return
	}

	// 全部取出.
	if limit >= total {
		count = total
		vs = o.cached[:]
		o.reset()
		return
	}

	// 取出片段.
	count = limit
	vs = o.cached[0:count]

	// 重置内存.
	o.cached = o.cached[count:]
	return
}

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *bucket) init() *bucket {
	o.reset()
	return o
}

func (o *bucket) reset() {
	o.cached = make([]*base.Payload, 0)
}
