// author: wsfuyibing <websearch@163.com>
// date: 2023-02-08

package md

import (
	"context"
	"fmt"
	"github.com/fuyibing/gmd/app"
	"github.com/fuyibing/gmd/app/md/base"
	"github.com/fuyibing/gmd/app/md/conf"
	"github.com/fuyibing/gmd/app/md/dispatchers"
	"github.com/fuyibing/log/v3"
	"github.com/fuyibing/log/v3/trace"
	"github.com/google/uuid"
	"strings"
	"sync/atomic"
	"time"
)

type (
	// ConsumerWorker
	// instance of consumer worker.
	ConsumerWorker interface {
		// Do
		// worker process.
		//
		// - Dispatch message
		// - Send notification if enabled.
		// - Store received message
		Do(t *base.Task, m *base.Message) (retry bool)

		// IsIdle
		// return idle status.
		//
		// Return false if any job is not completed, otherwise true returned.
		IsIdle() bool
	}

	worker struct {
		consuming, notifying, releasing int32
	}
)

// Do
// worker process.
func (o *worker) Do(t *base.Task, m *base.Message) (retry bool) {
	var ignored bool

	if m.TaskId == 0 {
		m.TaskId = t.Id
	}

	// Consume message
	// in sync coroutine.
	ignored, retry = o.DoConsume(t, m)

	// Send notification
	// in async coroutine if enabled and message not ignored.
	if !ignored && !retry {
		if m.GetError() != nil {
			if t.EnNotificationFailed() {
				o.DoNotify(m, conf.Config.Producer.NotificationTopic, conf.Config.Producer.NotificationTagFailed)
			}
		} else {
			if t.EnNotificationSucceed() {
				o.DoNotify(m, conf.Config.Producer.NotificationTopic, conf.Config.Producer.NotificationTagSucceed)
			}
		}
	}

	// Call
	// release process.
	o.DoRelease(m)
	return
}

// IsIdle
// return busy status.
func (o *worker) IsIdle() bool {
	return atomic.LoadInt32(&o.consuming) == 0 &&
		atomic.LoadInt32(&o.notifying) == 0 &&
		atomic.LoadInt32(&o.releasing) == 0
}

// /////////////////////////////////////////////////////////////
// Action methods.
// /////////////////////////////////////////////////////////////

// DoConsume
// consume method.
func (o *worker) DoConsume(t *base.Task, m *base.Message) (ignored, retry bool) {
	log.Infofc(m.GetContext(), "consumer worker: consume message, task-id=%d, try-count=%d, message-id=%s", t.Id, m.Dequeue, m.MessageId)

	var (
		c   = trace.Child(m.GetContext())
		err error
		raw string
		s   *base.Subscriber
	)

	// Called
	// when end.
	defer func() {
		// Set
		// ignored status.
		m.SetIgnored(ignored)

		// Execute
		// retry status.
		retry = !ignored && err != nil && m.Dequeue < t.MaxRetry
	}()

	// Subscriber selector
	// with task and message.
	if s, raw, err = o.getSubscriber(c, t, m); err != nil {
		return
	}

	// Condition filter
	// if enabled.
	if ignored, err = o.runCondition(c, s, raw); err != nil || ignored {
		return
	}

	// Dispatch process
	// in sync coroutine.
	err = o.runDispatcher(c, t, m, s, raw)
	return
}

// DoNotify
// send notification.
func (o *worker) DoNotify(m *base.Message, topic, tag string) {
	var (
		c context.Context
		p *base.Payload
		r *base.Registry
	)

	// Read registry
	// by topic name and tag.
	if r = base.Memory.GetRegistryByName(topic, tag); r == nil {
		log.Errorfc(m.GetContext(), "notify denied: topic and tag pair not registered")
		return
	}

	// Acquire payload
	// then assign fields.
	log.Infofc(m.GetContext(), "notify begin: topic=%s, tag=%s", topic, tag)
	c = trace.Child(m.GetContext())
	p = base.Pool.AcquirePayload().SetContext(c)
	p.Offset = 0
	p.Hash = strings.ToUpper(strings.ReplaceAll(uuid.NewString(), "-", ""))

	// Bind registry
	// to payload fields.
	p.RegistryId = r.Id
	p.TopicName = r.TopicName
	p.TopicTag = r.TopicTag
	p.FilterTag = r.FilterTag

	// Bind relation
	// to payload fields.
	p.MessageMessageId = m.MessageId
	p.MessageTaskId = m.TaskId

	// Publish directly
	// in async coroutine.
	go func(payload *base.Payload) {
		if err := Boot.Producer().PublishDirect(payload); err != nil {
			payload.SetError(err).Release()
		}
	}(p)
}

// DoRelease
// store message into database and release instance.
func (o *worker) DoRelease(m *base.Message) {
	go func() {
		atomic.AddInt32(&o.releasing, 1)
		defer atomic.AddInt32(&o.releasing, -1)
		m.Release()
	}()
}

// /////////////////////////////////////////////////////////////
// Dispatch methods.
// /////////////////////////////////////////////////////////////

func (o *worker) dispatchHttp(c context.Context, t *base.Task, m *base.Message, s *base.Subscriber, raw string) (body []byte, err error) {
	log.Infofc(c, "dispatcher call: type=http, method=%s, addr=%s, timeout=%d", s.Method, s.Addr, s.Timeout)

	// Acquire
	// http dispatcher and release when end.
	x := dispatchers.Pool.AcquireHttp()
	defer x.Release()

	// Set request
	// method and address.
	x.Request.SetRequestURI(s.Addr)
	x.Request.Header.SetMethod(s.Method)

	// Set user agent.
	x.Request.Header.SetUserAgent(app.Config.Software)

	// Set headers
	// with message properties.
	for k, v := range map[string]interface{}{
		"X-Gmd-Filter":       t.FilterTag,
		"X-Gmd-Message-Id":   m.MessageId,
		"X-Gmd-Message-Time": m.MessageTime,
		"X-Gmd-Topic":        t.TopicName,
		"X-Gmd-Tag":          t.TopicTag,
		"X-Gmd-Try":          m.Dequeue,
	} {
		x.Request.Header.Set(k, fmt.Sprintf("%v", v))
	}

	// Set
	// request body.
	if raw != "" {
		x.Request.SetBodyRaw([]byte(raw))
	}

	// Send request.
	body, err = x.Run(s.Timeout)
	return
}

func (o *worker) dispatchTcp(_ context.Context, _ *base.Task, _ *base.Message, _ *base.Subscriber, _ string) (body []byte, err error) {
	// todo : tcp dispatcher
	err = fmt.Errorf("tcp dispatcher not support")
	return
}

func (o *worker) dispatchRpc(_ context.Context, _ *base.Task, _ *base.Message, _ *base.Subscriber, _ string) (body []byte, err error) {
	// todo : rpc dispatcher
	err = fmt.Errorf("rpc dispatcher not support")
	return
}

func (o *worker) dispatchWebsocket(_ context.Context, _ *base.Task, _ *base.Message, _ *base.Subscriber, _ string) (body []byte, err error) {
	// todo : websocket dispatcher
	err = fmt.Errorf("websocket dispatcher not support")
	return
}

// /////////////////////////////////////////////////////////////
// Results validator.
// /////////////////////////////////////////////////////////////

func (o *worker) resultErrnoIsZero(body []byte) (code string, err error) {
	x := base.Result.Acquire()
	defer base.Result.Release(x)
	return x.Parse(body)
}

// /////////////////////////////////////////////////////////////
// Action methods.
// /////////////////////////////////////////////////////////////

func (o *worker) getSubscriber(c context.Context, t *base.Task, m *base.Message) (s *base.Subscriber, raw string, err error) {
	log.Infofc(c, "get subscriber")

	// Called
	// when end.
	defer func() {
		if err != nil {
			log.Errorfc(c, "get subscriber: %v", err)
		}
	}()

	// Normal subscription
	// which configured on task.handler column.
	if !t.IsNotification() {
		if t.HandlerSubscriber != nil {
			s = t.HandlerSubscriber
			raw = m.MessageBody
			return
		}

		err = fmt.Errorf("handler of task %d not defined", t.Id)
		return
	}

	// Notification subscription
	// which configured on task.failed or task.succeed column.
	var (
		n  = base.Pool.AcquireNotification()
		nt *base.Task
	)

	defer n.Release()

	// Return error
	// if message body not validated.
	if n.Parse(m.MessageBody) != nil || n.TaskId == 0 || n.MessageBody == "" {
		err = fmt.Errorf("invalid notification message body")
		return
	}

	// Return error
	// if source task disabled or deleted.
	if nt = base.Memory.GetTask(n.TaskId); nt == nil {
		err = fmt.Errorf("source task %d disabled or deleted", n.TaskId)
		return
	}

	// Subscription for dispatcher fail
	// notification.
	if t.IsNotificationFailed() {
		if nt.EnNotificationFailed() {
			s = nt.FailedSubscriber
			raw = n.MessageBody
			return
		}

		err = fmt.Errorf("failed handler of source task %d not defined", nt.Id)
		return
	}

	// Subscription for dispatcher succeed
	// notification.
	if t.IsNotificationSucceed() {
		if nt.EnNotificationSucceed() {
			s = nt.SucceedSubscriber
			raw = n.MessageBody
			return
		}

		err = fmt.Errorf("succeed handler of source task %d not defined", nt.Id)
		return
	}

	// Return error.
	err = fmt.Errorf("unknown notification type")
	return
}

func (o *worker) runCondition(c context.Context, s *base.Subscriber, raw string) (ignored bool, err error) {
	// Return
	// if not enabled.
	if s.Condition == nil {
		return
	}

	// Return
	// if error occurred on match called.
	if ignored, err = s.Condition.MatchJsonString(raw); err != nil {
		log.Errorfc(c, "condition match: expression=%s, error=%v", s.Condition.Expression(), err)
		return
	}

	// Completed
	// on match called.
	log.Infofc(c, "condition match: expression=%s, ignored=%v", s.Condition.Expression(), ignored)
	return
}

func (o *worker) runDispatcher(c context.Context, t *base.Task, m *base.Message, s *base.Subscriber, raw string) (err error) {
	var (
		body []byte
		code string
		ct   = time.Now()
	)

	// Called
	// when end.
	defer func() {
		dur := time.Now().Sub(ct).Seconds()

		// Logger result.
		if err != nil {
			if body == nil {
				body = []byte(err.Error())
			}
			log.Warnfc(c, "dispatcher failed: code=%s, reason=%v, result=%s", code, err, body)
		} else {
			log.Infofc(c, "dispatcher succeed: code=%s, result=%s", code, body)
		}

		// Set dispatcher result
		// to message.
		m.SetBody(body).SetDuration(dur).SetError(err)
	}()

	// Switch dispatcher
	// by protocol.
	switch s.Protocol {
	case base.SubscriberProtocolHttp:
		body, err = o.dispatchHttp(c, t, m, s, raw)
	case base.SubscriberProtocolTcp:
		body, err = o.dispatchTcp(c, t, m, s, raw)
	case base.SubscriberProtocolRpc:
		body, err = o.dispatchRpc(c, t, m, s, raw)
	case base.SubscriberProtocolWebsocket:
		body, err = o.dispatchWebsocket(c, t, m, s, raw)
	default:
		err = fmt.Errorf("unknown protocol on: %s", s.Addr)
	}

	// Return
	// if error occurred.
	if err != nil {
		return
	}

	// Validate
	// response body.
	switch s.ResponseType {
	case base.SubscriberResponseTypeErrnoIsZero:
		code, err = o.resultErrnoIsZero(body)
	}

	// Set dispatcher status as succeed
	// if response code ignored.
	if err != nil && code != "" && s.IgnoreCodes != nil {
		for _, ic := range s.IgnoreCodes {
			if ic == code {
				err = nil
				log.Infofc(c, "response ignore: code=%s, config=%v", code, s.IgnoreCodes)
				break
			}
		}
	}

	return
}

// /////////////////////////////////////////////////////////////
// Constructor methods.
// /////////////////////////////////////////////////////////////

func (o *worker) init() *worker {
	return o
}
