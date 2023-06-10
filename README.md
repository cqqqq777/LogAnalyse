# LogAnalyze
A distributed log analyze system.

## Project Implementation

### Technology Selection

- HTTP frame: [Hertz](https://www.cloudwego.io/zh/docs/hertz/)
- RPC frame: [Kitex](https://www.cloudwego.io/zh/docs/kitex/)
- Relational Database: Mysql
- Non-relational database: Redis
- The configuration center: [Nacos](https://nacos.io/zh-cn/docs/what-is-nacos.html)
- The service center: Nacos
- The message queue: Kafka
- The distributed storage system: Ceph
- Traceing: Yaeger
- The search engine: ElasticSearch

### Function

- [ ] Register
- [ ] Login
- [ ] Upload and download log file
- [ ] Analyze log file

### Service split

- User
- File
- Job
- Analysis

### Architecture

![image-20230601203357017](C:\Users\86132\AppData\Roaming\Typora\typora-user-images\image-20230601203357017.png)

