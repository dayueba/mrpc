package main

import (
	"net/http"

	"github.com/dayueba/mrpc"
	"github.com/dayueba/mrpc/plugin/jaeger"
	"github.com/dayueba/mrpc/examples/testdata"

	_ "net/http/pprof"
)


func main() {
	pprof()

	opts := []mrpc.ServerOption{
		mrpc.WithAddress(":8000"),
		mrpc.WithTracingSvrAddr("localhost:6831"),
		mrpc.WithTracingSpanName("helloworld.Greeter"),
		mrpc.WithPlugin(jaeger.Name),
	}
	s := mrpc.NewServer(opts ...)
	if err := s.RegisterService("helloworld.Greeter2", new(testdata.Service2)); err != nil {
		panic(err)
	}
	s.Serve()
}

func pprof() {
	go func() {
		http.ListenAndServe("0.0.0.0:8899", http.DefaultServeMux)
	}()
}