package main 

import (
	"time"

	"github.com/dayueba/mrpc"
	"github.com/dayueba/mrpc/examples/testdata"
	"github.com/dayueba/mrpc/plugin/consul"
)

func main() {
	opts := []mrpc.ServerOption{
		mrpc.WithAddress(":8000"),
		mrpc.WithTimeout(time.Millisecond * 20000),
		mrpc.WithPlugin(consul.Name),	
		mrpc.WithSelectorSvrAddr("localhost:8500"),
		mrpc.WithName("mrpc-consul-server"),
	}
	s := mrpc.NewServer(opts ...)

	if err := s.RegisterService("helloworld.Greeter2", new(testdata.Service2)); err != nil {
		panic(err)
	}
	s.Serve()
}