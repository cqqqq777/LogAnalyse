package mr

import "sync"

type Master struct {
	mu sync.Mutex

	intermediate chan []*KV
	files        []string
}

func NewMaster(files []string) *Master {
	return &Master{mu: sync.Mutex{}, intermediate: make(chan []*KV, len(files)), files: files}
}
