// Package cgerrors provides a way to return detailed information
// for an RPC request error. The error is normally JSON encoded.
package cgerrors

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kucjac/cleango/internal/uniqueid"
	"google.golang.org/grpc/status"
)

var g = uniqueid.NextGenerator("errors")

//go:generate protoc -I. --go_out=paths=source_relative:. errors.proto

// Error implements error interface.
func (x *Error) Error() string {
	b, _ := json.Marshal(x)
	return string(b)
}

// WithMeta sets the key, value metadata for given error.
func (x *Error) WithMeta(key, value string) *Error {
	if x.Meta == nil {
		x.Meta = make(map[string]string)
	}
	x.Meta[key] = value
	return x
}

// WithProcess sets the process for given error.
func (x *Error) WithProcess(process string) *Error {
	x.Process = process
	return x
}

// WithCode sets the code for given error.
func (x *Error) WithCode(code ErrorCode) *Error {
	x.Code = code
	return x
}

// GRPCError is an interface used to get grpcStatus
type GRPCError interface {
	GRPCStatus() *status.Status
}

// ToGRPCError converts an error to GRPC status.Status.
func ToGRPCError(err error) error {
	if e, ok := err.(GRPCError); ok {
		return e.GRPCStatus().Err()
	}
	return newError(Code(err), err.Error()).GRPCStatus().Err()
}

// GRPCStatus implements grpc client interface used to convert statuses.
func (x *Error) GRPCStatus() *status.Status {
	return status.New(x.Code.ToGRPCCode(), x.Error())
}

// Is compares the errors with their values.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// Is implements errors interface used by the
func (x *Error) Is(err error) bool {
	return Equal(x, err)
}

// To type check if given error is of *Error type or has encoded ErrorCode in it.
// Otherwise creates a new error with Unknown code.
func To(err error) *Error {
	if e, ok := err.(*Error); ok && e != nil {
		return e
	}
	if c := Code(err); c != ErrorCode_Unknown {
		return newError(c, err.Error())
	}

	e := new(Error)
	es := err.Error()
	errr := json.Unmarshal([]byte(es), e)
	if errr != nil {
		e.Detail = es
	} else {
		e.Code = ErrorCode_Unknown
	}
	return e
}

// IsNotFound checks if given input error is of code NotFound.
func IsNotFound(err error) bool {
	return Code(err) == ErrorCode_NotFound
}

// IsAlreadyExists checks if given error means that given entity already exists.
func IsAlreadyExists(err error) bool {
	return Code(err) == ErrorCode_AlreadyExists
}

// IsInvalidArgument checks if given error means that given entity already exists.
func IsInvalidArgument(err error) bool {
	return Code(err) == ErrorCode_InvalidArgument
}

// IsDeadlineExceeded checks if given error is of type Deadline Exceeded.
func IsDeadlineExceeded(err error) bool {
	return Code(err) == ErrorCode_Internal
}

// IsUnauthenticated checks if given error is an unauthenticated error.
func IsUnauthenticated(err error) bool {
	return Code(err) == ErrorCode_Unauthenticated
}

// IsInternal checks if the input is an internal error.
func IsInternal(err error) bool {
	return Code(err) == ErrorCode_Internal
}

// IsPermissionDenied checks if given error is of type PermissionDenied.
func IsPermissionDenied(err error) bool {
	return Code(err) == ErrorCode_PermissionDenied
}

// New generates a custom error.
func New(id, detail string, code ErrorCode) error {
	e := &Error{
		Id:     id,
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

// ErrInvalidArgument generates a 400 error.
func ErrInvalidArgument(a ...interface{}) *Error {
	return newError(ErrorCode_InvalidArgument, fmt.Sprint(a...))
}

// ErrInvalidArgumentf generates formatted 400 error.
func ErrInvalidArgumentf(format string, a ...interface{}) *Error {
	return newError(ErrorCode_InvalidArgument, fmt.Sprintf(format, a...))
}

// ErrUnauthenticated generates a 401 error.
func ErrUnauthenticated(a ...interface{}) *Error {
	return newError(ErrorCode_Unauthenticated, fmt.Sprint(a...))
}

// ErrUnauthenticatedf generates 401 error with formatted message.
func ErrUnauthenticatedf(format string, a ...interface{}) *Error {
	return newError(ErrorCode_Unauthenticated, fmt.Sprintf(format, a...))
}

// ErrPermissionDenied generates a 403 error.
func ErrPermissionDenied(a ...interface{}) *Error {
	return newError(ErrorCode_PermissionDenied, fmt.Sprint(a...))
}

// ErrPermissionDeniedf generates a 403 error.
func ErrPermissionDeniedf(format string, a ...interface{}) *Error {
	return newError(ErrorCode_PermissionDenied, fmt.Sprintf(format, a...))
}

// ErrNotFound generates a 404 error.
func ErrNotFound(a ...interface{}) *Error {
	return newError(ErrorCode_NotFound, fmt.Sprint(a...))
}

// ErrNotFoundf generates formatted 404 error.
func ErrNotFoundf(format string, a ...interface{}) *Error {
	return newError(ErrorCode_NotFound, fmt.Sprintf(format, a...))
}

// ErrDeadlineExceeded generates a 408 error.
func ErrDeadlineExceeded(a ...interface{}) *Error {
	return newError(ErrorCode_DeadlineExceeded, fmt.Sprint(a...))
}

// ErrDeadlineExceededf generates formatted 408 error.
func ErrDeadlineExceededf(format string, a ...interface{}) *Error {
	return newError(ErrorCode_DeadlineExceeded, fmt.Sprintf(format, a...))
}

// ErrAlreadyExists generates a 409 error.
func ErrAlreadyExists(a ...interface{}) *Error {
	return newError(ErrorCode_AlreadyExists, fmt.Sprint(a...))
}

// ErrAlreadyExistsf generates formatted 409 error.
func ErrAlreadyExistsf(format string, a ...interface{}) *Error {
	return newError(ErrorCode_AlreadyExists, fmt.Sprintf(format, a...))
}

// ErrInternal generates a 500 error.
func ErrInternal(a ...interface{}) *Error {
	return newError(ErrorCode_Internal, fmt.Sprint(a...))
}

// ErrInternalf generates formatted 500 error.
func ErrInternalf(format string, a ...interface{}) *Error {
	return newError(ErrorCode_Internal, fmt.Sprintf(format, a...))
}

// Equal tries to compare errors
func Equal(err1 error, err2 error) bool {
	verr1, ok1 := err1.(*Error)
	verr2, ok2 := err2.(*Error)

	if ok1 != ok2 {
		return false
	}

	if !ok1 {
		return err1 == err2
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
	if verr, ok := err.(*Error); ok && verr != nil {
		return verr
	}
	if c := Code(err); c != ErrorCode_Unknown {
		return newError(c, err.Error())
	}
	return Parse(err.Error())
}

func (x *Error) setDefaultID() {
	x.Id = g.NextId()
}
