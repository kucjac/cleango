package xpg

import (
	"context"

	"github.com/go-pg/pg/v10"

	"github.com/kucjac/cleango/errors"
)

const (
	errCancelled = "57014"
)

// ErrWrapf takes error returned from Postgres and returns structured error.
// Basically, we work with three database related kinds of errors:
// ErrNotFound for empty results,
// Coflict for insert/update violating data integrity,
// ErrInternal for all other kind of db errors.
func ErrWrapf(e error, process string, fmt string, args ...interface{}) error {
	// it is better to handle this directly in db functions
	// with more details about what was not found - see GetByID functions.
	if errors.Is(e, pg.ErrNoRows) {
		return errors.ErrNotFoundf(fmt, args...).WithProcess(process)
	}
	switch typed := e.(type) {
	case nil:
		return nil
	case pg.Error:
		if typed.IntegrityViolation() {
			return errors.ErrAlreadyExistsf(fmt, args...).WithProcess(process)
		}
		if typed.Field('C') == errCancelled {
			fmt += " %v"
			args = append(args, context.Canceled)
			return errors.ErrDeadlineExceededf(fmt, args).WithProcess(process)
		}
		args = append(args, e.Error())
		fmt = fmt + " %s"
		return errors.ErrInternalf(fmt, args...).WithProcess(process)
	case *errors.Error:
		return typed
	default:
		args = append(args, e.Error())
		fmt = fmt + " %s"
		return errors.ErrInternalf(fmt, args...).WithProcess(process)
	}
}

// ErrWrap takes error returned from Postgres and returns structured error.
// Basically, we work with three database related kinds of errors:
// ErrNotFound for empty results,
// Coflict for insert/update violating data integrity,
// ErrInternal for all other kind of db errors.
func ErrWrap(e error, process string, args ...interface{}) error {
	// it is better to handle this directly in db functions
	// with more details about what was not found - see GetByID functions.
	if errors.Is(e, pg.ErrNoRows) {
		return errors.ErrNotFound(args...)
	}
	switch typed := e.(type) {
	case nil:
		return nil
	case pg.Error:
		if typed.IntegrityViolation() {
			return errors.ErrAlreadyExists(args...)
		}
		if typed.Field('C') == errCancelled {
			return errors.ErrDeadlineExceeded(context.Canceled.Error())
		}
		return errors.ErrInternal(append([]interface{}{args}, e)...)
	case *errors.Error:
		return typed
	default:
		return errors.ErrInternal(append([]interface{}{args}, e)...)

	}
}
