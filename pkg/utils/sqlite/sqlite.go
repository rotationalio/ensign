/*
Package sqlite implements a connect hook around the sqlite3 driver so that the
underlying connection can be fetched from the driver for more advanced operations such
as backups. See: https://github.com/mattn/go-sqlite3/blob/master/_example/hook/hook.go

To use make sure you import this package so that the init code registers the driver:

	import _ github.com/rotationalio/ensign/pkg/utils/sqlite

Then you can use sql.Open in the same way you would with sqlite3:

	sql.Open("ensign_sqlite3", "path/to/database.db")
*/
package sqlite

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"sync"

	"github.com/mattn/go-sqlite3"
)

// Init creates the connections map and registers the driver with the SQL package.
func init() {
	conns = make(map[uint64]*Conn)
	sql.Register(DriverName, &Driver{})
}

// In order to use this driver, specify the DriverName to sql.Open.
const (
	DriverName = "ensign_sqlite3"
)

var (
	seq   uint64
	mu    sync.Mutex
	conns map[uint64]*Conn
)

// Driver embeds a sqlite3 driver but overrides the Open function to ensure the
// connection created is a local connection with a sequence ID. It then maintains the
// connection locally until it is closed so that the underlying sqlite3 connection can
// be returned on demand.
type Driver struct {
	sqlite3.SQLiteDriver
}

// Open implements the sql.Driver interface and returns a sqlite3 connection that can
// be fetched by the user using GetLastConn. The connection ensures it's cleaned up
// when it's closed. This method is not used by the user, but rather by sql.Open.
func (d *Driver) Open(dsn string) (_ driver.Conn, err error) {
	var inner driver.Conn
	if inner, err = d.SQLiteDriver.Open(dsn); err != nil {
		return nil, err
	}

	var (
		ok    bool
		sconn *sqlite3.SQLiteConn
	)

	if sconn, ok = inner.(*sqlite3.SQLiteConn); !ok {
		return nil, fmt.Errorf("unknown connection type %T", inner)
	}

	mu.Lock()
	seq++
	conn := &Conn{cid: seq, SQLiteConn: sconn}
	conns[conn.cid] = conn
	mu.Unlock()

	return conn, nil
}

// Wraps a sqlite3.SQLiteConn and maintains an ID so that the connection can be closed.
type Conn struct {
	cid uint64
	*sqlite3.SQLiteConn
}

// When the connection is closed, it is removed from the array of connections.
func (c *Conn) Close() error {
	mu.Lock()
	delete(conns, c.cid)
	mu.Unlock()
	return c.SQLiteConn.Close()
}

// The entire point of this package is to provide access to SQLite3 backup functionality
// on the sqlite3 connection. For more details on how to use the backup see the
// following links:
//
// https://www.sqlite.org/backup.html
// https://github.com/mattn/go-sqlite3/blob/master/_example/hook/hook.go
// https://github.com/mattn/go-sqlite3/blob/master/backup_test.go
//
// This is primarily used by the backups package and this method provides access
// directly to the underlying CGO call. This means the CGO call must be called correctly
// for example: the Finish() method MUST BE CALLED otherwise your code will panic.
func (c *Conn) Backup(dest string, srcConn *Conn, src string) (*sqlite3.SQLiteBackup, error) {
	return c.SQLiteConn.Backup(dest, srcConn.SQLiteConn, src)
}

// GetLastConn returns the last connection created by the driver. Unfortunately, there
// is no way to guarantee which connection will be returned since the sql.Open package
// does not provide any interface to the underlying connection object. The best a
// process can do is ping the server to open a new connection and then fetch the last
// connection immediately.
func GetLastConn() (*Conn, bool) {
	mu.Lock()
	defer mu.Unlock()
	conn, ok := conns[seq]
	return conn, ok
}

// For testing purposes, returns the number of active connections.
func NumConns() int {
	mu.Lock()
	defer mu.Unlock()
	return len(conns)
}
