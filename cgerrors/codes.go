package cgerrors

import (
	"gocloud.dev/gcerrors"
)

// ErrorCoder is an interface used to obtain an error code from given error within given implementation.
type ErrorCoder interface {
	ErrorCode(err error) ErrorCode
}

// Code gets the code defined in given
func Code(err error) ErrorCode {
	switch e := err.(type) {
	case *Error:
		return e.Code
	case ErrorCoder:
		return e.ErrorCode(err)
	}
	return fromGCErrors(gcerrors.Code(err))
}

func fromGCErrors(code gcerrors.ErrorCode) ErrorCode {
	switch code {
	case gcerrors.OK:
		return ErrorCode_OK
	case gcerrors.Unknown:
		return ErrorCode_Unknown
	case gcerrors.NotFound:
		return ErrorCode_NotFound
	case gcerrors.AlreadyExists:
		return ErrorCode_AlreadyExists
	case gcerrors.InvalidArgument:
		return ErrorCode_InvalidArgument
	case gcerrors.Internal:
		return ErrorCode_Internal
	case gcerrors.Unimplemented:
		return ErrorCode_Unimplemented
	case gcerrors.FailedPrecondition:
		return ErrorCode_FailedPrecondition
	case gcerrors.PermissionDenied:
		return ErrorCode_PermissionDenied
	case gcerrors.ResourceExhausted:
		return ErrorCode_ResourceExhausted
	case gcerrors.Canceled:
		return ErrorCode_Canceled
	case gcerrors.DeadlineExceeded:
		return ErrorCode_DeadlineExceeded
	}
	return ErrorCode_Unknown
}
