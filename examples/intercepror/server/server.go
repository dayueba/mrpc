package main

import (
	"time"
	"context"
	"fmt"

	"github.com/dayueba/mrpc"
	"github.com/dayueba/mrpc/examples/testdata"
	"github.com/dayueba/mrpc/interceptor"
)

func main() {
	var ceps []interceptor.ServerInterceptor

	timeSpend := func(ctx context.Context, req interface{}, handler interceptor.Handler) (interface{}, error) {
		stime := time.Now()
		fmt.Printf("req %+v\n", req)
		res, err := handler(ctx, req)
		fmt.Printf("res %+v\n", res)
		fmt.Println("request spend: ", time.Now().Sub(stime))
		return res, err
	}

	ceps = append(ceps, timeSpend)

	opts := []mrpc.ServerOption{
		mrpc.WithAddress("127.0.0.1:8000"),
		mrpc.WithTimeout(time.Millisecond * 20000),
		mrpc.WithInterceptor(ceps...),
	}
	s := mrpc.NewServer(opts ...)

	if err := s.RegisterService("helloworld.Greeter2", new(testdata.Service2)); err != nil {
		panic(err)
	}
	s.Serve()
}