package tasks_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/rotationalio/ensign/pkg/utils/tasks"
	"github.com/stretchr/testify/require"
)

func TestTasks(t *testing.T) {
	var wg sync.WaitGroup
	tm := tasks.New(8, 16)
	completed := 0

	for i := 0; i < 100; i++ {
		wg.Add(1)
		tm.Queue(tasks.TaskFunc(func(context.Context) {
			time.Sleep(1 * time.Millisecond)
			completed++
			wg.Done()
		}))
	}

	wg.Wait()
	tm.Stop()
	require.Equal(t, 100, completed)
}
