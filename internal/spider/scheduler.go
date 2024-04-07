package spider

type Scheduler interface {
	// AddTask adds a task to task queue in a thread-safe manner
	AddTask(task *Task, taskQ *TaskQueue)

	// FetchTask fetches a task to task queue in a thread-safe manner
	FetchTask(taskQ *TaskQueue) (*Task, bool)
}

func (s *Spider) AddTask(task *Task, taskQ *TaskQueue) {
	taskQ.PushTask(task)
}

func (s *Spider) FetchTask(taskQ *TaskQueue) (*Task, bool) {
	task, noMoreTask := taskQ.PopTask()
	return task, noMoreTask
}
