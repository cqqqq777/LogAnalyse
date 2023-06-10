package mr

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type WordCounter struct {
	*Master
}

var ErrNoFiles = errors.New("no files")

func NewWordCounter(files []string) *WordCounter {
	master := NewMaster(files)
	return &WordCounter{Master: master}
}

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
	scanner.Split(bufio.ScanLines)
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

func (w *WordCounter) Shuffle(intermediate []*KV) map[string][]string {
	groups := make(map[string][]string)

	for _, v := range intermediate {
		groups[v.Key] = append(groups[v.Key], v.value)
	}

	return groups
}

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
