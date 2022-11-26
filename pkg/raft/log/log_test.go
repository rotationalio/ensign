package log_test

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	pb "github.com/rotationalio/ensign/pkg/raft/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/raft/log"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestEmptyLog(t *testing.T) {

	log, err := log.New()
	require.NoError(t, err)
	require.Equal(t, uint64(0), log.LastApplied())
	require.Equal(t, uint64(0), log.CommitIndex())
	require.Equal(t, uint64(0), log.LastTerm())
	require.Equal(t, uint64(0), log.CommitTerm())
	require.Equal(t, pb.NullEntry, log.LastEntry())
	require.Equal(t, pb.NullEntry, log.LastCommit())

	entry, err := log.Get(0)
	require.NoError(t, err)
	require.Equal(t, pb.NullEntry, entry)

	_, err = log.Get(1)
	require.EqualError(t, err, "no entry at index 1")

	_, err = log.Prev(0)
	require.EqualError(t, err, "no entry before index 0")

	entry, err = log.Prev(1)
	require.NoError(t, err)
	require.Equal(t, pb.NullEntry, entry)

	entries, err := log.After(0)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, pb.NullEntry, entries[0])

	_, err = log.After(1)
	require.EqualError(t, err, "no entries after 1")

	meta := log.Meta()
	require.Equal(t, uint64(0), meta.LastApplied)
	require.Equal(t, uint64(0), meta.CommitIndex)
	require.Equal(t, uint64(0), meta.Length)
	require.NotEmpty(t, meta.Created.AsTime())
	require.NotEmpty(t, meta.Modified.AsTime())
	require.Empty(t, meta.Snapshot.AsTime())
}

// A normal sequence of operations - this sequence exercises most log methods.
func TestSequence(t *testing.T) {

	log, err := log.New()
	require.NoError(t, err)
	require.True(t, log.AsUpToDate(3, 1), "an empty log should be as up to date as a log with entries because it is farther ahead")

	// It should be able to create entries in the log
	for i := 0; i < 4; i++ {
		_, err := log.Create(cmdKey, makeValue(), 1)
		require.NoError(t, err, "could not create a log entry in term 1")
	}

	// Create a fifth entry for comparison
	entry, err := log.Create(cmdKey, makeValue(), 1)
	require.NoError(t, err, "could not create log entry in term 1")

	require.Equal(t, uint64(5), log.LastApplied())
	require.Equal(t, uint64(0), log.CommitIndex())
	require.Equal(t, uint64(1), log.LastTerm())
	require.Equal(t, uint64(0), log.CommitTerm())
	require.Equal(t, uint64(5), log.Length())
	require.Equal(t, entry, log.LastEntry())
	require.Equal(t, pb.NullEntry, log.LastCommit())

	// State check: log should have 5 entries in term 1
	require.False(t, log.AsUpToDate(3, 1), "the log should now be farther ahead then a log with 3 entries in term 1")
	require.True(t, log.AsUpToDate(5, 1), "the log should be as up to date as a log with 5 entries in term 1")
	require.True(t, log.AsUpToDate(6, 2), "the log should be as up to date as a log with 6 entries with the last entry in term 2")

	// We should be able to append an entry into the log
	err = log.Append(makeEntry(6, 2))
	require.NoError(t, err, "could not append entry to log")

	// We should not be able to append an entry to the log that is behind
	err = log.Append(makeEntry(7, 1))
	require.EqualError(t, err, "cannot append entry in earlier term (1 < 2)")

	err = log.Append(makeEntry(6, 2))
	require.EqualError(t, err, "cannot append entry with smaller index (6 <= 6)")

	err = log.Append(makeEntry(5, 3))
	require.EqualError(t, err, "cannot append entry with smaller index (5 <= 6)")

	// We should not be able to create an entry in a term that is behind
	_, err = log.Create(cmdKey, makeValue(), 1)
	require.EqualError(t, err, "cannot append entry in earlier term (1 < 2)")

	// We should be able to append multiple entries to the log
	entries := make([]*pb.LogEntry, 0, 5)
	for i := 0; i < 5; i++ {
		entries = append(entries, makeEntry(uint64(7+i), 3))
	}
	err = log.Append(entries...)
	require.NoError(t, err)

	require.Equal(t, uint64(11), log.LastApplied())
	require.Equal(t, uint64(0), log.CommitIndex())
	require.Equal(t, uint64(3), log.LastTerm())
	require.Equal(t, uint64(0), log.CommitTerm())
	require.Equal(t, uint64(11), log.Length())
	require.Equal(t, entries[len(entries)-1], log.LastEntry())
	require.Equal(t, pb.NullEntry, log.LastCommit())

	// Commit entry 5
	err = log.Commit(5)
	require.NoError(t, err, "could not commit entry 5")

	require.Equal(t, uint64(11), log.LastApplied())
	require.Equal(t, uint64(5), log.CommitIndex())
	require.Equal(t, uint64(3), log.LastTerm())
	require.Equal(t, uint64(1), log.CommitTerm())
	require.Equal(t, uint64(11), log.Length())
	require.Equal(t, entries[len(entries)-1], log.LastEntry())
	require.Equal(t, entry, log.LastCommit())

	// Cannot commit entry 5 again or anything earlier
	require.EqualError(t, log.Commit(5), "index at 5 already committed")
	require.EqualError(t, log.Commit(3), "index at 3 already committed")

	// Cannot commit an index that is not in the log
	require.EqualError(t, log.Commit(42), "cannot commit invalid index 42")

	// Cannot truncate an index that has already been committed or that does not exist
	require.EqualError(t, log.Truncate(4, 1), "cannot truncate already committed index 5")
	require.EqualError(t, log.Truncate(14, 1), "cannot truncate invalid index 14")

	// Can truncate all entries after the last committed index in same term
	require.NoError(t, log.Truncate(6, 2))

	require.Equal(t, uint64(6), log.LastApplied())
	require.Equal(t, uint64(5), log.CommitIndex())
	require.Equal(t, uint64(2), log.LastTerm())
	require.Equal(t, uint64(1), log.CommitTerm())
	require.Equal(t, uint64(6), log.Length())
	require.Equal(t, entry, log.LastCommit())

	// Cannot truncate entries not in the same term
	require.EqualError(t, log.Truncate(5, 3), "entry at index 5 does not match term 3")
}

func TestLoad(t *testing.T) {
	fixture := &fixture{
		metaPath:    "testdata/meta.pb.json",
		entriesPath: "testdata/entries.pb.json",
	}

	log, err := log.Load(log.WithSync(fixture))
	require.NoError(t, err, "could not load log from fixtures")

	require.Equal(t, uint64(10), log.LastApplied())
	require.Equal(t, uint64(5), log.CommitIndex())
	require.Equal(t, uint64(2), log.LastTerm())
	require.Equal(t, uint64(1), log.CommitTerm())
	require.Equal(t, uint64(10), log.Length())

	entry, err := log.Get(6)
	require.NoError(t, err)
	require.Equal(t, uint64(6), entry.Index)
	require.Equal(t, uint64(2), entry.Term)
	require.Equal(t, "testKeyBravo", string(entry.Key))
	require.Equal(t, "Thu Nov 24 14:55:13 CST 2022\n", string(entry.Value))
}

func TestLoadError(t *testing.T) {
	// Sync required for load
	_, err := log.Load()
	require.ErrorIs(t, err, log.ErrSyncRequired)

	// Test filepath errors; expects meta is loaded first then entries
	fixture := &fixture{
		metaPath:    "testdata/doesnotexist.pb.json",
		entriesPath: "testdata/doesnotexist.pb.json",
	}

	_, err = log.Load(log.WithSync(fixture))
	require.Error(t, err)

	fixture.metaPath = "testdata/meta.pb.json"

	_, err = log.Load(log.WithSync(fixture))
	require.Error(t, err)

	fixture.entriesPath = "testdata/entries.pb.json"

	_, err = log.Load(log.WithSync(fixture))
	require.NoError(t, err)
}

func TestStateMachine(t *testing.T) {
	t.Skip("not implemented yet")
}

func TestSync(t *testing.T) {
	t.Skip("not implemented yet")
}

func TestAccesses(t *testing.T) {
	// Individual Operations when log starts empty
	t.Run("FromEmpty", func(t *testing.T) {

		t.Run("Create", func(t *testing.T) {

			log, err := log.New()
			require.NoError(t, err)
			entry, err := log.Create(cmdKey, makeValue(), 8)
			require.NoError(t, err, "could not create entry in empty log")
			require.NotNil(t, entry)

			require.Equal(t, uint64(1), entry.Index)
			require.Equal(t, uint64(8), entry.Term)

			require.Equal(t, entry.Index, log.LastApplied())
			require.Equal(t, entry.Term, log.LastTerm())
			require.Equal(t, entry, log.LastEntry())
		})

		t.Run("Append", func(t *testing.T) {

			log, err := log.New()
			require.NoError(t, err)

			// Should not be able to append an entry at index 0
			entry := makeEntry(0, 0)
			require.EqualError(t, log.Append(entry), "cannot append entry with smaller index (0 <= 0)")

			// Should not be able to append an entry in the future
			entry = makeEntry(42, 8)
			require.EqualError(t, log.Append(entry), "cannot skip index (1 to 42)")

			// Should be able to append a valid entry
			entry = makeEntry(1, 8)
			require.NoError(t, log.Append(entry), "could not append entry to empty log")

			require.Equal(t, entry.Index, log.LastApplied())
			require.Equal(t, entry.Term, log.LastTerm())

			cmp := log.LastEntry()
			require.Equal(t, entry, cmp)

			cmp2, err := log.Get(entry.Index)
			require.NoError(t, err)
			require.Equal(t, entry, cmp2)

			prev, err := log.Prev(entry.Index)
			require.NoError(t, err)
			require.Equal(t, pb.NullEntry, prev)

			after, err := log.After(entry.Index)
			require.NoError(t, err)
			require.Len(t, after, 1)
			require.Contains(t, after, entry)
		})

		t.Run("AppendMany", func(t *testing.T) {

			log, err := log.New()
			require.NoError(t, err)

			// We should be able to append multiple entries to the log
			entries := make([]*pb.LogEntry, 0, 5)
			for i := 0; i < 5; i++ {
				entries = append(entries, makeEntry(uint64(i+1), 3))
			}
			err = log.Append(entries...)
			require.NoError(t, err)

			entry := entries[len(entries)-1]
			require.Equal(t, entry.Index, log.LastApplied())
			require.Equal(t, entry.Term, log.LastTerm())
			require.Equal(t, entry, log.LastEntry())

			for i, e := range entries {
				o, err := log.Get(e.Index)
				require.NoError(t, err)
				require.Equal(t, e, o)

				o, err = log.Prev(e.Index)
				require.NoError(t, err)
				if i == 0 {
					require.Equal(t, pb.NullEntry, o)
				} else {
					require.Equal(t, entries[i-1], o)
				}
			}

			cmp, err := log.After(uint64(1))
			require.NoError(t, err)
			require.Equal(t, entries, cmp)
		})

		t.Run("Commit", func(t *testing.T) {

			log, err := log.New()
			require.NoError(t, err)
			require.Error(t, log.Commit(log.LastApplied()))
		})

		t.Run("Truncate", func(t *testing.T) {

			log, err := log.New()
			require.NoError(t, err)
			require.NoError(t, log.Truncate(log.LastApplied(), log.LastTerm()))

			entry, err := log.Get(0)
			require.NoError(t, err)
			require.Equal(t, pb.NullEntry, entry)
		})
	})

	// Individual Operations when log starts with data partially committed
	// TODO: implement
	t.Run("WithData", func(t *testing.T) {

		fixture := &fixture{
			metaPath:    "testdata/meta.pb.json",
			entriesPath: "testdata/entries.pb.json",
		}

		t.Run("Create", func(t *testing.T) {

			log, err := log.Load(log.WithSync(fixture))
			require.NoError(t, err, "could not load log from fixtures")

			entry, err := log.Create(cmdKey, makeValue(), 2)
			require.NoError(t, err, "could not create entry in log")
			require.NotNil(t, entry)

			require.Equal(t, uint64(11), entry.Index)
			require.Equal(t, uint64(2), entry.Term)

			require.Equal(t, entry.Index, log.LastApplied())
			require.Equal(t, entry.Term, log.LastTerm())
			require.Equal(t, entry, log.LastEntry())
		})

		t.Run("Append", func(t *testing.T) {

			log, err := log.Load(log.WithSync(fixture))
			require.NoError(t, err, "could not load log from fixtures")

			// Should not be able to overwrite an existing entry
			entry := makeEntry(7, 2)
			require.EqualError(t, log.Append(entry), "cannot append entry with smaller index (7 <= 10)")

			entry = makeEntry(log.LastApplied(), log.LastTerm())
			require.EqualError(t, log.Append(entry), "cannot append entry with smaller index (10 <= 10)")

			// Should not be able to append an entry in the future
			entry = makeEntry(42, 8)
			require.EqualError(t, log.Append(entry), "cannot skip index (11 to 42)")

			// Should be able to append a valid entry
			entry = makeEntry(11, 3)
			require.NoError(t, log.Append(entry), "could not append entry to log with data in it")

			require.Equal(t, entry.Index, log.LastApplied())
			require.Equal(t, entry.Term, log.LastTerm())

			cmp := log.LastEntry()
			require.Equal(t, entry, cmp)

			cmp2, err := log.Get(entry.Index)
			require.NoError(t, err)
			require.Equal(t, entry, cmp2)

			prev, err := log.Prev(entry.Index)
			require.NoError(t, err)
			require.Equal(t, entry.Index-1, prev.Index)

			after, err := log.After(log.CommitIndex())
			require.NoError(t, err)
			require.Len(t, after, 7)
			require.Contains(t, after, entry)
		})

		t.Run("AppendMany", func(t *testing.T) {

			log, err := log.Load(log.WithSync(fixture))
			require.NoError(t, err, "could not load log from fixtures")

			// We should be able to append multiple entries to the log
			entries := make([]*pb.LogEntry, 0, 5)
			for i := 0; i < 5; i++ {
				entries = append(entries, makeEntry(uint64(i+11), 3))
			}
			err = log.Append(entries...)
			require.NoError(t, err)

			entry := entries[len(entries)-1]
			require.Equal(t, entry.Index, log.LastApplied())
			require.Equal(t, entry.Term, log.LastTerm())
			require.Equal(t, entry, log.LastEntry())

			for i, e := range entries {
				o, err := log.Get(e.Index)
				require.NoError(t, err)
				require.Equal(t, e, o)

				o, err = log.Prev(e.Index)
				require.NoError(t, err)
				if i == 0 {
					require.Equal(t, o.Index, e.Index-1)
				} else {
					require.Equal(t, entries[i-1], o)
				}
			}

			cmp, err := log.After(uint64(11))
			require.NoError(t, err)
			require.Equal(t, entries, cmp)
		})

		t.Run("Commit", func(t *testing.T) {

			log, err := log.Load(log.WithSync(fixture))
			require.NoError(t, err, "could not load log from fixtures")

			require.NoError(t, log.Commit(log.LastApplied()))
			require.Equal(t, uint64(10), log.CommitIndex())
		})

		t.Run("Truncate", func(t *testing.T) {

			log, err := log.Load(log.WithSync(fixture))
			require.NoError(t, err, "could not load log from fixtures")

			require.NoError(t, log.Truncate(5, 1))

			require.Equal(t, uint64(5), log.LastApplied())
			require.Equal(t, uint64(5), log.CommitIndex())
			require.Equal(t, uint64(5), log.Length())

			_, err = log.Get(7)
			require.Error(t, err)
		})

	})
}

func TestOptionError(t *testing.T) {
	erropt := func(*log.Log) error {
		return errors.New("bad thing happened")
	}

	_, err := log.New(erropt)
	require.Error(t, err)

	fixture := &fixture{
		metaPath:    "testdata/meta.pb.json",
		entriesPath: "testdata/entries.pb.json",
	}
	_, err = log.Load(log.WithSync(fixture), erropt)
	require.Error(t, err)

}

func TestHelpers(t *testing.T) {
	alpha := makeEntry(102221, 42)
	bravo := makeEntry(201, 12)

	require.NotEqual(t, alpha.Term, bravo.Term)
	require.NotEqual(t, alpha.Index, bravo.Index)
	require.True(t, bytes.Equal(alpha.Key, bravo.Key))
	require.False(t, bytes.Equal(alpha.Value, bravo.Value))
}

var (
	cmdKey       = []byte("cmd")
	seq    int64 = 0
	sqm    sync.Mutex
)

func makeEntry(index, term uint64) *pb.LogEntry {
	return &pb.LogEntry{
		Index: index,
		Term:  term,
		Key:   cmdKey,
		Value: makeValue(),
	}
}

func makeValue() []byte {
	buf := make([]byte, binary.MaxVarintLen64)

	sqm.Lock()
	seq++
	binary.PutVarint(buf, time.Now().UnixNano()+seq)
	sqm.Unlock()
	return buf
}

// Fixture implements the log.Sync interface for loading a log from JSON fixtures.
// There are two paths, one storing the JSON metadata and the other in JSON Lines format
// storing each log entry as protocol buffer serialized JSON.
type fixture struct {
	writer      log.Writer
	metaPath    string
	entriesPath string
}

func (f *fixture) Write(entries ...*pb.LogEntry) error {
	if f.writer != nil {
		return f.writer.Write(entries...)
	}
	return nil
}

func (f *fixture) WriteMeta(meta *pb.LogMeta) error {
	if f.writer != nil {
		return f.writer.WriteMeta(meta)
	}
	return nil
}

func (f *fixture) Trunc(startIndex uint64) error {
	if f.writer != nil {
		return f.writer.Trunc(startIndex)
	}
	return nil
}

func (f *fixture) Read(index uint64) (entry *pb.LogEntry, err error) {
	var entries []*pb.LogEntry
	if entries, err = f.ReadFrom(0); err != nil {
		return nil, err
	}

	if index < uint64(len(entries)) {
		return entries[index], nil
	}
	return nil, fmt.Errorf("could not find entry at index %d", index)
}

func (f *fixture) ReadFrom(index uint64) (entries []*pb.LogEntry, err error) {
	pbjson := protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	var file *os.File
	if file, err = os.Open(f.entriesPath); err != nil {
		return nil, err
	}
	defer file.Close()

	entries = make([]*pb.LogEntry, 0, 10)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entry := &pb.LogEntry{}
		if err = pbjson.Unmarshal(scanner.Bytes(), entry); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries[index:], nil
}

func (f *fixture) ReadMeta() (meta *pb.LogMeta, err error) {
	pbjson := protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	var data []byte
	if data, err = os.ReadFile(f.metaPath); err != nil {
		return nil, err
	}

	meta = &pb.LogMeta{}
	if err = pbjson.Unmarshal(data, meta); err != nil {
		return nil, err
	}
	return meta, nil
}

func (f *fixture) Close() error {
	return nil
}
