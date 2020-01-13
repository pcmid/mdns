# mdns

一个支持插件的dns代理工具

## 配置

### 全局配置

```json
{
  "Addr": ":53",
  "Upstream": {
    "Name": "Google DNS",
    "Address": "8.8.8.8:53",
    "Protocol": "tcp",
    "SOCKS5Address": "127.0.0.1:1080",
    "Timeout": 6
  },

  "PluginConfDir" : "config.sample.d/",
  "Plugins" : ["log", "cache", "dispatcher", "ipset"]
}
```

* Addr: 监听地址
* Upstream: 上游
   * Name: 名称
   * Address: 地址
   * Protocol: 协议, 支持`udp`,`tcp`,`tcp-tls`
* PluginConfDir: 插件配置目录
* Plugins: 加载的插件

