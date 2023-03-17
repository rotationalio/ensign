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
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
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
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return t(ctx)
}

// TaskManagers execute Tasks using a fixed number of workers that operate in their own
// go routines. The TaskManager also has a fixed task queue size, so that if there are
// more tasks added to the task manager than the queue size, back pressure is applied.
type TaskManager struct {
	sync.RWMutex
	wg       *sync.WaitGroup
	sg       *sync.WaitGroup
	queue    chan<- *TaskHandler
	retry    chan<- *TaskHandler
	stop     chan struct{}
	shutdown *sync.WaitGroup
	stopped  bool
}

// Option allows retries and backoff to be configured for individual tasks.
type Option func(*options)

type options struct {
	retries int
	backoff backoff.BackOff
	ctx     *gin.Context
	err     error
}

// Number of retries to attempt before giving up, default 0
func WithRetries(retries int) Option {
	return func(o *options) {
		o.retries = retries
	}
}

// Backoff strategy to use when retrying, default is an exponential backoff
func WithBackoff(backoff backoff.BackOff) Option {
	return func(o *options) {
		o.backoff = backoff
	}
}

// Log an error if all retries failed under the provided context
func WithError(ctx *gin.Context, err error) Option {
	return func(o *options) {
		o.ctx = ctx
		o.err = err
	}
}

// New returns TaskManager, running the specified number of workers in their own Go
// routines and creating a queue of the specified size. The task manager is now ready
// to perform routine tasks!
func New(workers, queueSize int, retryInterval time.Duration) *TaskManager {
	wg := &sync.WaitGroup{}
	queue := make(chan *TaskHandler, queueSize)
	retry := make(chan *TaskHandler, queueSize)
	stop := make(chan struct{})

	wg.Add(1)
	go TaskScheduler(wg, queue, retry, stop, retryInterval)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go TaskWorker(wg, queue, retry)
	}

	return &TaskManager{wg: wg, queue: queue, retry: retry, stop: stop}
}

// Stop the task manager waiting for all workers to stop their tasks before returning.
func (tm *TaskManager) Stop() {
	tm.Lock()
	defer tm.Unlock()

	// Don't close the queue multiple times (avoid panic)
	if tm.stopped {
		return
	}

	// Signal the workers to stop processing new tasks
	close(tm.stop)
	tm.wg.Wait()
	close(tm.retry)
	tm.stopped = true
}

// Check if the task manager has been stopped
func (tm *TaskManager) IsStopped() bool {
	tm.RLock()
	defer tm.RUnlock()
	return tm.stopped
}

// Queue a task defined by the handler to the specified channel. If the task returns an
// error, this method queues the task for retry by sending it to the retry channel.
func (tm *TaskManager) queueTask(ctx context.Context, handler *TaskHandler, dest chan<- *TaskHandler) {
	// Wrap the task with retry logic
	retry := func(ctx context.Context) (err error) {
		if err = handler.task.Do(ctx); err == nil {
			return nil
		}

		handler.attempts++
		if handler.attempts < handler.conf.retries {
			// Send the retry task to the scheduler
			handler.retry = time.Now().Add(handler.conf.backoff.NextBackOff())
			tm.queueTask(ctx, handler, tm.retry)
			return nil
		}

		// Log the error so we know that all the retries failed
		if handler.conf.err != nil {
			log.Error().Err(err).Int("attempts", handler.attempts).Msg("task failed after retries")

			if handler.conf.ctx != nil {
				// TODO: is this a thread-safe way to create a hub for capturing exceptions?
				hub := sentrygin.GetHubFromContext(handler.conf.ctx).Clone()
				if hub != nil {
					hub.CaptureException(handler.conf.err)
				}
			}
		}

		return err
	}

	// Create a new task to send to the queue
	dest <- &TaskHandler{
		task:     TaskFunc(retry),
		conf:     handler.conf,
		attempts: handler.attempts,
		retry:    handler.retry,
		ctx:      ctx,
	}
}

// Queue a task with the specified context. Blocks if the queue is full.
func (tm *TaskManager) QueueContext(ctx context.Context, task Task, opts ...Option) {
	conf := options{
		backoff: backoff.NewExponentialBackOff(),
	}
	for _, opt := range opts {
		opt(&conf)
	}

	handler := &TaskHandler{
		task: task,
		conf: conf,
		ctx:  ctx,
	}
	tm.queueTask(ctx, handler, tm.queue)
}

// Queue a task with a background context. Blocks if the queue is full.
func (tm *TaskManager) Queue(task Task, opts ...Option) {
	tm.QueueContext(context.Background(), task, opts...)
}

// TaskScheduler runs as a separate Go routine, listening for tasks on the retry
// channel and queueing them for a worker when their backoff period has expired.
func TaskScheduler(wg *sync.WaitGroup, queue chan<- *TaskHandler, retry chan *TaskHandler, stop <-chan struct{}, interval time.Duration) {
	defer wg.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Schedule tasks until signaled to stop
	tasks := make([]*TaskHandler, 0)
	for {
		select {
		case now := <-ticker.C:
			// Wake up and queue tasks that are ready to retry
			remain := make([]*TaskHandler, 0)
			for _, task := range tasks {
				if now.After(task.retry) {
					queue <- task
				} else {
					remain = append(remain, task)
				}
			}
			tasks = remain
		case task := <-retry:
			// Receive failed tasks from the workers
			tasks = append(tasks, task)
		case <-stop:
			// Flush remaining tasks to the queue
			for _, task := range tasks {
				queue <- task
			}
			close(queue)
			return
		}
	}
}

func TaskWorker(wg *sync.WaitGroup, queue <-chan *TaskHandler, retry <-chan *TaskHandler) {
	defer wg.Done()

	for handler := range queue {
		handler.task.Do(handler.ctx)
	}
}

type TaskHandler struct {
	task     Task
	conf     options
	attempts int
	retry    time.Time
	ctx      context.Context
}
