package radish_test

import (
	"context"
	"io"
	"math/rand"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/rotationalio/ensign/pkg/utils/radish"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
)

var discard = zerolog.New(io.Discard)

func TestScheduler(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping scheduler test in short mode")
		return
	}

	// Create a basic "task manager" that executes tasks until out is closed.
	var (
		wg        sync.WaitGroup
		completed uint32
	)

	out := make(chan Task)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for task := range out {
			task.Do(context.Background())
		}
		t.Log("task manager stopped")
	}()

	// Create a scheduler
	scheduler := NewScheduler(out, discard)
	require.False(t, scheduler.IsRunning(), "expected the scheduler to not be running when constructed")

	// Schedule a bunch of tasks before running the scheduler including tasks in the
	// past to ensure that all tasks are run correctly.
	for i := -5; i < 5; i++ {
		delay := time.Duration(i * 100 * int(time.Millisecond))
		scheduler.Delay(delay, Func(func(ctx context.Context) error {
			atomic.AddUint32(&completed, 1)
			t.Logf("task completed after %s delay", delay)
			return nil
		}))
	}

	// Run the scheduler in its own go routine
	scheduler.Start()
	require.True(t, scheduler.IsRunning(), "expected scheduler to be running when started")

	// Schedule more tasks, including tasks in the past while scheduler is running.
	for i := -5; i < 5; i++ {
		delay := time.Duration(i * 100 * int(time.Millisecond))
		scheduler.Delay(delay, Func(func(ctx context.Context) error {
			atomic.AddUint32(&completed, 1)
			t.Logf("task completed after %s delay", delay)
			return nil
		}))
	}

	// Schedule a final task far in the future to close the out channel
	// NOTE: the scheduler cannot be restarted after this!
	scheduler.Delay(1500*time.Millisecond, Func(func(ctx context.Context) error {
		close(out)
		t.Log("out channel closed")
		return nil
	}))

	wg.Wait()
	scheduler.Stop()
	require.False(t, scheduler.IsRunning(), "expected scheduler to be stopped")
	require.Equal(t, uint32(20), completed, "expected 20 tasks to be completed by scheduler")
}

func TestFutures(t *testing.T) {
	// Create a random time in the future.
	makeFuture := func() *Future {
		return &Future{
			Time: time.Now().Add(time.Duration(rand.Int63n(8.64e+13))),
			Task: Func(func(context.Context) error { return nil }),
		}
	}

	t.Run("RandomSort", func(t *testing.T) {
		futures := make(Futures, 0, 1000)
		for i := 0; i < 1000; i++ {
			futures = futures.Insert(makeFuture())
		}

		require.Len(t, futures, 1000)
		require.True(t, sort.IsSorted(futures))
	})

	t.Run("StableSort", func(t *testing.T) {
		// Sorted list of fixed timestamps (not random ones) with duplicates.
		timestamps := []string{
			"2023-10-14T09:36:40-05:00",
			"2023-10-14T09:37:06-05:00",
			"2023-10-14T09:39:21-05:00",
			"2023-10-14T09:39:35-05:00",
			"2023-10-14T09:39:35-05:00",
			"2023-10-14T09:40:04-05:00",
			"2023-10-14T09:40:05-05:00",
			"2023-10-14T14:40:34Z",
			"2023-10-14T14:40:34Z",
			"2023-10-14T10:40:48-04:00",
			"2023-10-14T07:41:08-07:00",
			"2023-10-14T07:41:08-07:00",
			"2023-10-14T14:41:20Z",
			"2023-10-14T09:41:34-05:00",
			"2023-10-14T22:41:53+08:00",
			"2023-10-14T16:42:18+02:00",
			"2023-10-14T09:42:34-05:00",
			"2023-10-15T01:42:55+11:00",
			"2023-10-15T01:42:55+11:00",
			"2023-10-14T14:43:12Z",
		}

		// Create a shuffled list of indexes to insert timestamps into futures in
		// a random order to ensure that the test is correct.
		index := make([]int, len(timestamps))
		for i := 0; i < len(timestamps); i++ {
			index[i] = i
		}

		rand := rand.New(rand.NewSource(time.Now().UnixNano()))
		rand.Shuffle(len(index), func(i, j int) { index[i], index[j] = index[j], index[i] })
		require.False(t, slices.IsSorted(index))

		// Create a list of futures from the timestamps.
		futures := make(Futures, 0)
		for _, i := range index {
			ts, _ := time.Parse(time.RFC3339, timestamps[i])
			futures = futures.Insert(&Future{Time: ts})
		}

		// Check that the futures are sorted correctly.
		require.Len(t, futures, len(timestamps))
		for i, f := range futures {
			ts, _ := time.Parse(time.RFC3339, timestamps[i])
			require.True(t, ts.Equal(f.Time))
		}
	})

	t.Run("GrowAndShrink", func(t *testing.T) {
		futures := make(Futures, 0)
		require.Equal(t, 0, len(futures))
		require.Equal(t, 0, cap(futures))

		// Grow futures by 7
		for i := 0; i < 7; i++ {
			futures = futures.Insert(makeFuture())
		}

		futures = futures.Resize()
		require.Equal(t, 7, len(futures))
		require.Equal(t, 16, cap(futures))

		// Shrink futures by 3
		futures = futures[3:]
		futures = futures.Resize()
		require.Equal(t, 4, len(futures))
		require.Equal(t, 16, cap(futures))

		// Grow futures by 24
		for i := 0; i < 24; i++ {
			futures = futures.Insert(makeFuture())
		}

		futures = futures.Resize()
		require.Equal(t, 28, len(futures))
		require.Equal(t, 28, cap(futures))

		// Shrink futures by 9
		futures = futures[9:]
		futures = futures.Resize()
		require.Equal(t, 19, len(futures))
		require.Equal(t, 19, cap(futures))
	})
}

func BenchmarkFutures(b *testing.B) {
	// Create a random time in the future.
	makeFuture := func() *Future {
		return &Future{
			Time: time.Now().Add(time.Duration(rand.Int63n(8.64e+13))),
			Task: Func(func(context.Context) error { return nil }),
		}
	}

	makeBenchmark := func(maxSize int) func(b *testing.B) {
		return func(b *testing.B) {
			futures := make(Futures, 0)
			futures = futures.Resize()
			b.ReportAllocs()
			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				b.StopTimer()
				task := makeFuture()
				if len(futures) > maxSize {
					futures = futures[:maxSize]
					futures = futures.Resize()
				}

				b.StartTimer()
				futures = futures.Insert(task)
			}
		}
	}

	b.Run("Small", makeBenchmark(16))
	b.Run("Medium", makeBenchmark(64))
	b.Run("Large", makeBenchmark(256))
	b.Run("XLarge", makeBenchmark(1024))
}
