## mrpc
called "msgpack rpc" or "my rpc"

**此项目为我边学习边写的项目，虽然可用，但是有些地方实现不够优雅。随着学习了更多的知识，发现有些地方可以实现的更好，但是因为精力等原因，改不动。**

**大家完全可以按照我下面给出的资料，自己实现一个更优雅的RPC项目**

## docs
- [Go高质量编码规范](https://xjip3se76o.feishu.cn/wiki/wikcnFYQhkMwXQ22kU9IynKrrbJ)
- 设计
  - [网络IO模型](https://xjip3se76o.feishu.cn/wiki/wikcnAmELWHaChSHkC5LYx5iD3q)
  - [消息协议设计](https://xjip3se76o.feishu.cn/wiki/wikcnwqy1WgahSaz1fOnNhoHrUc)
  - [连接池设计](https://xjip3se76o.feishu.cn/wiki/wikcnhhKMKTjAtiv1VCFqwD7fYt), [多路复用连接池](https://xjip3se76o.feishu.cn/wiki/wikcndJOo4Cz85V9EWm14Z0wjff)
  - 日志库设计: [docs](https://xjip3se76o.feishu.cn/wiki/wikcnLrnNKxMDe4xBotymH7HVqf)
  - 限流: [docs](https://xjip3se76o.feishu.cn/wiki/wikcnx5mMBOXaGYIeeM0uTXriTh)
  - 熔断: [docs](https://xjip3se76o.feishu.cn/wiki/wikcnawR2Gn782uhDUtinYUizNQ)
  - 负载均衡: [docs](https://xjip3se76o.feishu.cn/wiki/wikcnP8GuEVxgNl2qfa38GnSSCb)
  - 一个TCP连接被多个请求复用以减少开销，多个请求同时发往一个TCP连接: [异步处理](https://xjip3se76o.feishu.cn/wiki/wikcnQo4knFWi8xSRLzEBZe4Kib)
  - [重试](https://xjip3se76o.feishu.cn/wiki/wikcnj3PIxXhIiyaSKl8uebGqmh)
- 使用
  - 错误码设计: [docs](https://xjip3se76o.feishu.cn/wiki/wikcnlVQ9KKb1mqPDiVwuZxE3pb)
- [踩坑记录](https://xjip3se76o.feishu.cn/wiki/wikcnGY5Tpx9Izh8xmvTKrmKI7d)
- 其它
  - [字节序](https://xjip3se76o.feishu.cn/wiki/wikcnl9d6CIJ3nWNoZXA3zpmOsb)
  - [优雅重启](https://xjip3se76o.feishu.cn/wiki/wikcnRaAUML7cCujvTRBcGnW2zc)

## benckmark
- [如何做好压测](https://xjip3se76o.feishu.cn/wiki/wikcne3GYIP9i952pURS7Vxuhhe)
- [压测代码仓库](https://github.com/dayueba/mrpc-benchmark)

## todo
- [ ] errors
- [ ] 配置文件格式
- [ ] 服务注册与发现 重构
- [ ] 性能优化
- [ ] 超时控制完善

## 学习资料

- RPC 框架
    - [net/rpc](https://pkg.go.dev/net/rpc): go 官方的rpc框架，代码量很少
    - [go-kratos](https://go-kratos.dev/docs/), [go-zero](https://go-zero.dev/cn/): 不是rpc框架，是微服务框架，在rpc基础上实现了很多服务治理功能
    - [brpc](https://github.com/apache/brpc/blob/master/README_cn.md): 虽然是c++的写的，但是不用看源码，看文档已经能学到足够多的东西。
