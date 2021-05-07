package cgerrors

import (
	"gocloud.dev/gcerrors"
	"google.golang.org/grpc/codes"
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

// ToGRPCCode gets the related grpc code.
func (x ErrorCode) ToGRPCCode() codes.Code {
	switch x {
	case ErrorCode_OK:
		return codes.OK
	case ErrorCode_Canceled:
		return codes.Canceled
	case ErrorCode_Unknown:
		return codes.Unknown
	case ErrorCode_InvalidArgument:
		return codes.InvalidArgument
	case ErrorCode_DeadlineExceeded:
		return codes.DeadlineExceeded
	case ErrorCode_NotFound:
		return codes.NotFound
	case ErrorCode_AlreadyExists:
		return codes.AlreadyExists
	case ErrorCode_PermissionDenied:
		return codes.PermissionDenied
	case ErrorCode_ResourceExhausted:
		return codes.ResourceExhausted
	case ErrorCode_FailedPrecondition:
		return codes.FailedPrecondition
	case ErrorCode_Aborted:
		return codes.Aborted
	case ErrorCode_OutOfRange:
		return codes.OutOfRange
	case ErrorCode_Unimplemented:
		return codes.Unimplemented
	case ErrorCode_Internal:
		return codes.Internal
	case ErrorCode_Unavailable:
		return codes.Unavailable
	case ErrorCode_DataLoss:
		return codes.DataLoss
	case ErrorCode_Unauthenticated:
		return codes.Unauthenticated
	default:
		return codes.Unknown
	}
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
