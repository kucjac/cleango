package es

import (
	"time"
)

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
