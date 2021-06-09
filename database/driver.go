package database

import (
	"github.com/kucjac/cleango/cgerrors"
)

// Driver is an interface used by the database
type Driver interface {
	ErrorCode(err error) cgerrors.ErrorCode
	DriverName() string
	CanRetry(err error) bool
}
