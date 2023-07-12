package rlid

import "sync"

// Sequence generates totally ordered RLIDs but is not thread-safe.
type Sequence uint32

func (s *Sequence) Next() RLID {
	*s++
	return Make(uint32(*s))
}

// LockedSequence generates thread-safe totally ordered RLIDs.
type LockedSequence struct {
	sync.Mutex
	seq uint32
}

func (s *LockedSequence) Next() RLID {
	s.Lock()
	defer s.Unlock()
	s.seq++
	return Make(s.seq)
}
