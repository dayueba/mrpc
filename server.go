package mrpc

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"

	"github.com/dayueba/mrpc/interceptor"
	"github.com/dayueba/mrpc/log"
	"github.com/dayueba/mrpc/plugin"
	"github.com/dayueba/mrpc/plugin/jaeger"
	"github.com/dayueba/mrpc/protocol"
)

type Server struct {
	opts    *ServerOptions
	service Service
	name    string
	plugins []plugin.Plugin

	closing bool // whether the server is closing
}

func NewServer(opt ...ServerOption) *Server {
	s := &Server{
		opts: &ServerOptions{},
	}

	for _, o := range opt {
		o(s.opts)
	}

	s.service = NewService(s.opts)

	for pluginName, plugin := range plugin.PluginMap {
		if !containPlugin(pluginName, s.opts.pluginNames) {
			continue
		}
		s.plugins = append(s.plugins, plugin)
	}

	return s
}

func NewService(opts *ServerOptions) Service {
	return &service{
		opts: opts,
	}
}

func containPlugin(pluginName string, plugins []string) bool {
	for _, plugin := range plugins {
		if pluginName == plugin {
			return true
		}
	}
	return false
}

type emptyInterface interface{}

func (s *Server) RegisterService(serviceName string, svr interface{}) error {
	svrType := reflect.TypeOf(svr)
	svrValue := reflect.ValueOf(svr)

	sd := &ServiceDesc{
		ServiceName: serviceName,
		// for compatibility with code generation
		HandlerType: (*emptyInterface)(nil),
		Svr:         svr,
	}

	methods, err := getServiceMethods(svrType, svrValue)
	if err != nil {
		return err
	}

	sd.Methods = methods
	s.Register(sd, svr)
	return nil
}

func getServiceMethods(serviceType reflect.Type, serviceValue reflect.Value) ([]*MethodDesc, error) {
	var methods []*MethodDesc

	for i := 0; i < serviceType.NumMethod(); i++ {
		method := serviceType.Method(i)

		if err := checkMethod(method.Type); err != nil {
			return nil, err
		}

		methodHandler := func(ctx context.Context, svr interface{}, dec func(interface{}) error, ceps []interceptor.ServerInterceptor) (interface{}, error) {
			reqType := method.Type.In(2)

			// determine type
			req := reflect.New(reqType.Elem()).Interface()

			if err := dec(req); err != nil {
				return nil, err
			}

			if len(ceps) == 0 {
				values := method.Func.Call([]reflect.Value{serviceValue, reflect.ValueOf(ctx), reflect.ValueOf(req)})
				// determine error
				v1 := values[1].Interface()
				var err error
				if v1 == nil {
					err = nil
				} else if e, ok := v1.(protocol.RpcError); ok {
					err = e
				} else if e, ok := v1.(error); ok {
					err = protocol.RpcError{
						ErrorCode: -1,
						Message:   e.Error(),
					}
				}
				return values[0].Interface(), err
			}

			handler := func(ctx context.Context, reqbody interface{}) (interface{}, error) {
				values := method.Func.Call([]reflect.Value{serviceValue, reflect.ValueOf(ctx), reflect.ValueOf(req)})
				return values[0].Interface(), nil
			}

			return interceptor.ServerIntercept(ctx, req, ceps, handler)
		}

		methods = append(methods, &MethodDesc{
			MethodName: method.Name,
			Handler:    methodHandler,
		})
	}

	return methods, nil
}

func checkMethod(method reflect.Type) error {
	// params num must >= 2 , needs to be combined with itself
	if method.NumIn() < 3 {
		return fmt.Errorf("method %s invalid, the number of params < 2", method.Name())
	}

	// return values nums must be 2
	if method.NumOut() != 2 {
		return fmt.Errorf("method %s invalid, the number of return values != 2", method.Name())
	}

	// the first parameter must be context
	ctxType := method.In(1)
	var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()
	if !ctxType.Implements(contextType) {
		return fmt.Errorf("method %s invalid, first param is not context", method.Name())
	}

	// the second parameter type must be pointer
	argType := method.In(2)
	if argType.Kind() != reflect.Ptr {
		return fmt.Errorf("method %s invalid, req type is not a pointer", method.Name())
	}

	// the first return type must be a pointer
	replyType := method.Out(0)
	if replyType.Kind() != reflect.Ptr {
		return fmt.Errorf("method %s invalid, reply type is not a pointer", method.Name())
	}

	// The second return value must be an error
	errType := method.Out(1)
	var errorType = reflect.TypeOf((*error)(nil)).Elem()
	if !errType.Implements(errorType) {
		return fmt.Errorf("method %s invalid, returns %s , not error", method.Name(), errType.Name())
	}

	return nil
}

func (s *Server) Register(sd *ServiceDesc, svr interface{}) {
	if sd == nil || svr == nil {
		return
	}
	ht := reflect.TypeOf(sd.HandlerType).Elem()
	st := reflect.TypeOf(svr)
	if !st.Implements(ht) {
		log.Fatalf("handlerType %v not match service : %v ", ht, st)
	}

	if s.service == nil {
		s.service = &service{
			svr:      make(map[string]interface{}),
			handlers: make(map[string]Handler),
		}
	}

	s.service.AddSvr(strings.ToLower(sd.ServiceName), svr)

	for _, method := range sd.Methods {
		handlerName := strings.ToLower(sd.ServiceName + "." + method.MethodName)
		s.service.Register(handlerName, method.Handler)
	}
}

func (s *Server) Serve() {
	err := s.InitPlugins()
	if err != nil {
		panic(err)
	}

	s.service.Serve(s.opts)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGSEGV)
	<-ch

	s.Close()
}

func (s *Server) Close() {
	s.closing = false
	// todo 服务发现：注销
	s.service.Close()
}

func (s *Server) InitPlugins() error {
	// init plugins
	for _, p := range s.plugins {
		switch val := p.(type) {
		case plugin.ResolverPlugin:
			var services []string
			services = append(services, s.opts.name)

			pluginOpts := []plugin.Option{
				plugin.WithSelectorSvrAddr(s.opts.selectorSvrAddr),
				plugin.WithSvrAddr(s.opts.address),
				plugin.WithServices(services),
			}
			if err := val.Init(pluginOpts...); err != nil {
				log.Errorf("resolver init error, %v", err)
				return err
			}
		case plugin.TracingPlugin:

			pluginOpts := []plugin.Option{
				plugin.WithTracingSvrAddr(s.opts.tracingSvrAddr),
			}

			tracer, err := val.Init(pluginOpts...)
			if err != nil {
				log.Errorf("tracing init error, %v", err)
				return err
			}

			s.opts.interceptors = append(s.opts.interceptors, jaeger.OpenTracingServerInterceptor(tracer, s.opts.tracingSpanName))
		}
	}

	return nil
}
