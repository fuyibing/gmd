// author: wsfuyibing <websearch@163.com>
// date: 2023-02-17

// Package models
// Database model definition.
package models

import (
	"time"
)

const (
	DefaultParallels   = 1
	DefaultConcurrency = 10
	DefaultMaxRetry    = 5
)

const (
	ConnectionName = "db"
	GmtTimeLayout  = "2006-01-02 15:04:05"

	StatusDisabled = 0
	StatusEnabled  = 1

	StatusSucceed    = 1
	StatusFailed     = 2
	StatusWaiting    = 3
	StatusProcessing = 4
	StatusIgnored    = 9
)

type Timeline string

func (o Timeline) String() string {
	return string(o)
}

func (o Timeline) Time() time.Time {
	t, _ := time.Parse(GmtTimeLayout, string(o))
	return t
}

func NewTimeline() Timeline {
	return Timeline(time.Now().Format(GmtTimeLayout))
}
