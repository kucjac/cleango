package xpq

import (
	"database/sql"
	"errors"
	"net"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/database"
	"github.com/kucjac/cleango/database/internal/postgres"
	"github.com/lib/pq"
)

var defaultDriver = NewDriver()

// Compile time check for the xservice.Driver implementation.
var _ database.Driver = (*Driver)(nil)

// Driver is an implementation of the xservice.Driver for the go-pg.
type Driver struct {
	mp map[string]cgerrors.ErrorCode
}

// DriverName gets the name of the driver.
func (d *Driver) DriverName() string {
	return "postgres"
}

// DefaultDriver gets the default lib/pq driver.
func DefaultDriver() *Driver {
	return defaultDriver
}

// NewDriver creates a new driver implementation for the go-pg.
func NewDriver() *Driver {
	mp := map[string]cgerrors.ErrorCode{}
	for k, v := range postgres.ErrorMap {
		mp[k] = v
	}
	return &Driver{mp: mp}
}

// CustomErrorCode overwrites default error map for given class, which would result in given code.
func (d *Driver) CustomErrorCode(class string, code cgerrors.ErrorCode) {
	d.mp[class] = code
}

// ErrorCode implements xservice.Driver interface.
func (d *Driver) ErrorCode(err error) cgerrors.ErrorCode {
	if errors.Is(err, sql.ErrNoRows) {
		return cgerrors.CodeNotFound
	}
	if errors.Is(err, sql.ErrConnDone) || errors.Is(err, sql.ErrTxDone) {
		return cgerrors.CodeUnavailable
	}

	if code := cgerrors.Code(err); code != cgerrors.CodeUnknown {
		return code
	}

	e, ok := err.(*pq.Error)
	if !ok {
		return cgerrors.CodeUnknown
	}
	code, ok := d.mp[string(e.Code)]
	if ok {
		return code
	}
	class := e.Code.Class()
	if class == "" {
		return cgerrors.CodeUnknown
	}
	if len(class) < 2 {
		return cgerrors.CodeUnknown
	}
	code, ok = d.mp[string(class)]
	if !ok {
		return cgerrors.CodeInternal
	}
	return code
}

// CanRetry implements xservice.Driver interface.
func (d *Driver) CanRetry(err error) bool {
	switch e := err.(type) {
	case *pq.Error:
		class := e.Code.Class()
		if len(class) >= 2 {
			return class[:2] == "08"
		}
		return false
	case *net.OpError:
		return true
	}
	return false
}

// Err converts given error into a cgerrors.Error.
func (d *Driver) Err(err error) *cgerrors.Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*cgerrors.Error); ok {
		return e
	}
	return cgerrors.New("", err.Error(), d.ErrorCode(err))
}
