package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dayueba/mrpc/client"
	"github.com/dayueba/mrpc/interceptor"
	"github.com/dayueba/mrpc/breaker"
	"github.com/dayueba/mrpc/interceptor/breaker/sre"
)

type Response struct {
	Result int `mapstructure:"result" msgpack:"result"`
}

type Request struct {	
	A int `msgpack:"a"`
	B int `msgpack:"b"`
}

func main() {
	var ceps []interceptor.ClientInterceptor

	br := breaker.Client(
		breaker.WithBreaker(
			sre.NewBreaker(),
		),
	)

	ceps = append(ceps, br)

	opts := []client.Option {
		client.WithTarget("127.0.0.1:8000"),
		client.WithTimeout(20000 * time.Millisecond),
		client.WithInterceptor(ceps...),
	}
	c := client.DefaultClient

	req := &Request{
		A: 1,
		B: 2,
	}
	var err error

	for i := 0; i < 100; i++ {
		rsp2 := &Response{}
		err = c.Call(context.Background(), "helloworld.Greeter2.Breaker", req, rsp2, opts ...)
		fmt.Printf("%+v\n", rsp2)
		fmt.Println(err)
	}
}
