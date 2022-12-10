package log

import (
	"io"

	pb "github.com/rotationalio/ensign/pkg/raft/api/v1beta1"
)

type StateMachine interface {
	CommitEntry(*pb.LogEntry) error
	DropEntry(*pb.LogEntry) error
}

type Sync interface {
	Writer
	Reader
	io.Closer
}

type Writer interface {
	// Write should append all of the log entries to disk.
	Write(...*pb.LogEntry) error

	// WriteMeta should sync the log metadata to disk.
	WriteMeta(*pb.LogMeta) error

	// Trunc should delete all entries starting with the given index and all entries
	// that follow. Note that this is different than the log.Truncate() semantics.
	Trunc(startIndex uint64) error
}

type Reader interface {
	// Read should return the entry at the specified index.
	Read(index uint64) (*pb.LogEntry, error)

	// ReadFrom should return all entries starting at the given index and following.
	ReadFrom(index uint64) ([]*pb.LogEntry, error)

	// ReadMeta should return the log metadata from disk.
	ReadMeta() (*pb.LogMeta, error)
}
