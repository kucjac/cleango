package eventstate

import (
	"context"
	"time"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/database/es"
)

// StorageBase is an interface that contains base functionality of handling event state.
type StorageBase interface {
	es.StorageBase
	// MarkUnhandled marks the events as unhandled.
	MarkUnhandled(ctx context.Context, events ...*es.Event) error

	// StartHandling marks given event that it is being handled.
	StartHandling(ctx context.Context, e *es.Event, handlerName string, timestamp int64) error

	// FinishHandling marks an event that it is already handled.
	FinishHandling(ctx context.Context, e *es.Event, handlerName string, timestamp int64) error

	// HandlingFailed marks given handling as failure.
	HandlingFailed(ctx context.Context, failure *HandleFailure) error

	// RegisterHandlers registers the information about event handler.
	// This function should be done during migration of the event handler.
	RegisterHandlers(ctx context.Context, eventHandler ...Handler) error

	// ListHandlers list the handlers for the
	ListHandlers(ctx context.Context) ([]Handler, error)

	// FindUnhandled finds all unhandled events for given handler.
	FindUnhandled(ctx context.Context, query FindUnhandledQuery) ([]Unhandled, error)

	// FindFailures finds the handle failures for given handler name.
	FindFailures(ctx context.Context, query FindFailureQuery) ([]HandleFailure, error)
}

// Storage is an interface used for changing the state of given event with an ability of doing it in transaction.
type Storage interface {
	BeginTx(ctx context.Context) (TxStorage, error)
	StorageBase
}

// TxStorage is an interface used for changing the state of
type TxStorage interface {
	StorageBase
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// Handler is a structure that defines an event handler with its unique name and matched event type.
type Handler struct {
	// Name is the unique, human-readable name for the event handler.
	Name string
	// EventTypes defines the type of the event handled by given handler.
	EventTypes []string
}

// Unhandled is an event unhandled by the provided handler name.
type Unhandled struct {
	// Event contains the details of given event message.
	Event *es.Event
	// HandlerName defines the name of the handler where an event was not yet handled.
	HandlerName string
}

// HandleFailure contains the result of handling an event, with related error and number of retries.
type HandleFailure struct {
	// Event contains the details of given event message.
	Event *es.Event
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
