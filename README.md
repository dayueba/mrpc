## mrpc
called "msgpack rpc" or "my rpc"

## docs
- 消息协议设计: [docs](https://xjip3se76o.feishu.cn/wiki/wikcnwqy1WgahSaz1fOnNhoHrUc)
- 连接池设计: [docs](https://xjip3se76o.feishu.cn/wiki/wikcnhhKMKTjAtiv1VCFqwD7fYt)
- 日志库设计: [docs](https://xjip3se76o.feishu.cn/wiki/wikcnLrnNKxMDe4xBotymH7HVqf)
- 错误码设计: [docs](https://xjip3se76o.feishu.cn/wiki/wikcnlVQ9KKb1mqPDiVwuZxE3pb)
- 限流: [docs](https://xjip3se76o.feishu.cn/wiki/wikcnx5mMBOXaGYIeeM0uTXriTh)
- 熔断: [docs](https://xjip3se76o.feishu.cn/wiki/wikcnawR2Gn782uhDUtinYUizNQ)
- 负载均衡: [docs](https://xjip3se76o.feishu.cn/wiki/wikcnP8GuEVxgNl2qfa38GnSSCb)

## benckmark
[压测代码仓库](https://github.com/dayueba/mrpc-benchmark)
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
**rpc框架都免不了要与 [grpc](https://github.com/grpc/grpc-go) 比一下性能**

结果不太对，不应该这么低，后续要重新压测下
```
// grpc
INFO msg=took 196855 ms for 1000000 requests
INFO msg=sent     requests      : 1000000
INFO msg=received requests      : 1000000
INFO msg=received requests succ : 1000000
INFO msg=received requests fail : 0
INFO msg=throughput  (TPS)      : 5079
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
