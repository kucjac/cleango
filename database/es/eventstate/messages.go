package eventstate

import (
	"time"

	"github.com/kucjac/cleango/cgerrors"
)

//
// EventUnhandled Event
//

// newEventUnhandled is a constructor for the EventUnhandled event.
func newEventUnhandled(ts int64, eventType string, failures int, interval time.Duration, handlingTime time.Duration) (*EventUnhandled, error) {
	msg := &EventUnhandled{
		EventType:       eventType,
		Timestamp:       ts,
		MaxFailures:     int32(failures),
		MinFailInterval: int64(interval),
		MaxHandlingTime: int64(handlingTime),
	}
	if err := msg.Validate(); err != nil {
		return nil, err
	}
	return msg, nil
}

// EventUnhandledType is the type used by the Event aggregate on the Unhandled event.
const EventUnhandledType = "event_state:unhandled"

// MessageType gets the type of the event.
// Implements messages.Message interface.
func (x *EventUnhandled) MessageType() string {
	return EventUnhandledType
}

// EventUnhandledTopic is the topic used by the Event aggregate on the Unhandled event.
const EventUnhandledTopic = "eventsource.event_state.unhandled"

// MessageTopic returns messages.Topic from given message.
// Implements messages.Message interface.
func (x *EventUnhandled) MessageTopic() string {
	return EventUnhandledTopic
}

// Validate implements validator.Validator interface.
func (x *EventUnhandled) Validate() error {
	if x.Timestamp == 0 {
		return cgerrors.ErrInternal("no timestamp provided")
	}
	if x.EventType == "" {
		return cgerrors.ErrInternal("event type undefined")
	}
	return nil
}

//
// EventHandlingStarted Event
//

// newEventHandlingStarted is a constructor for the EventHandlingStarted event.
func newEventHandlingStarted(handlerName string) (*EventHandlingStarted, error) {
	msg := &EventHandlingStarted{HandlerName: handlerName}
	if err := msg.Validate(); err != nil {
		return nil, err
	}
	return msg, nil
}

// EventHandlingStartedType is the type used by the Event aggregate on the HandlingStarted event.
const EventHandlingStartedType = "event_state:handling_started"

// MessageType gets the type of the event.
// Implements messages.Message interface.
func (x *EventHandlingStarted) MessageType() string {
	return EventHandlingStartedType
}

// EventHandlingStartedTopic is the topic used by the Event aggregate on the HandlingStarted event.
const EventHandlingStartedTopic = "eventsource.event_state.handling_started"

// MessageTopic returns messages.Topic from given message.
// Implements messages.Message interface.
func (x *EventHandlingStarted) MessageTopic() string {
	return EventHandlingStartedTopic
}

// Validate implements validator.Validator interface.
func (x *EventHandlingStarted) Validate() error {
	if x.HandlerName == "" {
		return cgerrors.ErrInternal("undefined handler name")
	}
	return nil
}

//
// EventHandlingFinished Event
//

// newEventHandlingFinished is a constructor for the EventHandlingFinished event.
func newEventHandlingFinished(handlerName string) (*EventHandlingFinished, error) {
	msg := &EventHandlingFinished{HandlerName: handlerName}
	if err := msg.Validate(); err != nil {
		return nil, err
	}
	return msg, nil
}

// EventHandlingFinishedType is the type used by the Event aggregate on the HandlingFinished event.
const EventHandlingFinishedType = "event_state:handling_finished"

// MessageType gets the type of the event.
// Implements messages.Message interface.
func (x *EventHandlingFinished) MessageType() string {
	return EventHandlingFinishedType
}

// EventHandlingFinishedTopic is the topic used by the Event aggregate on the HandlingFinished event.
const EventHandlingFinishedTopic = "eventsource.event_state.handling_finished"

// MessageTopic returns messages.Topic from given message.
// Implements messages.Message interface.
func (x *EventHandlingFinished) MessageTopic() string {
	return EventHandlingFinishedTopic
}

// Validate implements validator.Validator interface.
func (x *EventHandlingFinished) Validate() error {
	if x.HandlerName == "" {
		return cgerrors.ErrInternal("undefined handler name")
	}
	return nil
}

//
// EventHandlingFailed Event
//

// newEventHandlingFailed is a constructor for the EventHandlingFailed event.
func newEventHandlingFailed(handlerName string, err error) (*EventHandlingFailed, error) {
	if err == nil {
		return nil, cgerrors.ErrInternal("no error message provided")
	}
	msg := &EventHandlingFailed{
		HandlerName: handlerName,
		Err:         err.Error(),
		ErrCode:     int32(cgerrors.Code(err)),
	}
	if err := msg.Validate(); err != nil {
		return nil, err
	}
	return msg, nil
}

// EventHandlingFailedType is the type used by the Event aggregate on the HandlingFailed event.
const EventHandlingFailedType = "event_state:handling_failed"

// MessageType gets the type of the event.
// Implements messages.Message interface.
func (x *EventHandlingFailed) MessageType() string {
	return EventHandlingFailedType
}

// EventHandlingFailedTopic is the topic used by the Event aggregate on the HandlingFailed event.
const EventHandlingFailedTopic = "eventsource.event_state.handling_failed"

// MessageTopic returns messages.Topic from given message.
// Implements messages.Message interface.
func (x *EventHandlingFailed) MessageTopic() string {
	return EventHandlingFailedTopic
}

// Validate implements validator.Validator interface.
func (x *EventHandlingFailed) Validate() error {
	if x.HandlerName == "" {
		return cgerrors.ErrInternal("handler name undefined")
	}
	if x.Err == "" {
		return cgerrors.ErrInternal("error message undefined")
	}
	if x.ErrCode == 0 {
		return cgerrors.ErrInternalf("undefined message error code")
	}
	return nil
}

//
// FailureCountReset Event
//

// newFailureCountReset is a constructor for the FailureCountReset event.
func newFailureCountReset(handlerName string) (*FailureCountReset, error) {
	msg := &FailureCountReset{HandlerName: handlerName}
	if err := msg.Validate(); err != nil {
		return nil, err
	}
	return msg, nil
}

// FailureCountResetType is the type used by the  aggregate on the FailureCountReset event.
const FailureCountResetType = "event_state:failure_count_reset"

// MessageType gets the type of the event.
// Implements messages.Message interface.
func (x *FailureCountReset) MessageType() string {
	return FailureCountResetType
}

// FailureCountResetTopic is the topic used by the  aggregate on the FailureCountReset event.
const FailureCountResetTopic = "eventsource.event_state.failure_count_reset"

// MessageTopic returns messages.Topic from given message.
// Implements messages.Message interface.
func (x *FailureCountReset) MessageTopic() string {
	return FailureCountResetTopic
}

// Validate implements validator.Validator interface.
func (x *FailureCountReset) Validate() error {
	if x.HandlerName == "" {
		return cgerrors.ErrInternal("no handler name defined")
	}
	return nil
}
