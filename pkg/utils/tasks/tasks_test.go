package tasks_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rotationalio/ensign/pkg/utils/tasks"
	"github.com/stretchr/testify/require"
)

func TestTasks(t *testing.T) {
	tm := tasks.New(8, 16)
	var completed int32

	for i := 0; i < 100; i++ {
		tm.Queue(tasks.TaskFunc(func(context.Context) {
			time.Sleep(1 * time.Millisecond)
			atomic.AddInt32(&completed, 1)
		}))
	}

	tm.Stop()
	require.Equal(t, int32(100), completed)
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
