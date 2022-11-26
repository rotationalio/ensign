package log

import (
	"fmt"
	"time"

	pb "github.com/rotationalio/ensign/pkg/raft/api/v1beta1"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Log implements the sequence of commands applied to the Raft state machine.
// This implementation uses an in-memory log that is backed by a disk based store for
// durability. The in-memory portion of the log contains all uncommitted entries and
// recently committed entries that need to be applied to the state machine. Once applied
// successfully, the in-memory log is truncated.
//
// The log's primarily responsibility is to ensure that the sequence of commands is
// consistent, e.g. that entries are appended in a monotonically increasing time order
// as defined by the Raft leader's term.
//
// Note that the log is not thread-safe and is not intended to be accessed from multiple
// go routines. Instead the log should be maintained by a single state machine that
// updates it sequentially when entries are committed.
//
// TODO: right now the log stores everything in-memory; refactor to store partial in-mem log
// TODO: implement snapshotting functionality
// TODO: implement log validation to ensure the log is in a correct state
type Log struct {
	sm          StateMachine   // State machine to apply commits to
	sync        Sync           // Synchronize the log to disk
	lastApplied uint64         // The index of the last applied log entry
	commitIndex uint64         // The index of the last committed log entry
	length      uint64         // The total number of entries in the log
	entries     []*pb.LogEntry // The in-memory array of log entries
	created     time.Time      // Timestamp the log was created
	modified    time.Time      // Timestamp of the last log modification
	snapshot    time.Time      // Timestamp of the last log snapshot
	meta        *pb.LogMeta    // Saved state; only updated on calls to Meta()
}

// New creates an empty log with a null entry at index 0. It is a fresh log ready for
// any log operations that may be applied to it. Generally, external users will want to
// load the log from disk using a sync command.
func New(opts ...Option) (*Log, error) {
	log := &Log{
		entries:  make([]*pb.LogEntry, 1),
		created:  time.Now(),
		modified: time.Now(),
		meta:     &pb.LogMeta{},
	}
	log.entries[0] = pb.NullEntry

	for _, opt := range opts {
		if err := opt(log); err != nil {
			return nil, err
		}
	}
	return log, nil
}

// Load a log from disk. This method creates a new log and reads entries and meta data
// from the sync reader. An error is returned if there is no WithSync() option.
func Load(opts ...Option) (log *Log, err error) {
	if log, err = New(opts...); err != nil {
		return nil, err
	}

	if log.sync == nil {
		return nil, ErrSyncRequired
	}

	if log.meta, err = log.sync.ReadMeta(); err != nil {
		return nil, fmt.Errorf("could not read meta: %w", err)
	}

	log.lastApplied = log.meta.LastApplied
	log.commitIndex = log.meta.CommitIndex
	log.length = log.meta.Length
	log.created = log.meta.Created.AsTime()
	log.modified = log.meta.Modified.AsTime()
	log.snapshot = log.meta.Snapshot.AsTime()

	if log.entries, err = log.sync.ReadFrom(0); err != nil {
		return nil, fmt.Errorf("could not read entries: %w", err)
	}
	return log, nil
}

//===========================================================================
// Index Management
//===========================================================================

// LastApplied returns the index of the last applied log entry.
func (l *Log) LastApplied() uint64 {
	return l.lastApplied
}

// CommitIndex returns the index of the last committed log entry.
func (l *Log) CommitIndex() uint64 {
	return l.commitIndex
}

// LastEntry returns the log entry at the last applied index.
func (l *Log) LastEntry() *pb.LogEntry {
	return l.entries[l.lastApplied]
}

// LastCommit returns the log entry at the commit index.
func (l *Log) LastCommit() *pb.LogEntry {
	return l.entries[l.commitIndex]
}

// LastTerm is a helper function to get the term of the entry at the last applied index.
func (l *Log) LastTerm() uint64 {
	return l.LastEntry().Term
}

// Length returns the number of entries in the log
func (l *Log) Length() uint64 {
	return l.length
}

// CommitTerm is a helper function to get the term of the entry at the commit index.
func (l *Log) CommitTerm() uint64 {
	return l.LastCommit().Term
}

// AsUpToDate returns true if the remote log specified by the last index and
// last term are at least as up to date (or farther ahead) than the local log.
func (l *Log) AsUpToDate(lastIndex, lastTerm uint64) bool {
	localTerm := l.LastTerm()

	// If we're in the same term as the remote host, our last applied index
	// should be at least as large as the remote's last applied index.
	if lastTerm == localTerm {
		return lastIndex >= l.lastApplied
	}

	// Otherwise ensure that the remote's term is greater than our own.
	return lastTerm > localTerm
}

//===========================================================================
// Entry Management
//===========================================================================

// Create an entry in the log and append it. This is essentially a helper method
// for quickly adding a command to the state machine consistent with the local log.
func (l *Log) Create(key, value []byte, term uint64) (*pb.LogEntry, error) {
	// Create the entry at the next log index
	entry := &pb.LogEntry{
		Index: l.lastApplied + 1,
		Term:  term,
		Key:   key,
		Value: value,
	}

	// Append the entry and perform invariant checks
	if err := l.Append(entry); err != nil {
		return nil, err
	}

	// Return the entry for use elsewhere
	return entry, nil
}

// Append one ore more entries and perform log invariant checks. If appending
// an entry creates a log inconsistency (out of order term or index), then an
// error is returned. A couple of important notes:
//
//  1. Append does not undo any successful appends even on error
//  2. Append will not compare entries that specify the same index
//
// These notes mean that all entries being appended to this log should be
// consistent with each other as well as the end of the log, and that the log
// needs to be truncated in order to "update" or splice two logs together.
func (l *Log) Append(entries ...*pb.LogEntry) error {
	// Append all entries one at a time, returning an error if an append fails.
	for _, entry := range entries {

		// Fetch the latest entry
		prev := l.LastEntry()

		// Ensure that the term is monotonically increasing
		if entry.Term < prev.Term {
			return fmt.Errorf("cannot append entry in earlier term (%d < %d)", entry.Term, prev.Term)
		}

		// Ensure that the index is monotonically increasing
		if entry.Index <= prev.Index {
			return fmt.Errorf("cannot append entry with smaller index (%d <= %d)", entry.Index, prev.Index)
		}

		// Ensure that the index is not skipped
		if entry.Index > prev.Index+1 {
			return fmt.Errorf("cannot skip index (%d to %d)", prev.Index+1, entry.Index)
		}

		// Append the entry and update metadata
		l.entries = append(l.entries, entry)
		l.lastApplied = entry.Index
		l.length++
	}

	// The log has been updated
	l.modified = time.Now()

	// Sync the log and metadata to disk
	if l.sync != nil {
		if err := l.sync.Write(entries...); err != nil {
			return err
		}

		if err := l.sync.WriteMeta(l.Meta()); err != nil {
			return err
		}
	}
	return nil
}

// Commit all entries up to and including the specified index.
func (l *Log) Commit(index uint64) error {
	// Ensure the index specified is in the log
	if index < 1 || index > l.lastApplied {
		return fmt.Errorf("cannot commit invalid index %d", index)
	}

	// Ensure that we haven't already committed this index
	if index <= l.commitIndex {
		return fmt.Errorf("index at %d already committed", index)
	}

	// Create a commit event for all entries now committed
	if l.sm != nil {
		for i := l.commitIndex + 1; i <= index; i++ {
			if err := l.sm.CommitEntry(l.entries[i]); err != nil {
				return err
			}
		}
	}

	// Update the commit index and the log
	l.commitIndex = index
	l.modified = time.Now()

	if l.sync != nil {
		if err := l.sync.WriteMeta(l.Meta()); err != nil {
			return err
		}
	}

	log.Debug().Uint64("index", l.commitIndex).Msg("committed index")
	return nil
}

// Truncate the log to the given position, conditioned by term.
// This method returns an error if the log has been committed after the
// specified index, there is an epoch mismatch, or there is some other log
// operation error.
//
// This method truncates everything after the given index, but keeps the
// entry at the specified index; e.g. truncate after.
func (l *Log) Truncate(index, term uint64) error {
	// Ensure the truncation matches an entry
	if index > l.lastApplied {
		return fmt.Errorf("cannot truncate invalid index %d", index)
	}

	// Specifies the index of the entry to be truncated
	nextIndex := index + 1

	// Do not allow committed entries to be truncted
	if nextIndex <= l.commitIndex {
		return fmt.Errorf("cannot truncate already committed index %d", nextIndex)
	}

	// Do not truncate if entry at index does not have matching term
	entry := l.entries[index]
	if entry.Term != term {
		return fmt.Errorf("entry at index %d does not match term %d", index, term)
	}

	// Only perform truncation if necessary
	if index < l.lastApplied {
		// Drop all entries that appear after the index
		if l.sm != nil {
			for _, droppedEntry := range l.entries[nextIndex:] {
				if err := l.sm.DropEntry(droppedEntry); err != nil {
					return err
				}
			}
		}

		// Update the entries and meta data
		l.entries = l.entries[0:nextIndex]
		l.length -= l.lastApplied - index
		l.lastApplied = index
		l.modified = time.Now()

		if l.sync != nil {
			if err := l.sync.Trunc(nextIndex); err != nil {
				return err
			}

			if err := l.sync.WriteMeta(l.Meta()); err != nil {
				return err
			}
		}

	}
	return nil
}

//===========================================================================
// Entry Access
//===========================================================================

// Get the entry at the specified index (whether or not it is committed).
// Returns an error if no entry exists at the index.
func (l *Log) Get(index uint64) (*pb.LogEntry, error) {
	if index > l.lastApplied {
		return nil, fmt.Errorf("no entry at index %d", index)
	}

	return l.entries[index], nil
}

// Prev returns the entry before the specified index (whether or not it is
// committed). Returns an error if no entry exists before.
func (l *Log) Prev(index uint64) (*pb.LogEntry, error) {
	if index < 1 || index > (l.lastApplied+1) {
		return nil, fmt.Errorf("no entry before index %d", index)
	}

	return l.entries[index-1], nil
}

// After returns all entries after the specified index, inclusive
func (l *Log) After(index uint64) ([]*pb.LogEntry, error) {
	if index > l.lastApplied {
		return make([]*pb.LogEntry, 0), fmt.Errorf("no entries after %d", index)
	}

	return l.entries[index:], nil
}

//===========================================================================
// Metadata Management
//===========================================================================

// Meta updates the state of the saved metadata on the log and returns a pointer to it.
// This means that the returned Meta is not safe for concurrent use and any writer that
// operates in another go routine should create a copy of it.
func (l *Log) Meta() *pb.LogMeta {
	l.meta.LastApplied = l.lastApplied
	l.meta.CommitIndex = l.commitIndex
	l.meta.Length = l.length
	l.meta.Created = timestamppb.New(l.created)
	l.meta.Modified = timestamppb.New(l.modified)
	l.meta.Snapshot = timestamppb.New(l.snapshot)
	return l.meta
}
