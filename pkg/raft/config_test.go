package raft_test

import (
	"testing"
	"time"

	"github.com/rotationalio/ensign/pkg/raft"
	"github.com/rotationalio/ensign/pkg/raft/peers"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Parallel()
	conf := &raft.Config{}
	require.Error(t, conf.Validate(), "empty config should be invalid")

	// Create a valid config with peers path
	conf = &raft.Config{
		ReplicaID: 13,
		Tick:      75 * time.Millisecond,
		Timeout:   25 * time.Millisecond,
		Aggregate: false,
		PeersPath: "testdata/quorum.json",
	}

	require.NoError(t, conf.Validate(), "expected valid configuration to start test")
	require.NotNil(t, conf.Quorum, "expected quorum to be loaded from disk on validation")

	// Config with quorum but no peers path should still be valid
	conf.PeersPath = ""
	require.NoError(t, conf.Validate(), "config with quorum but no peers path should be valid")

	conf.ReplicaID = 0
	require.ErrorIs(t, conf.Validate(), raft.ErrMissingReplicaID)

	conf.ReplicaID = 109
	require.ErrorIs(t, conf.Validate(), raft.ErrMissingReplica)

	conf.ReplicaID = 13
	conf.Tick = 0
	require.ErrorIs(t, conf.Validate(), raft.ErrMissingField)

	conf.Tick = 100 * time.Microsecond
	require.ErrorIs(t, conf.Validate(), raft.ErrTickTooSmall)

	conf.Tick = 1 * time.Second
	conf.Timeout = 0
	require.ErrorIs(t, conf.Validate(), raft.ErrMissingField)

	conf.Timeout = 30 * time.Second
	require.ErrorIs(t, conf.Validate(), raft.ErrTimeoutTooBig)

	conf.Timeout = 750 * time.Millisecond
	conf.Quorum.QID = 0
	require.ErrorIs(t, conf.Validate(), peers.ErrMissingQID)

	// Quorum needs either peers path or quorum
	conf.PeersPath = ""
	conf.Quorum = nil
	require.ErrorIs(t, conf.Validate(), raft.ErrMissingField)

	// Quorum needs to be able to be loaded from disk
	conf.PeersPath = "testdata/quorum.foo"
	require.Error(t, conf.Validate(), "expected error when quorum cannot be loaded from disk")
}
