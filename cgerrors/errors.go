// Package cgerrors provides a way to return detailed information
// for an RPC request error. The error is normally JSON encoded.
package cgerrors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kucjac/cleango/internal/uniqueid"
	"google.golang.org/grpc/status"
)

// MetaKeyWrapped adds the error metadata key that keeps the wrapped error value.
const MetaKeyWrapped = "x-err-wrapped"

// Error is the error message that has id, it's code and a detail.
type Error struct {
	ID      string            `json:"id,omitempty"`
	Code    ErrorCode         `json:"code,omitempty"`
	Detail  string            `json:"detail,omitempty"`
	Process string            `json:"process,omitempty"`
	Meta    map[string]string `json:"meta,omitempty"`

	wrapped error
}

// Is implements errors.Is function interface to check if input error matches this error or the one wrapped.
func (e *Error) Is(err error) bool {
	return Equal(e, err)
}

// Unwrap implements errors.Unwrap function internal interface.
func (e *Error) Unwrap() error {
	return e.wrapped
}

// Error implements error interface.
func (e *Error) Error() string {
	temp := *e

	if e.wrapped != nil {
		meta := map[string]string{}
		for k, v := range e.Meta {
			meta[k] = v
		}
		meta[MetaKeyWrapped] = e.wrapped.Error()
		temp.Meta = meta
	}

	data, _ := json.Marshal(temp)
	return string(data)
}

var g = uniqueid.NextGenerator("errors")

// FromString parses string error into an Error structure.
func FromString(err string) (*Error, bool) {
	var e Error
	if er := json.Unmarshal([]byte(err), &e); er != nil {
		return nil, false
	}
	return &e, true
}

// WithMeta sets the key, value metadata for given error.
func (e *Error) WithMeta(key, value string) *Error {
	if e.Meta == nil {
		e.Meta = make(map[string]string)
	}
	e.Meta[key] = value
	return e
}

// WithProcess sets the process for given error.
func (e *Error) WithProcess(process string) *Error {
	e.Process = process
	return e
}

// WithCode sets the code for given error.
func (e *Error) WithCode(code ErrorCode) *Error {
	e.Code = code
	return e
}

// GRPCError is an interface used to get grpcStatus
type GRPCError interface {
	GRPCStatus() *status.Status
}

// ToGRPCError converts an error to GRPC status.Status.
func ToGRPCError(err error) error {
	var ge GRPCError
	if errors.As(err, &ge) {
		return ge.GRPCStatus().Err()
	}
	return newError(Code(err), err.Error()).GRPCStatus().Err()
}

// GRPCStatus implements grpc client interface used to convert statuses.
func (e *Error) GRPCStatus() *status.Status {
	return status.New(e.Code.ToGRPCCode(), e.Error())
}

// Is compares the errors with their values.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// To type check if given error is of *Error type or has encoded ErrorCode in it.
// Otherwise creates a new error with Unknown code.
func To(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	if c := Code(err); c != CodeUnknown {
		return newError(c, err.Error())
	}

	e = new(Error)
	es := err.Error()
	if err := json.Unmarshal([]byte(es), e); err == nil {
		e.Detail = es
	} else {
		e.Code = CodeUnknown
	}
	return e
}

// IsNotFound checks if given input error is of code NotFound.
func IsNotFound(err error) bool {
	return Code(err) == CodeNotFound
}

// IsAlreadyExists checks if given error means that given entity already exists.
func IsAlreadyExists(err error) bool {
	return Code(err) == CodeAlreadyExists
}

// IsInvalidArgument checks if given error contains a CodeInvalidArgument.
func IsInvalidArgument(err error) bool {
	return Code(err) == CodeInvalidArgument
}

// IsDeadlineExceeded checks if given error is of type Deadline Exceeded.
func IsDeadlineExceeded(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	return Code(err) == CodeDeadlineExceeded
}

// IsUnauthenticated checks if given error is an unauthenticated error.
func IsUnauthenticated(err error) bool {
	return Code(err) == CodeUnauthenticated
}

// IsInternal checks if the input is an internal error.
func IsInternal(err error) bool {
	return Code(err) == CodeInternal
}

// IsPermissionDenied checks if given error is of type CodePermissionDenied.
func IsPermissionDenied(err error) bool {
	return Code(err) == CodePermissionDenied
}

// IsUnimplemented checks if given error contains CodeUnimplemented.
func IsUnimplemented(err error) bool {
	return Code(err) == CodeUnimplemented
}

// IsUnknown checks if given error contains CodeUnknown.
func IsUnknown(err error) bool {
	return Code(err) == CodeUnknown
}

// New generates a custom error.
func New(id, detail string, code ErrorCode) *Error {
	e := &Error{
		ID:     id,
		Code:   code,
		Detail: detail,
	}
	if id == "" {
		e.setDefaultID()
	}
	return e
}

// Parse tries to parse a JSON string into an error. If that
// fails, it will set the given string as the error detail.
func Parse(err string) *Error {
	e := new(Error)
	errr := json.Unmarshal([]byte(err), e)
	if errr != nil {
		e.Detail = err
	}
	return e
}

func newError(code ErrorCode, detail string) *Error {
	e := &Error{
		Code:   code,
		Detail: detail,
		Meta:   map[string]string{},
	}
	e.setDefaultID()
	return e
}

// Wrap wraps the error with given code and detail.
func Wrap(err error, code ErrorCode, detail string) *Error {
	e := newError(code, detail)
	e.wrapped = err
	return e
}

// Wrapf wraps the error with given code and formatted detail message.
func Wrapf(err error, code ErrorCode, format string, args ...interface{}) *Error {
	e := newError(code, fmt.Sprintf(format, args...))
	e.wrapped = err
	return e
}

// ErrInvalidArgument generates a 400 error.
func ErrInvalidArgument(a ...interface{}) *Error {
	return newError(CodeInvalidArgument, fmt.Sprint(a...))
}

// ErrInvalidArgumentf generates formatted 400 error.
func ErrInvalidArgumentf(format string, a ...interface{}) *Error {
	return newError(CodeInvalidArgument, fmt.Sprintf(format, a...))
}

// ErrFailedPrecondition generates a 400 error.
func ErrFailedPrecondition(a ...interface{}) *Error {
	return newError(CodeFailedPrecondition, fmt.Sprint(a...))
}

// ErrFailedPreconditionf generates a 400 error.
func ErrFailedPreconditionf(format string, a ...interface{}) *Error {
	return newError(CodeFailedPrecondition, fmt.Sprintf(format, a...))
}

// ErrUnauthenticated generates a 401 error.
func ErrUnauthenticated(a ...interface{}) *Error {
	return newError(CodeUnauthenticated, fmt.Sprint(a...))
}

// ErrUnauthenticatedf generates 401 error with formatted message.
func ErrUnauthenticatedf(format string, a ...interface{}) *Error {
	return newError(CodeUnauthenticated, fmt.Sprintf(format, a...))
}

// ErrPermissionDenied generates a 403 error.
func ErrPermissionDenied(a ...interface{}) *Error {
	return newError(CodePermissionDenied, fmt.Sprint(a...))
}

// ErrPermissionDeniedf generates a 403 error.
func ErrPermissionDeniedf(format string, a ...interface{}) *Error {
	return newError(CodePermissionDenied, fmt.Sprintf(format, a...))
}

// ErrNotFound generates a 404 error.
func ErrNotFound(a ...interface{}) *Error {
	return newError(CodeNotFound, fmt.Sprint(a...))
}

// ErrNotFoundf generates formatted 404 error.
func ErrNotFoundf(format string, a ...interface{}) *Error {
	return newError(CodeNotFound, fmt.Sprintf(format, a...))
}

// ErrDeadlineExceeded generates a 408 error.
func ErrDeadlineExceeded(a ...interface{}) *Error {
	return newError(CodeDeadlineExceeded, fmt.Sprint(a...))
}

// ErrDeadlineExceededf generates formatted 408 error.
func ErrDeadlineExceededf(format string, a ...interface{}) *Error {
	return newError(CodeDeadlineExceeded, fmt.Sprintf(format, a...))
}

// ErrAlreadyExists generates a 409 error.
func ErrAlreadyExists(a ...interface{}) *Error {
	return newError(CodeAlreadyExists, fmt.Sprint(a...))
}

// ErrAlreadyExistsf generates formatted 409 error.
func ErrAlreadyExistsf(format string, a ...interface{}) *Error {
	return newError(CodeAlreadyExists, fmt.Sprintf(format, a...))
}

// ErrInternal generates a 500 error.
func ErrInternal(a ...interface{}) *Error {
	return newError(CodeInternal, fmt.Sprint(a...))
}

// ErrInternalf generates formatted 500 error.
func ErrInternalf(format string, a ...interface{}) *Error {
	return newError(CodeInternal, fmt.Sprintf(format, a...))
}

// ErrUnimplemented generates Unimplemented error.
func ErrUnimplemented(a ...interface{}) *Error {
	return newError(CodeUnimplemented, fmt.Sprint(a...))
}

// ErrUnimplementedf generates Unimplemented error with formatting.
func ErrUnimplementedf(format string, a ...interface{}) *Error {
	return newError(CodeUnimplemented, fmt.Sprintf(format, a...))
}

// ErrUnavailable generates Unavailable error.
func ErrUnavailable(a ...interface{}) *Error {
	return newError(CodeUnavailable, fmt.Sprint(a...))
}

// ErrUnavailablef generates Unavailable error with formatting.
func ErrUnavailablef(format string, a ...interface{}) *Error {
	return newError(CodeUnavailable, fmt.Sprintf(format, a...))
}

// ErrUnknown generates Unknown error.
func ErrUnknown(a ...interface{}) *Error {
	return newError(CodeUnknown, fmt.Sprint(a...))
}

// ErrUnknownf generates Unknown error with formatting.
func ErrUnknownf(format string, a ...interface{}) *Error {
	return newError(CodeUnknown, fmt.Sprintf(format, a...))
}

// ErrOutOfRange generates an error with the CodeOutOfRange.
func ErrOutOfRange(a ...interface{}) *Error {
	return newError(CodeOutOfRange, fmt.Sprint(a...))
}

// ErrOutOfRangef generates formatted error with the CodeOutOfRange.
func ErrOutOfRangef(format string, a ...interface{}) *Error {
	return newError(CodeOutOfRange, fmt.Sprintf(format, a...))
}

// Equal tries to compare errors
func Equal(err1 error, err2 error) bool {
	var verr1, verr2 *Error
	ok1 := errors.As(err1, &verr1)
	ok2 := errors.As(err2, &verr2)

	if ok1 != ok2 {
		return false
	}
	if !ok1 {
		return verr1 == verr2
	}

	if verr1.Code != verr2.Code {
		return false
	}

	if verr1.Process != verr2.Process {
		return false
	}

	return true
}

// FromError try to convert go error to *Error
func FromError(err error) *Error {
	var vErr *Error
	if errors.As(err, &vErr) {
		return vErr
	}

	// Try to decode it from the GRPC Status.
	s, ok := status.FromError(err)
	if ok {
		if e, ok := FromString(s.Message()); ok {
			return e
		}
		return newError(ErrorCode(s.Code()), s.Message())
	}

	// Check if the error is in raw string form.
	if e, ok := FromString(err.Error()); ok {
		return e
	}

	return newError(Code(err), err.Error())
}

func (e *Error) setDefaultID() {
	e.ID = g.NextId()
}
