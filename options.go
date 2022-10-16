package mrpc

import (
	"time"

	"github.com/dayueba/mrpc/interceptor"
)

type ServerOptions struct {
	address           string        // 监听地址, e.g. :( ip://127.0.0.1:8080、 dns://www.google.com)
	timeout           time.Duration // timeout

	interceptors    []interceptor.ServerInterceptor
}

type ServerOption func(*ServerOptions)

func WithAddress(address string) ServerOption {
	return func(o *ServerOptions) {
		o.address = address
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
