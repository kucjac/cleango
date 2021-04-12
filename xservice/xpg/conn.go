package xpg

import (
	"fmt"

	"github.com/kucjac/cleango/errors"
)

// Logger defining format for various loggers
type Logger interface {
	pg.QueryHook
}

// Conn provides *pg.DB - a database handle
// representing a pool of zero or more underlying connections.
// It's safe for concurrent use by multiple goroutines.
type Conn struct {
	config *pg.Options
	conn   *pg.DB
}

// Parse will parse url and return options
func Parse(dial string) (*pg.Options, error) {
	opts, err := pg.ParseURL(dial)
	if err != nil {
		return opts, errors.ErrInternalf("failed to parse %s: %s", dial, err)
	}
	return opts, nil
}

// NewConn creates a new Conn.
func NewConn(config *pg.Options) *Conn {
	conn := &Conn{config: config, conn: pg.Connect(config)}
	return conn
}

// AddLogger will add xlog to database
func (c *Conn) AddLogger(log Logger) {
	c.conn.AddQueryHook(log)
}

// Get returns a database handle representing a pool of zero
// or more underlying connections.
// It's safe for concurrent use by multiple goroutines.
func (c *Conn) Get() *pg.DB {
	return c.conn
}

// Close closes the database client, releasing any open resources.
// It is rare to Close a DB, as the DB handle is meant
// to be long-lived and shared between many goroutines.
func (c *Conn) Close() {
	c.conn.Close()
}
