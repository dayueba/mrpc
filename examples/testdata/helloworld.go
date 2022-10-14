package testdata

import (
	"fmt"
	"context"
	// "errors"


	"github.com/dayueba/mrpc/protocol"
	"github.com/dayueba/mrpc/log"
)

type DB string

type Service struct {
	db DB
	log *log.Helper
}

type HelloRequest struct {
	Msg string
}

type HelloReply struct {
	Msg string
}

type AddRequest struct {
	A int32 `msgpack:"a"`
	B int32 `msgpack:"b"`
}

type AddReply struct {
	Result int32 `msgpack:"result"`
}

func NewService(db DB, logger log.Logger) *Service {
	return &Service{db: db, log: log.NewHelper(logger)}
} 

func (s *Service) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	rsp := &HelloReply{
		Msg : "world",
	}

	return rsp, nil
}

func (s *Service) Add(ctx context.Context, req *AddRequest) (*AddReply, error) {
	fmt.Println(s.db)
	rsp := &AddReply{
		Result: req.A + req.B,
	}

	return rsp, nil
}

func (s *Service) Oops(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	s.log.Log(log.LevelInfo, "message", "have request")
	return nil, protocol.RpcError{
		ErrorCode: -1,
		Message: "1231",
	}
}