/*
Package db establishes a connection with a Raft replicated sqlite3 database.
External packages can use this module to ensure that the database is at the most current
schema and can make thread-safe transactions against the database.

Users of the package have to call db.Connect() at least once to use the database, but
multiple calls to db.Connect() will not cause an error. A call to db.Close() will
require reconnecting before any additional queries are made. Arbitrary transactions to
the database can be executed by using db.BeginTx - the module guards a single connection
to the database from multiple go routines opening and closing access to the database.
*/
package db

import (
	"context"
	"database/sql"
	"errors"
	"net/url"
	"os"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	ro      bool         // if true, only allow database reads
	conn    *sql.DB      // connection to the database managed by the package
	connmu  sync.RWMutex // synchronize connect and close
	connect sync.Once    // ensure the database is only connected to once
)

var (
	ErrNotConnected   = errors.New("not connected to the database")
	ErrReadOnly       = errors.New("connected in read-only mode")
	ErrNotFound       = errors.New("record not found or no rows returned")
	ErrCannotParseDSN = errors.New("could not parse dsn, specify scheme:///path/to/data.db")
	ErrUnknownScheme  = errors.New("must specify a sqlite3 DSN")
)

// Connect to the sqlite3 database specified by the DSN. Connecting in readonly mode is
// managed by the package, not the database and is enforced by package functions.
// Subsequent calls to Connect will be ignored even if a different DSN or readonly mode
// is passed to the function.
func Connect(dsn string, readonly bool) (err error) {
	connmu.Lock()
	defer connmu.Unlock()

	connect.Do(func() {
		// Parse the DSN to get the path to the sqlite3 file
		var uri *DSN
		if uri, err = ParseDSN(dsn); err != nil {
			return
		}

		// TODO: do a better job of handling the DSN scheme
		if uri.Scheme != "sqlite3" {
			err = ErrUnknownScheme
		}

		// Check if the file exists, if it doesn't exist it will be created and all
		// migrations will be applied to the database. Otherwise the code will attempt
		// to only apply migrations that have not yet been applied.
		empty := false
		if _, err := os.Stat(uri.Path); os.IsNotExist(err) {
			empty = true
		}

		// Connect to the database
		ro = readonly
		if conn, err = sql.Open("sqlite3", uri.Path); err != nil {
			return
		}

		// Ensure the schema is initialized
		if err = InitializeSchema(empty); err != nil {
			return
		}
	})

	return err
}

// Close the database safely and allow for reconnect after close by resetting the
// package variables. No errors occur if the database is not connected.
func Close() (err error) {
	connmu.Lock()
	if conn != nil {
		err = conn.Close()
		conn = nil
		connect = sync.Once{}
	}
	connmu.Unlock()
	return err
}

// BeginTx creates a transaction with the connected database but returns an error if the
// database is not connected. If the database is set to readonly mode and the
// transaction options are not readonly, an error is returned.
func BeginTx(ctx context.Context, opts *sql.TxOptions) (tx *sql.Tx, err error) {
	connmu.RLock()
	defer connmu.RUnlock()
	if conn == nil {
		return nil, ErrNotConnected
	}

	if opts == nil {
		opts = &sql.TxOptions{ReadOnly: ro}
	} else if ro && !opts.ReadOnly {
		return nil, ErrReadOnly
	}

	return conn.BeginTx(ctx, opts)
}

// DSN represents the parsed components of an embedded database service.
type DSN struct {
	Scheme string
	Path   string
}

// DSN parsing and handling
func ParseDSN(uri string) (_ *DSN, err error) {
	var dsn *url.URL
	if dsn, err = url.Parse(uri); err != nil {
		return nil, err
	}

	if dsn.Scheme == "" || dsn.Path == "" {
		return nil, ErrCannotParseDSN
	}

	return &DSN{
		Scheme: dsn.Scheme,
		Path:   strings.TrimPrefix(dsn.Path, "/"),
	}, nil
}
