// Code generated by "stringer -output codes_string.go -trimprefix ErrorCode_ -type ErrorCode"; DO NOT EDIT.

package cgerrors

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ErrorCode_OK-0]
	_ = x[ErrorCode_Canceled-1]
	_ = x[ErrorCode_Unknown-2]
	_ = x[ErrorCode_InvalidArgument-3]
	_ = x[ErrorCode_DeadlineExceeded-4]
	_ = x[ErrorCode_NotFound-5]
	_ = x[ErrorCode_AlreadyExists-6]
	_ = x[ErrorCode_PermissionDenied-7]
	_ = x[ErrorCode_ResourceExhausted-8]
	_ = x[ErrorCode_FailedPrecondition-9]
	_ = x[ErrorCode_Aborted-10]
	_ = x[ErrorCode_OutOfRange-11]
	_ = x[ErrorCode_Unimplemented-12]
	_ = x[ErrorCode_Internal-13]
	_ = x[ErrorCode_Unavailable-14]
	_ = x[ErrorCode_DataLoss-15]
	_ = x[ErrorCode_Unauthenticated-16]
}

const _ErrorCode_name = "OKCanceledUnknownInvalidArgumentDeadlineExceededNotFoundAlreadyExistsPermissionDeniedResourceExhaustedFailedPreconditionAbortedOutOfRangeUnimplementedInternalUnavailableDataLossUnauthenticated"

var _ErrorCode_index = [...]uint8{0, 2, 10, 17, 32, 48, 56, 69, 85, 102, 120, 127, 137, 150, 158, 169, 177, 192}

func (c ErrorCode) String() string {
	if c < 0 || c >= ErrorCode(len(_ErrorCode_index)-1) {
		return "ErrorCode(" + strconv.FormatInt(int64(c), 10) + ")"
	}
	return _ErrorCode_name[_ErrorCode_index[c]:_ErrorCode_index[c+1]]
}
