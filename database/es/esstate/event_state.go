package esstate

import (
	"time"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/database/es"
	eventstate2 "github.com/kucjac/cleango/ddd/events/eventstate"
)

// State is a state value for the event handling.
type State int

const (
	// StateUndefined is an undefined event handle state.
	StateUndefined State = 0
	// StateUnhandled is a state for unhandled event.
	StateUnhandled State = 1
	// StateStarted is a state that handling of an event was already started.
	StateStarted State = 2
	// StateFinished is a state that handling an event is done.
	StateFinished State = 3
	// StateFailed is a state of an event that states handling an event had failed.
	StateFailed State = 4
)

// newEventHandleFailure creates a new event handling failure.
func newEventHandleFailure(state *EventState, eventID string, err error, handlerName string) *eventstate2.HandleFailure {
	return &eventstate2.HandleFailure{
		EventID:     eventID,
		Err:         err.Error(),
		ErrCode:     cgerrors.Code(err),
		HandlerName: handlerName,
		RetryNo:     state.getFailureRetryNo(handlerName),
		Timestamp:   state.getFailureTime(handlerName),
	}
}

// AggregateType is the type of the EventState aggregate.
const AggregateType = "eventsource.event_state"

// Compile time check if the EventState implements es.Aggregate.
var _ es.Aggregate = (*EventState)(nil)

type Options struct {
	// MaxFailures is the maximum number of failures for which the event would not allow to, start until it is reset.
	MaxFailures int
	// MinFailInterval is a minimum duration interval to wait after failure.
	// Next failure duration would be increased exponentially.
	MinFailInterval time.Duration
	// MaxHandlingTime is a handling time after which the message is treated as lost.
	MaxHandlingTime time.Duration
}

func (o *Options) Validate() error {
	if o.MaxFailures < 0 {
		return cgerrors.ErrInternal("invalid maximum failure number")
	}

	if o.MinFailInterval < 0 {
		return cgerrors.ErrInternal("invalid minimum failure interval")
	}

	if o.MaxHandlingTime < 0 {
		return cgerrors.ErrInternal("invalid maximum handling time interval")
	}

	return nil
}

// DefaultOptions creates default event state options.
func DefaultOptions() *Options {
	return &Options{
		MaxFailures:     5,
		MinFailInterval: 2 * time.Second,
		MaxHandlingTime: 10 * time.Second,
	}
}

// InitializeUnhandledEventState creates and initializes a new EventState model with an EventUnhandled message.
func InitializeUnhandledEventState(eventID, eventType string, timestamp time.Time, bs *es.AggregateBaseSetter, o *Options) (*EventState, error) {
	eh := &EventState{handlers: map[string]handles{}}
	if o == nil {
		o = DefaultOptions()
	}

	bs.SetAggregateBase(eh, eventID, AggregateType, 1)

	msg, err := newEventUnhandled(timestamp.UnixNano(), eventType, o.MaxFailures, o.MinFailInterval, o.MaxHandlingTime)
	if err != nil {
		return nil, err
	}
	if err = eh.base.SetEvent(msg); err != nil {
		return nil, err
	}
	return eh, nil
}

// NewEventState creates a new EventState aggregate model.
func NewEventState(eventID string, bs *es.AggregateBaseSetter) *EventState {
	eh := &EventState{handlers: map[string]handles{}}
	bs.SetAggregateBase(eh, eventID, AggregateType, 1)
	return eh
}

// EventState is an aggregate that stores the event handler changes.
type EventState struct {
	base                *es.AggregateBase
	timestamp           int64
	eventType           string
	handlers            map[string]handles
	minFailInterval     time.Duration
	maxFailures         int
	maxHandlingInterval time.Duration
}

// Apply implements es.Aggregate interface.
func (s *EventState) Apply(e *es.Event) (err error) {
	if s.base == nil {
		return cgerrors.ErrInternal("event handle aggregate base undefined")
	}
	if e.AggregateType != AggregateType {
		return cgerrors.ErrInternal("input event has different aggregate type").
			WithMeta("aggregate_type", e.AggregateType)
	}
	switch e.EventType {
	case EventUnhandledType:
		err = s.applyUnhandled(e)
	case EventHandlingStartedType:
		err = s.applyHandlingStarted(e)
	case EventHandlingFinishedType:
		err = s.applyHandlingFinished(e)
	case EventHandlingFailedType:
		err = s.applyHandlingFailed(e)
	case FailureCountResetType:
		err = s.applyFailureCountReset(e)
	default:
		return cgerrors.ErrInternal("undefined event type").WithMeta("event_type", e.EventType)
	}
	return err
}

// SetBase implements es.Aggregate interface.
func (s *EventState) SetBase(base *es.AggregateBase) {
	s.base = base
}

// AggBase implements es.Aggregate interface.
func (s *EventState) AggBase() *es.AggregateBase {
	return s.base
}

// Reset implements es.Aggregate interface.
func (s *EventState) Reset() {
	*s = EventState{base: s.base}
}

// StartHandling starts handling given event by the handlerName.
func (s *EventState) StartHandling(handlerName string) error {
	msg, err := newEventHandlingStarted(handlerName)
	if err != nil {
		return err
	}
	if err = s.base.SetEvent(msg); err != nil {
		return err
	}
	return nil
}

// FinishHandling finishes handling given event state by the handlerName.
func (s *EventState) FinishHandling(handlerName string) error {
	msg, err := newEventHandlingFinished(handlerName)
	if err != nil {
		return err
	}
	if err = s.base.SetEvent(msg); err != nil {
		return err
	}
	return nil
}

// HandlingFailed marks the event state that it's handling had failed with given error.
func (s *EventState) HandlingFailed(handlerName string, handlingErr error) error {
	msg, err := newEventHandlingFailed(handlerName, handlingErr)
	if err != nil {
		return err
	}
	if err = s.base.SetEvent(msg); err != nil {
		return err
	}
	return nil
}

// ResetFailures resets handling state failures.
func (s *EventState) ResetFailures(handlerName string) error {
	msg, err := newFailureCountReset(handlerName)
	if err != nil {
		return err
	}
	if err = s.base.SetEvent(msg); err != nil {
		return err
	}
	return nil
}

func (s *EventState) applyUnhandled(e *es.Event) error {
	var msg EventUnhandled
	if err := s.base.DecodeEventAs(e.EventData, &msg); err != nil {
		return err
	}
	s.timestamp = msg.Timestamp
	s.eventType = msg.EventType
	s.minFailInterval = time.Duration(msg.MinFailInterval)
	s.maxFailures = int(msg.MaxFailures)
	s.maxHandlingInterval = time.Duration(msg.MaxHandlingTime)
	return nil
}

func (s *EventState) applyHandlingStarted(e *es.Event) error {
	var msg EventHandlingStarted
	if err := s.base.DecodeEventAs(e.EventData, &msg); err != nil {
		return err
	}
	h := s.handlers[msg.HandlerName]
	if h.latestState == StateStarted &&
		time.Now().UTC().Before(h.lastStarted.Add(s.maxHandlingInterval)) {
		return cgerrors.ErrFailedPrecondition("given event handling had already started")
	}

	if h.latestState == StateFinished {
		return cgerrors.ErrAlreadyExists("event already handled")
	}

	if h.totalFailures > s.maxFailures {
		return cgerrors.New("", "too many handle tries for given handler", cgerrors.CodeResourceExhausted)
	}

	if h.latestState == StateFailed {
		// Check if there had passed at least minimum required time before last failure.
		failInterval := s.minFailInterval * time.Duration(1<<uint(h.totalFailures-1))
		minTime := h.lastFailure.Add(failInterval)
		if time.Now().UTC().Before(minTime) {
			return cgerrors.ErrFailedPrecondition("too many tries within time duration")
		}
	}

	h.handles = append(h.handles, handle{state: StateStarted, timestamp: e.Timestamp})
	h.latestState = StateStarted
	h.lastStarted = e.Time()
	s.handlers[msg.HandlerName] = h

	return nil
}

func (s *EventState) applyHandlingFinished(e *es.Event) error {
	var msg EventHandlingFinished
	if err := s.base.DecodeEventAs(e.EventData, &msg); err != nil {
		return err
	}
	h := s.handlers[msg.HandlerName]
	h.latestState = StateFinished
	h.finishedAt = e.Time()
	h.handles = append(h.handles, handle{state: StateFinished, timestamp: e.Timestamp})
	s.handlers[msg.HandlerName] = h

	return nil
}

func (s *EventState) applyHandlingFailed(e *es.Event) error {
	var msg EventHandlingFailed
	if err := s.base.DecodeEventAs(e.EventData, &msg); err != nil {
		return err
	}
	h := s.handlers[msg.HandlerName]
	h.latestState = StateFailed
	h.totalFailures++
	h.lastFailure = e.Time()
	h.handles = append(h.handles, handle{state: StateFailed, timestamp: e.Timestamp, failure: &handleFailure{
		err:         msg.Err,
		code:        cgerrors.ErrorCode(msg.ErrCode),
		retryNumber: h.totalFailures,
	}})
	s.handlers[msg.HandlerName] = h
	return nil
}

func (s *EventState) applyFailureCountReset(e *es.Event) error {
	var msg FailureCountReset
	if err := s.base.DecodeEventAs(e.EventData, &msg); err != nil {
		return err
	}
	h := s.handlers[msg.HandlerName]
	h.totalFailures = 0
	h.latestState = StateUnhandled
	h.lastFailure = time.Time{}
	s.handlers[msg.HandlerName] = h
	return nil
}

func (s *EventState) getFailureRetryNo(handlerName string) int {
	h := s.handlers[handlerName]
	return h.totalFailures
}

func (s *EventState) getFailureTime(handlerName string) time.Time {
	h := s.handlers[handlerName]
	return h.lastFailure
}

type handles struct {
	latestState   State
	handles       []handle
	totalFailures int
	lastFailure   time.Time
	lastStarted   time.Time
	finishedAt    time.Time
}

type handle struct {
	state     State
	timestamp int64
	failure   *handleFailure
}

type handleFailure struct {
	err         string
	code        cgerrors.ErrorCode
	retryNumber int
}
