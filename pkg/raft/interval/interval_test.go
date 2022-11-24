package interval_test

import (
	"testing"
	"time"

	"github.com/rotationalio/ensign/pkg/raft/interval"
	"github.com/stretchr/testify/require"
)

func TestInterval(t *testing.T) {
	t.Run("Fixed", func(t *testing.T) {
		t.Parallel()

		require.False(t, (&interval.FixedInterval{}).Start(), "should not be able to start a uninitialized ticker")

		ticker := interval.NewFixed(10 * time.Millisecond)
		timeout := time.NewTimer(100 * time.Millisecond)

		for i := 0; i < 100; i++ {
			require.Equal(t, 10*time.Millisecond, ticker.GetDelay())
		}

		ticks := 0
		require.True(t, ticker.Start(), "expected ticker to start")
		require.False(t, ticker.Start(), "should not be able to start a started ticker")
	clock:
		for {
			select {
			case <-ticker.C:
				ticks++
			case <-timeout.C:
				break clock
			}
		}

		require.True(t, ticker.Stop(), "expected ticker to stop")
		require.False(t, ticker.Stop(), "should not be able to stop a stopped ticker")
		require.False(t, ticker.Interrupt(), "cannot interrupt a stopped ticker")
		require.Greater(t, ticks, 8, "expected at least 8 ticks to occur in 100 milliseconds")
		require.LessOrEqual(t, ticks, 10, "expected up to 10 ticks to occur in 100 milliseconds")

		// Should be able to restart and interrupt a ticker
		wait := time.After(5 * time.Millisecond)
		timeout = time.NewTimer(20 * time.Millisecond)
		require.True(t, ticker.Start(), "expected ticker to start")

		ticks = 0
	clock2:
		for {
			select {
			case <-wait:
				ticker.Interrupt()
			case <-ticker.C:
				ticks++
			case <-timeout.C:
				break clock2
			}
		}

		require.True(t, ticker.Stop(), "expected ticker to stop")
		require.Equal(t, 1, ticks, "expected only 1 tick after interrupt")
	})

	t.Run("Random", func(t *testing.T) {
		t.Parallel()

		require.False(t, (&interval.RandomInterval{}).Start(), "should not be able to start a uninitialized ticker")

		ticker := interval.NewRandom(5*time.Millisecond, 15*time.Millisecond)
		timeout := time.NewTimer(100 * time.Millisecond)

		var prev time.Duration
		for i := 0; i < 100; i++ {
			delay := ticker.GetDelay()
			require.NotEqual(t, prev, delay, "should return a random delay")
			prev = delay
		}

		ticks := 0
		require.True(t, ticker.Start(), "expected ticker to start")
		require.False(t, ticker.Start(), "should not be able to start a started ticker")
	clock:
		for {
			select {
			case <-ticker.C:
				ticks++
			case <-timeout.C:
				break clock
			}
		}

		require.True(t, ticker.Stop(), "expected ticker to stop")
		require.False(t, ticker.Stop(), "should not be able to stop a stopped ticker")
		require.False(t, ticker.Interrupt(), "cannot interrupt a stopped ticker")
		require.Greater(t, ticks, 6, "expected at least 6 ticks to occur in 100 milliseconds")
		require.LessOrEqual(t, ticks, 20, "expected up to 20 ticks to occur in 100 milliseconds")

		// Should be able to restart and interrupt a ticker
		wait := time.After(5 * time.Millisecond)
		timeout = time.NewTimer(22 * time.Millisecond)
		require.True(t, ticker.Start(), "expected ticker to start")

		ticks = 0
	clock2:
		for {
			select {
			case <-wait:
				require.True(t, ticker.Interrupt(), "could not interrupt ticker")
			case <-ticker.C:
				ticks++
			case <-timeout.C:
				break clock2
			}
		}

		ticker.Stop()
		require.GreaterOrEqual(t, ticks, 1, "expected at least 1 tick after interrupt")
	})

	t.Run("Jitter", func(t *testing.T) {
		t.Parallel()

		require.False(t, (&interval.JitterInterval{}).Start(), "should not be able to start a uninitialized ticker")

		ticker := interval.NewJitter(15*time.Millisecond, 3*time.Millisecond)
		timeout := time.NewTimer(100 * time.Millisecond)

		var prev time.Duration
		for i := 0; i < 100; i++ {
			delay := ticker.GetDelay()
			require.NotEqual(t, prev, delay, "should return a random delay")
			prev = delay
		}

		ticks := 0
		require.True(t, ticker.Start(), "expected ticker to start")
		require.False(t, ticker.Start(), "should not be able to start a started ticker")
	clock:
		for {
			select {
			case <-ticker.C:
				ticks++
			case <-timeout.C:
				break clock
			}
		}

		require.True(t, ticker.Stop(), "expected ticker to stop")
		require.False(t, ticker.Stop(), "should not be able to stop a stopped ticker")
		require.False(t, ticker.Interrupt(), "cannot interrupt a stopped ticker")
		require.GreaterOrEqual(t, ticks, 5, "expected at least 5 ticks to occur in 100 milliseconds")
		require.LessOrEqual(t, ticks, 20, "expected up to 20 ticks to occur in 100 milliseconds")

		// Should be able to restart and interrupt a ticker
		wait := time.After(5 * time.Millisecond)
		timeout = time.NewTimer(32 * time.Millisecond)
		require.True(t, ticker.Start(), "expected ticker to start")

		ticks = 0
	clock2:
		for {
			select {
			case <-wait:
				require.True(t, ticker.Interrupt(), "could not interrupt ticker")
			case <-ticker.C:
				ticks++
			case <-timeout.C:
				break clock2
			}
		}

		ticker.Stop()
		require.GreaterOrEqual(t, ticks, 1, "expected at least 1 tick after interrupt")
	})
}
