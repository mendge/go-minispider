package spider

type Task struct {
	NowDepth int
	DstURL   string
}

// NewTask builds a new task
func NewTask(depth int, dstURL string) *Task {
	task := &Task{
		NowDepth: depth,
		DstURL:   dstURL,
	}
	return task
}
