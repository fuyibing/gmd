// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package core

import (
	"github.com/fuyibing/gmd/v8/app/md/base"
	"github.com/fuyibing/gmd/v8/app/md/conf"
	"github.com/fuyibing/gmd/v8/app/md/core/adapters/aliyunmns"
	"github.com/fuyibing/gmd/v8/app/md/core/adapters/rocketmq"
	"sync"
)

var (
	Boot   BootManager
	Config conf.Configuration

	builtInConsumer = map[base.Adapter]base.ConsumerCallable{
		base.AliyunMns: func(id, parallel int) base.ConsumerManager {
			return aliyunmns.NewConsumer(id, parallel)
		},
		base.RocketMq: func(id, parallel int) base.ConsumerManager {
			return rocketmq.NewConsumer(id, parallel)
		},
	}

	builtInProducer = map[base.Adapter]base.ProducerCallable{
		base.AliyunMns: func() base.ProducerManager {
			return aliyunmns.NewProducer()
		},
		base.RocketMq: func() base.ProducerManager {
			return rocketmq.NewProducer()
		},
	}

	builtInRemoting = map[base.Adapter]base.RemotingCallable{
		base.AliyunMns: func() base.RemotingManager {
			return aliyunmns.NewRemoting()
		},
		base.RocketMq: func() base.RemotingManager {
			return rocketmq.NewRemoting()
		},
	}

	buildInCondition  = map[base.Condition]func() base.ConditionManager{}
	buildInDispatcher = map[base.Dispatcher]func() base.DispatcherManager{}
	buildInResult     = map[base.Result]func() base.Result{}
)

func init() {
	new(sync.Once).Do(func() {
		Config = conf.Config
		Boot = (&boot{}).init()
	})
}
