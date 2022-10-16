package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dayueba/mrpc/client"
	"github.com/dayueba/mrpc/interceptor"
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

	timeSpend := func(ctx context.Context, req, rsp interface{}, ivk interceptor.Invoker) error {
		stime := time.Now()
		fmt.Printf("before %+v\n", rsp)
		err := ivk(ctx, req, rsp)
		fmt.Printf("end %+v\n", rsp)
		fmt.Println("request spend: ", time.Now().Sub(stime))
		return err
	}

	ceps = append(ceps, timeSpend)

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

	rsp2 := &Response{}
	err = c.Call(context.Background(), "helloworld.Greeter2.Add", req, rsp2, opts ...)
	fmt.Printf("%+v\n", rsp2)
	fmt.Println(err)
}
