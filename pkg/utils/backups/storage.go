package backups

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// Storage provides an interface for reading and writing compressed backups to disk.
type Storage interface {
	Open(name string) (io.WriteCloser, error)
	Remove(name string) error
	ListArchives() ([]string, error)
}

// FileStorage stores archived backups on a local disk. This is primarily used by
// mounting a second volume on a Kubernetes pod and saving the backups there or by
// writing to a network or RAID protected dis.
type FileStorage struct {
	root   string
	prefix string
}

// NewFileStorage creates the specified root directory if it does not exist. All storage
// archives will be opened relative to the root directory.
func NewFileStorage(root, prefix string) (_ *FileStorage, err error) {
	var stat os.FileInfo
	if stat, err = os.Stat(root); err != nil {
		if os.IsNotExist(err) {
			// Create the directory if it doesn't exist
			if err = os.MkdirAll(root, 0755); err != nil {
				return nil, fmt.Errorf("could not create backup storage directory: %w", err)
			}
		} else {
			return nil, fmt.Errorf("could not stat backup storage directory: %w", err)
		}
	}

	if !stat.IsDir() {
		return nil, ErrNotADirectory
	}

	return &FileStorage{root: root}, nil
}

// Open a file on disk and return the file for writting a gzip stream to.
func (s *FileStorage) Open(name string) (_ io.WriteCloser, err error) {
	var f *os.File
	if f, err = os.OpenFile(filepath.Join(s.root, name), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
		return nil, fmt.Errorf("could not open archive for writing: %w", err)
	}
	return f, nil
}

// Remove a file from the backup directory.
func (s *FileStorage) Remove(name string) error {
	if err := os.Remove(filepath.Join(s.root, name)); err != nil {
		return err
	}
	return nil
}

// ListArchives returns all backup archives in the FileStorage directory ordered by date
// ascending using string sorting that depends on the backup archive name format:
// prefix-YYYYmmddHHMM.tgz
func (s *FileStorage) ListArchives() (paths []string, err error) {
	pattern := fmt.Sprintf("%s-[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9].tgz", s.prefix)
	if paths, err = filepath.Glob(filepath.Join(s.root, pattern)); err != nil {
		return nil, err
	}

	// Sort the paths by timestamp ascending
	sort.Strings(paths)
	return paths, nil
}

// Memory storage is primarily used for testing and writes backups to a bytes buffer.
// If /dev/null is specified as the "root" path, then a noop-closer is returned.
type MemoryStorage struct {
	root    string
	backups map[string]*Buffer
}

func (s *MemoryStorage) Open(name string) (io.WriteCloser, error) {
	if s.root == "/dev/null" {
		return &Discard{}, nil
	}

	path := filepath.Join(s.root, name)
	s.backups[path] = &Buffer{Buffer: *bytes.NewBuffer(nil)}
	return s.backups[path], nil
}

func (s *MemoryStorage) Remove(name string) error {
	delete(s.backups, name)
	return nil
}

func (s *MemoryStorage) ListArchives() ([]string, error) {
	paths := make([]string, 0, len(s.backups))
	for key := range s.backups {
		paths = append(paths, key)
	}

	sort.Strings(paths)
	return paths, nil
}

func (s *MemoryStorage) Backup(name string) (*Buffer, bool) {
	data, ok := s.backups[name]
	return data, ok
}

type Discard struct{}

func (Discard) Write(p []byte) (int, error) {
	return io.Discard.Write(p)
}

func (Discard) Close() error { return nil }

type Buffer struct {
	bytes.Buffer
}

func (Buffer) Close() error { return nil }
