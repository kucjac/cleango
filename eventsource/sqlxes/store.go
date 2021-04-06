package sqlxes

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/kucjac/cleango/errors"
	"github.com/kucjac/cleango/eventsource"
)

var _ eventsource.Storage = (*sqlStorage)(nil)

// Config is the configuration for the event storage.
type Config struct {
	EventTable    string
	SnapshotTable string
	SchemaName    string // Optional
}

// Validate checks if the config is valid to use.
func (c *Config) Validate() error {
	if c.EventTable == "" {
		return errors.ErrInternal("no event table name provided")
	}
	if c.SnapshotTable == "" {
		return errors.ErrInternal("no snapshot table name provided")
	}
	return nil
}

func (c *Config) eventTableName() string {
	sb := strings.Builder{}
	if c.SchemaName != "" {
		sb.WriteString(c.SchemaName)
		sb.WriteRune('.')
	}
	sb.WriteString(c.EventTable)
	return sb.String()
}

func (c *Config) snapshotTableName() string {
	sb := strings.Builder{}
	if c.SchemaName != "" {
		sb.WriteString(c.SchemaName)
		sb.WriteRune('.')
	}
	sb.WriteString(c.SnapshotTable)
	return sb.String()
}

// New creates a new event storage based on provided sqlx connection.
func New(conn *sqlx.DB, cfg *Config, isErrDupFunc IsErrDuplicatedFunc) (eventsource.Storage, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if isErrDupFunc == nil {
		return nil, errors.ErrInternal("is error duplicated function - not defined")
	}
	return &sqlStorage{conn: conn, cfg: cfg, query: newQueries(conn, cfg), isErrDuplicated: isErrDupFunc}, nil
}

// IsErrDuplicatedFunc is a function that checks if given error is returing duplicated content error (unique constraint violation).
type IsErrDuplicatedFunc func(err error) bool

type sqlStorage struct {
	conn            *sqlx.DB
	cfg             *Config
	query           queries
	isErrDuplicated IsErrDuplicatedFunc
}

const (
	insertEventQuery           = `INSERT INTO %s (aggregate_id, aggregate_type, revision, timestamp, event_id, event_type, event_data) VALUES (?,?,?,?,?,?,?)`
	getEventStreamQuery        = `SELECT aggregate_id, aggregate_type, revision, timestamp, event_id, event_type, event_data FROM %s WHERE aggregate_id = ? AND aggregate_type = ?`
	saveSnapshotQuery          = `INSERT INTO %s (aggregate_id, aggregate_type, aggregate_version, revision, timestamp, snapshot_data) VALUES (?,?,?,?,?,?)`
	getSnapshotQuery           = `SELECT aggregate_id, aggregate_type, aggregate_version, revision, timestamp, snapshot_data FROM %s WHERE aggregate_id = ? AND aggregate_type = ? AND aggregate_version = ?`
	getStreamFromRevisionQuery = `SELECT aggregate_id, aggregate_type, revision, timestamp, event_id, event_type, event_data FROM %s WHERE aggregate_id = ? AND aggregate_type = ? AND revision > ?`
	batchInsertQueryBase       = `INSERT INTO %s (aggregate_id, aggregate_type, revision, timestamp, event_id, event_type, event_data) VALUES `
)

type queries struct {
	getEventStream        string
	saveSnapshot          string
	getSnapshot           string
	getStreamFromRevision string
	insertEvent           string
}

func (q queries) batchInsertEvent(length int) string {
	sb := strings.Builder{}
	sb.WriteString(batchInsertQueryBase)
	for i := 0; i < length; i++ {
		sb.WriteString("(?,?,?,?,?,?,?)")
		sb.WriteRune(',')
	}
	return sb.String()
}

func newQueries(conn *sqlx.DB, c *Config) queries {
	return queries{
		getEventStream:        conn.Rebind(fmt.Sprintf(getEventStreamQuery, c.eventTableName())),
		saveSnapshot:          conn.Rebind(fmt.Sprintf(saveSnapshotQuery, c.snapshotTableName())),
		getSnapshot:           conn.Rebind(fmt.Sprintf(getSnapshotQuery, c.snapshotTableName())),
		getStreamFromRevision: conn.Rebind(fmt.Sprintf(getStreamFromRevisionQuery, c.eventTableName())),
		insertEvent:           conn.Rebind(fmt.Sprintf(insertEventQuery, c.eventTableName())),
	}
}

// SaveEvents stores provided events in the database.
// Implements eventsource.Storage interface.
func (s *sqlStorage) SaveEvents(ctx context.Context, es []*eventsource.Event) error {
	var (
		query  string
		values []interface{}
	)
	switch len(es) {
	case 0:
		return nil
	case 1:
		e := es[0]
		query = s.query.insertEvent
		values = []interface{}{
			e.AggregateId,
			e.AggregateType,
			e.Revision,
			e.Timestamp,
			e.EventId,
			e.EventType,
			e.EventData,
		}
		_, err := s.conn.ExecContext(ctx, query, values...)
		if err != nil {
			if s.isErrDuplicated(err) {
				return errors.ErrAlreadyExists("event revision already exists")
			}
			if errors.Is(err, context.DeadlineExceeded) {
				return errors.ErrDeadlineExceeded(err.Error())
			}
			return errors.ErrInternal(err.Error())
		}
		return nil
	default:
		query = s.conn.Rebind(s.query.batchInsertEvent(len(es)))
		values = make([]interface{}, 7*len(es))
		for i, e := range es {
			values[(i * 7)] = e.AggregateId
			values[(i*7)+1] = e.AggregateType
			values[(i*7)+2] = e.Revision
			values[(i*7)+3] = e.Timestamp
			values[(i*7)+4] = e.EventId
			values[(i*7)+5] = e.EventType
			values[(i*7)+6] = e.EventData
		}
		tx, err := s.conn.BeginTxx(ctx, nil)
		if err != nil {
			return errors.ErrInternal(err.Error())
		}

		_, err = tx.ExecContext(ctx, query, values...)
		if err != nil {
			defer tx.Rollback()
			if s.isErrDuplicated(err) {
				return errors.ErrAlreadyExists("one of given event revision already exists")
			}
			if errors.Is(err, context.DeadlineExceeded) {
				return errors.ErrDeadlineExceeded(err.Error())
			}
			return errors.ErrInternal(err.Error())
		}

		// Commit changes.
		if err = tx.Commit(); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return errors.ErrDeadlineExceeded(err.Error())
			}
			return errors.ErrInternal(err.Error())
		}
		return nil
	}
}

// GetEventStream gets the event stream for provided aggregate.
// Implements eventsource.Storage interface.
func (s *sqlStorage) GetEventStream(ctx context.Context, aggId, aggType string) ([]*eventsource.Event, error) {
	rows, err := s.conn.QueryContext(ctx, s.query.getEventStream, aggId, aggType)
	if err != nil {
		return nil, errors.ErrInternal(err.Error())
	}
	defer rows.Close()

	var stream []*eventsource.Event
	for rows.Next() {
		e := &eventsource.Event{}
		// aggregate_id, aggregate_type, revision, timestamp, event_id, event_type, event_data
		if err = rows.Scan(&e.AggregateId, &e.AggregateType, &e.Revision, &e.Timestamp, &e.EventId, &e.EventType, &e.EventData); err != nil {
			return nil, errors.ErrInternalf("scanning  event row failed: %v", err.Error())
		}
		stream = append(stream, e)
	}
	if err := rows.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return stream, nil
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, errors.ErrDeadlineExceeded(err.Error())
		}
		return nil, errors.ErrInternal(err.Error())
	}
	return stream, nil
}

// SaveSnapshot stores the snapshot in the database.
// Implements eventsource.Storage interface.
func (s *sqlStorage) SaveSnapshot(ctx context.Context, snap *eventsource.Snapshot) error {
	// aggregate_id, aggregate_type, aggregate_version, revision, timestamp, snapshot_data
	_, err := s.conn.ExecContext(ctx, s.query.saveSnapshot, snap.AggregateId, snap.AggregateType, snap.AggregateVersion, snap.Revision, snap.Timestamp, snap.SnapshotData)
	if err != nil {
		if s.isErrDuplicated(err) {
			return errors.ErrAlreadyExists("one of given event revision already exists")
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return errors.ErrDeadlineExceeded(err.Error())
		}
		return errors.ErrInternal(err.Error())
	}
	return nil
}

// GetSnapshot gets the latest snapshot for given aggregate.
// Implements eventsource.Storage interface.
func (s *sqlStorage) GetSnapshot(ctx context.Context, aggId string, aggType string, aggVersion int64) (*eventsource.Snapshot, error) {
	// aggregate_id = ? AND aggregate_type = ? AND aggregate_version = ?
	row := s.conn.QueryRowContext(ctx, s.query.getSnapshot, aggId, aggType, aggVersion)
	if err := row.Err(); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, errors.ErrDeadlineExceeded(err.Error())
		}
		return nil, errors.ErrInternal(err.Error())
	}

	// aggregate_id, aggregate_type, aggregate_version, revision, timestamp, snapshot_data
	var snap eventsource.Snapshot
	if err := row.Scan(&snap.AggregateId, &snap.AggregateType, &snap.AggregateVersion, &snap.Revision, &snap.Timestamp, &snap.SnapshotData); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrNotFound("snapshot not found")
		}
		return nil, errors.ErrInternal(err.Error())
	}
	return &snap, nil
}

// GetStreamFromRevision gets the event stream for given aggregate where the revision is subsequent from provided.
func (s *sqlStorage) GetStreamFromRevision(ctx context.Context, aggId string, aggType string, from int64) ([]*eventsource.Event, error) {
	rows, err := s.conn.QueryContext(ctx, s.query.getStreamFromRevision, aggId, aggType, from)
	if err != nil {
		return nil, errors.ErrInternal(err.Error())
	}
	defer rows.Close()

	var stream []*eventsource.Event
	for rows.Next() {
		e := &eventsource.Event{}
		// aggregate_id, aggregate_type, revision, timestamp, event_id, event_type, event_data
		if err = rows.Scan(&e.AggregateId, &e.AggregateType, &e.Revision, &e.Timestamp, &e.EventId, &e.EventType, &e.EventData); err != nil {
			return nil, errors.ErrInternalf("scanning  event row failed: %v", err.Error())
		}
		stream = append(stream, e)
	}
	if err := rows.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return stream, nil
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, errors.ErrDeadlineExceeded(err.Error())
		}
		return nil, errors.ErrInternal(err.Error())
	}
	return stream, nil
}
