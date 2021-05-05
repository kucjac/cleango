package eventsource

import (
	"context"
)

//go:generate mockgen -destination=mockes/storage_gen.go -package=mockes . Storage
//go:generate mockgen -destination=mockes/cursor_gen.go -package=mockes . Cursor

// Storage is the interface used by the event store as a storage for its events and snapshots.
type Storage interface {
	SaveEvents(ctx context.Context, es []*Event) error
	GetEventStream(ctx context.Context, aggId string, aggType string) ([]*Event, error)
	SaveSnapshot(ctx context.Context, snap *Snapshot) error
	GetSnapshot(ctx context.Context, aggId string, aggType string, aggVersion int64) (*Snapshot, error)
	GetStreamFromRevision(ctx context.Context, aggId string, aggType string, from int64) ([]*Event, error)
	NewCursor(ctx context.Context, aggType string, aggVersion int64) (Cursor, error)
	StreamEvents(ctx context.Context, req *StreamEventsRequest) (<-chan *Event, error)
}

// Cursor is an interface used by the storages that enables listing the aggregates.
type Cursor interface {
	OpenChannel(withSnapshot bool) (<-chan *CursorAggregate, error)
}

// CursorAggregate is an aggregate events and snapshot taken by the cursor.
type CursorAggregate struct {
	AggregateID string
	Snapshot    *Snapshot
	Events      []*Event
}
