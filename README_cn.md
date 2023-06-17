# LogAnalyse

分布式日志分析系统。

## 技术选型

- HTTP 框架: [Hertz](https://www.cloudwego.io/zh/docs/hertz/)
- RPC 框架: [Kitex](https://www.cloudwego.io/zh/docs/kitex/)
- 关系型数据库: Mysql
- 非关系型数据库: Redis
- 配置中心: [Nacos](https://nacos.io/zh-cn/docs/what-is-nacos.html)
- 服务注册中心: Nacos
- 消息队列: Nsq
- 分布式存储系统: Minio
- 链路追踪: Jaeger

## 功能

- [x] 注册
- [x] 登录
- [x] 上传与下载日志文件
- [x] 用户提交与管理任务
- [x] 日志分析

## 服务划分

- User
- File
- Job
- Analysis

## 项目结构

![img_1.png](https://github.com/cqqqq777/LogAnalyse/blob/main/images/img_1.png?raw=true)

## 数据处理

1. 用户提交数据处理任务，Job Service 接受请求，创建 Job, 返回响应，异步调用 Analysis Service 进行数据处理

   ![img.png](https://github.com/cqqqq777/LogAnalyse/blob/main/images/img.png?raw=true)

2. Analysis Service 从 Minio 获取文件到本地进行切分

```go
func DownloadFile(url string, userId int64) ([]string, error) {
   // get file from minio
   resp, err := http.Get(url)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   path := GetFilePath(userId) + "task.log"
   tmp, err := os.Create(path)
   if err != nil {
      return nil, err
   }
   defer os.Remove(path)
   defer tmp.Close()

   _, err = io.Copy(tmp, resp.Body)
   if err != nil {
      return nil, err
   }

   // split file
   var wg sync.WaitGroup
   scanner := bufio.NewScanner(tmp)
   linesCh := make(chan string)
   fileCount := 1
   paths := make([]string, 0)

   path = fmt.Sprintf("%stask%d", GetFilePath(userId), fileCount)
   out, err := os.Create(path)
   if err != nil {
      return []string{path}, nil
   }
   paths = append(paths, path)

   wg.Add(1)
   go func() {
      defer wg.Done()
      for line := range linesCh {
         go possessFile(line, out, &wg)
      }
   }()

   for i := 1; scanner.Scan(); {
      linesCh <- scanner.Text()
      i++
      if i >= internal.ChunkLine {
         out.Close()
         fileCount++
         path = fmt.Sprintf("./task%d.log", fileCount)
         paths = append(paths, path)
         out, err = os.Create(path)
         if err != nil {
            log.Zlogger.Warn("create task file failed err:" + err.Error())
            return paths, nil
         }
         i = 1
      }
   }

   close(linesCh)
   out.Close()
   wg.Wait()
   return paths, nil
}
```

3. 系统将切分好的文件提交给 MapReduce 计算框架，并发执行计算任务
4. Analysis Service 将计算结果写入 Minio，向 MQ 发送消息，修改 Job 状态

## MapReduce

基于 MapReduce 的思想，实现了一个并发执行计算任务的计算框架

![img_2.png](https://github.com/cqqqq777/LogAnalyse/blob/main/images/img_2.png?raw=true)

1. 将已经切分好的文件交给 MapReduce 

2. Map 操作，将所有 Goroutine 的数据合并为 Intermediate Data

   - 具体实现：打开文件逐行读取，将每一行的日志反序列化到 `map[string]interface{}` 中，根据传入的参数生成 KV 结构体，最终由 channel 进行通信，合并为 Intermediate Data

     ```go
     func (w *WordCounter) Map(filed string) ([]*KV, error) {
        if len(w.files) == 0 {
           return nil, ErrNoFiles
        }
     
        for _, v := range w.files {
           go doMap(v, filed, w.intermediate)
        }
     
        intermediate := make([]*KV, 0, len(w.files))
     
        for i := 0; i < len(w.files); i++ {
           select {
           case kvs := <-w.intermediate:
              intermediate = append(intermediate, kvs...)
           }
        }
     
        return intermediate, nil
     }
     
     func doMap(path, filed string, intermediate chan<- []*KV) {
        file, err := os.Open(path)
        if err != nil {
           intermediate <- nil
           return
        }
        defer file.Close()
     
        fileInfo, err := file.Stat()
        if err != nil {
           intermediate <- nil
           return
        }
     
        outPut := make([]*KV, 0, fileInfo.Size()/64)
     
        scanner := bufio.NewScanner(file)
        logMap := make(map[string]interface{})
        for scanner.Scan() {
           err = json.Unmarshal([]byte(scanner.Text()), &logMap)
           if err != nil {
              continue
           }
           outPut = append(outPut, &KV{Key: logMap[filed].(string), value: ""})
        }
     
        intermediate <- outPut
     }
     ```

     

3. Shuffle 阶段，根据不同的 Key 将数据分为不同的 group ，交由 Reduce 阶段处理

   ```go
   func (w *WordCounter) Shuffle(intermediate []*KV) map[string][]string {
      groups := make(map[string][]string)
   
      for _, v := range intermediate {
         groups[v.Key] = append(groups[v.Key], v.value)
      }
   
      return groups
   }
   ```

4. Reduce 阶段，统计数量

   ```go
   func (w *WordCounter) Reduce(groups map[string][]string) []*KV {
      var wg sync.WaitGroup
      result := make([]*KV, 0, len(groups))
   
      for key, value := range groups {
         wg.Add(1)
         go func(key string, value []string) {
            defer wg.Done()
            kv := doReduce(key, value)
            w.mu.Lock()
            result = append(result, kv)
            w.mu.Unlock()
         }(key, value)
      }
      wg.Wait()
   
      return result
   }
   
   func doReduce(key string, v []string) *KV {
      return &KV{Key: key, Count: int64(len(v))}
   }
   ```

5. 最终结果返回给 Analysis Service，由其写入输出文件

