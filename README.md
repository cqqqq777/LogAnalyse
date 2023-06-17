# LogAnalyse

English | [中文](https://github.com/cqqqq777/LogAnalyse/blob/main/README_cn.md)

A distributed log analyze system.

## Technology Selection

- HTTP frame: [Hertz](https://www.cloudwego.io/zh/docs/hertz/)
- RPC frame: [Kitex](https://www.cloudwego.io/zh/docs/kitex/)
- Relational Database: Mysql
- Non-relational database: Redis
- The configuration center: [Nacos](https://nacos.io/zh-cn/docs/what-is-nacos.html)
- The service center: Nacos
- The message queue: Nsq
- The distributed storage system: Minio
- Traceing: Jaeger

## Function

- [x] Register
- [x] Login
- [x] Upload and download log file
- [x] Analyse log file

## Service split

- User
- File
- Job
- Analysis

## Architecture

![img_1.png](https://github.com/cqqqq777/LogAnalyse/blob/main/images/img_1.png?raw=true)

## Data Handle

1. The user submits the data processing task, the Job Service accepts the request, creates the Job, returns the response, and calls the Analysis Service asynchronously for data processing
   
    ![img.png](https://github.com/cqqqq777/LogAnalyse/blob/main/images/img.png?raw=true)

2. Analysis Service download file from Minio and split it info chunk
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

3. The system submits the split file to the MapReduce computing framework and performs the computing tasks concurrently

4. The Analysis Service writes the results of the calculation to Minio, sends a message to MQ, and modifies the Job status

## MapReduce

Based on the idea of MapReduce, a computing framework for concurrent execution of computing tasks is implemented

![img_2.png](https://github.com/cqqqq777/LogAnalyse/blob/main/images/img_2.png?raw=true)

1. Give MapReduce the files that have been shredded

2. The Map operation merges all Goroutine Data into Intermediate Data
    - Implementation: Open the file to read line by line, the log of each line will be deserialized to 'map [ string ] interface {} ' , according to the parameters passed into the generation of KV structure, and finally by the channel for communication, merged into Intermediate Data
    
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
    
      
    
3. The Shuffle stage, where data is divided into groups based on different keys and processed by the Reduce stage
   ```go
   func (w *WordCounter) Shuffle(intermediate []*KV) map[string][]string {
      groups := make(map[string][]string)
   
      for _, v := range intermediate {
         groups[v.Key] = append(groups[v.Key], v.value)
      }
   
      return groups
   }
   ```
   
4. Reduce phase, count the quantity
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
   
5. The final result is returned to Analysis Service, which writes it to the output file





