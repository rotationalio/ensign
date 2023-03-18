package tasks_test

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/rotationalio/ensign/pkg/utils/tasks"
	"github.com/stretchr/testify/require"
)

func TestTasks(t *testing.T) {
	// NOTE: ensure the queue size is zero so that queueing blocks until all tasks are
	// queued to prevent a race condition with the call to stop.
	tm := tasks.New(8, 0, 50*time.Millisecond)
	var completed int32

	// Queue basic tasks with no retries
	for i := 0; i < 100; i++ {
		tm.Queue(tasks.TaskFunc(func(context.Context) error {
			time.Sleep(1 * time.Millisecond)
			atomic.AddInt32(&completed, 1)
			return nil
		}))
	}

	require.False(t, tm.IsStopped())
	tm.Stop()

	// Should be able to call stop twice without panic
	tm.Stop()

	require.Equal(t, int32(100), completed)
	require.True(t, tm.IsStopped())

	// Should not be able to queue when the task manager is stopped
	err := tm.Queue(tasks.TaskFunc(func(context.Context) error { return nil }))
	require.ErrorIs(t, err, tasks.ErrTaskManagerStopped)
}

type ErroringTask struct {
	failUntil int
	attempts  int
	success   bool
	wg        *sync.WaitGroup
}

func (t *ErroringTask) Do(ctx context.Context) error {
	t.attempts++
	if t.attempts < t.failUntil {
		t.success = false
		return fmt.Errorf("task errored on attempt %d", t.attempts)
	}

	t.success = true
	t.wg.Done()
	return nil
}

func TestTasksRetry(t *testing.T) {
	// NOTE: ensure the queue size is zero so that queueing blocks until all tasks are
	// queued to prevent a race condition with the call to stop.
	tm := tasks.New(8, 0, 50*time.Millisecond)

	// Create a state of tasks that hold the number of attempts and success
	var wg sync.WaitGroup
	state := make([]*ErroringTask, 0, 100)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		state = append(state, &ErroringTask{failUntil: 3, wg: &wg})
	}

	// Queue state tasks with a retry limit that will ensure they all succeed
	for _, retryTask := range state {
		tm.Queue(retryTask, tasks.WithRetries(5), tasks.WithBackoff(&backoff.ZeroBackOff{}))
	}

	// Wait for all tasks to be completed and stop the task manager.
	wg.Wait()
	tm.Stop()

	// Analyze the results from the state
	var completed, attempts int
	for _, retryTask := range state {
		attempts += retryTask.attempts
		if retryTask.success {
			completed++
		}
	}

	require.Equal(t, 100, completed, "expected all tasks to have been completed")
	require.Equal(t, 300, attempts, "expected all tasks to have failed twice before success")
}

func TestTasksRetryFailure(t *testing.T) {
	// NOTE: ensure the queue size is zero so that queueing blocks until all tasks are
	// queued to prevent a race condition with the call to stop.
	tm := tasks.New(20, 0, 50*time.Millisecond)

	// Create a state of tasks that hold the number of attempts and success
	var wg sync.WaitGroup
	state := make([]*ErroringTask, 0, 100)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		state = append(state, &ErroringTask{failUntil: 5, wg: &wg})
	}

	// Queue state tasks with a retry limit that will ensure they all fail
	for _, retryTask := range state {
		tm.Queue(retryTask, tasks.WithRetries(1), tasks.WithBackoff(&backoff.ZeroBackOff{}))
	}

	// Wait for all tasks to be completed and stop the task manager.
	time.Sleep(500 * time.Millisecond)
	tm.Stop()

	// Analyze the results from the state
	var completed, attempts int
	for _, retryTask := range state {
		attempts += retryTask.attempts
		if retryTask.success {
			completed++
		}
	}

	require.Equal(t, 0, completed, "expected all tasks to have failed")
	require.Equal(t, 200, attempts, "expected all tasks to have failed twice before no more retries")
}

func TestTasksRetryBackoff(t *testing.T) {
	// NOTE: ensure the queue size is zero so that queueing blocks until all tasks are
	// queued to prevent a race condition with the call to stop.
	tm := tasks.New(20, 0, 5*time.Millisecond)

	// Create a state of tasks that hold the number of attempts and success
	var wg sync.WaitGroup
	state := make([]*ErroringTask, 0, 100)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		state = append(state, &ErroringTask{failUntil: 3, wg: &wg})
	}

	// Queue state tasks with a retry limit that will ensure they all succeed
	for _, retryTask := range state {
		tm.Queue(retryTask, tasks.WithRetries(5), tasks.WithBackoff(backoff.NewConstantBackOff(10*time.Millisecond)))
	}

	// Wait for all tasks to be completed and stop the task manager.
	wg.Wait()
	tm.Stop()

	// Analyze the results from the state
	var completed, attempts int
	for _, retryTask := range state {
		attempts += retryTask.attempts
		if retryTask.success {
			completed++
		}
	}

	// TODO: how to check if backoff was respected?
	require.Equal(t, 100, completed, "expected all tasks to have been completed")
	require.Equal(t, 300, attempts, "expected all tasks to have failed twice before success")
}

func TestTasksRetryContextCanceled(t *testing.T) {
	// NOTE: ensure the queue size is zero so that queueing blocks until all tasks are
	// queued to prevent a race condition with the call to stop.
	tm := tasks.New(20, 0, 50*time.Millisecond)
	var completed, attempts int32

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Queue tasks that are getting canceled
	for i := 0; i < 100; i++ {
		tm.QueueContext(ctx, tasks.TaskFunc(func(ctx context.Context) error {
			atomic.AddInt32(&attempts, 1)
			if err := ctx.Err(); err != nil {
				return err
			}

			atomic.AddInt32(&completed, 1)
			return nil
		}), tasks.WithRetries(1), tasks.WithBackoff(&backoff.ZeroBackOff{}))
	}

	// Wait for all tasks to be completed and stop the task manager.
	time.Sleep(500 * time.Millisecond)
	tm.Stop()

	require.Equal(t, int32(0), completed, "expected all tasks to have been canceled")
	require.Equal(t, int32(200), attempts, "expected all tasks to have failed twice before no more retries")
}

func TestTasksRetrySuccessAndFailure(t *testing.T) {
	// Test non-retry tasks alongside retry tasks
	// NOTE: ensure the queue size is zero so that queueing blocks until all tasks are
	// queued to prevent a race condition with the call to stop.
	tm := tasks.New(20, 0, 50*time.Millisecond)

	// Create a state of tasks that hold the number of attempts and success
	var wg sync.WaitGroup
	state := make([]*ErroringTask, 0, 100)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		state = append(state, &ErroringTask{failUntil: 2, wg: &wg})
	}

	// Queue state tasks with a retry limit that will ensure they all fail
	// First 50 have a retry, second 50 do not.
	for i, retryTask := range state {
		if i < 50 {
			tm.Queue(retryTask, tasks.WithRetries(2), tasks.WithBackoff(&backoff.ZeroBackOff{}))
		} else {
			tm.Queue(retryTask, tasks.WithBackoff(&backoff.ZeroBackOff{}))
		}
	}

	// Wait for all tasks to be completed and stop the task manager.
	time.Sleep(500 * time.Millisecond)
	tm.Stop()

	// Analyze the results from the state
	var completed, attempts int
	for _, retryTask := range state {
		attempts += retryTask.attempts
		if retryTask.success {
			completed++
		}
	}

	require.Equal(t, 50, completed, "expected all tasks to have failed")
	require.Equal(t, 150, attempts, "expected all tasks to have failed twice before no more retries")
}

func TestQueue(t *testing.T) {
	// A simple test to ensure that tm.Stop() will wait until all items in the queue are finished.
	var wg sync.WaitGroup
	queue := make(chan int32, 64)
	var final int32

	wg.Add(1)
	go func() {
		for num := range queue {
			time.Sleep(1 * time.Millisecond)
			atomic.SwapInt32(&final, num)
		}
		wg.Done()
	}()

	for i := int32(1); i < 101; i++ {
		queue <- i
	}

	close(queue)
	wg.Wait()
	require.Equal(t, int32(100), final)
}
