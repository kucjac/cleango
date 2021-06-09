package xsql

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/kucjac/cleango/cgerrors"
)

// Row is the row type wrapper of the sqlx.Row.
type Row sqlx.Row

// As extracts the *sqlx.Row from given Row type.
func (r *Row) As(in interface{}) error {
	switch it := in.(type) {
	case **sqlx.Row:
		*it = (*sqlx.Row)(r)
	default:
		return cgerrors.ErrInternalf("xsql.Row.As invalid input type: %T", in)
	}
	return nil
}

// Scan is a fixed implementation of sql.Row.Scan,
// which does not discard the underlying error from the
// internal rows object if it exists.
func (r *Row) Scan(dest ...interface{}) error {
	return (*sqlx.Row)(r).Scan(dest...)
}

// Columns returns the underlying sql.Rows.Columns(), or
// the deferred error usually returned by Row.Scan()
func (r *Row) Columns() ([]string, error) {
	return (*sqlx.Row)(r).Columns()
}

// ColumnTypes returns the underlying sql.Rows.ColumnTypes(), or
// the deferred error.
func (r *Row) ColumnTypes() ([]*sql.ColumnType, error) {
	return (*sqlx.Row)(r).ColumnTypes()
}

// Err returns the error encountered while scanning.
func (r *Row) Err() error {
	return (*sqlx.Row)(r).Err()
}

// SliceScan scans the row into a slice.
func (r *Row) SliceScan() ([]interface{}, error) {
	return (*sqlx.Row)(r).SliceScan()
}

// MapScan scans the row into an input dest map.
func (r *Row) MapScan(dest map[string]interface{}) error {
	return (*sqlx.Row)(r).MapScan(dest)
}

// StructScan scans the row as a structure provided as argument.
func (r *Row) StructScan(dest interface{}) error {
	return (*sqlx.Row)(r).StructScan(dest)
}
