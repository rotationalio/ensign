package mock

import (
	"fmt"

	pb "github.com/rotationalio/ensign/pkg/raft/api/v1beta1"
)

const (
	Write     = "Write"
	WriteMeta = "WriteMeta"
	Trunc     = "Trunc"
	Read      = "Read"
	ReadFrom  = "ReadFrom"
	ReadMeta  = "ReadMeta"
	Close     = "Close"
)

func NewSync() *Sync {
	return &Sync{
		Calls: make(map[string]int),
	}
}

type Sync struct {
	Calls       map[string]int
	OnWrite     func(...*pb.LogEntry) error
	OnWriteMeta func(*pb.LogMeta) error
	OnTrunc     func(uint64) error
	OnRead      func(uint64) (*pb.LogEntry, error)
	OnReadFrom  func(uint64) ([]*pb.LogEntry, error)
	OnReadMeta  func() (*pb.LogMeta, error)
	OnClose     func() error
}

func (m *Sync) UseError(method string, err error) {
	switch method {
	case Write:
		m.OnWrite = func(...*pb.LogEntry) error {
			return err
		}
	case WriteMeta:
		m.OnWriteMeta = func(*pb.LogMeta) error {
			return err
		}
	case Trunc:
		m.OnTrunc = func(uint64) error {
			return err
		}
	case Read:
		m.OnRead = func(uint64) (*pb.LogEntry, error) {
			return nil, err
		}
	case ReadFrom:
		m.OnReadFrom = func(uint64) ([]*pb.LogEntry, error) {
			return nil, err
		}
	case ReadMeta:
		m.OnReadMeta = func() (*pb.LogMeta, error) {
			return nil, err
		}
	case Close:
		m.OnClose = func() error {
			return err
		}
	default:
		panic(fmt.Errorf("unknown method %q", method))
	}
}

func (m *Sync) Reset() {
	m.OnWrite = nil
	m.OnWriteMeta = nil
	m.OnTrunc = nil
	m.OnRead = nil
	m.OnReadFrom = nil
	m.OnReadMeta = nil
	m.OnClose = nil
	for key := range m.Calls {
		m.Calls[key] = 0
	}
}

func (m *Sync) Write(entries ...*pb.LogEntry) error {
	m.incr(Write)
	if m.OnWrite != nil {
		return m.OnWrite(entries...)
	}
	return nil
}

func (m *Sync) WriteMeta(meta *pb.LogMeta) error {
	m.incr(WriteMeta)
	if m.OnWriteMeta != nil {
		return m.OnWriteMeta(meta)
	}
	return nil
}

func (m *Sync) Trunc(index uint64) error {
	m.incr(Trunc)
	if m.OnTrunc != nil {
		return m.OnTrunc(index)
	}
	return nil
}

func (m *Sync) Read(index uint64) (*pb.LogEntry, error) {
	m.incr(Read)
	if m.OnRead != nil {
		return m.OnRead(index)
	}
	return pb.NullEntry, nil
}

func (m *Sync) ReadFrom(index uint64) ([]*pb.LogEntry, error) {
	m.incr(ReadFrom)
	if m.OnReadFrom != nil {
		return m.OnReadFrom(index)
	}
	entries := make([]*pb.LogEntry, 1)
	entries[0] = pb.NullEntry
	return entries, nil
}

func (m *Sync) ReadMeta() (*pb.LogMeta, error) {
	m.incr(ReadMeta)
	if m.OnReadMeta != nil {
		return m.OnReadMeta()
	}
	return &pb.LogMeta{}, nil
}

func (m *Sync) Close() error {
	m.incr(Close)
	if m.OnClose != nil {
		return m.OnClose()
	}
	return nil
}

func (m *Sync) incr(name string) {
	if m.Calls == nil {
		m.Calls = make(map[string]int)
	}
	m.Calls[name]++
}
