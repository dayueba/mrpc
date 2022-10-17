package transport

import (
	"context"
)

type clientTransport struct {
	opts *ClientTransportOptions
}

var clientTransportMap = make(map[string]ClientTransport)

func init() {
	clientTransportMap["default"] = DefaultClientTransport
}

// RegisterClientTransport supports business custom registered ClientTransport
func RegisterClientTransport(name string, clientTransport ClientTransport) {
	if clientTransportMap == nil {
		clientTransportMap = make(map[string]ClientTransport)
	}
	clientTransportMap[name] = clientTransport
}

// Get the ServerTransport
func GetClientTransport(transport string) ClientTransport {
	if v, ok := clientTransportMap[transport]; ok {
		return v
	}

	return DefaultClientTransport
}

// The default ClientTransport
var DefaultClientTransport = New()

// Use the singleton pattern to create a ClientTransport
var New = func() ClientTransport {
	return &clientTransport{
		opts: &ClientTransportOptions{},
	}
}

func (c *clientTransport) Send(ctx context.Context, req []byte, opts ...ClientTransportOption) ([]byte, error) {
	for _, o := range opts {
		o(c.opts)
	}

	return c.SendTcpReq(ctx, req)
}

func (c *clientTransport) SendTcpReq(ctx context.Context, req []byte) ([]byte, error) {
	// 服务发现
	addr, err := c.opts.Selector.Select(c.opts.ServiceName)
	if err != nil {
		return nil, err
	}

	if addr == "" {
		addr = c.opts.Target
	}

	conn, err := c.opts.Pool.Get(ctx, c.opts.Network, addr)
	//	conn, err := net.DialTimeout("tcp", addr, c.opts.Timeout);
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	sendNum := 0
	num := 0
	for sendNum < len(req) {
		num, err = conn.Write(req[sendNum:])
		if err != nil {
			return nil, err
		}
		sendNum += num

		if err = isDone(ctx); err != nil {
			return nil, err
		}
	}

	// parse frame
	wrapperConn := wrapConn(conn)
	frame, err := wrapperConn.framer.ReadFrame(conn)
	if err != nil {
		return nil, err
	}

	return frame, err
}

func isDone(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return nil
}
