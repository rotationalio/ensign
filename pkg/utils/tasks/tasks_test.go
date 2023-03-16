package tasks_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/rotationalio/ensign/pkg/utils/tasks"
	"github.com/stretchr/testify/require"
)

func TestTasks(t *testing.T) {
	tm := tasks.New(8, 16)
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

	// Create a task that will fail
	retries := 0
	retryTask := tasks.TaskFunc(func(ctx context.Context) error {
		retries++
		if retries < 15 {
			return errors.New("retry")
		}
		return nil
	})

	// Task that hits the retry limit
	tm = tasks.New(1, 1)
	tm.Queue(retryTask, tasks.WithRetries(10), tasks.WithBackoff(&backoff.ZeroBackOff{}))
	tm.Stop()
	require.Equal(t, 10, retries)

	// Task that succeeds before the retry limit
	retries = 0
	tm = tasks.New(1, 1)
	tm.Queue(retryTask, tasks.WithRetries(20), tasks.WithBackoff(&backoff.ZeroBackOff{}))
	tm.Stop()
	require.Equal(t, 15, retries)

	// Task with a configured backoff
	retries = 0
	tm = tasks.New(1, 1)
	start := time.Now()
	tm.Queue(retryTask, tasks.WithRetries(20), tasks.WithBackoff(backoff.NewConstantBackOff(1*time.Millisecond)))
	tm.Stop()
	require.GreaterOrEqual(t, time.Since(start).Milliseconds(), int64(15))
	require.Equal(t, 15, retries)

	// Task with an expired context
	retries = 0
	tm = tasks.New(1, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()
	tm.QueueContext(ctx, retryTask, tasks.WithRetries(20), tasks.WithBackoff(backoff.NewConstantBackOff(1*time.Millisecond)))
	tm.Stop()
	require.LessOrEqual(t, retries, 15)

	// Test non-retry tasks alongside retry tasks
	retryCounts := make([]int, 10)
	queueRetryTask := func(i int) {
		t := tasks.TaskFunc(func(ctx context.Context) error {
			retryCounts[i]++
			if retryCounts[i] < 15 {
				return errors.New("retry")
			}
			return nil
		})

		tm.Queue(t, tasks.WithRetries(20), tasks.WithBackoff(&backoff.ZeroBackOff{}))
	}

	tm = tasks.New(8, 16)
	for i := 0; i < 10; i++ {
		queueRetryTask(i)
	}

	completed = 0
	for i := 0; i < 100; i++ {
		tm.Queue(tasks.TaskFunc(func(context.Context) error {
			time.Sleep(1 * time.Millisecond)
			atomic.AddInt32(&completed, 1)
			return nil
		}))
	}

	tm.Stop()
	require.Equal(t, int32(100), completed)
	for _, count := range retryCounts {
		require.Equal(t, 15, count)
	}
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
