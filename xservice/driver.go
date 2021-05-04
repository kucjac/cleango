package xservice

import (
	"github.com/kucjac/cleango/cgerrors"
)

// Driver is an interface used by the database
type Driver interface {
	ErrorCode(err error) cgerrors.ErrorCode
	CanRetry(err error) bool
}
