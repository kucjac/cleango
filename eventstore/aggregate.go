package eventstore

import (
	"time"

	"github.com/kucjac/cleango/messages/codec"
)

// Aggregate is an interface used for the aggregate models
type Aggregate interface {
	Apply(e *Event) error
	Validate() error
	SetBase(base *AggregateBase)
	Base() *AggregateBase
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
		Timestamp:     time.Now().UTC().Unix(),
		Revision:      revision + 1,
	}

	if err = a.agg.Apply(e); err != nil {
		return err
	}
	a.revision++
	a.uncommittedEvents = append(a.uncommittedEvents, e)

	return nil
}

// CommittedEvents gets the committed event messages.
func (a *AggregateBase) CommittedEvents() []*Event {
	return a.committedEvents
}

// DecodeEventAs decodes provided in put eventData into the structure of eventMsg.
// THe eventMsg is expected to be a pointer to the event msg.
func (a *AggregateBase) DecodeEventAs(eventData []byte, eventMsg interface{}) error {
	return a.eventCodec.Unmarshal(eventData, eventMsg)
}

func (a *AggregateBase) reset() {
	a.uncommittedEvents = nil
	a.revision = 0
	a.timestamp = 0
}
