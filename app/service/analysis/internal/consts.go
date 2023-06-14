package internal

const FilePrefix = "./tmp/"

const ChunkLine = 20000

const (
	StatusWait int8 = iota
	StatusSuccess
	StatusFailed
)
