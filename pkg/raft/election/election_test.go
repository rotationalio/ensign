package election_test

import (
	"fmt"
	"testing"

	"github.com/rotationalio/ensign/pkg/raft/election"
	"github.com/stretchr/testify/require"
)

func TestElection(t *testing.T) {
	// Create a generic test for quorums of size 3 to size 6.
	peers := []uint32{10, 15, 20, 25, 30, 42}
	makeQuorum := func(size, majority int) func(t *testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			election := election.New(peers[:size]...)
			require.Len(t, election, size, "expected quorum size to be %d", size)
			require.Equal(t, majority, election.Majority(), "expected majority to be %d", majority)

			election.Vote(99, true)
			require.Zero(t, election.Votes(), "should not allow non-members to vote")

			for i, peer := range peers[:size] {
				i++
				election.Vote(peer, true)
				require.Equal(t, i, election.Votes())
				require.Equal(t, i >= majority, election.Passed())
			}

		}
	}

	testCases := []struct {
		size     int
		majority int
	}{
		{3, 2},
		{4, 3},
		{5, 3},
		{6, 4},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Size-%d-Quorum", tc.size), makeQuorum(tc.size, tc.majority))
	}
}
