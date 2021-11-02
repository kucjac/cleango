package xpg

import (
	"github.com/kucjac/cleango/cgerrors"
)

var errorMap = map[string]cgerrors.ErrorCode{
	// Class 02 - No data
	"02":    cgerrors.CodeNotFound,
	"P0002": cgerrors.CodeNotFound,
	// Class 08 - Connection Exception
	"08": cgerrors.CodeUnavailable,

	// Class 22 - Data exception
	"22000": cgerrors.CodeInvalidArgument,
	"23502": cgerrors.CodeInvalidArgument, // NOT-NULL Violation

	// Class 23 - Integrity violation
	"23000": cgerrors.CodeAlreadyExists,
	"23503": cgerrors.CodeNotFound, // Foreign Key Violation
	"23505": cgerrors.CodeAlreadyExists,
}
