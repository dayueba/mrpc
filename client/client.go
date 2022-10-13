package client

import (
	"context"
	"math"
	"strconv"
	// "fmt"

	"github.com/dayueba/mrpc/codec"
	"github.com/dayueba/mrpc/codes"
	"github.com/dayueba/mrpc/interceptor"
	"github.com/dayueba/mrpc/pool/connpool"
	"github.com/dayueba/mrpc/protocol"
	"github.com/dayueba/mrpc/transport"
	"github.com/dayueba/mrpc/utils"
	"github.com/mitchellh/mapstructure"
)

// global client interface
type Client interface {
	// Invoke 这个方法表示向下游服务发起调用
	Invoke(ctx context.Context, req , rsp interface{}, path string, opts ...Option) error
}

// use a global client
var DefaultClient = New()

var New = func() *defaultClient {
	return &defaultClient{
		opts : &Options{
		},
	}
}

type defaultClient struct {
	opts *Options
	msgId int
}

// call by reflect
func (c *defaultClient) Call(ctx context.Context, method string, params interface{}, rsp interface{},
	opts ...Option) error {

	// reflection calls need to be serialized using msgpack
	callOpts := make([]Option, 0, len(opts)+1)
	callOpts = append(callOpts, opts ...)

	msgId := c.msgId
	if msgId == math.MaxInt32 {
		c.msgId = 1
	}

	req := &protocol.Request{
		Method: method,
		Type: "call",
		Params: params,
		MsgId: strconv.Itoa(msgId),
	}
	
	err := c.Invoke(ctx, req, rsp, method, callOpts ...)
	if err != nil {
		return err
	}

	return nil
}


func (c *defaultClient) Invoke(ctx context.Context, req , rsp interface{}, path string, opts ...Option) error {
	for _, o := range opts {
		o(c.opts)
	}

	if c.opts.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.opts.timeout)
		defer cancel()
	}

	serviceName, method := utils.ParseServicePath(path)

	c.opts.method = method
	c.opts.serviceName = serviceName

	// execute the interceptor first
	return interceptor.ClientIntercept(ctx, req, rsp, c.opts.interceptors, c.invoke)
}

func (c *defaultClient) invoke(ctx context.Context, req, rsp interface{}) error {
	serialization := codec.DefaultSerialization
	arr := make([]interface{}, 0)
	r := req.(*protocol.Request)
	arr = append(arr, r.MsgId)
	arr = append(arr, r.Type)
	arr = append(arr, r.Method)
	arr = append(arr, r.Params)
	payload, err := serialization.Marshal(arr)
	
	if err != nil {
		return codes.NewFrameworkError(codes.ClientMsgErrorCode, "request marshal failed ...")
	}

	clientCodec := codec.DefaultCodec

	reqbody, err := clientCodec.Encode(payload)
	if err != nil {
		return err
	}

	clientTransport := c.NewClientTransport()
	clientTransportOpts := []transport.ClientTransportOption {
		transport.WithServiceName(c.opts.serviceName),
		transport.WithClientTarget(c.opts.target),
		transport.WithClientNetwork("tcp"),
		transport.WithClientPool(connpool.GetPool("default")),
		transport.WithTimeout(c.opts.timeout),
	}
	frame, err := clientTransport.Send(ctx, reqbody, clientTransportOpts ...)
	if err != nil {
		return err
	}

	rspbuf, err := clientCodec.Decode(frame)
	if err != nil {
		return err
	}

	respp := make([]interface{}, 0)
	err = serialization.Unmarshal(rspbuf, &respp)
	if err != nil {
		return err
	}

	return mapstructure.Decode(respp[len(respp)-1],&rsp)
}

func (c *defaultClient) NewClientTransport() transport.ClientTransport {
	return transport.DefaultClientTransport
}
