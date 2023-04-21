package client

import (
	"context"
	"math"
	"strconv"

	"github.com/dayueba/mrpc/codec"
	"github.com/dayueba/mrpc/codes"
	"github.com/dayueba/mrpc/interceptor"
	"github.com/dayueba/mrpc/pool/connpool"
	"github.com/dayueba/mrpc/protocol"
	"github.com/dayueba/mrpc/selector"
	"github.com/dayueba/mrpc/transport"

	// "github.com/dayueba/mrpc/utils"
	"github.com/mitchellh/mapstructure"
)

type Client interface {
	// Invoke 这个方法表示向下游服务发起调用
	Invoke(ctx context.Context, req, rsp interface{}, path string, opts ...Option) error
}

// DefaultClient a global client
var DefaultClient = New()

var New = func() *defaultClient {
	return &defaultClient{
		opts: &Options{},
	}
}

type defaultClient struct {
	opts  *Options
	msgId int32
}

// call by reflect
func (c *defaultClient) Call(ctx context.Context, method string, params interface{}, rsp interface{},
	opts ...Option) error {

	// reflection calls need to be serialized using msgpack
	callOpts := make([]Option, 0, len(opts)+1)
	callOpts = append(callOpts, opts...)

	msgId := c.msgId
	if msgId == math.MaxInt32 {
		c.msgId = 1
	}

	req := &protocol.Request{
		Method: method,
		Type:   "call",
		Params: params,
		MsgId:  strconv.Itoa(int(msgId)),
	}

	err := c.Invoke(ctx, req, rsp, method, callOpts...)
	if err != nil {
		return err
	}

	return nil
}

func (c *defaultClient) Invoke(ctx context.Context, req, rsp interface{}, path string, opts ...Option) error {
	for _, o := range opts {
		o(c.opts)
	}

	if c.opts.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.opts.timeout)
		defer cancel()
	}

	// _, method := utils.ParseServicePath(path)

	// c.opts.method = method
	// c.opts.serviceName = serviceName

	// execute the interceptor first
	return interceptor.ClientIntercept(ctx, req, rsp, c.opts.interceptors, c.invoke)
}

func (c *defaultClient) invoke(ctx context.Context, req, rsp interface{}) error {
	// 对请求体序列化
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

	// 添加包头
	clientCodec := codec.DefaultCodec
	reqbody, err := clientCodec.Encode(payload)
	if err != nil {
		return err
	}

	// 发送请求
	clientTransport := c.NewClientTransport()
	clientTransportOpts := []transport.ClientTransportOption{
		transport.WithServiceName(c.opts.serviceName),
		transport.WithClientTarget(c.opts.target),
		transport.WithClientNetwork("tcp"),
		transport.WithClientPool(connpool.GetPool("default")),
		transport.WithTimeout(c.opts.timeout),
		transport.WithSelector(selector.GetSelector(c.opts.selectorName)),
	}
	frame, err := clientTransport.Send(ctx, reqbody, clientTransportOpts...)
	if err != nil {
		return err
	}

	// 对 server 回包进行解包
	rspbuf, err := clientCodec.Decode(frame)
	if err != nil {
		return err
	}

	respp := make([]interface{}, 0)
	err = serialization.Unmarshal(rspbuf, &respp)
	if err != nil {
		return err
	}

	if respp[1].(string) == "error" {
		// todo: to rpc error
		e := protocol.RpcError{}
		err = mapstructure.Decode(respp[len(respp)-1], &e)
		if err != nil {
			return err
		}
		return e
	}

	// 转结构体
	return mapstructure.Decode(respp[len(respp)-1], &rsp)
}

func (c *defaultClient) NewClientTransport() transport.ClientTransport {
	return transport.DefaultClientTransport
}
