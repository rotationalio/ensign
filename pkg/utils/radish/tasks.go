package radish

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/getsentry/sentry-go"
	"github.com/oklog/ulid/v2"
)

// Workers in the task manager handle Tasks which can hold state and other information
// needed by the task. You can also specify a simple function to execute by using the
// TaskFunc to create a Task to provide to the task manager.
type Task interface {
	Do(context.Context) error
}

// TaskFunc converts a function into a Task that can be queued or scheduled.
type TaskFunc func(context.Context) error

// Ensures a TaskFunc implements the Task interface.
func (f TaskFunc) Do(ctx context.Context) error {
	return f(ctx)
}

//===========================================================================
// Task Handler Implementation
//===========================================================================

type TaskHandler struct {
	id       ulid.ULID
	parent   *TaskManager
	task     Task
	ctx      context.Context
	attempts int
	retries  int
	backoff  backoff.BackOff
	timeout  time.Duration
	err      *Error
	queuedAt time.Time
}

func (tm *TaskManager) WrapTask(task Task, opts ...Option) *TaskHandler {
	// Don't re-wrap a task handler.
	if tmh, ok := task.(*TaskHandler); ok {
		return tmh
	}

	handler := &TaskHandler{
		id:       ulid.Make(),
		parent:   tm,
		task:     task,
		ctx:      context.Background(),
		err:      &Error{},
		queuedAt: time.Now().In(time.UTC),
	}

	for _, opt := range opts {
		opt(handler)
	}

	if handler.retries > 0 && handler.backoff == nil {
		handler.backoff = backoff.NewExponentialBackOff()
	}

	return handler
}

// Execute the wrapped task with the context. If the task fails, schedule the task to
// be retried using the backoff specified in the options.
func (h *TaskHandler) Exec() {
	// Create a new context for the task from the base context if a timeout is specified
	ctx := h.ctx
	if h.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(h.ctx, h.timeout)
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
	if h.attempts <= h.retries {
		// Schedule the retry be added back to the queue
		h.parent.logger.Warn().
			Err(err).
			Dict("radish", h.err.Dict()).
			Int("retries", h.retries-h.attempts).
			Msgf("%s failed, retrying", h)

		h.parent.scheduler.Delay(h.backoff.NextBackOff(), h)
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
