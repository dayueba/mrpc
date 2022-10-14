package main

import (
	"time"

	"go.uber.org/zap"
	"github.com/dayueba/mrpc"
	"github.com/dayueba/mrpc/examples/testdata"
	"github.com/dayueba/mrpc/log"
)


func main() {
	opts := []mrpc.ServerOption{
		mrpc.WithAddress("127.0.0.1:8000"),
		mrpc.WithTimeout(time.Millisecond * 20000),
	}
	s := mrpc.NewServer(opts ...)


	mockdb := testdata.DB("mockdb")
	logger, _ := zap.NewProduction()
	zapLog := log.NewZapLogger(logger)
	srv := testdata.NewService(mockdb, zapLog) // 依赖注入思想
	if err := s.RegisterService("helloworld.Greeter", srv); err != nil {
		panic(err)
	}

	if err := s.RegisterService("helloworld.Greeter2", new(testdata.Service2)); err != nil {
		panic(err)
	}
	s.Serve()
}