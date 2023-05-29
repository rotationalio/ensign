package backups_test

import (
	"crypto/sha512"
	"errors"
	"io"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rotationalio/ensign/pkg/utils/backups"
	"github.com/stretchr/testify/require"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	MaxBackupRecords = 1177
)

type MockBackup struct {
	err     error
	tmpdirs []string
}

func (m *MockBackup) Backup(tmpdir string) error {
	m.tmpdirs = append(m.tmpdirs, tmpdir)
	return m.err
}

func TestBackupManager(t *testing.T) {
	// Test setting up and running the backup manager
	conf := backups.Config{
		Enabled:    true,
		Interval:   50 * time.Millisecond,
		StorageDSN: "inmem:////dev/null",
	}

	mock := &MockBackup{}
	manager := backups.New(conf, mock)

	err := manager.Run()
	require.NoError(t, err, "could not run backup manager")
	require.NoError(t, manager.Run(), "running a running backup manager should not error")

	// Ensure at least one backup runs
	time.Sleep(150 * time.Millisecond)
	err = manager.Shutdown()
	require.NoError(t, err, "could not shutdown backup manager")

	// Shuting down a shutdown backup manager should not error
	require.NoError(t, manager.Shutdown(), "should be able to shutdown a shutdown backup manager without error")

	nbackups := len(mock.tmpdirs)
	require.GreaterOrEqual(t, nbackups, 1, "expected at least one backup to be run")

	// No more backups should be run after shutdown
	time.Sleep(150 * time.Millisecond)
	require.Equal(t, nbackups, len(mock.tmpdirs), "expected no more backups to be run after shutdown")

	// Should be able to restart the backup manager even if it errors
	mock.err = errors.New("something bad happened")
	require.NoError(t, manager.Run(), "could not run backup manager a second time")

	time.Sleep(150 * time.Millisecond)
	require.NoError(t, manager.Shutdown(), "could not shutdown manager")
	require.Greater(t, len(mock.tmpdirs), nbackups, "expected backup manager to run even with errors")

}

func TestDisabledBackupManager(t *testing.T) {
	// Test setting up and running the backup manager
	conf := backups.Config{
		Enabled:    false,
		StorageDSN: "inmem:////dev/null",
	}

	mock := &MockBackup{}
	manager := backups.New(conf, mock)

	err := manager.Run()
	require.ErrorIs(t, err, backups.ErrNotEnabled, "expected error on disabled backup")
}

func TestCanMkdTemp(t *testing.T) {
	conf := backups.Config{TempDir: "./testdata"}
	manager := backups.New(conf, nil)

	dir, err := manager.MkdirTemp()
	defer os.RemoveAll(dir)

	require.NoError(t, err, "could not create tempdir")
	require.DirExists(t, dir, "expected tmp dir to exist")
	require.True(t, strings.HasPrefix(dir, "./testdata"))
}

// A helper function that checks to make sure the fixture at the path exists; if it
// doesn't then the createFixture function is run and any errors returned.
func checkFixture(path string, createFixture func(path string) error) (err error) {
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return createFixture(path)
	}
	return nil
}

// A helper function to compute the SHA512 hash of a file for comparison purposes.
func fileHash(path string) (_ []byte, err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return nil, err
	}
	defer f.Close()

	hash := sha512.New()
	if _, err = io.Copy(hash, f); err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}
