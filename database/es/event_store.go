package es

import (
	"context"
	"time"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/codec"
)

//go:generate mockgen -destination=mock/event_store_gen.go -package=mockes . EventStore

// EventStore is an interface used by the event store to load, commit and create snapshot on aggregates.
type EventStore interface {
	// LoadEvents loads all events for given aggregate.
	LoadEvents(ctx context.Context, aggregate Aggregate) error
	// LoadEventsWithSnapshot loads the latest snapshot with the events that happened after it.
	LoadEventsWithSnapshot(ctx context.Context, aggregate Aggregate) error
	// Commit commits the event changes done in given aggregate.
	Commit(ctx context.Context, aggregate Aggregate) error
	// SaveSnapshot saves the snapshot of given aggregate.
	SaveSnapshot(ctx context.Context, aggregate Aggregate) error
	// StreamEvents opens stream events that matches given request.
	StreamEvents(ctx context.Context, req *StreamEventsRequest) (<-chan *Event, error)
	// SetAggregateBase sets the AggregateBase within an aggregate.
	SetAggregateBase(agg Aggregate, aggId, aggType string, version int64)
}

// StreamEventsRequest is a request for the stream events query.
type StreamEventsRequest struct {
	// AggregateTypes streams the events for selected aggregate types.
	AggregateTypes []string
	// AggregateIDs is the filter that streams events for selected aggregate ids.
	AggregateIDs []string
	// ExcludeEventTypes is the filter that provides a stream with excluded event types.
	ExcludeEventTypes []string
	// EventTypes is the filter that gets only selected event types.
	EventTypes []string
	// BuffSize defines the size of the stream channel buffer.
	BuffSize int
}

// New creates new EventStore implementation.
func New(cfg *Config, eventCodec EventCodec, snapCodec SnapshotCodec, storage StorageBase) (*Store, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Store{
		AggregateBaseSetter: NewAggregateBaseSetter(eventCodec, snapCodec, UUIDGenerator{}),
		snapCodec:           snapCodec,
		storage:             storage,
		bufferSize:          cfg.BufferSize,
	}, nil
}

// Store is the default implementation for the EventStore interface.
type Store struct {
	*AggregateBaseSetter
	snapCodec  codec.Codec
	storage    StorageBase
	bufferSize int
}

// WithStorage creates a copy of the event store with given storage base.
func (e *Store) WithStorage(base StorageBase) *Store {
	return &Store{
		AggregateBaseSetter: e.AggregateBaseSetter,
		snapCodec:           e.snapCodec,
		storage:             base,
		bufferSize:          e.bufferSize,
	}
}

// LoadEvents gets the event stream and applies on provided aggregate.
func (e *Store) LoadEvents(ctx context.Context, agg Aggregate) error {
	b := agg.AggBase()
	// Get the full event stream for given aggregate.
	events, err := e.storage.ListEvents(ctx, b.id, b.aggType)
	if err != nil {
		return e.err("listing aggregate events failed", err)
	}

	// If no events are found return an error.
	if len(events) == 0 {
		return cgerrors.ErrNotFoundf("aggregate: %s with id: %s not found", b.aggType, b.id)
	}

	// Apply all events from the stream on the aggregate.
	for _, event := range events {
		if err = agg.Apply(event); err != nil {
			return err
		}
		b.revision = event.Revision
		b.timestamp = event.Timestamp
	}
	return nil
}

// LoadEventsWithSnapshot gets the aggregate stream with the latest possible snapshot.
func (e *Store) LoadEventsWithSnapshot(ctx context.Context, agg Aggregate) error {
	// At first try to get the snapshot for given aggregate.
	b := agg.AggBase()
	snap, err := e.storage.GetSnapshot(ctx, b.id, b.aggType, b.version)
	isNotFound := false
	if err != nil {
		isNotFound = e.storage.ErrorCode(err) == cgerrors.CodeNotFound
		if !isNotFound {
			return e.err("getting aggregate snapshot failed", err)
		}
	}

	var events []*Event
	if !isNotFound {
		if err = e.snapCodec.Unmarshal(snap.SnapshotData, agg); err != nil {
			return err
		}
		b.timestamp = snap.Timestamp
		b.revision = snap.Revision
		// Get the event stream starting form the revision provided in the snapshot.
		events, err = e.storage.ListEventsAfterRevision(ctx, b.id, b.aggType, b.revision)
		if err != nil {
			return e.err("listing events after revision failed", err)
		}
	} else {
		// If the snapshot is not found gets the full event stream.
		events, err = e.storage.ListEvents(ctx, b.id, b.aggType)
		if err != nil {
			return e.err("listing events failed", err)
		}

		// If there are no events, and the snapshot was not found - then there is no aggregate matching given id.
		if len(events) == 0 {
			return cgerrors.ErrNotFoundf("aggregate: %s with id %s not found", b.aggType, b.id)
		}
	}

	// Iterate over each event and apply them on given aggregate.
	for _, event := range events {
		if err = agg.Apply(event); err != nil {
			return err
		}
		b.revision = event.Revision
		b.timestamp = event.Timestamp
	}
	return nil
}

// SaveSnapshot stores the snapshot
func (e *Store) SaveSnapshot(ctx context.Context, agg Aggregate) error {
	// Create a snapshot and store it in the storage.
	data, err := e.snapCodec.Marshal(agg)
	if err != nil {
		return err
	}
	b := agg.AggBase()
	snap := &Snapshot{
		AggregateId:      b.id,
		AggregateType:    b.aggType,
		AggregateVersion: b.version,
		Revision:         b.revision,
		Timestamp:        b.timestamp,
		SnapshotData:     data,
	}
	if err = e.storage.SaveSnapshot(ctx, snap); err != nil {
		return e.err("saving snapshot failed", err)
	}
	return nil
}

// Commit commits all uncommitted events within given aggregate.
func (e *Store) Commit(ctx context.Context, agg Aggregate) error {
	b := agg.AggBase()
	events := b.uncommittedEvents
	if len(events) == 0 {
		return nil
	}
	for {
		// Try to save the events.
		err := e.storage.SaveEvents(ctx, events)
		if err == nil {
			b.committedEvents, b.uncommittedEvents = b.uncommittedEvents, nil
			return nil
		}

		// Everytime when the save fails due to the already exists error.
		// It means that there already is an event with provided revision.
		// In given case in order to apply the events
		if e.storage.ErrorCode(err) != cgerrors.CodeAlreadyExists {
			return e.err("saving events failed", err)
		}

		// Reset aggregate, and it's base.
		agg.Reset()
		b.reset()
		agg.SetBase(b)
		if err = e.LoadEventsWithSnapshot(ctx, agg); err != nil {
			return e.err("loading events with snapshot failed", err)
		}

		for _, event := range events {
			// Make a copy of given event.
			b.revision++
			event.Revision = b.revision
			event.Timestamp = time.Now().UTC().UnixNano()
			if err = agg.Apply(event); err != nil {
				return err
			}
		}
	}
}

// StreamEvents opens an event stream that matches given request.
func (e *Store) StreamEvents(ctx context.Context, req *StreamEventsRequest) (<-chan *Event, error) {
	c, err := e.storage.StreamEvents(ctx, req)
	if err != nil {
		return nil, e.err("opening storage stream events failed", err)
	}
	return c, nil
}

func (e *Store) err(msg string, err error) error {
	if cgerrors.Is(err, context.DeadlineExceeded) {
		return err
	}
	return cgerrors.Wrap(err, e.storage.ErrorCode(err), msg)
}
