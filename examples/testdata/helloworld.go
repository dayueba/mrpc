package testdata

import (
	"fmt"
	"context"
	"errors"
)

type DB string

type Service struct {
	db DB
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

func NewService(db DB) *Service {
	return &Service{db: db}
} 

func (s *Service) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	rsp := &HelloReply{
		Msg : "world",
	}

	return rsp, nil
}

func (s *Service) Add(ctx context.Context, req *HelloRequest) (*AddReply, error) {
	fmt.Println(s.db)
	rsp := &AddReply{
		Result: 1,
	}

	return rsp, nil
}

func (s *Service) Oops(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	return nil, errors.New("oops")
}