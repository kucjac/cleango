package eventsource

import (
	"context"
	"time"

	"github.com/kucjac/cleango/errors"
	"github.com/kucjac/cleango/pubsub/codec"
)

//go:generate mockgen -destination=internal/storemock/store_gen.go -package=storemock . Store

// EventStore is an interface used by the event store to load, commit and create snapshot on aggregates.
type EventStore interface {
	SetAggregateBase(agg Aggregate, aggId, aggType string, version int64)
	LoadEventStream(ctx context.Context, aggregate Aggregate) error
	LoadEventStreamWithSnapshot(ctx context.Context, aggregate Aggregate) error
	Commit(ctx context.Context, aggregate Aggregate) error
	SaveSnapshot(ctx context.Context, aggregate Aggregate) error
}

// New creates new EventStore implementation.
func New(eventCodec codec.Codec, snapCodec codec.Codec, storage Storage) EventStore {
	return &eventStore{
		aggBaseSetter: newAggregateBaseSetter(eventCodec, snapCodec, UUIDGenerator{}),
		snapCodec:     snapCodec,
		storage:       storage,
	}
}

type eventStore struct {
	*aggBaseSetter
	snapCodec codec.Codec
	storage   Storage
}

// LoadEventStream gets the event stream and applies on provided aggregate.
func (e *eventStore) LoadEventStream(ctx context.Context, agg Aggregate) error {
	b := agg.AggBase()
	// Get the full event stream for given aggregate.
	events, err := e.storage.GetEventStream(ctx, b.id, b.aggType)
	if err != nil {
		return err
	}

	// If no events are found return an error.
	if len(events) == 0 {
		return errors.ErrNotFoundf("aggregate: %s with id: %s not found", b.aggType, b.id)
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
func (e *eventStore) LoadEventStreamWithSnapshot(ctx context.Context, agg Aggregate) error {
	// At first try to get the snapshot for given aggregate.
	b := agg.AggBase()
	snap, err := e.storage.GetSnapshot(ctx, b.id, b.aggType, b.version)
	isNotFound := errors.IsNotFound(err)
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
			return errors.ErrNotFoundf("aggregate: %s with id %s not found", b.aggType, b.id)
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
func (e *eventStore) SaveSnapshot(ctx context.Context, agg Aggregate) error {
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
func (e *eventStore) Commit(ctx context.Context, agg Aggregate) error {
	b := agg.AggBase()
	events := b.uncommittedEvents
	for {
		err := e.storage.SaveEvents(ctx, events)
		if err == nil {
			return nil
		}
		if !errors.IsAlreadyExists(err) {
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
			event.Timestamp = time.Now().Unix()
			if err = agg.Apply(event); err != nil {
				return err
			}
		}
	}
}
