package keymu_test

import (
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/rotationalio/ensign/pkg/utils/keymu"
	"github.com/stretchr/testify/require"
)

func TestMutex(t *testing.T) {
	mu := keymu.New()

	keys := 20
	iters := 10000
	out := make(chan string, iters*2)

	var wg sync.WaitGroup
	wg.Add(iters)

	for i := 0; i < iters; i++ {
		go func(rn int) {
			defer wg.Done()

			// To confirm the tests work, comment the locking and the tests should fail.
			key := strconv.Itoa(rn)
			lock := mu.Lock(key)
			defer lock.Unlock()

			out <- key + " A"
			time.Sleep(time.Microsecond)
			out <- key + " B"

		}(rand.Intn(keys))
	}

	wg.Wait()
	close(out)

	// The map should be empty now that all the work is done
	require.Equal(t, 0, mu.Len(), "expected the map to be empty")

	// Confirm the output always produced the correct sequence
	keyops := make([][]string, keys)
	for seq := range out {
		parts := strings.Fields(seq)
		require.Len(t, parts, 2)

		key, err := strconv.Atoi(parts[0])
		require.NoError(t, err, "couldn't parse the key")

		keyops[key] = append(keyops[key], parts[1])
	}

	// For every key, the sequence should be AB AB AB AB AB ...
	for key, ops := range keyops {
		for i := 0; i < len(ops); i += 2 {
			require.Equal(t, "A", ops[i], "expected A for the %dth sequence of key %s", i, key)
			require.Equal(t, "B", ops[i+1], "expected B for the %dth sequence of key %d", i+1, key)
		}
	}
}

func BenchmarkMutex(b *testing.B) {
	b.Run("Uncontested", func(b *testing.B) {
		mu := keymu.New()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mu.Lock(i).Unlock()
		}
	})
}
