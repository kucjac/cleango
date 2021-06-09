package xpq

import (
	"database/sql"
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
func (p *Driver) CustomErrorCode(class string, code cgerrors.ErrorCode) {
	p.mp[class] = code
}

// ErrorCode implements xservice.Driver interface.
func (p *Driver) ErrorCode(err error) cgerrors.ErrorCode {
	switch err {
	case sql.ErrNoRows:
		return cgerrors.ErrorCode_NotFound
	case sql.ErrConnDone, sql.ErrTxDone:
		return cgerrors.ErrorCode_Internal
	}
	e, ok := err.(*pq.Error)
	if !ok {
		return cgerrors.ErrorCode_Unknown
	}
	code, ok := p.mp[string(e.Code)]
	if ok {
		return code
	}
	class := e.Code.Class()
	if class == "" {
		return cgerrors.ErrorCode_Unknown
	}
	if len(class) < 2 {
		return cgerrors.ErrorCode_Unknown
	}
	code, ok = p.mp[string(class)]
	if !ok {
		return cgerrors.ErrorCode_Internal
	}
	return code
}

// CanRetry implements xservice.Driver interface.
func (p *Driver) CanRetry(err error) bool {
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
func (p *Driver) Err(err error) *cgerrors.Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*cgerrors.Error); ok {
		return e
	}
	return cgerrors.New("", err.Error(), p.ErrorCode(err))
}
