package xsql

import (
	"context"
	"database/sql"

	"github.com/kucjac/cleango/cgerrors"
)

// DB is the common interface for both the Conn and Tx.
type DB interface {
	// QueryContext queries the database and returns an *xsql.Rows. Any placeholder parameters are replaced with supplied args.
	QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error)
	// QueryRowContext queries the database and returns an *xsql.Row. Any placeholder parameters are replaced with supplied args.
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *Row

	// Query queries the database and returns an *xsql.Rows. Any placeholder parameters are replaced with supplied args.
	Query(query string, args ...interface{}) (*Rows, error)
	// QueryRow queries the database and returns an *xsql.Row. Any placeholder parameters are replaced with supplied args.
	QueryRow(query string, args ...interface{}) *Row

	// ExecContext executes provided query with the input arguments.
	// The connection is aware of given context.
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	// Exec execute provided query with the input arguments.
	Exec(query string, args ...interface{}) (sql.Result, error)

	// PrepareContext creates a prepared statement.
	// Provided context is used for the preparation of the statement, not for the execution of the statement.
	PrepareContext(ctx context.Context, query string) (*Stmt, error)
	// Prepare creates a prepared statement.
	Prepare(query string) (*Stmt, error)

	// As extracts the types on which given implementation is based on.
	// I.e: Tx accepts: **sqlx.Tx or **sql.Tx.
	As(in interface{}) error

	// ErrorCode gets the error code for related database error.
	ErrorCode(err error) cgerrors.ErrorCode
	// CanRetry checks if the query could be retried for given error.
	CanRetry(err error) bool

	// Rebind changes argument format in provided query.
	Rebind(query string) string
}
