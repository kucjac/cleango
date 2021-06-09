package xsql

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/kucjac/cleango/cgerrors"
)

// Rows is the simple type wrapper for the sqlx.Rows that provides minor enhancement.
type Rows sqlx.Rows

// As extracts the *sqlx.Rows or *sql.Rows implementation by providing a pointer to it.
func (r *Rows) As(in interface{}) error {
	switch it := in.(type) {
	case **sqlx.Rows:
		*it = (*sqlx.Rows)(r)
	case **sql.Rows:
		*it = (*sqlx.Rows)(r).Rows
	default:
		return cgerrors.ErrInternalf("xsql.Rows.As invalid input type: %T", in)
	}
	return nil
}

// Next prepares the next result row for reading with the Scan method. It returns true on success, or false if
// there is no next result row, or an error happened while preparing it. Err should be consulted to distinguish
// between the two cases.
// Every call to Scan, even the first one, must be preceded by a call to Next.
func (r *Rows) Next() bool {
	return (*sqlx.Rows)(r).Next()
}

// NextResultSet prepares the next result set for reading. It reports whether there is further result sets, or false
// if there is no further result set or if there is an error advancing to it. The Err method should be consulted to
// distinguish between the two cases.
// After calling NextResultSet, the Next method should always be called before scanning.
// If there are further result sets they may not have rows in the result set.
func (r *Rows) NextResultSet() bool {
	return (*sqlx.Rows)(r).NextResultSet()
}

// Err returns the error, if any, that was encountered during iteration.
// It may be called after an explicit or implicit Close.
func (r *Rows) Err() error {
	return (*sqlx.Rows)(r).Err()
}

// Columns returns the column names. It returns an error if the rows are closed
func (r *Rows) Columns() ([]string, error) {
	return (*sqlx.Rows)(r).Columns()
}

// ColumnTypes returns column information such as column type, length, and nullable.
// Some information may not be available from some drivers.
func (r *Rows) ColumnTypes() ([]*sql.ColumnType, error) {
	return (*sqlx.Rows)(r).ColumnTypes()
}

// Scan copies the columns in the current row into the values pointed at by dest.
// The number of values in dest must be the same as the number of columns in Rows.
// It converts columns read from the database into the following common Go types and
// special types provided by the sql package:
// - *string
// - *[]byte
// - *int, *int8, *int16, *int32, *int64
// - *uint, *uint8, *uint16, *uint32, *uint64
// - *bool
// - *float32, *float64
// - *interface{}
// - *RawBytes
// - *Rows (cursor value)
// any type implementing Scanner (see Scanner docs)
// In the most simple case, if the type of the value from the source column is an integer, bool
// or string type T and dest is of type *T, Scan simply assigns the value through the pointer.
// Scan also converts between string and numeric types, as long as no information would be lost.
// While Scan stringifies all numbers scanned from numeric database columns into *string, scans into
// numeric types are checked for overflow. For example, a float64 with value 300, or a string with value "300"
// can scan into a uint16, but not into a uint8, though float64(255) or "255" can scan into an uint8.
// One exception is that scans of some float64 numbers to strings may lose information when stringifying.
// In general, scan floating point columns into *float64.
// If a dest argument has type *[]byte, Scan saves in that argument a copy of the corresponding data.
// The copy is owned by the caller and can be modified and held indefinitely. The copy can be avoided by using
// an argument of type *RawBytes instead; see the documentation for RawBytes for restrictions on its use.
// If an argument has type *interface{}, Scan copies the value provided by the underlying driver without conversion.
// When scanning from a source value of type []byte to *interface{}, a copy of the slice is made,
// and the caller owns the result.
// Source values of type time.Time may be scanned into values of type *time.Time, *interface{}, *string, or *[]byte.
// When converting to the latter two, time.RFC3339Nano is used.
// Source values of type bool may be scanned into types *bool, *interface{}, *string, *[]byte, or *RawBytes.
// For scanning into *bool, the source may be true, false, 1, 0, or string inputs parseable by strconv.ParseBool.
// Scan can also convert a cursor returned from a query, such as "select cursor(select * from my_table) from dual",
// into a *Rows value that can itself be scanned from. The parent select query will close any cursor *Rows if
// the parent *Rows is closed.
// If any of the first arguments implementing Scanner returns an error, that error will be wrapped in the returned error.
func (r *Rows) Scan(dest ...interface{}) error {
	return (*sqlx.Rows)(r).Scan(dest...)
}

// Close closes the Rows, preventing further enumeration.
// If Next is called and returns false and there are no further result sets,
// the Rows are closed automatically, and it will suffice to check the result of Err.
// Close is idempotent and does not affect the result of Err.
func (r *Rows) Close() error {
	return (*sqlx.Rows)(r).Close()
}

// SliceScan using this Rows.
func (r *Rows) SliceScan() ([]interface{}, error) {
	return (*sqlx.Rows)(r).SliceScan()
}

// MapScan using this Rows.
func (r *Rows) MapScan(dest map[string]interface{}) error {
	return (*sqlx.Rows)(r).MapScan(dest)
}

// StructScan is like sql.Rows.Scan, but scans a single Row into a single Struct.
// Use this and iterate over Rows manually when the memory load of Select() might be prohibitive.
// *Rows.StructScan caches the reflection work of matching up column positions to fields to avoid that overhead
// per scan, which means it is not safe to run StructScan on the same Rows instance with different struct types.
func (r *Rows) StructScan(dest interface{}) error {
	return (*sqlx.Rows)(r).StructScan(dest)
}
