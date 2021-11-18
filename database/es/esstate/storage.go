package esstate

import (
	"context"

	"github.com/kucjac/cleango/database/es"
	eventstate2 "github.com/kucjac/cleango/ddd/events/eventstate"
)

// StorageBase is an interface that contains base functionality of handling event state.
type StorageBase interface {
	es.StorageBase
	// MarkUnhandled marks the events as unhandled.
	MarkUnhandled(ctx context.Context, eventID, eventType string, timestamp int64) error
	// StartHandling marks given event that it is being handled.
	StartHandling(ctx context.Context, eventID string, handlerName string, timestamp int64) error
	// FinishHandling marks an event that it is already handled.
	FinishHandling(ctx context.Context, eventID string, handlerName string, timestamp int64) error
	// HandlingFailed marks given handling as failure.
	HandlingFailed(ctx context.Context, failure *eventstate2.HandleFailure) error
	// RegisterHandlers registers the information about event handler.
	// This function should be done during migration of the event handler.
	RegisterHandlers(ctx context.Context, eventHandler ...eventstate2.Handler) error
	// ListHandlers list the handlers for the
	ListHandlers(ctx context.Context) ([]eventstate2.Handler, error)
	// FindUnhandled finds all unhandled events for given handler.
	FindUnhandled(ctx context.Context, query eventstate2.FindUnhandledQuery) ([]eventstate2.Unhandled, error)
	// FindFailures finds the handle failures for given handler name.
	FindFailures(ctx context.Context, query eventstate2.FindFailureQuery) ([]eventstate2.HandleFailure, error)
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
