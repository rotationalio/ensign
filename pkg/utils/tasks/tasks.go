/*
Package tasks provides functionality for services to run a fixed number of workers to
conduct generic asynchronous tasks. This is an intentionally simple package to make
sure that routine, non-critical work happens in a non-blocking fashion.
*/
package tasks

import (
	"context"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/rs/zerolog/log"
)

// Workers in the task manager handle Tasks which can hold state and other information
// needed by the task. You can also specify a simple function to execute by using the
// TaskFunc to create a Task to provide to the task manager.
type Task interface {
	Do(context.Context) error
}

// TaskFunc is an adapter to allow ordinary functions to be used as tasks.
type TaskFunc func(context.Context) error

func (t TaskFunc) Do(ctx context.Context) error {
	return t(ctx)
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

// Option allows retries and backoff to be configured for individual tasks.
type Option func(*options)

type options struct {
	Retries int
	Backoff backoff.BackOff
	err     error
}

// Number of retries to attempt before giving up, default 0
func WithRetries(retries int) Option {
	return func(o *options) {
		o.Retries = retries
	}
}

// Backoff strategy to use when retrying, default is no backoff
func WithBackoff(backoff backoff.BackOff) Option {
	return func(o *options) {
		o.Backoff = backoff
	}
}

// Log an error if all the retries fail, by default nothing is logged
func WithError(err error) Option {
	return func(o *options) {
		o.err = err
	}
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
func (tm *TaskManager) QueueContext(ctx context.Context, task Task, opts ...Option) {
	conf := options{
		Backoff: &backoff.ZeroBackOff{},
	}
	for _, opt := range opts {
		opt(&conf)
	}

	// Wrap the task with retry logic
	retry := func(ctx context.Context) (err error) {
		var retries int
	retryLoop:
		for {
			if err = task.Do(ctx); err == nil {
				return nil
			}

			retries++
			if retries >= conf.Retries {
				break retryLoop
			}

			// Wait for the backoff duration before retrying
			// Note: This blocks the worker thread so queue sizes should be large
			// enough to avoid blocking new tasks from being queued.
			wait := time.After(conf.Backoff.NextBackOff())

			select {
			case <-ctx.Done():
				err = ctx.Err()
				break retryLoop
			case <-wait:
			}
		}

		// Log the error so we know that all the retries failed
		if conf.err != nil {
			log.Error().Err(err).Int("retries", retries).Msg("task failed after retries")
		}

		return err
	}

	// Queue the task
	tm.queue <- &TaskHandler{TaskFunc(retry), ctx}
}

// Queue a task with a background context. Blocks if the queue is full.
func (tm *TaskManager) Queue(task Task, opts ...Option) {
	tm.QueueContext(context.Background(), task, opts...)
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
