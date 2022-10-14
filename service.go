package mrpc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dayueba/mrpc/codec"
	"github.com/dayueba/mrpc/interceptor"
	"github.com/dayueba/mrpc/utils"
	"github.com/dayueba/mrpc/protocol"
	"github.com/dayueba/mrpc/transport"
	"github.com/mitchellh/mapstructure"
)

type Service interface {
	Register(string, Handler)
	Serve(*ServerOptions)
	Close()
	Name() string
	AddSvr(string, interface{})
}

type service struct {
	svr         map[string]interface{}      // server
	ctx         context.Context    // Each service is managed in one context
	cancel      context.CancelFunc // controller of context
	serviceName string             // service name
	handlers    map[string]Handler
	opts        *ServerOptions // parameter options

	closing bool // whether the service is closing
}

// ServiceDesc is a detailed description of a service
type ServiceDesc struct {
	Svr         interface{}
	ServiceName string
	Methods     []*MethodDesc
	HandlerType interface{}
}

// MethodDesc is a detailed description of a method
type MethodDesc struct {
	MethodName string
	Handler    Handler
}

// Handler is the handler of a method
type Handler func(context.Context, interface{}, func(interface{}) error, []interceptor.ServerInterceptor) (interface{}, error)

func (s *service) Register(handlerName string, handler Handler) {
	if s.handlers == nil {
		s.handlers = make(map[string]Handler)
	}
	s.handlers[handlerName] = handler
}

func (s *service) AddSvr(serviceName string, svr interface{}) {
	if s.svr == nil {
		s.svr = make(map[string]interface{})
	}
	s.svr[serviceName] = svr
}

func (s *service) Serve(opts *ServerOptions) {
	s.opts = opts

	transportOpts := []transport.ServerTransportOption{
		transport.WithServerAddress(s.opts.address),
		transport.WithHandler(s),
		transport.WithServerTimeout(s.opts.timeout),
	}

	serverTransport := transport.DefaultServerTransport
	s.ctx, s.cancel = context.WithCancel(context.Background())
	if err := serverTransport.ListenAndServe(s.ctx, transportOpts...); err != nil {
		fmt.Printf("tcp serve error, %v", err)
		return
	}

	fmt.Printf("service serving at %s ... \n", s.opts.address)

	<-s.ctx.Done()
}

func (s *service) Close() {
	s.closing = true
	if s.cancel != nil {
		s.cancel()
	}
	fmt.Println("service closing ...")
}

func (s *service) Name() string {
	return s.serviceName
}

func (s *service) Handle(ctx context.Context, reqbuf []byte) ([]byte, error) {
	request := []interface{}{}
	serverSerialization := codec.DefaultSerialization
	err := serverSerialization.Unmarshal(reqbuf, &request)
	if err != nil {
		return nil, err
	}
	msgId := request[0].(string)
	payload := request[len(request)-1]
	pathArr := make([]string, 0)

	for i := 2; i < len(request) - 1; i++ {
		pathArr = append(pathArr, request[i].(string))
	}
	path := strings.ToLower(strings.Join(pathArr, "."))
	
	srvName, _ := utils.ParseServicePath(path)

	dec := func(req interface{}) error {
		if err := mapstructure.Decode(payload, req); err != nil {
			return protocol.RpcError{
				Message: err.Error(),
			}
		}
		return nil
	}

	if s.opts.timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.opts.timeout)
		defer cancel()
	}

	// _, method, err := utils.ParseServicePath(string(request.ServicePath))
	// if err != nil {
	// 	return nil, codes.New(codes.ClientMsgErrorCode, "method is invalid")
	// }


	// 首字母大写
	handler := s.handlers[path]
	if handler == nil {
		return nil, errors.New("handlers is nil")
	}

	rsp, err := handler(ctx, s.svr[srvName], dec, s.opts.interceptors)
	result := []interface{}{}

	if err != nil {
		result = append(result, msgId)
		// todo 
		// result = append(result, "error")
		result = append(result, "reply")
		result = append(result, err)
		// return nil, err
	} else {
		result = append(result, msgId)
		result = append(result, "reply")
		result = append(result, rsp)
	}

	// result = append(result, msgId)
	// result = append(result, "reply")
	// result = append(result, rsp)

	rspbuf, err := serverSerialization.Marshal(result)
	if err != nil {
		return nil, err
	}

	return rspbuf, nil
}
