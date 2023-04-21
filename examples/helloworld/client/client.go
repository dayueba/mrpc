package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dayueba/mrpc/client"
)

type Response struct {
	Result int `mapstructure:"result" msgpack:"result"`
}

type Request struct {
	A int `msgpack:"a"`
	B int `msgpack:"b"`
}

func main() {
	opts := []client.Option{
		client.WithTarget("127.0.0.1:8000"),
		client.WithTimeout(20000 * time.Millisecond),
	}
	c := client.DefaultClient

	req := &Request{
		A: 1,
		B: 2,
	}
	var err error
	var success int64
	var fail int64

	var wg sync.WaitGroup
	wg.Add(10000)
	for i := 0; i < 10000; i++ {
		go func() {
			defer wg.Done()
			rsp := &Response{}
			err = c.Call(context.Background(), "helloworld.Greeter.Add", req, rsp, opts...)
			//fmt.Printf("%+v\n", rsp)
			if err != nil {
				//fmt.Println("rsp: ", err)
				atomic.AddInt64(&fail, 1)
			} else {
				atomic.AddInt64(&success, 1)
			}
			//time.Sleep(time.Second * 10)
		}()
	}
	wg.Wait()

	fmt.Println("success: ", success, "fail: ", fail)

	// rsp2 := &Response{}
	// err = c.Call(context.Background(), "helloworld.Greeter2.Add", req, rsp2, opts ...)
	// fmt.Printf("%+v\n", rsp2)
	// fmt.Println(err)

	//rsp3 := &Response{}
	//err = c.Call(context.Background(), "helloworld.greeter.oops", req, rsp3, opts ...)
	//// fmt.Printf("%+v, %d, %s\n", rsp3, err.(protocol.RpcError).ErrorCode, err.(protocol.RpcError).Message)
	//fmt.Printf("%+v\n", rsp3)
	//if e, ok := err.(*protocol.RpcError); ok {
	//	fmt.Printf("%+v\n", e)
	//} else if err != nil{
	//	fmt.Println("oops, have error: ", err)
	//}

}
