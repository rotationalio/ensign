package raft_test

import (
	"testing"
	"time"

	"github.com/rotationalio/ensign/pkg/raft"
	"github.com/stretchr/testify/suite"
)

type raftTestSuite struct {
	suite.Suite
	replica *raft.Replica
}

func (s *raftTestSuite) SetupSuite() {
	var err error
	require := s.Require()

	conf := raft.Config{
		ReplicaID: 58,
		Tick:      75 * time.Millisecond,
		Timeout:   50 * time.Millisecond,
		Aggregate: false,
		PeersPath: "testdata/quorum.json",
	}

	s.replica, err = raft.New(conf)
	require.NoError(err, "could not create new raft replica")
}

func (s *raftTestSuite) TearDownSuite() {

}

func TestRaft(t *testing.T) {
	suite.Run(t, new(raftTestSuite))
}

func (s *raftTestSuite) TestElection() {
	require := s.Require()
	votes := s.replica.Election()
	require.Len(votes, 3)
	require.Equal(1, votes.Votes(), "expected only one vote for the leader")
}
