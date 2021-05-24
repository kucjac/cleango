package sqlxes

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kucjac/cleango/xservice"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/eventsource"
)

var _ eventsource.Storage = (*storage)(nil)

// Config is the configuration for the event storage.
type Config struct {
	EventTable     string
	SnapshotTable  string
	SchemaName     string // Optional
	AggregateTable string
	WorkersCount   int
}

// DefaultConfig creates a new default config.
func DefaultConfig() *Config {
	return &Config{
		EventTable:     "event",
		SnapshotTable:  "snapshot",
		AggregateTable: "aggregate",
		WorkersCount:   10,
	}
}

// Validate checks if the config is valid to use.
func (c *Config) Validate() error {
	if c.EventTable == "" {
		return cgerrors.ErrInternal("no event table name provided")
	}
	if c.SnapshotTable == "" {
		return cgerrors.ErrInternal("no snapshot table name provided")
	}
	if c.AggregateTable == "" {
		return cgerrors.ErrInternalf("no aggregate table name provided")
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

func (c *Config) aggregateTableName() string {
	sb := strings.Builder{}
	if c.SchemaName != "" {
		sb.WriteString(c.SchemaName)
		sb.WriteRune('.')
	}
	sb.WriteString(c.AggregateTable)
	return sb.String()
}

// New creates a new event storage based on provided sqlx connection.
func New(conn *sqlx.DB, cfg *Config, d xservice.Driver) (*Storage, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if d == nil {
		return nil, cgerrors.ErrInternal("sqlxes driver not defined")
	}
	if cfg.WorkersCount == 0 {
		cfg.WorkersCount = 10
	}
	return &Storage{storage: storage{conn: conn, cfg: cfg, query: newQueries(conn, cfg), d: d}}, nil
}

type Storage struct {
	storage
}

type storage struct {
	conn       sqlx.ExtContext
	cfg        *Config
	query      queries
	d          xservice.Driver
	maxRetries int
}

// NewCursor creates a new cursor.
func (s *storage) NewCursor(ctx context.Context, aggType string, aggVersion int64) (eventsource.Cursor, error) {
	return s.newCursor(ctx, aggType, aggVersion), nil
}

// SaveEvents stores provided events in the database.
// Implements eventsource.Storage interface.
func (s *storage) SaveEvents(ctx context.Context, es []*eventsource.Event) error {
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

		// If this is initial aggregate revision insert new entry in the aggregate table.
		if e.Revision == 1 {

		}

		var err error
		for i := 1; i <= s.maxRetries; i++ {
			_, err = s.conn.ExecContext(ctx, query, values...)
			if err != nil {
				if s.d.CanRetry(err) {
					time.Sleep(time.Millisecond * 500)
					continue
				}
			}
			break
		}
		if err != nil {
			c := s.d.ErrorCode(err)
			if c == cgerrors.ErrorCode_AlreadyExists {
				return cgerrors.ErrAlreadyExists("event revision already exists")
			}
			if cgerrors.Is(err, context.DeadlineExceeded) {
				return err
			}
			return cgerrors.New("", err.Error(), c)
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

		// Begin transaction.
		tx, began, err := s.getTx(ctx)
		if err != nil {
			return err
		}
		if began {
			defer func() {
				if err == nil {
					tx.Commit()
				} else {
					tx.Rollback()
				}
			}()
		}

		// Execute the query.
		err = s.tryTx(ctx, tx, func(ctx context.Context, tx *sqlx.Tx) error {
			_, err = tx.ExecContext(ctx, query, values...)
			return err
		})
		if err != nil {
			return s.Err(err)
		}

		// Commit changes.
		err = s.tryTx(ctx, tx, func(ctx context.Context, tx *sqlx.Tx) error {
			return tx.Commit()
		})
		if err != nil {
			s.Err(err)
		}
		return nil
	}
}

// Err handles error message with given driver.
func (s *storage) Err(err error) error {
	c := s.d.ErrorCode(err)
	if c == cgerrors.ErrorCode_AlreadyExists {
		return cgerrors.ErrAlreadyExists("event revision already exists")
	}
	if cgerrors.Is(err, context.DeadlineExceeded) {
		return err
	}
	return cgerrors.New("", err.Error(), c)
}

func (s *storage) getTx(ctx context.Context) (tx *sqlx.Tx, began bool, err error) {
	switch db := s.conn.(type) {
	case *sqlx.Tx:
		tx = db
		began = true
	case *sqlx.DB:
		for i := 1; i <= s.maxRetries; i++ {
			tx, err = db.BeginTxx(ctx, nil)
			if err != nil {
				if s.d.CanRetry(err) {
					continue
				}
				return nil, false, err
			}
			break
		}
		return nil, false, err
	default:
		return nil, false, cgerrors.ErrInternal("undefined connection type: %T", s.conn)
	}
	return tx, began, nil
}

// GetEventStream gets the event stream for provided aggregate.
// Implements eventsource.Storage interface.
func (s *storage) GetEventStream(ctx context.Context, aggId, aggType string) ([]*eventsource.Event, error) {
	var rows *sqlx.Rows
	err := s.try(ctx, s.conn, func(ctx context.Context, db sqlx.ExtContext) error {
		var err error
		rows, err = db.QueryxContext(ctx, s.query.getEventStream, aggId, aggType)
		return err
	})
	if err != nil {
		return nil, cgerrors.New("", err.Error(), s.d.ErrorCode(err))
	}
	defer rows.Close()

	var stream []*eventsource.Event
	for rows.Next() {
		e := &eventsource.Event{}
		// aggregate_id, aggregate_type, revision, timestamp, event_id, event_type, event_data
		if err = rows.Scan(&e.AggregateId, &e.AggregateType, &e.Revision, &e.Timestamp, &e.EventId, &e.EventType, &e.EventData); err != nil {
			return nil, cgerrors.ErrInternalf("scanning  event row failed: %v", err.Error())
		}
		stream = append(stream, e)
	}
	if err = rows.Err(); err != nil {
		if cgerrors.Is(err, sql.ErrNoRows) {
			return stream, nil
		}
		if cgerrors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		return nil, cgerrors.New("", err.Error(), s.d.ErrorCode(err))
	}
	return stream, nil
}

// SaveSnapshot stores the snapshot in the database.
// Implements eventsource.Storage interface.
func (s *storage) SaveSnapshot(ctx context.Context, snap *eventsource.Snapshot) error {
	// aggregate_id, aggregate_type, aggregate_version, revision, timestamp, snapshot_data
	err := s.try(ctx, s.conn, func(ctx context.Context, db sqlx.ExtContext) error {
		_, err := s.conn.ExecContext(ctx, s.query.saveSnapshot, snap.AggregateId, snap.AggregateType, snap.AggregateVersion, snap.Revision, snap.Timestamp, snap.SnapshotData)
		return err
	})
	if err != nil {
		c := s.d.ErrorCode(err)
		if c == cgerrors.ErrorCode_AlreadyExists {
			return cgerrors.ErrAlreadyExists("one of given event revision already exists")
		}
		if cgerrors.Is(err, context.DeadlineExceeded) {
			return err
		}
		return cgerrors.New("", err.Error(), c)
	}
	return nil
}

// GetSnapshot gets the latest snapshot for given aggregate.
// Implements eventsource.Storage interface.
func (s *storage) GetSnapshot(ctx context.Context, aggId string, aggType string, aggVersion int64) (*eventsource.Snapshot, error) {
	// aggregate_id = ? AND aggregate_type = ? AND aggregate_version = ?
	var row *sqlx.Row
	err := s.try(ctx, s.conn, func(ctx context.Context, db sqlx.ExtContext) error {
		row = db.QueryRowxContext(ctx, s.query.getSnapshot, aggId, aggType, aggVersion)
		return row.Err()
	})
	if err != nil {
		if cgerrors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		c := s.d.ErrorCode(err)
		return nil, cgerrors.New("", err.Error(), c)
	}

	// aggregate_id, aggregate_type, aggregate_version, revision, timestamp, snapshot_data
	var snap eventsource.Snapshot
	if err := row.Scan(&snap.AggregateId, &snap.AggregateType, &snap.AggregateVersion, &snap.Revision, &snap.Timestamp, &snap.SnapshotData); err != nil {
		if cgerrors.Is(err, sql.ErrNoRows) {
			return nil, cgerrors.ErrNotFound("snapshot not found")
		}
		return nil, cgerrors.New("", err.Error(), s.d.ErrorCode(err))
	}
	return &snap, nil
}

// GetStreamFromRevision gets the event stream for given aggregate where the revision is subsequent from provided.
func (s *storage) GetStreamFromRevision(ctx context.Context, aggId string, aggType string, from int64) ([]*eventsource.Event, error) {
	rows, err := s.conn.QueryContext(ctx, s.query.getStreamFromRevision, aggId, aggType, from)
	if err != nil {
		return nil, cgerrors.ErrInternal(err.Error())
	}
	defer rows.Close()

	var stream []*eventsource.Event
	for rows.Next() {
		e := &eventsource.Event{}
		// aggregate_id, aggregate_type, revision, timestamp, event_id, event_type, event_data
		if err = rows.Scan(&e.AggregateId, &e.AggregateType, &e.Revision, &e.Timestamp, &e.EventId, &e.EventType, &e.EventData); err != nil {
			return nil, cgerrors.ErrInternalf("scanning  event row failed: %v", err.Error())
		}
		stream = append(stream, e)
	}
	if err := rows.Err(); err != nil {
		if cgerrors.Is(err, sql.ErrNoRows) {
			return stream, nil
		}
		if cgerrors.Is(err, context.DeadlineExceeded) {
			return nil, cgerrors.ErrDeadlineExceeded(err.Error())
		}
		return nil, cgerrors.ErrInternal(err.Error())
	}
	return stream, nil
}

func (s *storage) StreamEvents(ctx context.Context, req *eventsource.StreamEventsRequest) (<-chan *eventsource.Event, error) {
	c := s.newStreamCursor(ctx, req)
	return c.openChannel()
}

func (s *storage) try(ctx context.Context, db sqlx.ExtContext, fn func(context.Context, sqlx.ExtContext) error) error {
	var err error
	for i := 1; i <= s.maxRetries; i++ {
		if err = fn(ctx, db); err != nil {
			if s.d.CanRetry(err) {
				continue
			}
			return err
		}
		break
	}
	return err
}

func (s *storage) tryTx(ctx context.Context, tx *sqlx.Tx, fn func(ctx context.Context, tx *sqlx.Tx) error) error {
	var err error
	for i := 1; i <= s.maxRetries; i++ {
		if err = fn(ctx, tx); err != nil {
			if s.d.CanRetry(err) {
				continue
			}
			return err
		}
		break
	}
	return err
}
