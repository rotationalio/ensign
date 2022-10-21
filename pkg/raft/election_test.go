package raft_test

import (
	"fmt"
	"testing"

	"github.com/rotationalio/ensign/pkg/raft"
	"github.com/stretchr/testify/require"
)

func TestElection(t *testing.T) {
	peers := []string{"Quintus", "Gaius", "Fabius", "Julius", "Marcus", "Pontius"}
	makeQuorum := func(size, majority int) func(t *testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			election := raft.NewElection(peers[:size]...)
			require.Len(t, election, size, "expected quorum size to be %d", size)
			require.Equal(t, majority, election.Majority(), "expected majority to be %d", majority)

			election.Vote("Lucifer", true)
			require.Zero(t, election.Votes(), "should not allow non-members to vote")

			for i, peer := range peers[:size] {
				i++
				election.Vote(peer, true)
				require.Equal(t, i, election.Votes())
				require.Equal(t, i >= majority, election.Passed())
			}

		}
	}

	for i := 3; i < 7; i++ {
		t.Run(fmt.Sprintf("Size%dQuorum", i), makeQuorum(i, i/2+1))
	}
}
