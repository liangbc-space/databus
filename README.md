# mysql(binlog)监听同步数据到elasticsearch

## 服务依赖



###### 注意：考虑效率问题，以下中间件最好安装在mysql和elasticsearch的同一个局域网内安装，否则可能会因为网路请求问题导致性能降低

* 安装java jdk 须大于1.8版本

* 安装zookeeper   版本无特殊要求，不要太低

* 安装kafka		版本无特殊要求，不要太低

* 安装canal服务端	最新版1.4.4
```
    https://github.com/alibaba/canal/wiki/Canal-Kafka-RocketMQ-QuickStart

    https://github.com/alibaba/canal/wiki/aliyun-RDS-QuickStart
```

* 配置canal、kafka、zookeeper
