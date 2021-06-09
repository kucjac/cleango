package es

import (
	"context"
	"time"
)

//go:generate protoc -I=. --go_out=. event.proto --go_opt=paths=source_relative

// EventMessage is an interface that defines event messages.
type EventMessage interface {
	MessageType() string
}

// Copy creates a copy of given event.
func (x *Event) Copy() *Event {
	return &Event{
		EventId:       x.EventId,
		EventType:     x.EventType,
		AggregateType: x.AggregateType,
		AggregateId:   x.AggregateId,
		EventData:     x.EventData,
		Timestamp:     x.Timestamp,
		Revision:      x.Revision,
	}
}

func (x *Event) Time() time.Time {
	return time.Unix(0, x.Timestamp)
}

// EventHandler is an interface used for handling events.
type EventHandler interface {
	Handle(ctx context.Context, e *Event)
}
