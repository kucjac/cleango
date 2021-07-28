package esstate

import (
	"context"

	"github.com/kucjac/cleango/database/es"
	eventstate2 "github.com/kucjac/cleango/ddd/events/eventstate"
	"github.com/kucjac/cleango/xlog"
)

// EventStore is an interface that allows to operate on top of the standard es.EventStore, but also
// allows to control the state of the event handling.
type EventStore interface {
	es.EventStore

	// ChangeTypeOptions changes event type handling state options.
	ChangeTypeOptions(eventType string, options *Options)
}

// Compile time check if the Store implements EventStore interface.
var _ EventStore = (*Store)(nil)

// Store is an implementation of the event store that is also handling the event state on each commit.
type Store struct {
	*es.Store
	storage     Storage
	typeOptions map[string]*Options
}

// ChangeTypeOptions changes the default options for given event type.
func (s *Store) ChangeTypeOptions(eventType string, options *Options) {
	s.typeOptions[eventType] = options
}

// RegisterHandlers registers unique event handlers.
// This function should be used on the event migration only once per service.
func (s *Store) RegisterHandlers(ctx context.Context, eventHandlers ...eventstate2.Handler) error {
	return s.storage.RegisterHandlers(ctx, eventHandlers...)
}

// ListHandlers list the handlers that matches given event type.
func (s *Store) ListHandlers(ctx context.Context) ([]eventstate2.Handler, error) {
	return s.storage.ListHandlers(ctx)
}

// FindUnhandledEvents finds all unhandled events that matches given query.
func (s *Store) FindUnhandledEvents(ctx context.Context, query eventstate2.FindUnhandledQuery) ([]eventstate2.Unhandled, error) {
	return s.storage.FindUnhandled(ctx, query)
}

// FindEventHandleFailures find all handling failures that matches given query.
func (s *Store) FindEventHandleFailures(ctx context.Context, query eventstate2.FindFailureQuery) ([]eventstate2.HandleFailure, error) {
	return s.storage.FindFailures(ctx, query)
}

// Commit overwrites the default method of the es.Store and atomically commits given aggregate events, but also
// creates a new EventState per each committed event.
// This way no event is lost in handling, and the handlers now are able to control its status.
func (s *Store) Commit(ctx context.Context, aggregate es.Aggregate) error {
	// In case if the aggregate is the EventState use a regular commit flow.
	if aggregate.AggBase().Type() == AggregateType {
		return s.Store.Commit(ctx, aggregate)
	}
	tx, err := s.storage.BeginTx(ctx)
	if err != nil {
		return err
	}

	txStore := s.Store.WithStorage(tx)
	if err = txStore.Commit(ctx, aggregate); err != nil {
		if er := tx.Rollback(ctx); er != nil {
			xlog.WithContext(ctx).
				WithField("err", err).
				Errorf("Rolling back transaction failed")
		}
		return err
	}

	// Iterate over all committed events and create a new event state for each.
	events := aggregate.AggBase().CommittedEvents()
	for _, e := range events {
		// Check if there are some custom options for given event type.
		options := s.getEventOptions(e.EventType)

		// define new event state aggregate.
		state, err := InitializeUnhandledEventState(e.EventId, e.EventType, e.Time(), s.Store.AggregateBaseSetter, options)
		if err != nil {
			if er := tx.Rollback(ctx); er != nil {
				xlog.WithContext(ctx).
					WithField("err", err).
					Errorf("Rolling back transaction failed")
			}
			return err
		}

		// Commit it immediately.
		if err = txStore.Commit(ctx, state); err != nil {
			if er := tx.Rollback(ctx); er != nil {
				xlog.WithContext(ctx).
					WithField("err", err).
					Errorf("Rolling back transaction failed")
			}
			return err
		}
	}

	// Mark the events unhandled.
	for _, e := range events {
		if err = tx.MarkUnhandled(ctx, e.EventId, e.EventType, e.Timestamp); err != nil {
			return err
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

// StartHandling starts handling given event by the handler with a name = handlerName.
func (s *Store) StartHandling(ctx context.Context, eventID, handlerName string) error {
	state := NewEventState(eventID, s.AggregateBaseSetter)

	if err := s.LoadEvents(ctx, state); err != nil {
		return err
	}

	if err := state.StartHandling(handlerName); err != nil {
		return err
	}

	if err := s.Commit(ctx, state); err != nil {
		return err
	}

	if err := s.storage.StartHandling(ctx, eventID, handlerName, state.handlers[handlerName].lastStarted.UnixNano()); err != nil {
		return err
	}
	return nil
}

// FinishHandling finishes handling given event by the handlerName.
func (s *Store) FinishHandling(ctx context.Context, eventID, handlerName string) error {
	state := NewEventState(eventID, s.AggregateBaseSetter)

	if err := s.LoadEvents(ctx, state); err != nil {
		return err
	}

	if err := state.FinishHandling(handlerName); err != nil {
		return err
	}

	if err := s.Commit(ctx, state); err != nil {
		return err
	}

	if err := s.storage.FinishHandling(ctx, eventID, handlerName, state.handlers[handlerName].finishedAt.UnixNano()); err != nil {
		return err
	}
	return nil
}

// HandlingFailed finishes handling given event by the handlerName.
func (s *Store) HandlingFailed(ctx context.Context, eventID, handlerName string, handleErr error) error {
	// Load up the event state.
	state := NewEventState(eventID, s.AggregateBaseSetter)
	if err := s.LoadEvents(ctx, state); err != nil {
		return err
	}

	// Add the event that handling given event had failed.
	if err := state.HandlingFailed(handlerName, handleErr); err != nil {
		return err
	}

	// Commit given state.
	if err := s.Commit(ctx, state); err != nil {
		return err
	}

	// Create a failure for given event.
	failure := newEventHandleFailure(state, eventID, handleErr, handlerName)

	// Store the failure information in the storage.
	if err := s.storage.HandlingFailed(ctx, failure); err != nil {
		return err
	}
	return nil
}

// NewStore creates a new store that works in the same way as the es.Store with enhanced feature of tracking event state on each commit.
func NewStore(cfg *es.Config, eventCodec es.EventCodec, snapCodec es.SnapshotCodec, storage Storage) (*Store, error) {
	eventStore, err := es.New(cfg, eventCodec, snapCodec, storage)
	if err != nil {
		return nil, err
	}
	return &Store{Store: eventStore, typeOptions: map[string]*Options{}, storage: storage}, nil
}

func (s *Store) getEventOptions(eventType string) *Options {
	options := s.typeOptions[eventType]
	if options == nil {
		options = DefaultOptions()
	}
	return options
}
