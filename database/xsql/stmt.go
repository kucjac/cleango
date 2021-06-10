package xsql

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kucjac/cleango/cgerrors"
)

// Stmt is the xsql type wrapper for the sqlx.Stmt.
type Stmt struct {
	stmt   *sqlx.Stmt
	txID   string
	query  string
	config *Config
}

// As extracts direct implementation of the stmt in the *sqlx.Stmt or *sql.Stmt.
func (s *Stmt) As(in interface{}) error {
	switch it := in.(type) {
	case **sqlx.Stmt:
		*it = s.stmt
	case **sql.Stmt:
		*it = s.stmt.Stmt
	default:
		return cgerrors.ErrInternalf("invalid *Stmt.As type: %T - expected: %T", in, (**sql.Stmt)(nil))
	}
	return nil
}

// ExecContext executes the statement with provided arguments.
func (s *Stmt) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	ts := time.Now()
	defer logQuery(s.txID, s.query, ts, s.config, args...)

	return s.stmt.ExecContext(ctx, args...)
}

// Exec executes the statement with provided arguments.
func (s *Stmt) Exec(args ...interface{}) (sql.Result, error) {
	return s.ExecContext(context.Background(), args...)
}

// QueryContext executes statement query with provided arguments.
func (s *Stmt) QueryContext(ctx context.Context, args ...interface{}) (*Rows, error) {
	ts := time.Now()
	defer logQuery(s.txID, s.query, ts, s.config, args...)

	rows, err := s.stmt.QueryxContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	return (*Rows)(rows), nil
}

// Query executes statement query with provided arguments.
func (s *Stmt) Query(args ...interface{}) (*Rows, error) {
	return s.QueryContext(context.Background(), args...)
}

// QueryRowContext executes a query with provided arguments and creates a new Row.
// The connection is based on given context.
func (s *Stmt) QueryRowContext(ctx context.Context, args ...interface{}) *Row {
	ts := time.Now()
	defer logQuery(s.txID, s.query, ts, s.config, args...)

	return (*Row)(s.stmt.QueryRowxContext(ctx, args...))
}

// QueryRow executes a query with provided arguments and creates a new Row.
func (s *Stmt) QueryRow(args ...interface{}) *Row {
	return s.QueryRowContext(context.Background(), args...)
}

// Close closes given statement.
func (s *Stmt) Close() error {
	return s.stmt.Close()
}

func newStmt(stmt *sqlx.Stmt, query string, config *Config) *Stmt {
	return &Stmt{stmt: stmt, query: query, config: config}
}

func newTxStmt(stmt *sqlx.Stmt, id, query string, config *Config) *Stmt {
	return &Stmt{stmt: stmt, query: query, txID: id, config: config}
}
