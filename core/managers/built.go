// author: wsfuyibing <websearch@163.com>
// date: 2023-02-19

package managers

import (
	"github.com/fuyibing/gmd/v8/core/adapters/aliyunmns"
	"github.com/fuyibing/gmd/v8/core/adapters/rocketmq"
	"github.com/fuyibing/gmd/v8/core/base"
	"github.com/fuyibing/gmd/v8/core/conditions/el"
	"github.com/fuyibing/gmd/v8/core/dispatchers/http_get"
	"github.com/fuyibing/gmd/v8/core/dispatchers/http_post_form"
	"github.com/fuyibing/gmd/v8/core/dispatchers/http_post_json"
	"github.com/fuyibing/gmd/v8/core/dispatchers/rpc"
	"github.com/fuyibing/gmd/v8/core/dispatchers/tcp"
	"github.com/fuyibing/gmd/v8/core/dispatchers/wss"
	"github.com/fuyibing/gmd/v8/core/results/http_ok"
	"github.com/fuyibing/gmd/v8/core/results/json_errno_zero"
)

var (
	builtInConsumer = map[base.Adapter]base.ConsumerCallable{
		base.AliyunMns: func(id, parallel int, name string, process base.ConsumerProcess) base.ConsumerManager {
			return aliyunmns.NewConsumer(id, parallel, name, process)
		},
		base.RocketMq: func(id, parallel int, name string, process base.ConsumerProcess) base.ConsumerManager {
			return rocketmq.NewConsumer(id, parallel, name, process)
		},
	}

	builtInProducer = map[base.Adapter]base.ProducerCallable{
		base.AliyunMns: func() base.ProducerManager { return aliyunmns.NewProducer() },
		base.RocketMq:  func() base.ProducerManager { return rocketmq.NewProducer() },
	}

	builtInRemoting = map[base.Adapter]base.RemotingCallable{
		base.AliyunMns: func() base.RemotingManager { return aliyunmns.NewRemoting() },
		base.RocketMq:  func() base.RemotingManager { return rocketmq.NewRemoting() },
	}

	buildInConditions = map[string]base.ConditionCallable{
		base.ConditionEl: func() base.ConditionManager { return el.New() },
	}

	buildInDispatchers = map[string]base.DispatcherCallable{
		base.DispatchHttpGet:      func() base.DispatcherManager { return http_get.New() },
		base.DispatchHttpPostForm: func() base.DispatcherManager { return http_post_form.New() },
		base.DispatchHttpPostJson: func() base.DispatcherManager { return http_post_json.New() },
		base.DispatchRpc:          func() base.DispatcherManager { return rpc.New() },
		base.DispatchTcp:          func() base.DispatcherManager { return tcp.New() },
		base.DispatchWebsocket:    func() base.DispatcherManager { return wss.New() },
	}

	buildInResults = map[string]base.ResultCallable{
		base.ResultHttpOk:        func() base.ResultManager { return http_ok.New() },
		base.ResultJsonErrnoZero: func() base.ResultManager { return json_errno_zero.New() },
	}
)
