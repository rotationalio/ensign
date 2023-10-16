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
	handler := tm.WrapTask(task, opts...)

	tm.RLock()
	defer tm.RUnlock()

	if !tm.running {
		tm.logger.Warn().Err(ErrTaskManagerStopped).Msgf("cannot queue %s", handler)
		return ErrTaskManagerStopped
	}

	tm.add <- handler
	return nil
}

// Queue a task with the specified context. Note that the context should not contain a
// deadline that might be sooner than backoff retries or the task will always fail. To
// specify a timeout for each retry, use WithTimeout. Blocks if the queue is full.
//
// Deprecated: use tm.Queue(task, WithContext(ctx)) instead.
func (tm *TaskManager) QueueContext(ctx context.Context, task Task, opts ...Option) error {
	opts = append(opts, WithContext(ctx))
	return tm.Queue(task, opts...)
}

// Delay a task to be scheduled the specified duration from now.
func (tm *TaskManager) Delay(delay time.Duration, task Task, opts ...Option) error {
	return tm.scheduler.Delay(delay, tm.WrapTask(task, opts...))
}

// Schedule a task to be executed at the specific timestamp.
func (tm *TaskManager) Schedule(at time.Time, task Task, opts ...Option) error {
	return tm.scheduler.Schedule(at, tm.WrapTask(task, opts...))
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
				queue <- tm.WrapTask(task)
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
