package testdata

import (
	"context"
	"time"
	"sync/atomic"
	"errors"
	"fmt"
)

type Service2 struct {
	count int32
}

type AddReply2 struct {
	Result int32 `msgpack:"result"`
}

func (s *Service2) Add(ctx context.Context, req *AddRequest) (*AddReply2, error) {
	time.Sleep(300 * time.Millisecond)
	rsp := &AddReply2{
		Result: req.A + req.B,
	}

	return rsp, nil
}

func (s *Service2) Breaker(ctx context.Context, req *HelloRequest) (*CountResponse, error) {
	atomic.AddInt32(&s.count, 1)

	if s.count > 50 {
		return nil, errors.New(fmt.Sprintf("count is %d", s.count))
	}
	rsp := &CountResponse{
		Count: 100,
	}

	return rsp, nil
}
