/*
Package keymu provides a map of mutexes that allows users to lock and unlock single keys
in the data structure. This functionality can be used to provide key-based "transactions"
in the form of exclusive locks on a given key.

See: https://stackoverflow.com/questions/40931373/how-to-gc-a-map-of-mutexes-in-go
*/
package keymu

import (
	"fmt"
	"sync"
)

// TODO: create an RWMutex
type Mutex struct {
	mu      sync.Mutex
	entries map[any]*entry
}

type Unlocker interface {
	Unlock()
}

func New() *Mutex {
	return &Mutex{entries: make(map[any]*entry)}
}

type entry struct {
	mu     sync.Mutex
	parent *Mutex
	refs   uint64
	key    any
}

func (m *Mutex) Lock(key any) Unlocker {
	m.mu.Lock()
	e, ok := m.entries[key]
	if !ok {
		e = &entry{parent: m, key: key}
		m.entries[key] = e
	}

	e.refs++
	m.mu.Unlock()

	// Acquire lock here, will block until entry.refs == 1
	e.mu.Lock()
	return e
}

func (e *entry) Unlock() {
	e.parent.mu.Lock()
	entry, ok := e.parent.entries[e.key]
	if !ok {
		// entry must exist
		e.parent.mu.Unlock()
		panic(fmt.Errorf("cannot unlock key=%v no entry found", e.key))
	}
	entry.refs--
	if entry.refs < 1 {
		delete(e.parent.entries, e.key)
	}
	e.parent.mu.Unlock()

	// Now that the map stuff is handled, we unlock
	e.mu.Unlock()
}

func (m *Mutex) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.entries)
}
