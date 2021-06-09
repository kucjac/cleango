package xsql

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/xlog"
	"github.com/sirupsen/logrus"
)

// RunInTransaction executes given function based on provided 'db' within a transaction.
func RunInTransaction(ctx context.Context, db DB, fn func(tx *Tx) error) (err error) {
	switch dbt := db.(type) {
	case *Conn:
		err = dbt.RunInTransaction(ctx, fn)
	case *Tx:
		err = fn(dbt)
	default:
		return cgerrors.ErrInternalf("unknown xsql.DB implementation: %T", db)
	}
	return err
}

// Tx is the database connection.
type Tx struct {
	id string
	tx *sqlx.Tx
}

// RunInTransaction runs a function in the transaction. If function
// returns an error transaction is rolled back, otherwise transaction
// is committed.
func (tx *Tx) RunInTransaction(fn func(tx *Tx) error) error {
	defer func() {
		if err := recover(); err != nil {
			if err := tx.Rollback(); err != nil {
				xlog.Errorf("tx.Rollback panicked: %s", err)
			}
			panic(err)
		}
	}()
	if err := fn(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			xlog.Errorf("tx.Rollback failed: %v", err)
		}
		return err
	}
	return tx.Commit()
}

// As extracts the transaction base types like:
// - *sqlx.Tx
// - *sql.Tx
func (tx *Tx) As(in interface{}) error {
	switch it := in.(type) {
	case **sqlx.Tx:
		*it = tx.tx
	case **sql.Tx:
		*it = tx.tx.Tx
	default:
		return cgerrors.ErrInternalf("invalid xsql.Tx.As type: %T", in)
	}
	return nil
}

// DriverName gets the name of the driver provided during establishing connection.
func (tx *Tx) DriverName() string {
	return tx.tx.DriverName()
}

// BindNamed binds the named arguments.
func (tx *Tx) BindNamed(s string, i interface{}) (string, []interface{}, error) {
	return tx.tx.BindNamed(s, i)
}

// Rebind change input query bindings to the ones that matches given database driver.
func (tx *Tx) Rebind(s string) string {
	n := tx.tx.Rebind(s)
	return n
}

// Query queries the database within given transaction and returns an *xsql.Rows. Any placeholder parameters are replaced with supplied args.
func (tx *Tx) Query(query string, args ...interface{}) (*Rows, error) {
	return tx.QueryContext(context.Background(), query, args...)
}

// QueryContext queries the database within given transaction and returns an *xsql.Rows. Any placeholder parameters are replaced with supplied args.
func (tx *Tx) QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	if xlog.IsLevelEnabled(logrus.DebugLevel) {
		ts := time.Now()
		defer logQuery(tx.id, query, ts, args...)
	}
	rows, err := tx.tx.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return (*Rows)(rows), nil
}

// QueryRow queries the database within given transaction and returns an *xsql.Row. Any placeholder parameters are replaced with supplied args.
func (tx *Tx) QueryRow(query string, args ...interface{}) *Row {
	return tx.QueryRowContext(context.Background(), query, args...)
}

// QueryRowContext queries the database within given transaction and returns an *xsql.Row. Any placeholder parameters are replaced with supplied args.
func (tx *Tx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *Row {
	if xlog.IsLevelEnabled(logrus.DebugLevel) {
		ts := time.Now()
		defer logQuery(tx.id, query, ts, args...)
	}
	return (*Row)(tx.tx.QueryRowxContext(ctx, query, args...))
}

// ExecContext execute provided query with the input arguments.
func (tx *Tx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if xlog.IsLevelEnabled(logrus.DebugLevel) {
		ts := time.Now()
		defer logQuery(tx.id, query, ts, args...)
	}
	return tx.tx.ExecContext(ctx, query, args...)
}

// Exec execute provided query with the input arguments.
func (tx *Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return tx.ExecContext(context.Background(), query, args...)
}

// PrepareContext creates a prepared statement.
// Provided context is used for the preparation of the statement, not for the execution of the statement.
func (tx *Tx) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	xlog.Debugf("query: %s prepared", query)
	stmt, err := tx.tx.PreparexContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return newTxStmt(stmt, tx.id, query), nil
}

// Prepare creates a prepared statement.
func (tx *Tx) Prepare(query string) (*Stmt, error) {
	return tx.PrepareContext(context.Background(), query)
}

// Commit commits this transaction.
func (tx *Tx) Commit() error {
	if xlog.IsLevelEnabled(logrus.DebugLevel) {
		ts := time.Now()
		defer logQuery(tx.id, "COMMIT", ts)
	}
	return tx.tx.Commit()
}

// Rollback aborts this transaction.
func (tx *Tx) Rollback() error {
	if xlog.IsLevelEnabled(logrus.DebugLevel) {
		ts := time.Now()
		defer logQuery(tx.id, "ROLLBACK", ts)
	}
	return tx.tx.Rollback()
}
