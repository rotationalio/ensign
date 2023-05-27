/*
Package sqlite implements a connect hook around the sqlite3 driver so that the
underlying connection can be fetched from the driver for more advanced operations such
as backups. See: https://github.com/mattn/go-sqlite3/blob/master/_example/hook/hook.go
*/
package sqlite

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"sync"

	"github.com/mattn/go-sqlite3"
)

func init() {
	conns = make(map[uint64]*Conn)
	sql.Register(DriverName, &Driver{})
}

const (
	DriverName = "ensign_sqlite3"
)

var (
	seq   uint64
	mu    sync.Mutex
	conns map[uint64]*Conn
)

type Driver struct {
	sqlite3.SQLiteDriver
}

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

type Conn struct {
	cid uint64
	*sqlite3.SQLiteConn
}

func (c *Conn) Close() error {
	mu.Lock()
	delete(conns, c.cid)
	mu.Unlock()
	return c.SQLiteConn.Close()
}

func GetLastConn() (*Conn, bool) {
	mu.Lock()
	defer mu.Unlock()
	conn, ok := conns[seq]
	return conn, ok
}

func NumConns() int {
	mu.Lock()
	defer mu.Unlock()
	return len(conns)
}
