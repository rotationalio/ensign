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
	return t(ctx)
}

// TaskManagers execute Tasks using a fixed number of workers that operate in their own
// go routines. The TaskManager also has a fixed task queue size, so that if there are
// more tasks added to the task manager than the queue size, back pressure is applied.
type TaskManager struct {
	sync.RWMutex
	wg      *sync.WaitGroup
	queue   chan<- *TaskHandler
	stop    chan struct{}
	stopped bool
}

// Option allows retries and backoff to be configured for individual tasks.
type Option func(*options)

type options struct {
	attempts int
	retries  int
	backoff  backoff.BackOff
	ctx      *gin.Context
	err      error
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
func New(workers, queueSize int) *TaskManager {
	wg := &sync.WaitGroup{}
	queue := make(chan *TaskHandler, queueSize)
	stop := make(chan struct{})

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go TaskWorker(wg, queue, stop)
	}

	return &TaskManager{wg: wg, queue: queue, stop: stop}
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
	close(tm.queue)
	tm.stopped = true
}

// Check if the task manager has been stopped
func (tm *TaskManager) IsStopped() bool {
	tm.RLock()
	defer tm.RUnlock()
	return tm.stopped
}

// Queue a task with the current retry configuration.
func (tm *TaskManager) queueRetry(ctx context.Context, task Task, conf options) {
	// Wrap the task with retry logic
	retry := func(ctx context.Context) (err error) {
		if err = task.Do(ctx); err == nil {
			return nil
		}

		conf.attempts++
		if conf.attempts < conf.retries {
			// Wait for the backoff duration before retrying
			wait := time.After(conf.backoff.NextBackOff())

			// Queue the task again to avoid blocking the current worker
			select {
			case <-ctx.Done():
				err = ctx.Err()
				break
			case <-wait:
				tm.queueRetry(ctx, task, conf)
			}
		}

		// Log the error so we know that all the retries failed
		if conf.err != nil {
			log.Error().Err(err).Int("attempts", conf.attempts).Msg("task failed after retries")

			if conf.ctx != nil {
				// TODO: is this a thread-safe way to create a hub for capturing exceptions?
				hub := sentrygin.GetHubFromContext(conf.ctx).Clone()
				if hub != nil {
					hub.CaptureException(conf.err)
				}
			}
		}

		return err
	}

	tm.queue <- &TaskHandler{task: TaskFunc(retry), ctx: ctx}
}

// Queue a task with the specified context. Blocks if the queue is full.
func (tm *TaskManager) QueueContext(ctx context.Context, task Task, opts ...Option) {
	conf := options{
		backoff: backoff.NewExponentialBackOff(),
	}
	for _, opt := range opts {
		opt(&conf)
	}

	tm.queueRetry(ctx, task, conf)
}

// Queue a task with a background context. Blocks if the queue is full.
func (tm *TaskManager) Queue(task Task, opts ...Option) {
	tm.QueueContext(context.Background(), task, opts...)
}

func TaskWorker(wg *sync.WaitGroup, queue <-chan *TaskHandler, stop <-chan struct{}) {
	defer wg.Done()

	// Handle tasks until signaled to stop
taskLoop:
	for {
		select {
		case <-stop:
			break taskLoop
		case handler := <-queue:
			handler.task.Do(handler.ctx)
		}
	}

	// At this point no new tasks are being queued, but there may be retry tasks still
	// in the queue or yet to be queued. The default case ensures that the worker will
	// eventually exit.
	for {
		select {
		case handler := <-queue:
			handler.task.Do(handler.ctx)
		default:
			return
		}
	}
}

type TaskHandler struct {
	task Task
	ctx  context.Context
}
