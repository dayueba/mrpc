package main

import (
	"fmt"
	"context"

	"github.com/dayueba/mrpc/client"
	"github.com/dayueba/mrpc/plugin/jaeger"
)

type Response struct {
	Result int `mapstructure:"result" msgpack:"result"`
}

type Request struct {	
	A int `msgpack:"a"`
	B int `msgpack:"b"`
}

func main() {
	tracer, err := jaeger.Init("localhost:6831")
	if err != nil {
		panic(err)
	}

	opts := []client.Option {
		client.WithTarget("127.0.0.1:8000"),
		client.WithInterceptor(jaeger.OpenTracingClientInterceptor(tracer, "helloworld.Greeter.Add")),
	}
	c := client.DefaultClient

	req := &Request{
		A: 1,
		B: 2,
	}
	rsp := &Response{}
	
	for i := 0; i < 500; i++ {
		req.A = i
		err = c.Call(context.Background(), "helloworld.Greeter2.Add", req, rsp, opts ...)
		fmt.Printf("%+v\n", rsp)
		fmt.Println(err)	
	}
}