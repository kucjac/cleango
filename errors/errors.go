// Package errors provides a way to return detailed information
// for an RPC request error. The error is normally JSON encoded.
package errors

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
)

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
func (x *Error) WithCode(code codes.Code) *Error {
	x.Code = uint32(code)
	return x
}

// Is compares the errors with their values.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// Is implements errors interface used by the
func (x *Error) Is(err error) bool {
	return Equal(x, err)
}

// IsNotFound checks if given input error is of code NotFound.
func IsNotFound(err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	return e.Code == uint32(codes.NotFound)
}

// IsAlreadyExists checks if given error means that given entity already exists.
func IsAlreadyExists(err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	return e.Code == uint32(codes.AlreadyExists)
}

// IsInvalidArgument checks if given error means that given entity already exists.
func IsInvalidArgument(err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	return e.Code == uint32(codes.InvalidArgument)
}

// IsDeadlineExceeded checks if given error is of type Deadline Exceeded.
func IsDeadlineExceeded(err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	return e.Code == uint32(codes.DeadlineExceeded)
}

// IsUnauthenticated checks if given error is an unauthenticated error.
func IsUnauthenticated(err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	return e.Code == uint32(codes.Unauthenticated)
}

// IsInternal checks if the input is an internal error.
func IsInternal(err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	return e.Code == uint32(codes.Internal)
}

// IsPermissionDenied checks if given error is of type PermissionDenied.
func IsPermissionDenied(err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	return e.Code == uint32(codes.PermissionDenied)
}

// New generates a custom error.
func New(id, detail string, code codes.Code) error {
	e := &Error{
		Id:     id,
		Code:   uint32(code),
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

func newError(code codes.Code, detail string) *Error {
	e := &Error{
		Code:   uint32(code),
		Detail: detail,
		Meta:   map[string]string{},
	}
	e.setDefaultID()
	return e
}

// ErrInvalidArgument generates a 400 error.
func ErrInvalidArgument(a ...interface{}) *Error {
	return newError(codes.InvalidArgument, fmt.Sprint(a...))
}

// ErrInvalidArgumentf generates formatted 400 error.
func ErrInvalidArgumentf(format string, a ...interface{}) *Error {
	return newError(codes.InvalidArgument, fmt.Sprintf(format, a...))
}

// ErrUnauthenticated generates a 401 error.
func ErrUnauthenticated(a ...interface{}) *Error {
	return newError(codes.Unauthenticated, fmt.Sprint(a...))
}

// ErrUnauthenticatedf generates 401 error with formatted message.
func ErrUnauthenticatedf(format string, a ...interface{}) *Error {
	return newError(codes.Unauthenticated, fmt.Sprintf(format, a...))
}

// ErrPermissionDenied generates a 403 error.
func ErrPermissionDenied(a ...interface{}) *Error {
	return newError(codes.PermissionDenied, fmt.Sprint(a...))
}

// ErrPermissionDeniedf generates a 403 error.
func ErrPermissionDeniedf(format string, a ...interface{}) *Error {
	return newError(codes.PermissionDenied, fmt.Sprintf(format, a...))
}

// ErrNotFound generates a 404 error.
func ErrNotFound(a ...interface{}) *Error {
	return newError(codes.NotFound, fmt.Sprint(a...))
}

// ErrNotFoundf generates formatted 404 error.
func ErrNotFoundf(format string, a ...interface{}) *Error {
	return newError(codes.NotFound, fmt.Sprintf(format, a...))
}

// ErrDeadlineExceeded generates a 408 error.
func ErrDeadlineExceeded(a ...interface{}) *Error {
	return newError(codes.DeadlineExceeded, fmt.Sprint(a...))
}

// ErrDeadlineExceededf generates formatted 408 error.
func ErrDeadlineExceededf(format string, a ...interface{}) *Error {
	return newError(codes.DeadlineExceeded, fmt.Sprintf(format, a...))
}

// ErrAlreadyExists generates a 409 error.
func ErrAlreadyExists(a ...interface{}) *Error {
	return newError(codes.AlreadyExists, fmt.Sprint(a...))
}

// ErrAlreadyExistsf generates formatted 409 error.
func ErrAlreadyExistsf(format string, a ...interface{}) *Error {
	return newError(codes.AlreadyExists, fmt.Sprintf(format, a...))
}

// ErrInternal generates a 500 error.
func ErrInternal(a ...interface{}) *Error {
	return newError(codes.Internal, fmt.Sprint(a...))
}

// ErrInternalf generates formatted 500 error.
func ErrInternalf(format string, a ...interface{}) *Error {
	return newError(codes.Internal, fmt.Sprintf(format, a...))
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

	return Parse(err.Error())
}

func (x *Error) setDefaultID() {
	var id string
	temp := make([]byte, 16)
	_, err := rand.Read(temp)
	if err == nil {
		id = hex.EncodeToString(temp)
	}
	x.Id = id
}
