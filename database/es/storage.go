package es

import (
	"context"

	"github.com/kucjac/cleango/cgerrors"
)

//go:generate mockgen -destination=mockes/storage_gen.go -package=mockes . Storage,TxStorage
//go:generate mockgen -destination=mockes/cursor_gen.go -package=mockes . Cursor

// StorageBase is the interface used by the event store as a storage for its events and snapshots.
type StorageBase interface {
	// SaveEvents all input events atomically in the storage.
	SaveEvents(ctx context.Context, es []*Event) error

	// ListEvents lists all events for given aggregate type with given id.
	ListEvents(ctx context.Context, aggId string, aggType string) ([]*Event, error)

	// SaveSnapshot stores a snapshot.
	SaveSnapshot(ctx context.Context, snap *Snapshot) error

	// GetSnapshot gets the snapshot of the aggregate with it's id, type and version.
	GetSnapshot(ctx context.Context, aggId string, aggType string, aggVersion int64) (*Snapshot, error)

	// ListEventsAfterRevision gets the event stream for given aggregate id, type starting after given revision.
	ListEventsAfterRevision(ctx context.Context, aggId string, aggType string, from int64) ([]*Event, error)

	// NewCursor creates a new cursor of given aggregate type and version.
	NewCursor(ctx context.Context, aggType string, aggVersion int64) (Cursor, error)

	// StreamEvents streams the events that matching given request.
	StreamEvents(ctx context.Context, req *StreamEventsRequest) (<-chan *Event, error)

	// As allows drivers to expose driver-specific types.
	As(dst interface{}) error

	// ErrorCode gets the error code from the storage.
	ErrorCode(err error) cgerrors.ErrorCode
}

// Cursor is an interface used by the storages that enables listing the aggregates.
type Cursor interface {
	// GetAggregateStream gets the stream of aggregates.
	GetAggregateStream(withSnapshot bool) (<-chan *CursorAggregate, error)
}

// CursorAggregate is an aggregate events and snapshot taken by the cursor.
type CursorAggregate struct {
	AggregateID string
	Snapshot    *Snapshot
	Events      []*Event
}

// Storage is a transaction beginner.
type Storage interface {
	BeginTx(ctx context.Context) (TxStorage, error)
	StorageBase
}

type TxStorage interface {
	StorageBase
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
