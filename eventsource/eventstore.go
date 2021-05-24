package eventsource

import (
	"context"
	"time"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/codec"
)

//go:generate mockgen -destination=mock.go -package=eventsource . Store

// EventStore is an interface used by the event store to load, commit and create snapshot on aggregates.
type EventStore interface {
	SetAggregateBase(agg Aggregate, aggId, aggType string, version int64)
	LoadEventStream(ctx context.Context, aggregate Aggregate) error
	LoadEventStreamWithSnapshot(ctx context.Context, aggregate Aggregate) error
	Commit(ctx context.Context, aggregate Aggregate) error
	SaveSnapshot(ctx context.Context, aggregate Aggregate) error
	StreamEvents(ctx context.Context, req *StreamEventsRequest) (<-chan *Event, error)
	StreamAggregates(ctx context.Context, aggType string, aggVersion int64, factory AggregateFactory) (<-chan Aggregate, error)
	StreamProjections(ctx context.Context, aggType string, aggVersion int64, factory ProjectionFactory) (<-chan Projection, error)
}

// StreamEventsRequest is a request for the stream events query.
type StreamEventsRequest struct {
	AggregateTypes    []string
	AggregateIDs      []string
	ExcludeEventTypes []string
	EventTypes        []string
	BuffSize          int
}

// New creates new EventStore implementation.
func New(cfg *Config, eventCodec EventCodec, snapCodec SnapshotCodec, storage Storage) (*Store, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Store{
		aggBaseSetter: newAggregateBaseSetter(eventCodec, snapCodec, UUIDGenerator{}),
		snapCodec:     snapCodec,
		storage:       storage,
		bufferSize:    cfg.BufferSize,
	}, nil
}

// Store is the default implementation for the EventStore interface.
type Store struct {
	*aggBaseSetter
	snapCodec  codec.Codec
	storage    Storage
	bufferSize int
}

// WithStorage creates a copy of the Store structure that has a different storage.
// This function could be used to create transaction implementations.
func (e *Store) WithStorage(storage Storage) *Store {
	return &Store{aggBaseSetter: e.aggBaseSetter, snapCodec: e.snapCodec, bufferSize: e.bufferSize, storage: storage}
}

// LoadEventStream gets the event stream and applies on provided aggregate.
func (e *Store) LoadEventStream(ctx context.Context, agg Aggregate) error {
	b := agg.AggBase()
	// Get the full event stream for given aggregate.
	events, err := e.storage.GetEventStream(ctx, b.id, b.aggType)
	if err != nil {
		return err
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

// LoadEventStreamWithSnapshot gets the aggregate stream with the latest possible snapshot.
func (e *Store) LoadEventStreamWithSnapshot(ctx context.Context, agg Aggregate) error {
	// At first try to get the snapshot for given aggregate.
	b := agg.AggBase()
	snap, err := e.storage.GetSnapshot(ctx, b.id, b.aggType, b.version)
	isNotFound := cgerrors.IsNotFound(err)
	if err != nil && !isNotFound {
		return err
	}

	var events []*Event
	if !isNotFound {
		if err = e.snapCodec.Unmarshal(snap.SnapshotData, agg); err != nil {
			return err
		}
		b.timestamp = snap.Timestamp
		b.revision = snap.Revision
		// Get the event stream starting form the revision provided in the snapshot.
		events, err = e.storage.GetStreamFromRevision(ctx, b.id, b.aggType, b.revision)
		if err != nil {
			return err
		}
	} else {
		// If the snapshot is not found gets the full event stream.
		events, err = e.storage.GetEventStream(ctx, b.id, b.aggType)
		if err != nil {
			return err
		}

		// If there is no events and the snapshot was not found - then there is no aggregate matching given id.
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
		return err
	}
	return err
}

// Commit commits provided aggregate events. In case when the
func (e *Store) Commit(ctx context.Context, agg Aggregate) error {
	b := agg.AggBase()
	events := b.uncommittedEvents
	for {
		err := e.storage.SaveEvents(ctx, events)
		if err == nil {
			return nil
		}
		if !cgerrors.IsAlreadyExists(err) {
			return err
		}

		// Everytime when the save fails due to the already exists error - it means that there already is an event with provided revision.
		// In given case in order to apply the events

		// Reset aggregate and it's base.
		agg.Reset()
		b.reset()
		if err = e.LoadEventStreamWithSnapshot(ctx, agg); err != nil {
			return err
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

// StreamAggregates opens up the aggregate streaming channel. The channel would got closed when there is no more aggregate to read
// or when the context is done.
// Closing resulting channel would result with a panic.
func (e *Store) StreamAggregates(ctx context.Context, aggType string, aggVersion int64, factory AggregateFactory) (<-chan Aggregate, error) {
	c, err := e.storage.NewCursor(ctx, aggType, aggVersion)
	if err != nil {
		return nil, err
	}

	l := newAggregateLoader(c, aggType, aggVersion, factory, e.bufferSize)
	return l.readAggregateChannel()
}

func (e *Store) StreamProjections(ctx context.Context, aggType string, aggVersion int64, factory ProjectionFactory) (<-chan Projection, error) {
	c, err := e.storage.NewCursor(ctx, aggType, aggVersion)
	if err != nil {
		return nil, err
	}

	l := newProjectionLoader(e.eventCodec, c, aggType, aggVersion, factory, e.bufferSize)
	return l.readProjectionChannel()
}

func (e *Store) StreamEvents(ctx context.Context, req *StreamEventsRequest) (<-chan *Event, error) {
	return e.storage.StreamEvents(ctx, req)
}