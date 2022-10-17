package main 

import (
	"fmt"
	"context"

	"github.com/dayueba/mrpc/client"
	"github.com/dayueba/mrpc/plugin/consul"
	"github.com/dayueba/mrpc/examples/testdata"
)

func main() {
	opts := []client.Option {
		 client.WithSelectorName(consul.Name),
		 client.WithServiceName("mrpc-consul-server"),
	}
	c := client.DefaultClient
	consul.Init("localhost:8500")

	req := &testdata.Request{
		A: 1,
		B: 2,
	}
	var err error

	rsp2 := &testdata.Response{}
	err = c.Call(context.Background(), "helloworld.Greeter2.Add", req, rsp2, opts ...)
	fmt.Printf("%+v\n", rsp2)
	fmt.Println(err)
}

