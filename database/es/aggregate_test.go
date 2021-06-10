package es_test

import (
	"time"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/database/es"
)

const aggregateType = "aggregate"

var _ es.Aggregate = (*testAggregate)(nil)

type testAggregate struct {
	Base      *es.AggregateBase `json:"-"`
	Name      string            `json:"name"`
	CreatedAt time.Time         `json:"createdAt"`
}

func (t *testAggregate) Apply(e *es.Event) error {
	if e.AggregateType != aggregateType {
		return cgerrors.ErrInternal("event doesnt match given aggregate type")
	}
	switch e.EventType {
	case aggregateCreatedType:
		var msg aggregateCreated
		if err := t.Base.DecodeEventAs(e.EventData, &msg); err != nil {
			return cgerrors.ErrInternalf("decoding aggregate created event failed: %v", err)
		}
		if !t.CreatedAt.IsZero() {
			return cgerrors.ErrInvalidArgument("aggregate already created")
		}
		t.CreatedAt = e.Time()
	case aggregateNameChangedType:
		var msg aggregateNameChanged
		if err := t.Base.DecodeEventAs(e.EventData, &msg); err != nil {
			return cgerrors.ErrInternalf("decoding aggregate name changed event failed: %v", err)
		}
		t.Name = msg.Name
	default:
		return cgerrors.ErrInternalf("unsupported event type: %v", e.EventType)
	}
	return nil
}

func (t *testAggregate) SetBase(base *es.AggregateBase) {
	t.Base = base
}

func (t *testAggregate) AggBase() *es.AggregateBase {
	return t.Base
}

func (t *testAggregate) Reset() {
	*t = testAggregate{}
}

//
// aggregateCreated Event
//

// newAggregateCreated is a constructor for the aggregateCreated event.
func newAggregateCreated() (*aggregateCreated, error) {
	msg := &aggregateCreated{}
	if err := msg.Validate(); err != nil {
		return nil, err
	}
	return msg, nil
}

type aggregateCreated struct{}

// aggregateCreatedType is the type used by the aggregate aggregate on the Created event.
const aggregateCreatedType = "aggregate:created"

// MessageType gets the type of the event.
// Implements messages.Message interface.
func (x *aggregateCreated) MessageType() string {
	return aggregateCreatedType
}

// aggregateCreatedTopic is the topic used by the aggregate aggregate on the Created event.
const aggregateCreatedTopic = "cleango.aggregate.created"

// MessageTopic returns messages.Topic from given message.
// Implements messages.Message interface.
func (x *aggregateCreated) MessageTopic() string {
	return aggregateCreatedTopic
}

func (x *aggregateCreated) Validate() error {
	return nil
}

type aggregateNameChanged struct {
	Name string
}

//
// aggregateNameChanged Event
//

// newAggregateNameChanged is a constructor for the aggregateNameChanged event.
func newAggregateNameChanged(newName string) (*aggregateNameChanged, error) {
	msg := &aggregateNameChanged{Name: newName}
	if err := msg.Validate(); err != nil {
		return nil, err
	}
	return msg, nil
}

// aggregateNameChangedType is the type used by the aggregate aggregate on the NameChanged event.
const aggregateNameChangedType = "aggregate:name_changed"

// MessageType gets the type of the event.
// Implements messages.Message interface.
func (x *aggregateNameChanged) MessageType() string {
	return aggregateNameChangedType
}

// aggregateNameChangedTopic is the topic used by the aggregate aggregate on the NameChanged event.
const aggregateNameChangedTopic = "cleango.aggregate.name_changed"

// MessageTopic returns messages.Topic from given message.
// Implements messages.Message interface.
func (x *aggregateNameChanged) MessageTopic() string {
	return aggregateNameChangedTopic
}

func (x *aggregateNameChanged) Validate() error {
	if x.Name == "" {
		return cgerrors.ErrInternalf("%s - name not set", x.MessageType())
	}
	return nil
}
