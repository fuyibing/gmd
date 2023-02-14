// author: wsfuyibing <websearch@163.com>
// date: 2023-02-09

package base

import (
	"encoding/json"
	"fmt"
	"sync"
)

var (
	// Result
	// instance of result.
	Result ResultInterface
)

type (
	// ResultInterface
	// interface of result.
	ResultInterface interface {
		// Acquire
		// acquire validator instance from registered pool.
		Acquire() ResultValidator

		// Register
		// register validator instance override default.
		Register(v func() ResultValidator)

		// Release
		// release validator instance to pool.
		Release(x ResultValidator)
	}

	// ResultValidator
	// instance of result validator.
	ResultValidator interface {
		// After
		// called when release to pool.
		After()

		// Before
		// called when acquired from pool.
		Before()

		// Parse
		// parse dispatched result and return status.
		Parse(body []byte) (code string, err error)
	}

	result struct {
		p *sync.Pool
		v func() ResultValidator
	}

	validator struct {
		Errno interface{} `json:"errno"`
		Error interface{} `json:"error"`
	}
)

func (o *result) Acquire() ResultValidator {
	// Return validator
	// if acquired succeed.
	if x := o.p.Get(); x != nil {
		if v, ok := x.(ResultValidator); ok {
			v.Before()
			return v
		}
	}

	// Create validator
	// and return if not acquired.
	v := o.v()
	v.Before()
	return v
}

func (o *result) Register(v func() ResultValidator) { o.v = v }
func (o *result) Release(x ResultValidator)         { x.After(); o.p.Put(x) }

func (o *validator) After()  { o.Errno = nil; o.Errno = nil }
func (o *validator) Before() {}
func (o *validator) Parse(body []byte) (code string, err error) {
	// Parse body
	// on interface.
	if err = json.Unmarshal(body, o); err != nil {
		return
	}

	// Return
	// if code is zero string.
	if code = fmt.Sprintf("%v", o.Errno); code == "0" {
		return
	}

	// Reset
	// error reason.
	err = fmt.Errorf("%v", o.Error)
	return
}

// /////////////////////////////////////////////////////////////
// Access methods.
// /////////////////////////////////////////////////////////////

func (o *result) init() *result {
	o.p = &sync.Pool{}
	o.Register(func() ResultValidator { return &validator{} })
	return o
}
