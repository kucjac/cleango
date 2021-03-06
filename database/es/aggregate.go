package es

import (
	"time"

	"github.com/kucjac/cleango/codec"
	"github.com/kucjac/cleango/pkg/xlog"
)

// Aggregate is an interface used for the aggregate models
type Aggregate interface {
	Apply(e *Event) error
	SetBase(base *AggregateBase)
	AggBase() *AggregateBase
	Reset()
}

// AggregateBase is the base of the aggregate which
type AggregateBase struct {
	id                string
	aggType           string
	agg               Aggregate
	eventCodec        codec.Codec
	snapCodec         codec.Codec
	idGen             IdGenerator
	timestamp         int64
	uncommittedEvents []*Event
	committedEvents   []*Event
	revision          int64
	version           int64
}

// SetID sets aggregate id.
func (a *AggregateBase) SetID(id string) {
	a.id = id
}

// ID gets the aggregate identifier.
func (a *AggregateBase) ID() string {
	return a.id
}

// Type gets the aggregate type.
func (a *AggregateBase) Type() string {
	return a.aggType
}

// Revision gets aggregate current revision.
func (a *AggregateBase) Revision() int64 {
	return a.revision
}

// Timestamp gets the aggregate base timestamp.
func (a *AggregateBase) Timestamp() int64 {
	return a.timestamp
}

// Version gets aggregate version.
func (a *AggregateBase) Version() int64 {
	return a.version
}

// SetEvent sets new event message into given aggregate.
func (a *AggregateBase) SetEvent(eventMsg EventMessage) error {
	eventData, err := a.eventCodec.Marshal(eventMsg)
	if err != nil {
		return err
	}

	revision := a.revision
	e := &Event{
		EventId:       a.idGen.GenerateId(),
		EventType:     eventMsg.MessageType(),
		AggregateType: a.aggType,
		AggregateId:   a.id,
		EventData:     eventData,
		Timestamp:     time.Now().UTC().UnixNano(),
		Revision:      revision + 1,
	}

	if err = a.agg.Apply(e); err != nil {
		return err
	}
	a.revision++
	a.timestamp = e.Timestamp
	a.uncommittedEvents = append(a.uncommittedEvents, e)

	return nil
}

// NewEvent creates a new event with given identifier, at given timestamp and with given message.
// An event is added to uncommitted events of the aggregate base.
// NOTE: Created event is not applied to given aggregate base.
func (a *AggregateBase) NewEvent(id string, msg EventMessage, timestamp time.Time) (*Event, error) {
	eventData, err := a.eventCodec.Marshal(msg)
	if err != nil {
		return nil, err
	}
	if id == "" {
		id = a.idGen.GenerateId()
	}
	if timestamp.IsZero() {
		timestamp = time.Now().UTC()
	}
	revision := a.revision
	e := &Event{
		EventId:       id,
		EventType:     msg.MessageType(),
		AggregateType: a.aggType,
		AggregateId:   a.id,
		EventData:     eventData,
		Timestamp:     timestamp.UnixNano(),
		Revision:      revision + 1,
	}

	a.revision++
	a.timestamp = e.Timestamp
	a.uncommittedEvents = append(a.uncommittedEvents, e)
	return e, nil
}

// CommittedEvents gets the committed event messages.
func (a *AggregateBase) CommittedEvents() []*Event {
	return a.committedEvents
}

// UncommittedEvents gets the slice of uncommitted event messages.
func (a *AggregateBase) UncommittedEvents() []*Event {
	return a.uncommittedEvents
}

// MustLatestCommittedEvent gets the latest committed event message or panics.
func (a *AggregateBase) MustLatestCommittedEvent() *Event {
	if len(a.committedEvents) == 0 {
		xlog.Panicf("no committed events found for the aggregate: %s - id: %s", a.aggType, a.id)
	}
	return a.committedEvents[len(a.committedEvents)-1]
}

// LatestCommittedEvent gets the latest committed event message.
func (a *AggregateBase) LatestCommittedEvent() (*Event, bool) {
	if len(a.committedEvents) > 0 {
		return a.committedEvents[len(a.committedEvents)-1], true
	}
	return nil, false
}

// MarkEventsCommitted marks the aggregate events as committed.
// NOTE: Use this function carefully, as the event store wouldn't try to commit events, already marked as committed.
func (a *AggregateBase) MarkEventsCommitted() {
	a.committedEvents, a.uncommittedEvents = a.uncommittedEvents, nil
}

// DecodeEventAs decodes provided in put eventData into the structure of eventMsg.
// THe eventMsg is expected to be a pointer to the event msg.
func (a *AggregateBase) DecodeEventAs(eventData []byte, eventMsg interface{}) error {
	return a.eventCodec.Unmarshal(eventData, eventMsg)
}

func (a *AggregateBase) reset() {
	a.revision = 0
	a.timestamp = 0
}
