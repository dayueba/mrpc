package main

import (
	"context"
	"fmt"
	"time"
	"sync"

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
	opts := []client.Option {
		client.WithTarget("127.0.0.1:8000"),
		client.WithTimeout(20000 * time.Millisecond),
	}
	c := client.DefaultClient

	var wg sync.WaitGroup
	wg.Add(1000)
	req := &Request{
		A: 1,
		B: 2,
	}

	for i := int64(0); i < 1000; i++ {
		go func() {
			rsp2 := &Response{}
			err := c.Call(context.Background(), "helloworld.Greeter2.Add", req, rsp2, opts ...)
			fmt.Printf("%+v\n", rsp2)
			fmt.Println(err)
			wg.Done()
		}()
	}

	wg.Wait()
}
