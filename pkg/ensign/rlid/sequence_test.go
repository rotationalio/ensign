package rlid_test

import (
	"sync"
	"testing"

	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/stretchr/testify/require"
)

func TestSequence(t *testing.T) {
	seq := rlid.Sequence(0)
	for i := uint32(1); i < 10000; i++ {
		id := seq.Next()
		require.Equal(t, i, uint32(seq))
		require.Equal(t, i, id.Sequence())
	}

	require.Equal(t, uint32(9999), uint32(seq))
}

func TestLockedSequence(t *testing.T) {
	var wg sync.WaitGroup
	nOps := 20
	nRoutines := 100

	seq := &rlid.LockedSequence{}
	for i := 0; i < nRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < nOps; i++ {
				seq.Next()
			}
		}()
	}

	wg.Wait()
	id := seq.Next()
	require.Equal(t, uint32(nRoutines*nOps+1), id.Sequence())
}
