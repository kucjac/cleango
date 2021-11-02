package xmongo

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/kucjac/cleango/cgerrors"
)

// Driver is the implementation of the database.Driver for the MongoDB.
type Driver struct{}

// Err gets the standard error conversion.
func (d Driver) Err(err error) error {
	if errors.Is(err, context.Canceled) {
		return err
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	if mongo.IsTimeout(err) {
		return cgerrors.Wrap(err, cgerrors.CodeDeadlineExceeded, "context deadline exceeded")
	}
	if mongo.IsDuplicateKeyError(err) {
		return cgerrors.Wrap(err, cgerrors.CodeAlreadyExists, "duplicated key")
	}
	if mongo.IsNetworkError(err) {
		return cgerrors.New("", err.Error(), cgerrors.CodeUnavailable)
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return cgerrors.ErrNotFound(err.Error())
	}

	return cgerrors.ErrUnknown(err)
}

// ErrWrap wraps an error, checks its related code and sets the detail as the args.
func (d Driver) ErrWrap(err error, args ...interface{}) error {
	code := d.ErrorCode(err)
	return cgerrors.Wrap(err, code, fmt.Sprint(args...))
}

// ErrWrapf wraps an error and sets current error formatted details.
func (d Driver) ErrWrapf(err error, format string, args ...interface{}) error {
	code := d.ErrorCode(err)
	return cgerrors.Wrapf(err, code, format, args...)
}

// ErrorCode implements database.Driver interface.
func (d Driver) ErrorCode(err error) cgerrors.ErrorCode {
	if errors.Is(err, context.Canceled) {
		return cgerrors.CodeCanceled
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return cgerrors.CodeDeadlineExceeded
	}
	if mongo.IsTimeout(err) {
		return cgerrors.CodeDeadlineExceeded
	}
	if mongo.IsDuplicateKeyError(err) {
		return cgerrors.CodeAlreadyExists
	}
	if mongo.IsNetworkError(err) {
		return cgerrors.ErrorCode_Unavailable
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return cgerrors.ErrorCode_NotFound
	}

	if code := cgerrors.Code(err); code != cgerrors.ErrorCode_Unknown {
		return code
	}
	return cgerrors.ErrorCode_Unknown
}

// DriverName implements database.Driver interface.
func (d Driver) DriverName() string {
	return "mongo"
}

// CanRetry implements database.Driver interface.
func (d Driver) CanRetry(err error) bool {
	if mongo.IsNetworkError(err) {
		return true
	}
	return false
}
