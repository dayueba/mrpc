package testdata

import "context"

type Service struct {

}

type HelloRequest struct {
	Msg string
}

type HelloReply struct {
	Msg string
}

type AddReply struct {
	Result int32 `msgpack:"result"`
}

func (s *Service) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	rsp := &HelloReply{
		Msg : "world",
	}

	return rsp, nil
}

func (s *Service) Add(ctx context.Context, req *HelloRequest) (*AddReply, error) {
	rsp := &AddReply{
		Result: 1,
	}

	return rsp, nil
}
