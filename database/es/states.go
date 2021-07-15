package es

import (
	"context"
)

// EventHandler is an interface used for handling events.
type EventHandler interface {
	Handle(ctx context.Context, e *Event)
}

type EventState int

const (
	// EventStateUndefined is an undefined state for event.
	EventStateUndefined EventState = 0
	// EventStateUnhandled is a state for unhandled event.
	EventStateUnhandled EventState = 1
	// EventStateHandlingStarted is a state that handling of an event was already started.
	EventStateHandlingStarted EventState = 2
	// EventStateHandlingDone is a state that handling an event is done.
	EventStateHandlingDone EventState = 3
)

type StorageStater interface {
	MarkUnhandled(ctx context.Context, events []*Event) error
	StartHandling(ctx context.Context, e *Event) error
	FinishHandling(ctx context.Context, e *Event) error
}


