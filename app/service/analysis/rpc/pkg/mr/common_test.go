package mr

import (
	"fmt"
	"log"
	"testing"
)

func TestStartWordCount(t *testing.T) {
	kvs, err := StartWordCount([]string{"log1.log", "log.log", "log3.log", "log2.log"}, "level")
	if err != nil {
		log.Println(err)
		return
	}
	for _, v := range kvs {
		fmt.Println("key:", v.Key, "value:", v.Count)
	}

	fmt.Println()

	SortByASC(kvs)

	for _, v := range kvs {
		fmt.Println("key:", v.Key, "count:", v.Count)
	}
	fmt.Println()

	SortByDESC(kvs)

	for _, v := range kvs {
		fmt.Println("key:", v.Key, "count:", v.Count)
	}
}
