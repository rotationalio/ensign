package mock

import (
	"fmt"

	pb "github.com/rotationalio/ensign/pkg/raft/api/v1beta1"
)

const (
	CommitEntry = "CommitEntry"
	DropEntry   = "DropEntry"
)

func NewStateMachine() *StateMachine {
	return &StateMachine{
		Calls: make(map[string]int),
	}
}

type StateMachine struct {
	Calls         map[string]int
	OnCommitEntry func(*pb.LogEntry) error
	OnDropEntry   func(*pb.LogEntry) error
}

func (m *StateMachine) UseError(method string, err error) {
	switch method {
	case CommitEntry:
		m.OnCommitEntry = func(*pb.LogEntry) error {
			return err
		}
	case DropEntry:
		m.OnDropEntry = func(*pb.LogEntry) error {
			return err
		}
	default:
		panic(fmt.Errorf("unknown method %q", method))
	}
}

func (m *StateMachine) Reset() {
	m.OnCommitEntry = nil
	m.OnDropEntry = nil
	for key := range m.Calls {
		m.Calls[key] = 0
	}
}

func (m *StateMachine) CommitEntry(entry *pb.LogEntry) error {
	m.incr(CommitEntry)
	if m.OnCommitEntry != nil {
		return m.OnCommitEntry(entry)
	}
	return nil
}

func (m *StateMachine) DropEntry(entry *pb.LogEntry) error {
	m.incr(DropEntry)
	if m.OnDropEntry != nil {
		return m.OnDropEntry(entry)
	}
	return nil
}

func (m *StateMachine) incr(name string) {
	if m.Calls == nil {
		m.Calls = make(map[string]int)
	}
	m.Calls[name]++
}
