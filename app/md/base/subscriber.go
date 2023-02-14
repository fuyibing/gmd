// author: wsfuyibing <websearch@163.com>
// date: 2023-02-08

package base

import (
	"fmt"
	"github.com/fuyibing/gmd/app/md/conf"
	"github.com/fuyibing/gmd/app/models"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type (
	SubscriberProtocol     int
	SubscriberResponseType int
	SubscriberType         int
)

const (
	_ SubscriberProtocol = iota

	SubscriberProtocolHttp
	SubscriberProtocolRpc
	SubscriberProtocolTcp
	SubscriberProtocolWebsocket
)

var (
	regexSubscriberProtocol = regexp.MustCompile(`^([^:]+)://`)
	regexSubscriberPort     = regexp.MustCompile(`^([^:]+)://([_a-zA-Z0-9-.]+):?(\d*)/?`)

	// SubscriberProtocolDefault
	// default subscription protocol.
	SubscriberProtocolDefault = "http"

	// SubscriberProtocolList
	// subscription protocol list.
	SubscriberProtocolList = map[string]SubscriberProtocol{
		"http":  SubscriberProtocolHttp,
		"https": SubscriberProtocolHttp,
		"rpc":   SubscriberProtocolRpc,
		"grpc":  SubscriberProtocolRpc,
		"tcp":   SubscriberProtocolTcp,
		"ws":    SubscriberProtocolWebsocket,
		"wss":   SubscriberProtocolWebsocket,
	}

	// SubscriberProtocolPort
	// subscription port mapping.
	SubscriberProtocolPort = map[string]int{
		"http":  80,
		"https": 443,
	}
)

const (
	_ SubscriberResponseType = iota
	SubscriberResponseTypeErrnoIsZero
)

const (
	_ SubscriberType = iota

	SubscriberTypeHandler
	SubscriberTypeFailed
	SubscriberTypeSucceed
)

type (
	// Subscriber
	// struct for subscription.
	Subscriber struct {
		Host, Addr, Method string
		Port, Timeout      int

		Condition    ConditionManager
		IgnoreCodes  []string
		Protocol     SubscriberProtocol
		ResponseType SubscriberResponseType
	}
)

// NewSubscriber
// create and return subscriber instance.
func NewSubscriber(m *models.Task, t SubscriberType) (s *Subscriber) {
	switch t {
	case SubscriberTypeHandler:
		if m.Handler != "" {
			s = (&Subscriber{
				Addr: m.Handler, Method: m.HandlerMethod, Timeout: m.HandlerTimeout,
				ResponseType: SubscriberResponseType(m.HandlerResponseType),
			}).init(m.HandlerCondition, m.HandlerIgnoreCodes)
		}
	case SubscriberTypeFailed:
		if m.Failed != "" {
			s = (&Subscriber{
				Addr: m.Failed, Method: m.FailedMethod, Timeout: m.FailedTimeout,
				ResponseType: SubscriberResponseType(m.FailedResponseType),
			}).init(m.FailedCondition, m.FailedIgnoreCodes)
		}
	case SubscriberTypeSucceed:
		if m.Succeed != "" {
			s = (&Subscriber{
				Addr: m.Succeed, Method: m.SucceedMethod, Timeout: m.SucceedTimeout,
				ResponseType: SubscriberResponseType(m.SucceedResponseType),
			}).init(m.SucceedCondition, m.SucceedIgnoreCodes)
		}
	}

	return
}

// /////////////////////////////////////////////////////////////
// Access methods.
// /////////////////////////////////////////////////////////////

func (o *Subscriber) init(sc, ic string) *Subscriber {
	// Condition definition.
	if sc = strings.TrimSpace(sc); sc != "" {
		o.Condition = (&condition{s: sc}).init()
	}

	// Ignored codes definition.
	if ic = strings.TrimSpace(ic); ic != "" {
		o.IgnoreCodes = make([]string, 0)
		for _, s := range strings.Split(ic, ",") {
			if s = strings.TrimSpace(s); s != "" {
				o.IgnoreCodes = append(o.IgnoreCodes, s)
			}
		}
	}

	// Extension fields.
	o.initDefaults()
	o.initProtocol()
	o.initHostAndPort()
	return o
}

func (o *Subscriber) initDefaults() {
	// Set default protocol
	// as HTTP.
	o.Protocol = SubscriberProtocolHttp

	// Set default method
	// as POST.
	if o.Method == "" {
		o.Method = http.MethodPost
	}

	// Set default timeout
	// seconds.
	if o.Timeout == 0 {
		o.Timeout = conf.Config.Consumer.DispatchTimeout
	}
}

func (o *Subscriber) initHostAndPort() {
	var k string

	// Parse host and port
	// from address.
	if m := regexSubscriberPort.FindStringSubmatch(o.Addr); len(m) == 4 {
		k = strings.ToLower(m[1])

		// Collect host.
		o.Host = m[2]

		// Collect port.
		if m[3] != "" {
			if i, ie := strconv.ParseInt(m[3], 0, 32); ie == nil && i > 0 {
				o.Port = int(i)
			}
		}
	}

	// Reset port.
	if o.Port == 0 && k != "" {
		if i, ok := SubscriberProtocolPort[k]; ok && i > 0 {
			o.Port = i
		}
	}
}

func (o *Subscriber) initProtocol() {
	// Parse protocol from address.
	if m := regexSubscriberProtocol.FindStringSubmatch(o.Addr); len(m) == 2 {
		if k := strings.TrimSpace(strings.ToLower(m[1])); k != "" {
			if v, ok := SubscriberProtocolList[k]; ok {
				o.Protocol = v
			}
			return
		}
	}

	// Fill default protocol.
	o.Addr = fmt.Sprintf("%v://%s", SubscriberProtocolDefault, o.Addr)
}
