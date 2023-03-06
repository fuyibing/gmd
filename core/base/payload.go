// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package base

import (
	"context"
	"time"
)

// Payload
// 入参结构.
//
// 发布到MQ的消息结构 或 本地存储的重试消息.
type Payload struct {
	// 上下文
	// 应用于 Span 跟踪.
	ctx context.Context

	// 发布毫时.
	duration time.Duration

	// 错误原因.
	err error

	// 消息ID.
	messageId string

	Hash   string
	Offset int

	RegistryId int
	TopicName  string
	TopicTag   string
	FilterTag  string

	MessageBody string
}

func (o *Payload) GetContext() context.Context { return o.ctx }

func (o *Payload) Release() { Pool.ReleasePayload(o) }

func (o *Payload) SetContext(ctx context.Context) *Payload { o.ctx = ctx; return o }
func (o *Payload) SetDuration(dur time.Duration) *Payload  { o.duration = dur; return o }
func (o *Payload) SetError(err error) *Payload             { o.err = err; return o }
func (o *Payload) SetMessageId(s string) *Payload          { o.messageId = s; return o }

// /////////////////////////////////////////////////////////////////////////////
// Access and constructor
// /////////////////////////////////////////////////////////////////////////////

func (o *Payload) after() {}

func (o *Payload) before() {}

func (o *Payload) init() *Payload {
	return o
}

func (o *Payload) save() {
}
