package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type WordCountManger struct{}

type KV struct {
	Key   string
	Value string
}

func NewWordCountManger() *WordCountManger {
	return &WordCountManger{}
}

func (w *WordCountManger) Map(inputs []string) []KV {
	logMap := make(map[string]interface{})
	kvs := make([]KV, 0)
	for _, input := range inputs {
		err := json.Unmarshal([]byte(input), &logMap)
		if err != nil {
			log.Fatalln("unmarshal log string failed err:", err)
		}
		kvs = append(kvs, KV{Key: logMap["level"].(string), Value: ""})
	}
	return kvs
}

func (w *WordCountManger) Reduce(key string, values []string) KV {
	return KV{Key: key, Value: strconv.FormatInt(int64(len(values)), 10)}
}

func (w *WordCountManger) MapReduce(path string) []KV {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Println(err)
		return nil
	}
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	inputs := make([]string, 0)
	for fileScanner.Scan() {
		inputs = append(inputs, fileScanner.Text())
	}
	intermediate := w.Map(inputs)

	groups := make(map[string][]string)
	for _, kv := range intermediate {
		groups[kv.Key] = append(groups[kv.Key], kv.Value)
	}

	var results []KV
	for key, value := range groups {
		result := w.Reduce(key, value)
		results = append(results, result)
	}

	return results
}

func main() {
	start := time.Now()
	m := NewWordCountManger()
	results := m.MapReduce("./log.log")
	results = append(results, m.MapReduce("./log1.log")...)
	results = append(results, m.MapReduce("./log2.log")...)
	results = append(results, m.MapReduce("./log3.log")...)
	for _, v := range results {
		fmt.Println("key:", v.Key, "value:", v.Value)
	}
	spend := time.Since(start)
	fmt.Println(spend.Seconds())
}
