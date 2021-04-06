package eventsource

import (
	"context"
)

// Storage is the interface used by the event store as a storage for its events and snapshots.
type Storage interface {
	SaveEvents(ctx context.Context, es []*Event) error
	GetEventStream(ctx context.Context, aggId string, aggType string) ([]*Event, error)
	SaveSnapshot(ctx context.Context, snap *Snapshot) error
	GetSnapshot(ctx context.Context, aggId string, aggType string, aggVersion int64) (*Snapshot, error)
	GetStreamFromRevision(ctx context.Context, aggId string, aggType string, from int64) ([]*Event, error)
}
