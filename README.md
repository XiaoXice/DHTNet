基于DHT的简易互联网辅助连接程序
===

说白了就是可以p2p的方式匿名向整个网络提供某种服务。

## 当前实现

- A主机运行程序得到一个指定的域名，eg：`-qmc-s59r-p-ruv-y-v4-pcgx-t-jqyw8-yq-k3v-lvb92-j-k-ar-z3-gk-uhfi`
- 其他位置运行本程序，透过socks5访问 `-qmc-s59r-p-ruv-y-v4-pcgx-t-jqyw8-yq-k3v-lvb92-j-k-ar-z3-gk-uhfi.qm`的流量将被转发到A主机

## TODO

- [ ] 自动网络发现服务
- [ ] 提供介绍
- [ ] 安全限制
- [ ] 流量转发服务