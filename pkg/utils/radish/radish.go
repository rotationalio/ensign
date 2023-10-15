/*
Package tasks provides functionality for services to run a fixed number of workers to
conduct generic asynchronous tasks. This is an intentionally simple package to make
sure that routine, non-critical work happens in a non-blocking fashion.
*/
package radish

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// TaskManagers execute Tasks using a fixed number of workers that operate in their own
// go routines. The TaskManager also has a fixed task queue size, so that if there are
// more tasks added to the task manager than the queue size, back pressure is applied.
type TaskManager struct {
	sync.RWMutex
	conf      Config
	logger    zerolog.Logger
	scheduler *Scheduler
	wg        *sync.WaitGroup
	add       chan Task
	stop      chan struct{}
	running   bool
}

// Create a new task manager with the specified configuration.
func New(conf Config) *TaskManager {
	if conf.IsZero() {
		conf.Workers = 4
		conf.QueueSize = 64
		conf.ServerName = "radish"
	}

	add := make(chan Task, conf.QueueSize)
	logger := log.With().Str("task_manager", conf.ServerName).Logger()
	scheduler := NewScheduler(add, logger)

	return &TaskManager{
		conf:      conf,
		logger:    logger,
		scheduler: scheduler,
		wg:        &sync.WaitGroup{},
		add:       add,
		stop:      make(chan struct{}, 1),
		running:   false,
	}
}

// Queue a task to be executed asynchronously as soon as a worker is available. Options
// can be specified to influence the handling of the task. Blocks if queue is full.
func (tm *TaskManager) Queue(task Task, opts ...Option) error {
	return tm.QueueContext(context.Background(), task, opts...)
}

// Queue a task with the specified context. Note that the context should not contain a
// deadline that might be sooner than backoff retries or the task will always fail. To
// specify a timeout for each retry, use WithTimeout. Blocks if the queue is full.
func (tm *TaskManager) QueueContext(ctx context.Context, task Task, opts ...Option) error {
	options := makeOptions(opts...)
	handler := &TaskHandler{
		id:       ulid.Make(),
		parent:   tm,
		task:     task,
		opts:     options,
		ctx:      ctx,
		err:      options.err,
		queuedAt: time.Now().In(time.UTC),
	}

	tm.RLock()
	defer tm.RUnlock()
	if !tm.running {
		tm.logger.Warn().Err(ErrTaskManagerStopped).Msgf("cannot queue %s", handler)
		return ErrTaskManagerStopped
	}

	tm.add <- handler
	return nil
}

// Start the task manager and scheduler in their own go routines (no-op if already started)
func (tm *TaskManager) Start() {
	tm.Lock()
	defer tm.Unlock()

	// Start the scheduler (also a no-op if already started)
	tm.scheduler.Start(tm.wg)

	if tm.running {
		return
	}

	tm.running = true
	go tm.run()
}

func (tm *TaskManager) run() {
	tm.wg.Add(1)
	defer tm.wg.Done()
	tm.logger.Info().Int("workers", tm.conf.Workers).Int("queue_size", tm.conf.QueueSize).Msg("task manager running")

	queue := make(chan *TaskHandler, tm.conf.QueueSize)
	for i := 0; i < tm.conf.Workers; i++ {
		tm.wg.Add(1)
		go worker(tm.wg, queue)
	}

	for {
		select {
		case task := <-tm.add:
			if handler, ok := task.(*TaskHandler); ok {
				queue <- handler
			} else {
				queue <- &TaskHandler{
					id:       ulid.Make(),
					parent:   tm,
					task:     task,
					opts:     makeOptions(),
					ctx:      context.Background(),
					err:      &Error{},
					queuedAt: time.Now().In(time.UTC),
				}
			}

		case <-tm.stop:
			close(queue)
			tm.logger.Info().Msg("task manager stopped")
			return
		}
	}
}

func worker(wg *sync.WaitGroup, tasks <-chan *TaskHandler) {
	defer wg.Done()
	for handler := range tasks {
		handler.Exec()
	}
}

// Stop the task manager and scheduler if running (otherwise a no-op). This method
// blocks until all pending tasks have been completed, however future scheduled tasks
// will likely be dropped and not scheduled for execution.
func (tm *TaskManager) Stop() {
	tm.Lock()

	// Stop the scheduler (also a no-op if already stopped)
	tm.scheduler.Stop()

	if tm.running {
		// Send the stop signal to the task manager
		tm.stop <- struct{}{}
		tm.running = false

		tm.Unlock()

		// Wait for all tasks to be completed and workers closed
		tm.wg.Wait()
	} else {
		tm.Unlock()
	}
}

func (tm *TaskManager) IsRunning() bool {
	tm.RLock()
	defer tm.RUnlock()
	return tm.running
}

type TaskHandler struct {
	id       ulid.ULID
	parent   *TaskManager
	task     Task
	opts     *options
	ctx      context.Context
	err      *Error
	queuedAt time.Time
	attempts int
}

// Execute the wrapped task with the context. If the task fails, schedule the task to
// be retried using the backoff specified in the options.
func (h *TaskHandler) Exec() {
	// Create a new context for the task from the base context if a timeout is specified
	ctx := h.ctx
	if h.opts.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(h.ctx, h.opts.timeout)
		defer cancel()
	}

	// Attempt to execute the task
	var err error
	if err = h.task.Do(ctx); err == nil {
		// Success!
		h.parent.logger.Debug().
			Str("task_id", h.id.String()).
			Dur("duration", time.Since(h.queuedAt)).
			Int("attempts", h.attempts+1).
			Msgf("%s completed", h)
		return
	}

	// Deal with the error
	h.attempts++
	h.err.Append(err)
	h.err.Since(h.queuedAt)

	// Check if we have retries left
	if h.attempts <= h.opts.retries {
		// Schedule the retry be added back to the queue
		h.parent.logger.Warn().
			Err(err).
			Dict("radish", h.err.Dict()).
			Int("retries", h.opts.retries-h.attempts).
			Msgf("%s failed, retrying", h)

		h.parent.scheduler.Delay(h.opts.backoff.NextBackOff(), h)
		return
	}

	// At this point we've exhausted all possible retries, so log the error.
	h.parent.logger.Error().Err(h.err).Dict("radish", h.err.Dict()).Msgf("%s failed", h)
	h.err.Capture(sentry.GetHubFromContext(h.ctx))
}

// TaskHandler implements Task so that it can be scheduled, but it should never be
// called as a Task rather than a Handler (to avoid re-wrapping) so this method simply
// panics if called -- it is a developer error.
func (h *TaskHandler) Do(context.Context) error {
	panic("a task handler should not wrap another task handler")
}

// String implements fmt.Stringer and checks if the underlying task does as well; if so
// the task name is fetched from the task stringer, otherwise a default name is returned.
func (h *TaskHandler) String() string {
	if s, ok := h.task.(fmt.Stringer); ok {
		return s.String()
	}
	return "async task"
}
