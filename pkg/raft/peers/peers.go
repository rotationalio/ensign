package peers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v3"
)

// Quorum represents the configuration of a quorum on disk.
type Quorum struct {
	QID             uint32  `json:"quorum_id" yaml:"quorum_id"`                                  // the unique ID of the quorum of peers
	BootstrapLeader uint32  `json:"boostrap_leader,omitempty" yaml:"bootstrap_leader,omitempty"` // the PID of the leader to bootstrap at term 0
	Peers           []*Peer `json:"peers" yaml:"peers"`                                          // the peers that are in the quorum
}

// Peer is a configuration of a peer and contains all connection and helper information.
type Peer struct {
	PID      uint32 `json:"peer_id" yaml:"peer_id"`     // the unique ID of the peer in the system
	Name     string `json:"name" yaml:"name"`           // a human readable name for the peer (e.g. hostname)
	BindAddr string `json:"bind_addr" yaml:"bind_addr"` // the address to bind the peer on to listen for requests
	Endpoint string `json:"endpoint" yaml:"endpoint"`   // the domain or ip address and port to connect to the peer on
}

//===========================================================================
// Validation
//===========================================================================

// Validation Errors
var (
	ErrMissingQID       = errors.New("invalid quorum configuration: quorum id is required")
	ErrNoPeers          = errors.New("invalid quorum configuration: no peers assigned to quorum")
	ErrUniquePID        = errors.New("invalid quorum configuration: peer ids must be unique")
	ErrMissingPID       = errors.New("invalid peer configuration: peer id is required")
	ErrPeerMissingField = errors.New("invalid peer configuration: name, bind_addr, and endpoint are required")
)

func (q *Quorum) Validate() (err error) {
	if q.QID == 0 {
		err = multierror.Append(err, ErrMissingQID)
	}

	if len(q.Peers) == 0 {
		err = multierror.Append(err, ErrNoPeers)
	}

	pids := make(map[uint32]int)
	for _, peer := range q.Peers {
		if perr := peer.Validate(); perr != nil {
			err = multierror.Append(err, perr)
		}
		pids[peer.PID]++
	}

	for _, count := range pids {
		if count > 1 {
			err = multierror.Append(err, ErrUniquePID)
			break
		}
	}

	return err
}

func (p *Peer) Validate() error {
	if p.PID == 0 {
		return ErrMissingPID
	}

	if p.Name == "" || p.BindAddr == "" || p.Endpoint == "" {
		return ErrPeerMissingField
	}
	return nil
}

// Contains checks if the quorum contains the peer specified by PID or name.
func (q *Quorum) Contains(pidorname interface{}) bool {
	switch t := pidorname.(type) {
	case uint32:
		for _, peer := range q.Peers {
			if peer.PID == t {
				return true
			}
		}
	case string:
		for _, peer := range q.Peers {
			if peer.Name == t {
				return true
			}
		}
	}
	return false
}

//===========================================================================
// Serialization
//===========================================================================

// Load the quorum configuration from disk.
func Load(path string) (quorum *Quorum, err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return nil, err
	}
	defer f.Close()

	switch ext := filepath.Ext(path); ext {
	case ".json":
		if err = json.NewDecoder(f).Decode(&quorum); err != nil {
			return nil, err
		}
	case ".yaml":
		if err = yaml.NewDecoder(f).Decode(&quorum); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown file extension %q", ext)
	}
	return quorum, nil
}

// Dump the quorum configuration to disk.
func (q *Quorum) Dump(path string) (err error) {
	var f *os.File
	if f, err = os.Create(path); err != nil {
		return err
	}
	defer f.Close()

	switch ext := filepath.Ext(path); ext {
	case ".json":
		if err = json.NewEncoder(f).Encode(q); err != nil {
			return err
		}
	case ".yaml":
		if err = yaml.NewEncoder(f).Encode(q); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown file extension %q", ext)
	}

	return nil
}
