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
	require.Equal(t, int32(100), completed)
	require.True(t, tm.IsStopped())
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
	t.Skip()

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

	// // Task that hits the retry limit
	// tm = tasks.New(1, 1, time.Millisecond)
	// tm.Queue(retryTask, tasks.WithRetries(3), tasks.WithBackoff(&backoff.ZeroBackOff{}))
	// time.Sleep(6 * time.Millisecond)
	// tm.Stop()
	// require.Equal(t, 3, retries)

	// // Task that succeeds before the retry limit
	// retries = 0
	// tm = tasks.New(1, 1, time.Millisecond)
	// tm.Queue(retryTask, tasks.WithRetries(10), tasks.WithBackoff(&backoff.ZeroBackOff{}))
	// time.Sleep(6 * time.Millisecond)
	// tm.Stop()
	// require.Equal(t, 3, retries)

	// // Task with a configured backoff
	// retries = 0
	// tm = tasks.New(1, 1, time.Millisecond)
	// tm.Queue(retryTask, tasks.WithRetries(10), tasks.WithBackoff(backoff.NewConstantBackOff(1*time.Millisecond)))
	// time.Sleep(6 * time.Millisecond)
	// tm.Stop()
	// require.Equal(t, 3, retries)

	// // Task with a canceled context
	// retries = 0
	// tm = tasks.New(1, 1, time.Millisecond)
	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Millisecond)
	// cancel()
	// tm.QueueContext(ctx, retryTask, tasks.WithRetries(10), tasks.WithBackoff(backoff.NewConstantBackOff(1*time.Millisecond)))
	// tm.Stop()
	// require.Equal(t, 0, retries)

	// // Test non-retry tasks alongside retry tasks
	// retryCounts := make([]int, 10)
	// queueRetryTask := func(i int) {
	// 	t := tasks.TaskFunc(func(ctx context.Context) error {
	// 		retryCounts[i]++
	// 		if retryCounts[i] < 5 {
	// 			return errors.New("retry")
	// 		}
	// 		return nil
	// 	})

	// 	tm.Queue(t, tasks.WithRetries(10), tasks.WithBackoff(&backoff.ZeroBackOff{}))
	// }

	// tm = tasks.New(8, 16, time.Millisecond)
	// for i := 0; i < 10; i++ {
	// 	queueRetryTask(i)
	// }

	// completed = 0
	// for i := 0; i < 100; i++ {
	// 	tm.Queue(tasks.TaskFunc(func(context.Context) error {
	// 		time.Sleep(1 * time.Millisecond)
	// 		atomic.AddInt32(&completed, 1)
	// 		return nil
	// 	}))
	// }

	// time.Sleep(10 * time.Millisecond)
	// tm.Stop()
	// require.Equal(t, int32(100), completed)
	// for _, count := range retryCounts {
	// 	require.Equal(t, 5, count)
	// }
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
