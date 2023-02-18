// author: wsfuyibing <websearch@163.com>
// date: 2023-02-03

package aliyunmns

import (
	"context"
	"fmt"
	mns "github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/fuyibing/gmd/app/md/base"
	"github.com/fuyibing/gmd/app/md/conf"
	"github.com/fuyibing/log/v8"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	Agent AgentManager

	RegexMessageNotExist      = regexp.MustCompile(`code\s*:\s*MessageNotExist`)
	RegexQueueNotExist        = regexp.MustCompile(`code\s*:\s*QueueNotExist`)
	RegexSubscriptionNotExist = regexp.MustCompile(`code\s*:\s*SubscriptionNotExist`)
	RegexTopicNotExist        = regexp.MustCompile(`code\s*:\s*TopicNotExist`)
)

const (
	DefaultLoggingEnable          = true
	DefaultMaxMessageSize         = 65536
	DefaultMessageRetentionPeriod = 604800
	DefaultNotifyContentFormat    = "JSON"
	DefaultNotifyStrategy         = mns.BACKOFF_RETRY
	DefaultPollingWaitSeconds     = 30
	DefaultSlices                 = 0
	DefaultVisibilityTimeout      = 30
)

type (
	// AgentManager
	// interface of aliyunmns agent.
	AgentManager interface {
		// Build
		// create remote relation if not exists.
		Build(ctx context.Context, task *base.Task) error

		// Destroy
		// delete remote relation if exists.
		Destroy(ctx context.Context, task *base.Task) error

		// GenQueueName
		// generate and return queue name.
		//
		//   Agent.GenQueueName(1) // return "X-Q1"
		GenQueueName(id int) string

		// GenSubscriptionName
		// generate and return subscription name.
		//
		//   Agent.GenSubscriptionName(1) // return "X-S1"
		GenSubscriptionName(id int) string

		// GenTopicName
		// generate and return topic name.
		//
		//   Agent.GenTopicName("Topic") // return "X-Topic"
		GenTopicName(name string) string

		// GetQueueClient
		// return queue client instance.
		GetQueueClient(id int) mns.AliMNSQueue

		// GetQueueManager
		// return queue manager interface.
		GetQueueManager() mns.AliQueueManager

		// GetTopicClient
		// return topic client instance.
		GetTopicClient(name string) mns.AliMNSTopic

		// GetTopicManager
		// return topic manager interface.
		GetTopicManager() mns.AliTopicManager
	}

	agent struct {
		mu           *sync.RWMutex
		queueClients map[int]mns.AliMNSQueue
		queueManager mns.AliQueueManager
		topicClients map[string]mns.AliMNSTopic
		topicManager mns.AliTopicManager
	}
)

// /////////////////////////////////////////////////////////////
// Interface methods.
// /////////////////////////////////////////////////////////////

func (o *agent) Build(ctx context.Context, task *base.Task) (err error) {
	var (
		qm = o.GetQueueManager()
		tm = o.GetTopicManager()
	)

	// Range callables
	// to build element.
	for _, f := range []func() error{
		func() error { return o.buildQueue(ctx, qm, task) },
		func() error { return o.buildTopic(ctx, tm, task) },
		func() error { return o.buildSubscribe(ctx, task) },
	} {
		if err = f(); err != nil {
			return
		}
	}

	return
}

func (o *agent) Destroy(ctx context.Context, task *base.Task) (err error) {
	var (
		tm = o.GetTopicManager()
		qm = o.GetQueueManager()
	)

	// Range callables
	// to build element.
	for _, f := range []func() error{
		func() error { return o.destroySubscription(ctx, task) },
		func() error { return o.destroyQueue(ctx, qm, task) },
		func() error { return o.destroyTopic(ctx, tm, task) },
	} {
		if err = f(); err != nil {
			return
		}
	}

	return
}

func (o *agent) GenQueueName(id int) string {
	return fmt.Sprintf("%sQ%d", conf.Config.Account.Aliyunmns.Prefix, id)
}

func (o *agent) GenSubscriptionName(id int) string {
	return fmt.Sprintf("%sS%d", conf.Config.Account.Aliyunmns.Prefix, id)
}

func (o *agent) GenTopicName(name string) string {
	return fmt.Sprintf("%s%s", conf.Config.Account.Aliyunmns.Prefix, name)
}

func (o *agent) GetQueueClient(id int) mns.AliMNSQueue {
	o.mu.Lock()
	defer o.mu.Unlock()

	if c, ok := o.queueClients[id]; ok {
		return c
	}

	c := mns.NewMNSQueue(o.GenQueueName(id), o.client())
	o.queueClients[id] = c
	return c
}

func (o *agent) GetQueueManager() mns.AliQueueManager {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.queueManager == nil {
		o.queueManager = mns.NewMNSQueueManager(o.client())
	}
	return o.queueManager
}

func (o *agent) GetTopicClient(name string) mns.AliMNSTopic {
	o.mu.Lock()
	defer o.mu.Unlock()

	if c, ok := o.topicClients[name]; ok {
		return c
	}

	c := mns.NewMNSTopic(o.GenTopicName(name), o.client())
	o.topicClients[name] = c
	return c
}

func (o *agent) GetTopicManager() mns.AliTopicManager {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.topicManager == nil {
		o.topicManager = mns.NewMNSTopicManager(o.client())
	}
	return o.topicManager
}

// /////////////////////////////////////////////////////////////
// Synchronize methods.
// /////////////////////////////////////////////////////////////

func (o *agent) buildQueue(ctx context.Context, m mns.AliQueueManager, task *base.Task) (err error) {
	var (
		attr mns.QueueAttribute
		name = o.GenQueueName(task.Id)
	)

	// Read queue attribute
	// from aliyunmns server.
	if attr, err = m.GetQueueAttributes(name); err == nil {
		// Send update request
		// if properties not matched on local task. Error
		// returned if request failed.
		if attr.DelaySeconds != int32(task.DelaySeconds) ||
			attr.MaxMessageSize != DefaultMaxMessageSize ||
			attr.MessageRetentionPeriod != DefaultMessageRetentionPeriod ||
			attr.VisibilityTimeout != DefaultVisibilityTimeout ||
			attr.PollingWaitSeconds != DefaultPollingWaitSeconds {
			if err = m.SetQueueAttributes(name,
				int32(task.DelaySeconds),
				DefaultMaxMessageSize,
				DefaultMessageRetentionPeriod,
				DefaultVisibilityTimeout,
				DefaultPollingWaitSeconds,
				DefaultSlices); err != nil {
				return
			}
		}

		// Succeed return
		// if queue exists and matched on local.
		log.Infofc(ctx, "aliyunmns-client: queue synced, name=%s, delay=%d, messages=%d, updated=%v",
			attr.QueueName, task.DelaySeconds,
			attr.ActiveMessages+attr.InactiveMessages+attr.DelayMessages,
			time.Unix(attr.LastModifyTime, 0).Format("2006-01-02/15:04:05"),
		)
		return
	}

	// Return error
	// if aliyunmns server response not recognized.
	if !RegexQueueNotExist.MatchString(err.Error()) {
		return
	}

	// Return error
	// if send create request to aliyunmns server failed.
	if err = m.CreateQueue(name,
		int32(task.DelaySeconds),
		DefaultMaxMessageSize,
		DefaultMessageRetentionPeriod,
		DefaultVisibilityTimeout,
		DefaultPollingWaitSeconds,
		DefaultSlices); err != nil {
		return
	}

	// Return succeed.
	log.Infofc(ctx, "aliyunmns-client: queue created, name=%s, delay=%d",
		name,
		task.DelaySeconds,
	)
	return
}

func (o *agent) buildSubscribe(ctx context.Context, task *base.Task) (err error) {
	var (
		attr mns.SubscriptionAttribute
		name = o.GenSubscriptionName(task.Id)
		tc   = Agent.GetTopicClient(task.TopicName)
	)

	// Read subscription attribute
	// from aliyunmns server.
	if attr, err = tc.GetSubscriptionAttributes(name); err == nil {
		// Send update request
		// if properties not matched on local task. Error
		// returned if request failed.
		if attr.NotifyStrategy != DefaultNotifyStrategy {
			if err = tc.SetSubscriptionAttributes(name, DefaultNotifyStrategy); err != nil {
				return
			}
		}

		// Succeed return
		// if subscription exists and matched on local.
		log.Infofc(ctx, "aliyunmns-client: subscription synced, name=%s, strategy=%v, updated=%s",
			name,
			attr.NotifyStrategy,
			time.Unix(attr.LastModifyTime, 0).Format("2006-01-02/15:04:05"),
		)
		return
	}

	// Return
	// if error type is subscription exists.
	if !RegexSubscriptionNotExist.MatchString(err.Error()) {
		return
	}

	// Return error
	// if aliyunmns server response not recognized.
	if err = tc.Subscribe(name, mns.MessageSubsribeRequest{
		Endpoint:            tc.GenerateQueueEndpoint(o.GenQueueName(task.Id)),
		FilterTag:           task.FilterTag,
		NotifyStrategy:      DefaultNotifyStrategy,
		NotifyContentFormat: DefaultNotifyContentFormat,
	}); err != nil {
		return
	}

	// Return succeed.
	log.Infofc(ctx, "aliyunmns-client: subscription created, name=%s, strategy=%v",
		name,
		DefaultNotifyStrategy,
	)
	return
}

func (o *agent) buildTopic(ctx context.Context, m mns.AliTopicManager, task *base.Task) (err error) {
	var (
		attr mns.TopicAttribute
		name = o.GenTopicName(task.TopicName)
	)

	// Read topic attribute
	// from aliyunmns server.
	if attr, err = m.GetTopicAttributes(name); err == nil {
		// Send update request
		// if properties not matched on local task. Error
		// returned if request failed.
		if attr.MaxMessageSize != int32(DefaultMaxMessageSize) ||
			attr.LoggingEnabled != DefaultLoggingEnable {
			if err = m.SetTopicAttributes(name,
				DefaultMaxMessageSize,
				DefaultLoggingEnable); err != nil {
				return
			}
		}

		// Succeed return
		// if topic exists and matched on local.
		log.Infofc(ctx, "aliyunmns-client: topic synced, name=%s, period=%d, messages=%d, updated=%v",
			name,
			attr.MessageRetentionPeriod,
			attr.MessageCount,
			time.Unix(attr.LastModifyTime, 0).Format("2006-01-02/15:04:05"),
		)
		return
	}

	// Return error
	// if aliyunmns server response not recognized.
	if !RegexTopicNotExist.MatchString(err.Error()) {
		return
	}

	// Return error
	// if send create request to aliyunmns server failed.
	if err = m.CreateTopic(name,
		DefaultMaxMessageSize,
		DefaultLoggingEnable); err != nil {
		return
	}

	// Return succeed.
	log.Infofc(ctx, "aliyunmns-client: topic created, name=%s, period=%d",
		name,
		DefaultMessageRetentionPeriod,
	)
	return
}

func (o *agent) destroyQueue(ctx context.Context, m mns.AliQueueManager, task *base.Task) (err error) {
	var (
		attr mns.QueueAttribute
		name = o.GenQueueName(task.Id)
	)

	// Return queue attribute
	// from aliyunmns server.
	if attr, err = m.GetQueueAttributes(name); err != nil {
		// Reset error
		// if queue not exists.
		if RegexQueueNotExist.MatchString(err.Error()) {
			err = nil
			log.Infofc(ctx, "aliyunmns-client: queue not exists, name=%s",
				name,
			)
		}
		return
	}

	// Return error
	// if send create request to aliyunmns server failed.
	if err = m.DeleteQueue(name); err != nil {
		return
	}

	// Return succeed.
	log.Infofc(ctx, "aliyunmns-client: queue deleted, name=%s, delay=%d, messages=%d, updated=%s",
		name,
		attr.DelaySeconds,
		attr.ActiveMessages+attr.InactiveMessages+attr.DelayMessages,
		time.Unix(attr.LastModifyTime, 0).Format("2006-01-02/15:04:05"),
	)
	return
}

func (o *agent) destroySubscription(ctx context.Context, task *base.Task) (err error) {
	var (
		attr mns.SubscriptionAttribute
		name = o.GenSubscriptionName(task.Id)
		tc   = o.GetTopicClient(task.TopicName)
	)

	// Return subscription attribute
	// from aliyunmns server.
	if attr, err = tc.GetSubscriptionAttributes(name); err != nil {
		// Reset error
		// if topic or subscription not exists.
		if RegexTopicNotExist.MatchString(err.Error()) ||
			RegexSubscriptionNotExist.MatchString(err.Error()) {
			err = nil
			log.Infofc(ctx, "aliyunmns-client: subscription not exists, name=%s",
				name,
			)
		}
		return
	}

	// Return error
	// if send unsubscribe request to aliyunmns server failed.
	if err = tc.Unsubscribe(name); err != nil {
		return
	}

	// Return succeed.
	log.Infofc(ctx, "aliyunmns-client: subscription deleted, name=%s, strategy=%s, updated=%s",
		name,
		attr.NotifyStrategy,
		time.Unix(attr.LastModifyTime, 0).Format("2006-01-02/15:04:05"),
	)
	return
}

func (o *agent) destroyTopic(ctx context.Context, m mns.AliTopicManager, task *base.Task) (err error) {
	var (
		attr mns.TopicAttribute
		name = o.GenTopicName(task.TopicName)
	)

	// Read topic attribute
	// from aliyunmns server.
	if attr, err = m.GetTopicAttributes(name); err != nil {
		// Reset error
		// if topic not exists.
		if RegexTopicNotExist.MatchString(err.Error()) {
			err = nil
			log.Infofc(ctx, "aliyunmns-client: topic not exists, name=%s", name)
		}
		return
	}

	// Validate process.
	var (
		ss mns.Subscriptions
		tc = o.GetTopicClient(task.TopicName)
	)

	// Return error
	// if send list request to aliyunmns server failed.
	if ss, err = tc.ListSubscriptionByTopic("", 0, conf.Config.Account.Aliyunmns.Prefix); err != nil {
		return
	}

	// Return succeed
	// if any subscription found under topic.
	if n := len(ss.Subscriptions); n > 0 {
		log.Infofc(ctx, "aliyunmns-client: topic exists and subscription exists, name=%s, subscriptions=%d",
			name, n)
		return
	}

	// Return error
	// if send delete request to aliyunmns server failed.
	if err = m.DeleteTopic(name); err != nil {
		return
	}

	// Return succeed.
	log.Infofc(ctx, "aliyunmns-client: topic deleted, name=%s, period=%d, messages=%d",
		name,
		attr.MessageRetentionPeriod,
		attr.MessageCount,
	)
	return nil
}

// /////////////////////////////////////////////////////////////
// Access methods.
// /////////////////////////////////////////////////////////////

func (o *agent) client() mns.MNSClient {
	return mns.NewAliMNSClient(
		conf.Config.Account.Aliyunmns.Endpoint,
		conf.Config.Account.Aliyunmns.AccessId,
		conf.Config.Account.Aliyunmns.AccessKey,
	)
}

func (o *agent) removePrefix(s string) string {
	if conf.Config.Account.Aliyunmns.Prefix == "" {
		return s
	}
	return strings.TrimPrefix(s, conf.Config.Account.Aliyunmns.Prefix)
}

// /////////////////////////////////////////////////////////////
// Construct methods.
// /////////////////////////////////////////////////////////////

func (o *agent) init() *agent {
	o.mu = &sync.RWMutex{}
	o.queueClients = make(map[int]mns.AliMNSQueue)
	o.topicClients = make(map[string]mns.AliMNSTopic)
	return o
}
