package raft

import (
	"errors"
	"fmt"
	"time"

	"github.com/rotationalio/ensign/pkg/raft/peers"
)

// Config is intended to be loaded from the environment using confire and embedded as
// an subconfiguration of the service configuration (not loaded directly).
type Config struct {
	ReplicaID uint32        `required:"true"`    // the process ID of the local replica to identify from the peers config
	Tick      time.Duration `default:"1s"`       // the tick interval used to compute the heartbeat and election timeouts
	Timeout   time.Duration `default:"500ms"`    // the timeout to wait for a response from a peer
	Aggregate bool          `default:"true"`     // aggregate multiple commands into one append entries
	PeersPath string        `split_words:"true"` // the path to the peers configuration (usually loaded from a config map), only loaded if quorum is nil
	Quorum    *peers.Quorum `ignored:"true"`     // the peers configuration, will not be loaded from the environment
}

// Validation Errors
var (
	ErrMissingReplicaID = errors.New("invalid raft configuration: local replica id is required")
	ErrMissingField     = errors.New("invalid raft configuration: tick and timeout are required and one of peers path or quorum")
	ErrTickTooSmall     = errors.New("invalid raft configuration: tick must be greater than 10ms")
	ErrTimeoutTooBig    = errors.New("invalid raft configuration: timeout must be smaller than the tick")
	ErrMissingReplica   = errors.New("invalid raft configuration: local replica is not defined in the quorum")
)

// Validate the raft configuration. This also loads the peers from disk and validates
// the quorum and peers configuration if not directly specified.
func (c *Config) Validate() (err error) {
	if c.ReplicaID == 0 {
		return ErrMissingReplicaID
	}

	if c.Tick == 0 || c.Timeout == 0 || (c.Quorum == nil && c.PeersPath == "") {
		return ErrMissingField
	}

	if c.Tick < 10*time.Millisecond {
		return ErrTickTooSmall
	}

	if c.Timeout >= c.Tick {
		return ErrTimeoutTooBig
	}

	// Load the quorum from the peers path if it hasn't been specified by the user
	if c.Quorum == nil {
		if c.Quorum, err = peers.Load(c.PeersPath); err != nil {
			return fmt.Errorf("invalid raft configuration: could not load quorum from peers path: %w", err)
		}
	}

	if err = c.Quorum.Validate(); err != nil {
		return err
	}

	if !c.Quorum.Contains(c.ReplicaID) {
		return ErrMissingReplica
	}
	return nil
}
