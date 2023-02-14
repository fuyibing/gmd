// author: wsfuyibing <websearch@163.com>
// date: 2023-02-08

package md

import (
	"fmt"
	"github.com/fuyibing/gmd/app/md/base"
	"sync"
)

type (
	// ProducerBucket
	// interface of producer bucket.
	ProducerBucket interface {
		// IsEmpty
		// return bucket is empty or not.
		//
		// Return true if no payload in bucket, otherwise false
		// returned.
		IsEmpty() (yes bool)

		// IsFull
		// return bucket is full or not.
		//
		// Return true if payload count is greater or equal to size,
		// otherwise false returned.
		IsFull() (yes bool)

		// Length
		// return payloads count in bucket.
		Length() (count int)

		// Pop
		// get one payload from left cached.
		Pop() (payload *base.Payload)

		// Popn
		// get specified count payloads from left cached.
		Popn(n int) (list []*base.Payload, count int)

		// Push
		// add payload to right cached.
		Push(ps ...*base.Payload) error
	}

	bucket struct {
		cached []*base.Payload
		mu     *sync.Mutex
		size   int
	}
)

// IsEmpty
// return bucket is empty or not.
func (o *bucket) IsEmpty() bool {
	return o.Length() == 0
}

// IsFull
// return bucket is full or not.
func (o *bucket) IsFull() bool {
	return o.Length() >= o.size
}

// Length
// return payloads count in bucket.
func (o *bucket) Length() int {
	o.mu.Lock()
	defer o.mu.Unlock()
	return len(o.cached)
}

// Pop
// get one payload from left.
func (o *bucket) Pop() *base.Payload {
	if s, n := o.Popn(1); n == 1 {
		return s[0]
	}
	return nil
}

// Popn
// get specified count payloads from left.
func (o *bucket) Popn(n int) (list []*base.Payload, count int) {
	o.mu.Lock()
	defer o.mu.Unlock()

	var total int

	if total = len(o.cached); total == 0 {
		return
	}

	if n >= total {
		count = total
		list = o.cached[0:]
		o.cached = []*base.Payload{}
		return
	}

	count = n
	list = o.cached[0:n]
	o.cached = o.cached[n:]
	return
}

// Push
// add payload to right.
func (o *bucket) Push(ps ...*base.Payload) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	n := len(ps)
	if n == 0 {
		return nil
	}

	if (n + len(o.cached)) > o.size {
		return fmt.Errorf("bucket is full")
	}

	o.cached = append(o.cached, ps...)
	return nil
}

// /////////////////////////////////////////////////////////////
// Access methods.
// /////////////////////////////////////////////////////////////

func (o *bucket) init() *bucket {
	o.cached = make([]*base.Payload, 0)
	o.mu = &sync.Mutex{}
	return o
}
