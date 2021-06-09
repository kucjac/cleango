package xpg

import (
	"github.com/kucjac/cleango/cgerrors"
)

var errorMap = map[string]cgerrors.ErrorCode{
	// Class 02 - No data
	"02":    cgerrors.ErrorCode_NotFound,
	"P0002": cgerrors.ErrorCode_NotFound,
	// Class 08 - Connection Exception
	"08": cgerrors.ErrorCode_Unavailable,

	// Class 22 - Data exception
	"22000": cgerrors.ErrorCode_InvalidArgument,
	"23502": cgerrors.ErrorCode_InvalidArgument, // NOT-NULL Violation

	// Class 23 - Integrity violation
	"23000": cgerrors.ErrorCode_AlreadyExists,
	"23503": cgerrors.ErrorCode_NotFound, // Foreign Key Violation
	"23505": cgerrors.ErrorCode_AlreadyExists,
}
