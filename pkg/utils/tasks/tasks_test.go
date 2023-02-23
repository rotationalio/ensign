package tasks_test

import (
	"context"
	"testing"
	"time"

	"github.com/rotationalio/ensign/pkg/utils/tasks"
	"github.com/stretchr/testify/require"
)

func TestTasks(t *testing.T) {
	tm := tasks.New(8, 16)
	completed := 0

	for i := 0; i < 100; i++ {
		tm.Queue(tasks.TaskFunc(func(context.Context) {
			time.Sleep(1 * time.Millisecond)
			completed++
		}))
	}

	tm.Stop()
	require.Equal(t, 100, completed)
}
