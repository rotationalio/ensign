/*
Package tasks provides functionality for services to run a fixed number of workers to
conduct generic asynchronous tasks. This is an intentionally simple package to make
sure that routine, non-critical work happens in a non-blocking fashion.
*/
package tasks

import (
	"context"
	"sync"
)

// Workers in the task manager handle Tasks which can hold state and other information
// needed by the task. You can also specify a simple function to execute by using the
// TaskFunc to create a Task to provide to the task manager.
type Task interface {
	Do(context.Context)
}

// TaskFunc is an adapter to allow ordinary functions to be used as tasks.
type TaskFunc func(context.Context)

func (t TaskFunc) Do(ctx context.Context) {
	t(ctx)
}

// TaskManagers execute Tasks using a fixed number of workers that operate in their own
// go routines. The TaskManager also has a fixed task queue size, so that if there are
// more tasks added to the task manager than the queue size, back pressure is applied.
type TaskManager struct {
	sync.RWMutex
	wg      *sync.WaitGroup
	queue   chan<- *TaskHandler
	stopped bool
}

// New returns TaskManager, running the specified number of workers in their own Go
// routines and creating a queue of the specified size. The task manager is now ready
// to perform routine tasks!
func New(workers, queueSize int) *TaskManager {
	wg := &sync.WaitGroup{}
	queue := make(chan *TaskHandler, queueSize)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go TaskWorker(wg, queue)
	}

	return &TaskManager{wg: wg, queue: queue}
}

// Stop the task manager waiting for all workers to stop their tasks before returning.
func (tm *TaskManager) Stop() {
	tm.Lock()
	defer tm.Unlock()

	// Don't close the queue multiple times (avoid panic)
	if tm.stopped {
		return
	}

	close(tm.queue)
	tm.wg.Wait()
	tm.stopped = true
}

// Check if the task manager has been stopped
func (tm *TaskManager) IsStopped() bool {
	tm.RLock()
	defer tm.RUnlock()
	return tm.stopped
}

// Queue a task with the specified context. Blocks if the queue is full.
func (tm *TaskManager) QueueContext(ctx context.Context, task Task) {
	tm.queue <- &TaskHandler{task, ctx}
}

// Queue a task with a background context. Blocks if the queue is full.
func (tm *TaskManager) Queue(task Task) {
	tm.QueueContext(context.Background(), task)
}

func TaskWorker(wg *sync.WaitGroup, queue <-chan *TaskHandler) {
	defer wg.Done()
	for handler := range queue {
		handler.task.Do(handler.ctx)
	}
}

type TaskHandler struct {
	task Task
	ctx  context.Context
}
