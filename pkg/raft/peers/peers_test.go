package peers_test

import (
	"path/filepath"
	"testing"

	"github.com/rotationalio/ensign/pkg/raft/peers"
	"github.com/stretchr/testify/require"
)

func TestQuorumValidation(t *testing.T) {
	q := &peers.Quorum{}
	require.Error(t, q.Validate(), "empty quorums should not be valid")

	// Load quorum fixture
	q, err := peers.Load("testdata/quorum.json")
	require.NoError(t, err, "could not load quorum fixture")
	require.NoError(t, q.Validate(), "expected quorum fixture to be valid")

	// Require ID
	q.QID = 0
	require.ErrorIs(t, q.Validate(), peers.ErrMissingQID)

	// Require Peers
	q.QID = 42
	cache := q.Peers
	q.Peers = nil
	require.ErrorIs(t, q.Validate(), peers.ErrNoPeers)

	// Require peers to be valid
	q.Peers = cache
	q.Peers[1].PID = 0
	require.ErrorIs(t, q.Validate(), peers.ErrMissingPID)

	// Require unique PIDs
	q.Peers[1].PID = 58
	require.ErrorIs(t, q.Validate(), peers.ErrUniquePID)
}

func TestPeerValidation(t *testing.T) {
	peer := &peers.Peer{}
	require.ErrorIs(t, peer.Validate(), peers.ErrMissingPID)

	peer.PID = 42
	peer.Name = ""
	peer.BindAddr = ":7000"
	peer.Endpoint = "localhost:7000"
	require.ErrorIs(t, peer.Validate(), peers.ErrPeerMissingField)

	peer.Name = "example"
	peer.BindAddr = ""
	peer.Endpoint = "localhost:7000"
	require.ErrorIs(t, peer.Validate(), peers.ErrPeerMissingField)

	peer.Name = "example"
	peer.BindAddr = ":7000"
	peer.Endpoint = ""
	require.ErrorIs(t, peer.Validate(), peers.ErrPeerMissingField)

	peer.Name = "example"
	peer.BindAddr = ":7000"
	peer.Endpoint = "localhost:7000"
	require.NoError(t, peer.Validate(), "expected peer to be valid")
}

func TestContains(t *testing.T) {
	// Load quorum fixture
	quorum, err := peers.Load("testdata/quorum.json")
	require.NoError(t, err, "could not load quorum fixture")
	require.NoError(t, quorum.Validate(), "expected quorum fixture to be valid")

	testCases := []struct {
		pidorname interface{}
		assert    require.BoolAssertionFunc
	}{
		{uint32(13), require.True},
		{uint32(21), require.True},
		{uint32(58), require.True},
		{"alpha", require.True},
		{"bravo", require.True},
		{"charlie", require.True},
		{uint32(109), require.False},
		{"omega", require.False},
		{[]string{}, require.False},
	}

	for _, tc := range testCases {
		tc.assert(t, quorum.Contains(tc.pidorname))
	}

}

func TestSerialization(t *testing.T) {
	_, err := peers.Load("testdata/quorum.foo")
	require.EqualError(t, err, "unknown file extension \".foo\"", "expected error for .foo extension")

	q := &peers.Quorum{}
	err = q.Dump("testdata/quorum.foo")
	require.EqualError(t, err, "unknown file extension \".foo\"", "expected error for .foo extension")

	t.Run("YAML", func(t *testing.T) {
		t.Parallel()

		quorum, err := peers.Load("testdata/quorum.yaml")
		require.NoError(t, err, "could not load quorum")
		require.Len(t, quorum.Peers, 3, "unexpected number of peers")
		require.Equal(t, uint32(42), quorum.QID, "unexpected quorum id")
		require.Equal(t, uint32(13), quorum.BootstrapLeader, "unexpected bootstrap leader")

		path := filepath.Join(t.TempDir(), "out.yaml")
		err = quorum.Dump(path)
		require.NoError(t, err, "could not dump quorum")

		other, err := peers.Load(path)
		require.NoError(t, err, "could not load dumped quorum")
		require.Equal(t, quorum, other, "quorums did not match as expected")

	})

	t.Run("JSON", func(t *testing.T) {
		t.Parallel()

		quorum, err := peers.Load("testdata/quorum.json")
		require.NoError(t, err, "could not load quorum")
		require.Len(t, quorum.Peers, 3, "unexpected number of peers")
		require.Equal(t, uint32(42), quorum.QID, "unexpected quorum id")
		require.Equal(t, uint32(13), quorum.BootstrapLeader, "unexpected bootstrap leader")

		path := filepath.Join(t.TempDir(), "out.json")
		err = quorum.Dump(path)
		require.NoError(t, err, "could not dump quorum")

		other, err := peers.Load(path)
		require.NoError(t, err, "could not load dumped quorum")
		require.Equal(t, quorum, other, "quorums did not match as expected")
	})
}
