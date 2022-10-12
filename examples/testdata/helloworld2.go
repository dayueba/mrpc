package testdata

import "context"

type Service2 struct {

}

type HelloRequest2 struct {
	Msg string
}

type AddReply2 struct {
	Result int32 `msgpack:"result"`
}

func (s *Service2) Add(ctx context.Context, req *HelloRequest2) (*AddReply2, error) {
	rsp := &AddReply2{
		Result: 2,
	}

	return rsp, nil
}
