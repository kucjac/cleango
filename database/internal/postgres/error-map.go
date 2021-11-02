package postgres

import (
	"github.com/kucjac/cleango/cgerrors"
)

// ErrorMap is the postgres error mapping.
var ErrorMap = map[string]cgerrors.ErrorCode{
	// Class 02 - No data
	"02":    cgerrors.CodeNotFound,
	"P0002": cgerrors.CodeNotFound,
	// Class 08 - Connection Exception
	"08": cgerrors.CodeUnavailable,

	// Class 22 - Data exception
	"22000": cgerrors.CodeInvalidArgument,
	"23502": cgerrors.CodeInvalidArgument, // NOT-NULL Violation

	// Class 23 - Integrity violation
	"23000": cgerrors.ErrorCode_AlreadyExists,
	"23503": cgerrors.ErrorCode_NotFound, // Foreign Key Violation
	"23505": cgerrors.ErrorCode_AlreadyExists,

	// Class 42 - Table
	"42P06": cgerrors.ErrorCode_AlreadyExists,
	"42P07": cgerrors.ErrorCode_AlreadyExists,
	"42712": cgerrors.ErrorCode_AlreadyExists,
}
