## mrpc
called "msgpack rpc" or "my rpc"

## docs
- 限流：[docs](https://xjip3se76o.feishu.cn/wiki/wikcnx5mMBOXaGYIeeM0uTXriTh)

## benckmark
```
go run examples/helloworld/server/server.go
go run -v benchmark/client.go -concurrency=500 -total=1000000
INFO msg=took 13273 ms for 1000000 requests
INFO msg=sent     requests      : 1000000
INFO msg=received requests      : 1000000
INFO msg=received requests succ : 999685
INFO msg=received requests fail : 315
INFO msg=throughput  (TPS)      : 75340
```


## todo
- [ ] rpc接口异常处理
- [ ] 日志
- [ ] 配置文件格式
- [ ] 自适应限流 
- [ ] 熔断
- [ ] 负载均衡
- [ ] handlers 目前使用map存储，可以改为radix tree
- [ ] 服务注册与发现 重构
- [ ] 压测的过程中 `read: connection reset by peer`