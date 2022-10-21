package raft

import (
	"math/rand"
	"time"

	api "github.com/rotationalio/ensign/pkg/raft/api/v1beta1"
)

// Initialize the package and random numbers, etc.
func init() {
	// Set the random seed to something different each time.
	rand.Seed(time.Now().UnixNano())
}

func New(conf Config) (replica *Replica, err error) {
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
	Peer

	conf Config // the configuration of the local replica

	// Consensus state
	state     State    // the current state of the local replica
	leader    string   // the name of the leader of the quorum
	term      uint64   // current term of the replica
	votes     Election // the current leader election, if any
	votedFor  string   // the peer we voted for in the current term
	heartbeat Interval // the heartbeat ticker
	candidacy Interval // the candidate timeout
}
