package es

import (
	"context"
	"fmt"
	"time"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/codec"
	"github.com/kucjac/cleango/xlog"
	"github.com/kucjac/cleango/xpubsub"
)

// Generate the mock store.
//go:generate mockgen -destination=mock.go -package=eventsource . EventStore

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

	// StreamAggregates opens aggregate stream for given type and version.
	StreamAggregates(ctx context.Context, aggType string, aggVersion int64, factory AggregateFactory) (<-chan Aggregate, error)

	// StreamProjections streams the projections based on given aggregate type and version.
	StreamProjections(ctx context.Context, aggType string, aggVersion int64, factory ProjectionFactory) (<-chan Projection, error)

	// SetAggregateBase sets the AggregateBase within an aggregate.
	SetAggregateBase(agg Aggregate, aggId, aggType string, version int64)
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
	storage    Storage
	bufferSize int
	topics     map[string]xpubsub.Topic
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
		isNotFound = e.storage.ErrorCode(err) == cgerrors.ErrorCode_NotFound
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
	return e.doAndCommit(ctx, agg, func([]*Event) error { return nil })
}

// DoAndCommit atomically saves all uncommitted events within given aggregate, executes a 'do' function
// and on success commits the changes.
func (e *Store) DoAndCommit(ctx context.Context, agg Aggregate, do func(events []*Event) error) error {
	return e.doAndCommit(ctx, agg, do)
}

func (e *Store) doAndCommit(ctx context.Context, agg Aggregate, do func(events []*Event) error) error {
	b := agg.AggBase()
	events := b.uncommittedEvents
	if len(events) == 0 {
		return nil
	}
	for {
		// Begin a new transaction.
		tx, err := e.storage.BeginTx(ctx)
		if err != nil {
			return err
		}

		// Try to save the events.
		if err = tx.SaveEvents(ctx, events); err == nil {
			// If everything went well execute a 'do' function.
			if err = do(events); err != nil {
				// Rollback the transaction on error.
				if er := tx.Rollback(ctx); er != nil {
					xlog.WithContext(ctx).
						WithField("err", er).
						Error("rolling back transaction failed")
				}
				return err
			}

			// Commit the changes atomically with 'do' function.
			if er := tx.Commit(ctx); er != nil {
				xlog.WithContext(ctx).
					WithField("err", err).
					Error("committing transaction failed")
			}

			// Mark the events committed.
			b.committedEvents, b.uncommittedEvents = b.uncommittedEvents, nil
			return nil
		}

		// Rollback the transaction as the error occurred.
		if err = tx.Rollback(ctx); err != nil {
			return e.err("rolling back transaction failed", err)
		}

		// Everytime when the save fails due to the already exists error.
		// It means that there already is an event with provided revision.
		// In given case in order to apply the events
		if e.storage.ErrorCode(err) != cgerrors.ErrorCode_AlreadyExists {
			return e.err("saving events failed", err)
		}

		// Reset aggregate, and it's base.
		agg.Reset()
		b.reset()
		agg.SetBase(b)

		// Load all the events, try with snapshot.
		if err = e.LoadEventsWithSnapshot(ctx, agg); err != nil {
			return e.err("loading events with snapshot failed", err)
		}

		// Again try to apply the events.
		for _, event := range events {
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
		return nil, e.err("creating new cursor", err)
	}

	l := newAggregateLoader(c, aggType, aggVersion, factory, e.bufferSize)
	return l.readAggregateChannel()
}

// StreamProjections streams the projection of given aggregate.
func (e *Store) StreamProjections(ctx context.Context, aggType string, aggVersion int64, factory ProjectionFactory) (<-chan Projection, error) {
	c, err := e.storage.NewCursor(ctx, aggType, aggVersion)
	if err != nil {
		return nil, err
	}

	l := newProjectionLoader(e.eventCodec, c, aggType, aggVersion, factory, e.bufferSize)
	return l.readProjectionChannel()
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
	return cgerrors.New("", fmt.Sprintf("%s: %v", msg, err), e.storage.ErrorCode(err))
}
