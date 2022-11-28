package raft

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

// Raft server states (not part of the state machine)
const (
	Stopped State = iota // stopped should be the zero value and default
	Initialized
	Running
	Follower
	Candidate
	Leader
)

// Names of the states for serialization
var stateStrings = [...]string{
	"stopped", "initialized", "running", "follower", "candidate", "leader",
}

//===========================================================================
// State Enumeration
//===========================================================================

// State is an enumeration of the possible status of a replica.
type State uint8

// String returns a human readable representation of the state.
func (s State) String() string {
	return stateStrings[s]
}

//===========================================================================
// State Transitions
//===========================================================================

// State returns the current state of the replica for testing.
func (r *Replica) State() State {
	return r.state
}

// SetState updates the state of the local replica, performing any actions
// related to multiple states, modifying internal private variables as
// needed and calling the correct internal state setting function.
//
// NOTE: These methods are not thread-safe.
func (r *Replica) setState(state State) (err error) {
	switch state {
	case Stopped:
		err = r.setStoppedState()
	case Initialized:
		err = r.setInitializedState()
	case Running:
		err = r.setRunningState()
	case Follower:
		err = r.setFollowerState()
	case Candidate:
		err = r.setCandidateState()
	case Leader:
		err = r.setLeaderState()
	default:
		err = fmt.Errorf("unknown state '%s'", state)
	}

	if err == nil {
		r.state = state
	}

	return err
}

// Stops all timers that might be running.
func (r *Replica) setStoppedState() error {
	if r.state == Stopped {
		log.Trace().Str("replica", r.Name).Msg("already stopped")
		return nil
	}

	r.heartbeat.Stop()
	r.candidacy.Stop()
	log.Debug().Str("replica", r.Name).Msg("replica stopped")
	return nil
}

// Resets any volatile variables on the local replica and is called when the replica
// becomes a follower or a candidate.
func (r *Replica) setInitializedState() error {
	r.votes = nil
	r.votedFor = 0

	// TODO: reset nextIndex and matchIndex on remotes

	log.Debug().Str("replica", r.Name).Msg("replica initialized")
	return nil
}

// Should be called once after initialization to bootstrap the quorum by starting the
// leader's heartbeat or starting the election timeout for all other replicas.
func (r *Replica) setRunningState() error {
	if r.state != Initialized {
		return ErrCannotSetRunningState
	}

	// TODO: determine leader PID from quorum
	// if r.conf.Leader == r.Name {
	// 	return r.setLeaderState()
	// }

	// Start the election timeout
	r.candidacy.Start()
	log.Debug().Str("replica", r.Name).Msg("replica running")
	return nil
}

func (r *Replica) setFollowerState() error {
	// Reset volatile state
	r.setInitializedState()

	// Update the tickers
	r.heartbeat.Stop()
	r.candidacy.Start()

	log.Info().Str("replica", r.Name).Uint64("term", r.term).Msg("replica is now a follower")
	return nil
}

func (r *Replica) setCandidateState() error {
	// Reset volatile state
	r.setInitializedState()

	// Create the election for the next term and vote for self
	// TODO: create election and vote for self
	r.term++

	// TODO: broadcast vote request

	log.Info().Str("replica", r.Name).Uint64("term", r.term).Msg("replica is now a candidate")
	return nil
}

func (r *Replica) setLeaderState() error {
	if r.state == Leader {
		return nil
	}

	// Stop the election timeout if we're leader
	r.candidacy.Stop()
	r.leader = r.PID

	// TODO: set the volatile state for known followers

	// TODO: broadcast heartbeat message

	// Start the heartbeat interval
	r.heartbeat.Start()
	log.Info().Str("replica", r.Name).Uint64("term", r.term).Msg("replica is now the leader")
	return nil
}
