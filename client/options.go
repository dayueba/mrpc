package client

import (
	"time"

	"github.com/dayueba/mrpc/interceptor"
	// "github.com/dayueba/mrpc/transport"
)

// Options defines the client call parameters
type Options struct {
	serviceName  string        // service name
	method       string        // method name
	target       string        // format e.g.:  ip:port 127.0.0.1:8000
	timeout      time.Duration // timeout
	selectorName string        // 服务发现插件名, e.g. : consul、zookeeper、etcd
	// transportOpts transport.ClientTransportOptions
	interceptors []interceptor.ClientInterceptor
}

type Option func(*Options)

func WithServiceName(serviceName string) Option {
	return func(o *Options) {
		o.serviceName = serviceName
	}
}

func WithMethod(method string) Option {
	return func(o *Options) {
		o.method = method
	}
}

func WithTarget(target string) Option {
	return func(o *Options) {
		o.target = target
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.timeout = timeout
	}
}

func WithInterceptor(interceptors ...interceptor.ClientInterceptor) Option {
	return func(o *Options) {
		o.interceptors = append(o.interceptors, interceptors...)
	}
}

func WithSelectorName(selectorName string) Option {
	return func(o *Options) {
		o.selectorName = selectorName
	}
}
