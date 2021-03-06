package xsql

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/database"
	"github.com/kucjac/cleango/internal/uniqueid"
	"github.com/kucjac/cleango/pkg/xlog"
)

var txIdGen = uniqueid.NextBaseGenerator("xsql")

// Connect establish a new database connection using provided driverName and given dataSourceName (DSN).
func Connect(driver database.Driver, dataSourceName string) (*Conn, error) {
	conn, err := sqlx.Connect(driver.DriverName(), dataSourceName)
	return &Conn{
		db:     conn,
		driver: driver,
		config: &Config{
			LongQueriesTime: time.Millisecond * 100,
			WarnLongQueries: true,
		},
	}, err
}

// Compile time check if the Conn implements DB interface.
var _ DB = (*Conn)(nil)

// Conn is the database connection.
type Conn struct {
	db     *sqlx.DB
	driver database.Driver
	config *Config
}

// SetLongQueriesTime sets the duration at which the queries would become treated as long-running.
func (c *Conn) SetLongQueriesTime(d time.Duration) {
	c.config.LongQueriesTime = d
}

// WarnOnLongQueries logs warnings if the queries are long-running.
func (c *Conn) WarnOnLongQueries(warn bool) {
	c.config.WarnLongQueries = warn
}

// ErrorCode gets the error code related to given database result error.
func (c *Conn) ErrorCode(err error) cgerrors.ErrorCode {
	return c.driver.ErrorCode(err)
}

// CanRetry checks if the query done by given connection could be retried.
func (c *Conn) CanRetry(err error) bool {
	return c.driver.CanRetry(err)
}

// RunInTransaction runs provided function within a new database transaction.
func (c *Conn) RunInTransaction(ctx context.Context, fn func(tx *Tx) error) error {
	tx, err := c.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	return tx.RunInTransaction(fn)
}

// As sets the input in to one of the following types:
//	- *sqlx.DB
//	- *sql.DB
func (c *Conn) As(in interface{}) error {
	switch it := in.(type) {
	case **sqlx.DB:
		*it = c.db
	case **sql.DB:
		*it = c.db.DB
	default:
		return cgerrors.ErrInternalf("xsql.Conn.As provided invalid input type: %T", in)
	}
	return nil
}

// Begin starts a new transaction.
func (c *Conn) Begin() (*Tx, error) {
	id := txIdGen.NextId()

	ts := time.Now()
	logQuery(id, "BEGIN", ts, c.config, nil)

	tx, err := c.db.Beginx()
	if err != nil {
		return nil, err
	}

	return c.beginTx(tx, id), nil
}

// BeginTx starts a new transaction.
// The provided context is used until the transaction is committed or rolled
// back. If the context is canceled, the sql package will roll back the
// transaction. Tx.Commit will return an error if this context is canceled.
func (c *Conn) BeginTx(ctx context.Context, options *sql.TxOptions) (*Tx, error) {
	id := txIdGen.NextId()
	ts := time.Now()
	logQuery(id, "BEGIN", ts, c.config, nil)

	tx, err := c.db.BeginTxx(ctx, options)
	if err != nil {
		return nil, err
	}
	return c.beginTx(tx, id), nil
}

// Ping verifies a connection to the database is still alive, establishing a connection if necessary.
func (c *Conn) Ping() error {
	return c.db.Ping()
}

// PingContext verifies a connection to the database is still alive, establishing a connection if necessary.
func (c *Conn) PingContext(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

// DriverName gets the name of the driver provided during establishing connection.
func (c *Conn) DriverName() string {
	return c.db.DriverName()
}

// BindNamed binds the named arguments.
func (c *Conn) BindNamed(q string, arg interface{}) (query string, args []interface{}, err error) {
	return c.db.BindNamed(q, arg)
}

// Rebind change input query bindings to the ones that matches given database driver.
func (c *Conn) Rebind(query string) string {
	n := c.db.Rebind(query)
	return n
}

// QueryContext queries the database and returns an *xsql.Rows. Any placeholder parameters are replaced with supplied args.
func (c *Conn) QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	ts := time.Now()
	defer logQuery("", query, ts, c.config, nil, args...)

	rows, err := c.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return (*Rows)(rows), nil
}

// QueryRowContext queries the database and returns an *xsql.Row. Any placeholder parameters are replaced with supplied args.
func (c *Conn) QueryRowContext(ctx context.Context, query string, args ...interface{}) *Row {
	ts := time.Now()
	defer logQuery("", query, ts, c.config, nil, args...)
	return (*Row)(c.db.QueryRowxContext(ctx, query, args...))
}

// Query queries the database and returns an *xsql.Rows. Any placeholder parameters are replaced with supplied args.
func (c *Conn) Query(query string, args ...interface{}) (*Rows, error) {
	return c.QueryContext(context.Background(), query, args...)
}

// QueryRow queries the database and returns an *xsql.Row. Any placeholder parameters are replaced with supplied args.
func (c *Conn) QueryRow(query string, args ...interface{}) *Row {
	return c.QueryRowContext(context.Background(), query, args...)
}

// ExecContext executes provided query with the input arguments.
// The connection is aware of given context.
func (c *Conn) ExecContext(ctx context.Context, query string, args ...interface{}) (res sql.Result, err error) {
	ts := time.Now()
	res, err = c.db.ExecContext(ctx, query, args...)
	logQuery("", query, ts, c.config, res, args...)
	return res, err
}

// Exec execute provided query with the input arguments.
func (c *Conn) Exec(query string, args ...interface{}) (sql.Result, error) {
	return c.ExecContext(context.Background(), query, args...)
}

// PrepareContext creates a prepared statement.
// Provided context is used for the preparation of the statement, not for the execution of the statement.
func (c *Conn) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	xlog.Debugf("query: %s prepared", query)
	stmt, err := c.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return newStmt(stmt, query, c.config), nil
}

// Prepare creates a prepared statement.
func (c *Conn) Prepare(query string) (*Stmt, error) {
	return c.PrepareContext(context.Background(), query)
}

func (c *Conn) beginTx(tx *sqlx.Tx, id string) *Tx {
	return &Tx{id: id, tx: tx, driver: c.driver, config: c.config}
}
