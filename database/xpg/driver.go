package xpg

import (
	"net"

	"github.com/go-pg/pg/v10"
	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/database"
)

var defaultDriver = NewDriver()

// Compile time check for the xservice.Driver implementation.
var _ database.Driver = (*PGDriver)(nil)

// PGDriver is an implementation of the xservice.Driver for the go-pg.
type PGDriver struct {
	mp map[string]cgerrors.ErrorCode
}

// DriverName gets the name of the PGDriver.
func (p *PGDriver) DriverName() string {
	return "go-pg"
}

// NewDriver creates a new driver implementation for the go-pg.
func NewDriver() *PGDriver {
	mp := map[string]cgerrors.ErrorCode{}
	for k, v := range errorMap {
		mp[k] = v
	}
	return &PGDriver{mp: mp}
}

// CustomErrorCode overwrites default error map for given class, which would result in given code.
func (p *PGDriver) CustomErrorCode(class string, code cgerrors.ErrorCode) {
	p.mp[class] = code
}

// ErrorCode implements xservice.Driver interface.
func (p *PGDriver) ErrorCode(err error) cgerrors.ErrorCode {
	switch err {
	case pg.ErrNoRows:
		return cgerrors.ErrorCode_NotFound
	case pg.ErrMultiRows, pg.ErrTxDone:
		return cgerrors.ErrorCode_Internal
	}
	e, ok := err.(pg.Error)
	if !ok {
		return cgerrors.ErrorCode_Unknown
	}
	class := e.Field('C')
	if class == "" {
		return cgerrors.ErrorCode_Unknown
	}
	code, ok := p.mp[class]
	if ok {
		return code
	}
	if len(class) < 2 {
		return cgerrors.ErrorCode_Unknown
	}
	code, ok = p.mp[class[:2]]
	if !ok {
		return cgerrors.ErrorCode_Internal
	}
	return code
}

// CanRetry implements xservice.Driver interface.
func (p *PGDriver) CanRetry(err error) bool {
	switch e := err.(type) {
	case pg.Error:
		class := e.Field('C')
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
func (p *PGDriver) Err(err error) *cgerrors.Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*cgerrors.Error); ok {
		return e
	}
	return cgerrors.New("", err.Error(), p.ErrorCode(err))
}
