package eventstate

import (
	"context"
)

// Handler is a structure that defines an event handler with its unique name and matched event type.
type Handler struct {
	// Name is the unique, human-readable name for the event handler.
	Name string
	// EventTypes defines the type of the event handled by given handler.
	EventTypes []string
}

// StateHandler is an interface that allows to handle event state.
type StateHandler interface {
	// StartHandling starts handling given event by the handler with a name = handlerName.
	StartHandling(ctx context.Context, eventID, handlerName string) error

	// FinishHandling finishes handling given event by the handlerName.
	FinishHandling(ctx context.Context, eventID, handlerName string) error

	// HandlingFailed finishes handling given event by the handlerName.
	HandlingFailed(ctx context.Context, eventID, handlerName string, handleErr error) error
}

// HandlerRegistry is an interface used to register different handler state.
type HandlerRegistry interface {
	// RegisterHandlers registers the information about event handler.
	// This function should be done during migration of the event handler.
	RegisterHandlers(ctx context.Context, eventHandlers ...Handler) error

	// ListHandlers list the handlers for the
	ListHandlers(ctx context.Context) ([]Handler, error)
}

// HandleSearcher is an interface used to find event handle state.
type HandleSearcher interface {
	// FindUnhandledEvents finds all unhandled events for given handler.
	FindUnhandledEvents(ctx context.Context, query FindUnhandledQuery) ([]Unhandled, error)

	// FindEventHandleFailures finds the handle failures for given handler name.
	FindEventHandleFailures(ctx context.Context, query FindFailureQuery) ([]HandleFailure, error)
}
