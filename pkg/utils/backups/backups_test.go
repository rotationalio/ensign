package backups_test

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha512"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
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

// Extract a gzip file to the specified destination directory.
func extract(file, destDir string, skipHidden bool) (root string, err error) {
	var (
		f  *os.File
		gr *gzip.Reader
	)

	// Read the gzip file.
	if f, err = os.Open(file); err != nil {
		return "", err
	}
	defer f.Close()
	if gr, err = gzip.NewReader(f); err != nil {
		return "", err
	}
	defer gr.Close()

	// Write the contents to the temporary directory.
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err = os.MkdirAll(filepath.Join(destDir, hdr.Name), os.FileMode(hdr.Mode)); err != nil {
				return "", err
			}
			if root == "" {
				root = filepath.Join(destDir, hdr.Name)
			}
		case tar.TypeReg:
			var reg *os.File
			if skipHidden && hdr.Name[0] == '.' {
				// Skip hidden files if requested.
				continue
			}
			if reg, err = os.Create(filepath.Join(destDir, hdr.Name)); err != nil {
				return "", err
			}
			if _, err = io.Copy(reg, tr); err != nil {
				reg.Close()
				return "", err
			}
			reg.Close()
		default:
			return "", fmt.Errorf("extracting %s: unknown type flag: %c", hdr.Name, hdr.Typeflag)
		}
	}
	return root, nil
}

// Archive a directory to the specified gzip file
func archive(dir, file string) (err error) {
	var (
		f *os.File
	)
	// Create a gzip file.
	if f, err = os.Create(file); err != nil {
		return err
	}
	defer f.Close()
	gw := gzip.NewWriter(f)
	defer gw.Close()

	// Create a tar file.
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Write the DB to the tar file.
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		var hdr *tar.Header
		if hdr, err = tar.FileInfoHeader(info, ""); err != nil {
			return err
		}
		hdr.Name = path[len(dir):]
		if err = tw.WriteHeader(hdr); err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		var tmp *os.File
		if tmp, err = os.Open(path); err != nil {
			return err
		}
		defer tmp.Close()
		if _, err = io.Copy(tw, tmp); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
