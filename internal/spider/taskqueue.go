package spider

import (
	"context"
	"sync"
)

import (
	"minispider/pkg/queue"
	"minispider/pkg/semaphore"
)

const (
	weight     = 1
	initweight = 0
)

// TaskQueue is a thread-safa Queue to add and fetch tasks
type TaskQueue struct {
	dq   *queue.Deque[Task]
	muDq sync.Mutex
	// marks the working spider goroutine on task queue.
	runningCount int
	muRC         sync.Mutex
	sema         *semaphore.Weighted
}

// NewTaskDeque builds new task queue based on initial goroutine count.
func NewTaskDeque(threadCount int) *TaskQueue {
	sq := &TaskQueue{
		dq:           queue.NewDeque[Task](),
		sema:         semaphore.NewWeighted(initweight),
		runningCount: threadCount,
	}
	return sq
}

// PopTask return the popped task and a mark whether there is no more task
func (sq *TaskQueue) PopTask() (*Task, bool) {
	var v Task
	var noMoreTask bool

	// mark noMoreTask
	sq.muRC.Lock()
	sq.runningCount--
	// if all spider goroutine on task queue are blocked and task queue is empty,
	// it means that all tasks are done and no tasks will put in queue.
	noMoreTask = (sq.runningCount == 0) && (sq.dq.Len() == 0)
	sq.muRC.Unlock()
	if noMoreTask {
		return nil, true
	}
	sq.sema.Acquire(context.Background(), weight)
	sq.muRC.Lock()
	sq.runningCount++
	sq.muRC.Unlock()

	// pop task
	sq.muDq.Lock()
	v = sq.dq.PopFront()
	sq.muDq.Unlock()

	return &v, false
}

// PushTask pushes a task to task queue
func (sq *TaskQueue) PushTask(task *Task) {
	sq.muDq.Lock()
	sq.dq.PushBack(*task)
	sq.muDq.Unlock()
	sq.sema.Release(weight)
}
