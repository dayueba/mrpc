package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dayueba/mrpc/client"
	"github.com/dayueba/mrpc/protocol"
)

type Response struct {
	Result int `mapstructure:"result" msgpack:"result"`
}

type Request struct {	
	A int `msgpack:"a"`
	B int `msgpack:"b"`
}

func main() {
	opts := []client.Option {
		client.WithTarget("127.0.0.1:8000"),
		client.WithTimeout(20000 * time.Millisecond),
	}
	c := client.DefaultClient

	req := &Request{
		A: 1,
		B: 2,
	}
	var err error
	rsp := &Response{}
	err = c.Call(context.Background(), "helloworld.Greeter.Add", req, rsp, opts ...)
	fmt.Printf("%+v\n", rsp)
	fmt.Println(err)

	// rsp2 := &Response{}
	// err = c.Call(context.Background(), "helloworld.Greeter2.Add", req, rsp2, opts ...)
	// fmt.Printf("%+v\n", rsp2)
	// fmt.Println(err)

	rsp3 := &Response{}
	err = c.Call(context.Background(), "helloworld.greeter.oops", req, rsp3, opts ...)
	// fmt.Printf("%+v, %d, %s\n", rsp3, err.(protocol.RpcError).ErrorCode, err.(protocol.RpcError).Message)
	fmt.Printf("%+v\n", rsp3)
	if e, ok := err.(*protocol.RpcError); ok {
		fmt.Printf("%+v\n", e)
	} else if err != nil{
		fmt.Println("oops, have error: ", err)
	}

}
