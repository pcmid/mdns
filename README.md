# mdns

一个支持插件的dns代理工具

## 配置

```json
{
  "addr": ":53",
  "upstream": {
    "name": "Google DNS",
    "address": "8.8.8.8:53",
    "protocol": "tcp",
    "socks5_address": "127.0.0.1:1080",
    "timeout": 6
  },
  "plugins": [
    {
      "name": "log",
      "config": {
        "log_file": ""
      }
    },
    {
      "name": "cache",
      "config": {
        "capacity": 1024,
        "MTTL": 3600
      }
    },
    {
      "name": "dispatcher",
      "config": {
        "areas": {
          "TEST": {
            "upstream": {
              "name": "114 DNS",
              "address": "114.114.114.114:53",
              "protocol": "udp",
              "socks5_address": "",
              "timeout": 6
            },
            "domain_file": "config.sample.d/domain_test.txt"
          }
        }
      }
    },
    {
      "name": "ipset",
      "config": {
        "sets": {
          "TEST": {
            "domain_file": "config.sample.d/domain_test.txt",
            "ip_file": "config.sample.d/ip_test.txt"
          }
        }
      }
    }
  ]
}

```

* addr: 监听地址
* upstream: 上游
   * name: 名称
   * address: 地址
   * protocol: 协议, 支持`udp`,`tcp`,`tcp-tls`
* plugins: 开启的插件及其配置，name为空表示不开启
---

## 插件


### log
支持简单的查询日志记录

### cache
缓存 最大缓存数量`capacity`和最小ttl `MTTL`

### dispatcher
分流器，通过匹配`domain_file`执行分流策略，选择不同的上游

### ipset
根据`domain_file`将查询到的ip插入ipset中，暂未支持`ip_file`
