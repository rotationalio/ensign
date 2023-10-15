/*
Package tasks provides functionality for services to run a fixed number of workers to
conduct generic asynchronous tasks. This is an intentionally simple package to make
sure that routine, non-critical work happens in a non-blocking fashion.
*/
package radish

import (
	"context"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog/log"
)

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

// New returns TaskManager, running the specified number of workers in their own Go
// routines and creating a queue of the specified size. The task manager is now ready
// to perform routine tasks!
func New(workers, queueSize int, retryInterval time.Duration) *TaskManager {
	wg := &sync.WaitGroup{}                     // Waits for all go routines (scheduler and workers) to stop
	queue := make(chan *TaskHandler, queueSize) // Queue sends tasks from the manager to the scheduler
	tasks := make(chan *TaskHandler, queueSize) // Tasks is used by the scheduler to send tasks to the workers (including retries)
	stop := make(chan struct{})

	wg.Add(1)
	go TaskScheduler(wg, queue, tasks, stop, retryInterval)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go TaskWorker(wg, tasks)
	}

	return &TaskManager{wg: wg, queue: queue, stop: stop}
}

// Stop the task manager waiting for all workers to stop their tasks before returning.
func (tm *TaskManager) Stop() {
	tm.Lock()

	// Don't close the queue multiple times (avoid panic)
	if tm.stopped {
		tm.Unlock()
		return
	}

	// Signal the scheduler to stop sending retry tasks to the workers
	// Note that the write lock will prevent users from queuing new tasks
	// Note also that the write lock will prevent workers from queuing retries
	tm.stop <- struct{}{}

	// Mark the task manager as stopped, unlock so we don't deadlock with workers.
	tm.stopped = true
	tm.Unlock()

	// Wait until the workers and the scheduler terminate
	tm.wg.Wait()
	close(tm.queue)
	close(tm.stop)
}

// Check if the task manager has been stopped (blocks until fully stopped).
func (tm *TaskManager) IsStopped() bool {
	tm.RLock()
	defer tm.RUnlock()
	return tm.stopped
}

// Queue a task with a background context. Blocks if the queue is full.
func (tm *TaskManager) Queue(task Task, opts ...Option) error {
	return tm.QueueContext(context.Background(), task, opts...)
}

// Queue a task with the specified context. Blocks if the queue is full.
func (tm *TaskManager) QueueContext(ctx context.Context, task Task, opts ...Option) error {
	options := makeOptions(opts...)
	handler := &TaskHandler{
		parent:    tm,
		task:      task,
		opts:      options,
		ctx:       ctx,
		err:       &Error{err: options.err},
		attempts:  0,
		retryAt:   time.Time{},
		scheduled: time.Now(),
	}
	return tm.queueTask(handler)
}

// Queue a task defined by the handler to the specified channel. If the task returns an
// error, this method queues the task for retry by sending it to the scheduler.
func (tm *TaskManager) queueTask(handler *TaskHandler) error {
	// The read-lock allows us to check tm.stopped concurrently. If Stop() has been
	// called it holds a write lock that prevents this lock from being acquired until
	// the scheduler has closed the queue.
	tm.RLock()
	defer tm.RUnlock()

	if tm.stopped {
		// Dropping the task because the task manager is not running
		log.Warn().Err(ErrTaskManagerStopped).Msg("cannot queue async task when task manager is stopped")
		return ErrTaskManagerStopped
	}

	// Queue the handler
	tm.queue <- handler
	return nil
}

type TaskHandler struct {
	parent    *TaskManager
	task      Task
	opts      *options
	ctx       context.Context
	err       *Error
	attempts  int
	retryAt   time.Time
	scheduled time.Time
}

// Execute the wrapped task with the context. If the task fails, schedule the task to
// be retried using the backoff specified in the options.
func (h *TaskHandler) Exec() {
	// Attempt to execute the task
	var err error
	if err = h.task.Do(h.ctx); err == nil {
		// Success!
		log.Debug().
			Dur("duration", time.Since(h.scheduled)).
			Int("attempts", h.attempts).
			Msg("async tasks completed")
		return
	}

	// Deal with the error
	h.attempts++
	h.err.Append(err)

	// Check if we have retries left
	if h.attempts <= h.opts.retries {
		// Schedule the retry be added back to the queue
		h.retryAt = time.Now().Add(h.opts.backoff.NextBackOff())
		log.Warn().
			Err(err).
			Time("retry_at", h.retryAt).
			Int("attempts", h.attempts).
			Int("retries_remaining", h.opts.retries-h.attempts).
			Msg("async task failed, retrying")

		h.parent.queueTask(h)
		return
	}

	// At this point we've exhausted all possible retries, so log the error.
	h.err.Since(h.scheduled)
	log.Error().Err(h.err).Dict("radish", h.err.Dict()).Msg("task failed")
	h.err.Capture(sentry.GetHubFromContext(h.ctx))
}

// TaskScheduler runs as a separate Go routine, listening for tasks on the retry
// channel and queueing them for a worker when their backoff period has expired.
func TaskScheduler(wg *sync.WaitGroup, queue <-chan *TaskHandler, tasks chan<- *TaskHandler, stop <-chan struct{}, interval time.Duration) {
	defer wg.Done()

	// Hold tasks awaiting retry and queue them every tick if ready
	// TODO: how do we test to ensure there is no memory leak?
	pending := make([]*TaskHandler, 0, 64)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case task := <-queue:
			// Check if the task is a retry task that needs to be held
			if !task.retryAt.IsZero() && time.Now().Before(task.retryAt) {
				pending = append(pending, task)
				continue
			}

			// Otherwise send the task to the worker queue immediately
			// Do not block the send; if no workers are available, append the task to
			// the pending data structure. This will reduce the backpressure on the
			// queue but also prevent deadlocks where there are more retries than
			// workers available and no one can make progress.
			select {
			case tasks <- task:
				continue
			default:
				pending = append(pending, task)
			}

		case now := <-ticker.C:
			// Do not modify pending if it contains no tasks.
			if len(pending) == 0 {
				continue
			}

			// Check all of the pending tasks to see if any are ready to be queued
			for i, task := range pending {
				if task.retryAt.IsZero() || task.retryAt.Before(now) {
					// The task is ready to retry; queue it up and delete it from pending
					// Note: this is a non-blocking write to tasks in case there are no
					// workers available to handle the current task.
					select {
					case tasks <- task:
						pending[i] = nil
					default:
						continue
					}
				}
			}

			// Prevent memory leaks by shifting tasks to deleted spots without allocation
			i := 0
			for _, task := range pending {
				if task != nil {
					pending[i] = task
					i++
				}
			}

			// Compute the new capacity, shrinking it if necessary to prevent leaks.
			newcap := cap(pending)
			if i+64 < newcap {
				newcap = i + 64
			}

			pending = pending[:i:newcap]
			log.Trace().Int("pending_length", len(pending)).Int("pending_capacity", cap(pending)).Msg("async task scheduler memory usage")

		case <-stop:
			// Flush remaining tasks to the workers
			for _, task := range pending {
				tasks <- task
			}
			close(tasks)
			return
		}
	}
}

func TaskWorker(wg *sync.WaitGroup, tasks <-chan *TaskHandler) {
	defer wg.Done()
	for handler := range tasks {
		handler.Exec()
	}
}
