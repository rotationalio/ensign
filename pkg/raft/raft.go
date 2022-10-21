package raft

import (
	"math/rand"
	"time"

	api "github.com/rotationalio/ensign/pkg/raft/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/raft/election"
	"github.com/rotationalio/ensign/pkg/raft/interval"
	"github.com/rotationalio/ensign/pkg/raft/peers"
)

// Initialize the package and random numbers, etc.
func init() {
	// Set the random seed to something different each time.
	rand.Seed(time.Now().UnixNano())
}

func New(conf Config) (replica *Replica, err error) {
	if err = conf.Validate(); err != nil {
		return nil, err
	}

	replica = &Replica{
		conf: conf,
	}

	if err = replica.setState(Initialized); err != nil {
		return nil, err
	}
	return replica, err
}

type Replica struct {
	api.UnimplementedRaftServer
	peers.Peer

	conf Config // the configuration of the local replica

	// Consensus state
	state     State             // the current state of the local replica
	leader    uint32            // the PID of the leader of the quorum
	term      uint64            // current term of the replica
	votes     election.Election // the current leader election, if any
	votedFor  uint32            // the PID of the replica we voted for in the current term
	heartbeat interval.Interval // the heartbeat ticker
	candidacy interval.Interval // the candidate timeout
}

// Election returns a new election from the replica's internal configuration with the
// replica voting yes for itself automatically (e.g. to start its candidacy).
func (r *Replica) Election() election.Election {
	peers := make([]uint32, 0, len(r.conf.Quorum.Peers))
	for _, peer := range r.conf.Quorum.Peers {
		peers = append(peers, peer.PID)
	}

	votes := election.New(peers...)
	votes.Vote(r.conf.ReplicaID, true)
	return votes
}
