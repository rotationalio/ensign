package raft

import (
	api "github.com/rotationalio/ensign/pkg/raft/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/raft/election"
	"github.com/rotationalio/ensign/pkg/raft/interval"
	"github.com/rotationalio/ensign/pkg/raft/log"
	"github.com/rotationalio/ensign/pkg/raft/peers"
)

// New creates a new replica from the configuration, validating it and setting the
// replica to its initialized state. If the configuration is invalid or the replica
// cannot be correctly initialized then an error is returned.
func New(conf Config) (replica *Replica, err error) {
	if err = conf.Validate(); err != nil {
		return nil, err
	}

	// TODO: document the heartbeat and candidacy interval tick rates.
	replica = &Replica{
		conf:      conf,
		heartbeat: interval.NewFixed(conf.Tick),
		candidacy: interval.NewRandom(2*conf.Tick, 4*conf.Tick),
	}

	// TODO: set state-machine and sync from configuration
	if replica.log, err = log.New(); err != nil {
		return nil, err
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
	log       *log.Log          // state machine command log maintained by consensus
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
