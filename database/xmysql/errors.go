package xmysql

import (
	"github.com/kucjac/cleango/cgerrors"
)

// IsDuplicatedError checks if given error states for the mysql duplicated input error.
func IsDuplicatedError(err error) bool {
	return ErrorCode(err) == cgerrors.ErrorCode_AlreadyExists
}

// ErrorCode gets the mysql based error code from input error.
func ErrorCode(err error) cgerrors.ErrorCode {
	return defaultDriver.ErrorCode(err)
}

// CanRetry checks if a query could be retried on the base of given error.
func CanRetry(err error) bool {
	return defaultDriver.CanRetry(err)
}
