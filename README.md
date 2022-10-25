## mrpc
called "msgpack rpc" or "my rpc"

## docs
- 连接池设计: [docs](https://xjip3se76o.feishu.cn/wiki/wikcnhhKMKTjAtiv1VCFqwD7fYt)
- 限流: [docs](https://xjip3se76o.feishu.cn/wiki/wikcnx5mMBOXaGYIeeM0uTXriTh)
- 熔断: [docs](https://xjip3se76o.feishu.cn/wiki/wikcnawR2Gn782uhDUtinYUizNQ)
- 负载均衡: [docs](https://xjip3se76o.feishu.cn/wiki/wikcnP8GuEVxgNl2qfa38GnSSCb)

## benckmark
```
go run examples/helloworld/server/server.go
go run -v benchmark/client.go -concurrency=100 -total=1000000
INFO msg=took 58976 ms for 1000000 requests
INFO msg=sent     requests      : 1000000
INFO msg=received requests      : 1000000
INFO msg=received requests succ : 1000000
INFO msg=received requests fail : 0
INFO msg=throughput  (TPS)      : 16956
```


## todo
- [ ] errors
- [ ] 日志
- [ ] 配置文件格式
- [ ] 服务注册与发现 重构
- [ ] 性能优化
- [ ] examples 完善
- [ ] 超时控制完善
- [ ] 代码生成工具
