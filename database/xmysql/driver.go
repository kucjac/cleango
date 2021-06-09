package xmysql

import (
	"database/sql"
	"database/sql/driver"
	"net"

	"github.com/go-sql-driver/mysql"
	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/database"
)

var defaultDriver = newDriver()

// Compile time check if the MySQLDriver implements database.Driver interface.
var _ database.Driver = (*MySQLDriver)(nil)

// NewDriver creates a new error driver.
func NewDriver() *MySQLDriver {
	return newDriver()
}

func newDriver() *MySQLDriver {
	cm := make(map[interface{}]cgerrors.ErrorCode)
	sm := make(map[uint16]string)
	for k, v := range mysqlErrMap {
		cm[k] = v
	}
	for k, v := range codeSQLState {
		sm[k] = v
	}
	return &MySQLDriver{
		codeMap:  cm,
		stateMap: sm,
	}
}

// MySQLDriver is the implementation of the MySQL based database.Driver.
type MySQLDriver struct {
	codeMap  map[interface{}]cgerrors.ErrorCode
	stateMap map[uint16]string
}

// DriverName implements database.Driver interface.
func (m *MySQLDriver) DriverName() string {
	return "mysql"
}

// CanRetry implements database.Driver interface.
func (m *MySQLDriver) CanRetry(err error) bool {
	switch e := err.(type) {
	case *mysql.MySQLError:
		switch e.Number {
		case 1053, 1317, 1290, 1836:
			return true
		default:
			return false
		}
	case *net.OpError:
		return true
	}

	switch err {
	case mysql.ErrInvalidConn, driver.ErrBadConn:
		return true
	}
	return false
}

// ErrorCode implements cgerrors.ErrorCoder interface.
func (m *MySQLDriver) ErrorCode(err error) cgerrors.ErrorCode {
	mySQLErr, ok := err.(*mysql.MySQLError)
	if !ok {
		// Otherwise check if it sql.Err* or other errors from mysql package
		switch err {
		case mysql.ErrInvalidConn, mysql.ErrNoTLS, mysql.ErrOldProtocol,
			mysql.ErrMalformPkt, sql.ErrTxDone:
			return cgerrors.ErrorCode_Internal
		case sql.ErrNoRows:
			return cgerrors.ErrorCode_NotFound
		default:
			return cgerrors.ErrorCode_Unknown
		}
	}

	// Check if Error Number is in recogniser
	c, ok := m.codeMap[mySQLErr.Number]
	if ok {
		// Return if found
		return c
	}

	// Otherwise check if given sqlstate is in the codeMap
	sqlState, ok := m.stateMap[mySQLErr.Number]
	if !ok || len(sqlState) != 5 {
		return cgerrors.ErrorCode_Unknown
	}
	c, ok = m.codeMap[sqlState]
	if ok {
		return c
	}

	// First two letter from sqlState represents error class
	// Check if class is in error map
	sqlStateClass := sqlState[0:2]
	c, ok = m.codeMap[sqlStateClass]
	if ok {
		return c
	}
	return cgerrors.ErrorCode_Unknown
}
