package main

import (
	"context"
	"flag"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dayueba/mrpc/client"
	"github.com/dayueba/mrpc/examples/testdata"
	"github.com/dayueba/mrpc/log"
)

var concurrency = flag.Int64("concurrency", 100, "concurrency")
var target = flag.String("target", "127.0.0.1:8000", "target")
var total = flag.Int64("total", 1000000, "total requests")

func main() {
	flag.Parse()
	request(*total, *concurrency, *target)
	// requestOnce(*target)
}

type Counter struct {
	Succ        int64 // 成功量
	Fail        int64 // 失败量
	Total       int64 // 总量
	Concurrency int64 // 并发量
	Cost        int64 // 总耗时 ms
}

func request(totalReqs int64, concurrency int64, target string) {
	perClientReqs := totalReqs / concurrency

	counter := &Counter{
		Total:       perClientReqs * concurrency,
		Concurrency: concurrency,
	}

	opts := []client.Option{
		client.WithTarget(target),
		//  client.WithTimeout(2 * time.Second),
	}
	c := client.DefaultClient
	req := &testdata.HelloRequest{
		Msg: "hello",
	}

	var wg sync.WaitGroup
	wg.Add(int(concurrency))

	startTime := time.Now().UnixNano()

	for i := int64(0); i < counter.Concurrency; i++ {
		go func(i int64) {
			for j := int64(0); j < perClientReqs; j++ {
				rsp := &testdata.CountResponse{}
				err := c.Call(context.Background(), "helloworld.greeter.count", req, rsp, opts...)
				fmt.Println(rsp, err)
				if err == nil {
					// log.Infof("req success : %v", rsp)
					atomic.AddInt64(&counter.Succ, 1)
				} else {
					// log.Infof("rsp fail : %v", err)
					atomic.AddInt64(&counter.Fail, 1)
				}
			}

			wg.Done()
		}(i)
	}

	wg.Wait()

	counter.Cost = (time.Now().UnixNano() - startTime) / 1000000

	log.Infof("took %d ms for %d requests", counter.Cost, counter.Total)
	log.Infof("sent     requests      : %d\n", counter.Total)
	log.Infof("received requests      : %d\n", atomic.LoadInt64(&counter.Succ)+atomic.LoadInt64(&counter.Fail))
	log.Infof("received requests succ : %d\n", atomic.LoadInt64(&counter.Succ))
	log.Infof("received requests fail : %d\n", atomic.LoadInt64(&counter.Fail))
	log.Infof("throughput  (TPS)      : %d\n", totalReqs*1000/counter.Cost)
}

func requestOnce(target string) {
	opts := []client.Option{
		client.WithTarget(target),
		client.WithTimeout(2 * time.Second),
	}
	c := client.DefaultClient
	req := &testdata.HelloRequest{
		Msg: "hello",
	}
	rsp := &testdata.CountResponse{}
	err := c.Call(context.Background(), "helloworld.greeter.count", req, rsp, opts...)
	fmt.Println(rsp, err)
}
