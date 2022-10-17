package mrpc

import (
	"time"

	"github.com/dayueba/mrpc/interceptor"
)

type ServerOptions struct {
	address           string        // 监听地址, e.g. :( ip://127.0.0.1:8080、 dns://www.google.com)
	timeout           time.Duration // timeout
	name string

	interceptors    []interceptor.ServerInterceptor
	pluginNames []string         // 插件名字
	selectorSvrAddr string       // 服务发现地址
}

type ServerOption func(*ServerOptions)

func WithAddress(address string) ServerOption {
	return func(o *ServerOptions) {
		o.address = address
	}
}

func WithName(name string) ServerOption {
	return func(o *ServerOptions) {
		o.name = name
	}
}

func WithTimeout(timeout time.Duration) ServerOption {
	return func(o *ServerOptions) {
		o.timeout = timeout
	}
}

func WithInterceptor(interceptors ...interceptor.ServerInterceptor) ServerOption {
	return func(o *ServerOptions) {
		o.interceptors = append(o.interceptors, interceptors...)
	}
}

func WithPlugin(pluginName ... string) ServerOption {
	return func(o *ServerOptions) {
		o.pluginNames = append(o.pluginNames, pluginName ...)
	}
}

func WithSelectorSvrAddr(addr string) ServerOption {
	return func(o *ServerOptions) {
		o.selectorSvrAddr = addr
	}
}
