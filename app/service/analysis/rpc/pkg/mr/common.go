package mr

import (
	"fmt"
	"sort"
)

type KV struct {
	Key   string
	value string
	Count int64
}

type ByDESC []*KV

func (a ByDESC) Len() int { return len(a) }

func (a ByDESC) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a ByDESC) Less(i, j int) bool { return a[i].Count < a[j].Count }

type ByASC []*KV

func (a ByASC) Len() int { return len(a) }

func (a ByASC) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a ByASC) Less(i, j int) bool { return a[i].Count > a[j].Count }

func SortByASC(kvs []*KV) {
	sort.Sort(ByASC(kvs))
}

func SortByDESC(kvs []*KV) {
	sort.Sort(ByDESC(kvs))
}

func StartWordCount(files []string, filed string) ([]*KV, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no files to handle")
	}
	wc := NewWordCounter(files)
	intermediate, err := wc.Map(filed)
	if err != nil {
		return nil, err
	}

	groups := wc.Shuffle(intermediate)

	result := wc.Reduce(groups)

	return result, nil
}
