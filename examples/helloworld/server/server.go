package main

import (
	"time"

	"github.com/dayueba/mrpc"
	"github.com/dayueba/mrpc/examples/testdata"
)


func main() {
	opts := []mrpc.ServerOption{
		mrpc.WithAddress("127.0.0.1:8000"),
		mrpc.WithTimeout(time.Millisecond * 20000),
	}
	s := mrpc.NewServer(opts ...)

	mockdb := testdata.DB("mockdb")
	srv := testdata.NewService(mockdb) // 依赖注入思想
	if err := s.RegisterService("helloworld.Greeter", srv); err != nil {
		panic(err)
	}

	if err := s.RegisterService("helloworld.Greeter2", new(testdata.Service2)); err != nil {
		panic(err)
	}
	s.Serve()
}