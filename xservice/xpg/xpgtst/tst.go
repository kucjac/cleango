package xpgtst

import (
	"os"
	"testing"

	"github.com/go-pg/pg/v10"
	"github.com/stretchr/testify/require"

	"github.com/kucjac/cleango/xservice/xpg"
)

// Conn creates new real connection to database
// and t.SkipNow() if is not parsable or missing
func Conn(t testing.TB) *xpg.Conn {
	env := os.Getenv("PAL_POSTGRES_URI")
	if env == "" {
		t.Logf("skipping test invalid uri %q", env)
		t.SkipNow()
	}
	opts, err := xpg.Parse(env)
	if err != nil {
		t.Logf("failed to parse postgres uri %q: %s", env, err.Error())
		t.SkipNow()
	}
	conn := xpg.NewConn(opts)
	if testing.Verbose() {
		conn.AddLogger(xpg.NewLogger(0))
	}
	return conn
}

// TestingConn is struct helping to keep
// conn to database and help to close it when needed.
type TestingConn struct {
	test testing.TB
	conn *xpg.Conn
	db   *pg.DB
}

// NewConn creates new testing connection.
func NewConn(t testing.TB) *TestingConn {
	return &TestingConn{test: t, conn: Conn(t)}
}

// Get return conn to database which is type orm.DB
func (t *TestingConn) Get() *pg.DB {
	t.db = t.conn.Get()
	return t.db
}

// Close will close connection to db.
func (t *TestingConn) Close() {
	require.NoError(t.test, t.db.Close())
}

// TestingConnTx is struct helping to keep
// conn to database and help to close it when needed.
type TestingConnTx struct {
	test testing.TB
	conn *xpg.Conn
	db   *pg.DB
	tx   *pg.Tx
}

// NewConnTx creates new testing connection. Also create
// transaction which will be closed one .Close() is called
func NewConnTx(t testing.TB) *TestingConnTx {
	return &TestingConnTx{test: t, conn: Conn(t)}
}

// Get return conn to database which is type orm.DB
func (t *TestingConnTx) Get() *pg.Tx {
	t.db = t.conn.Get()
	var err error
	t.tx, err = t.db.Begin()
	require.NoError(t.test, err)
	return t.tx
}

// Close will close connection to db and rollback transaction
func (t *TestingConnTx) Close() {
	// if the 'tx' connection was used with db.RunInTransaction
	// then it would be automatically closed and set to nil.
	// It causes panics
	if t.tx == nil {
		return
	}
	err := t.tx.Rollback()
	require.NoError(t.test, err)
	err = t.db.Close()
	require.NoError(t.test, err)
}
