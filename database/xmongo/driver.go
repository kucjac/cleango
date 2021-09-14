package xmongo

import (
	"context"
	"errors"

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
		return cgerrors.Wrap(err, cgerrors.ErrorCode_DeadlineExceeded, "context deadline exceeded")
	}
	if mongo.IsDuplicateKeyError(err) {
		return cgerrors.ErrAlreadyExists(err)
	}
	if mongo.IsNetworkError(err) {
		return cgerrors.New("", err.Error(), cgerrors.ErrorCode_Unavailable)
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return cgerrors.ErrNotFound(err.Error())
	}
	return cgerrors.ErrUnknown(err)
}

// ErrorCode implements database.Driver interface.
func (d Driver) ErrorCode(err error) cgerrors.ErrorCode {
	if errors.Is(err, context.Canceled) {
		return cgerrors.ErrorCode_Canceled
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return cgerrors.ErrorCode_DeadlineExceeded
	}
	if mongo.IsTimeout(err) {
		return cgerrors.ErrorCode_DeadlineExceeded
	}
	if mongo.IsDuplicateKeyError(err) {
		return cgerrors.ErrorCode_AlreadyExists
	}
	if mongo.IsNetworkError(err) {
		return cgerrors.ErrorCode_Unavailable
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return cgerrors.ErrorCode_NotFound
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
