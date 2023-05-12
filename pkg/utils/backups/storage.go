package backups

import "io"

// Storage provides an interface for reading and writing compressed backups to disk.
type Storage interface {
	Open(name string) (io.WriteCloser, error)
	ListArchives() ([]string, error)
}
