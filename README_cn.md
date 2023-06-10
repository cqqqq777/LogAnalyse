# LogAnalyze

分布式日志分析系统。

## 实现

### 技术选型

- HTTP 框架: [Hertz](https://www.cloudwego.io/zh/docs/hertz/)
- RPC 框架: [Kitex](https://www.cloudwego.io/zh/docs/kitex/)
- 关系型数据库: Mysql
- 非关系型数据库: Redis
- 配置中心: [Nacos](https://nacos.io/zh-cn/docs/what-is-nacos.html)
- 服务注册中心: Nacos
- 消息队列: Kafka
- 分布式存储系统: Minio
- 链路追踪: Yaeger
- 搜索引擎: ElasticSearch

### 功能

- [ ] 注册
- [ ] 登录
- [ ] 上传与下载日志文件
- [ ] 用户提交与管理任务
- [ ] 日志收集
- [ ] 日志分析

### 服务划分

- User
- File
- Job
- Analysis

### 项目结构

![image-20230601203357017](C:\Users\86132\AppData\Roaming\Typora\typora-user-images\image-20230601203357017.png)

### 日志系统具体分析

对于日志系统的实现，大致分为以下几步：

1. 日志规范化：对日志的格式同一进行规范
2. 日志收集：提供接口由用户调用来上传日志到系统或由用户提供日志文件
   1. 客户端调用日志收集接口
   2. 将日志写入 kafka ，返回响应
   3. 消费 kafka 消息，写入 ES
3. 日志存储：将日志信息存储在 ElasticSearch 与 Minio 中以实现日志信息的备份与快速检索
4. 日志分析：基于 MapReduce 进行计算





