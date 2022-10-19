package main

import (
	"time"

	"github.com/dayueba/mrpc"
	"github.com/dayueba/mrpc/examples/testdata"
	"github.com/dayueba/mrpc/interceptor"
	"github.com/dayueba/mrpc/ratelimit"
	"github.com/dayueba/mrpc/interceptor/ratelimit/maxconcurrency"
)


func main() {
	var ceps []interceptor.ServerInterceptor

	// ceps = append(ceps, ratelimit.Server()) // 默认bbr限流
	ceps = append(ceps, ratelimit.Server(
		ratelimit.WithLimiter(
			maxconcurrency.NewLimiter(maxconcurrency.WithMaxconcurrency(1000)),
		),
	)) // 最大并发数 限流

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