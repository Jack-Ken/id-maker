# id-maker

id-maker 是使用golang开发的生成分布式Id系统，是基于美团开源的基于数据库号段算法以及雪花算法实现的Leaf分布式ID生成系统。

## 使用

有两种方式来调用接口

1. HTTP 方式
2. gRPC 方式

### HTTP方式

1、健康检查

```she
curl http://127.0.0.1:8080/v1/ping
```

2、获取ID

- 获取 tag 是 test 的 ID

```shell
curl http://127.0.0.1:8080/v1/id?tag=test
```

- 雪花ID

```shell
curl http://127.0.0.1:8080/v1/snowid
```

3、创建tag

```shell
curl.exe -i -X POST http://127.0.0.1:8080/v1/tag -H "'Content-type':'application/json'" -d  "{\`"biz_tag\`":\`"bz\`",\`"max_id\`":1,\`"step\`":100,\`"remark\`":\`"first\`"}" 
```

### gRPC方式

目录/cmd下的客户端代码样例



## 特性

- 路由注册的解耦，使用时直接注册，不需要去修改主体代码

- 全局唯一的int64型id

- 分配ID只访问内存

- 可无限横向扩展

- 依赖mysql恢复服务迅速



## 安装

- 创建数据库

```
create database id_maker;
```

- 修改配置文件

```yaml
app:
  name: "id-maker"
  mode: "dev"
  port: "8080"

log:
  level: "debug"
  filename: "../log/id-maker.log"
  max_size: 200
  max_age: 30
  max_backups: 5

mysql:
  host: "127.0.0.1"
  port: 3306
  user: "root"
  password: "123456"
  dbname: "id_maker"
  charset: "utf8mb4"
  max_open_connections: 200
  max_idle_connections: 50

## Redis未使用，可忽略
redis:
  host: "127.0.0.1"
  port: 6379
  password: "123456"
  db: 0
  pool_size: 100

grpc:
  port: "50051"
```

- 运行项目

  ```shell
      git clone https://github.com/Jack-Ken/id-maker.git
      cd id-maker/cmd
      go run main.go
  ```

## 文献

[Leaf——美团点评分布式ID生成系统](https://tech.meituan.com/2017/04/21/mt-leaf.html)

[hwholiday/gid](https://github.com/hwholiday/gid)

