package eventstate

import (
	"time"

	"github.com/kucjac/cleango/cgerrors"
)

// Unhandled is an event unhandled by the provided handler name.
type Unhandled struct {
	// EventUD is an identifier of the related event.
	EventID string
	// HandlerName defines the name of the handler where an event was not yet handled.
	HandlerName string
}

// HandleFailure contains the result of handling an event, with related error and number of retries.
type HandleFailure struct {
	// EventID is an identifier of the related event.
	EventID string
	// HandlerName defines the name of the handler for which the handling failed.
	HandlerName string
	// Err keeps the error information.
	Err string
	// ErrCode keeps the code of the error.
	ErrCode cgerrors.ErrorCode
	// RetryNo is the number of the retries to handle given event.
	RetryNo int
	// Timestamp of the failure.
	Timestamp time.Time
}

// FindUnhandledQuery is a query messa
type FindUnhandledQuery struct {
	// HandlerNames defines the filter for the handler names in a query for unhandled events.
	HandlerNames []string
}

// FindFailureQuery is a query messa
type FindFailureQuery struct {
	// HandlerNames defines the filter for the handler names in a query for unhandled events.
	HandlerNames []string
}
