package esxsql

import (
	"context"
	"database/sql"
	"strings"

	"github.com/kucjac/cleango/database"
	"github.com/kucjac/cleango/database/xsql"
	"github.com/kucjac/cleango/xlog"
	uuid "github.com/satori/go.uuid"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/database/es"
)

var _ es.Storage = (*Storage)(nil)

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
func New(conn xsql.DB, cfg *Config, d database.Driver) (*Storage, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if d == nil {
		return nil, cgerrors.ErrInternal("esxsql driver not defined")
	}
	if cfg.WorkersCount == 0 {
		cfg.WorkersCount = 10
	}
	return &Storage{storage: storage{conn: conn, cfg: cfg, query: newQueries(conn, cfg), d: d}}, nil
}

// Storage is the implementation of the eventsource.Storage interface for the sqlx driver.
type Storage struct {
	storage
}

// BeginTx creates and begins a new transaction, which exposes *sqlx.Tx and allows atomic commits.
func (s *Storage) BeginTx(ctx context.Context) (*Transaction, error) {
	tx, _, err := s.getTx(ctx)
	if err != nil {
		return nil, err
	}
	st := s.storage
	st.conn = tx
	return &Transaction{id: uuid.NewV4().String(), storage: st}, nil
}

// As exposes driver specific implementation.
func (s *Storage) As(dst interface{}) error {
	db, ok := s.conn.(*xsql.Conn)
	if !ok {
		return cgerrors.ErrInternalf("invalid esxsql.Storage conn type: %T", s.conn)
	}
	ddb, ok := dst.(**xsql.Conn)
	if !ok {
		return cgerrors.ErrInternalf("invalid input type: %T, wanted **sqlx.DB", dst)
	}
	*ddb = db
	return nil
}

type storage struct {
	conn  xsql.DB
	cfg   *Config
	query queries
	d     database.Driver
}

func (s *storage) ErrorCode(err error) cgerrors.ErrorCode {
	return s.d.ErrorCode(err)
}

// NewCursor creates a new cursor.
func (s *storage) NewCursor(ctx context.Context, aggType string, aggVersion int64) (es.Cursor, error) {
	return s.newCursor(ctx, aggType, aggVersion), nil
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

// SaveEvents stores provided events in the database.
// Implements eventsource.Storage interface.
func (s *storage) SaveEvents(ctx context.Context, es []*es.Event) error {
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

		var err error
		// If this is initial aggregate revision insert new entry in the aggregate table.
		if e.Revision == 1 {
			err = xsql.RunInTransaction(ctx, s.conn, func(tx *xsql.Tx) error {
				_, err := tx.ExecContext(ctx, query, values...)
				if err != nil {
					return err
				}
				_, err = tx.ExecContext(ctx, s.query.insertAggregate, e.AggregateId, e.AggregateType, e.Timestamp)
				if err != nil {
					return err
				}
				return nil
			})
		} else {
			_, err = s.conn.ExecContext(ctx, query, values...)
		}
		if err != nil {
			return err
		}
		return nil
	default:
		query = s.conn.Rebind(s.query.batchInsertEvent(len(es)))
		values = make([]interface{}, 7*len(es))
		var aggregates []aggregate
		for i, e := range es {
			// Check if the aggregate needs to be inserted.
			if e.Revision == 1 {
				aggregates = append(aggregates, aggregate{
					ID:        e.AggregateId,
					Type:      e.AggregateType,
					Timestamp: e.Timestamp,
				})
			}

			values[(i * 7)] = e.AggregateId
			values[(i*7)+1] = e.AggregateType
			values[(i*7)+2] = e.Revision
			values[(i*7)+3] = e.Timestamp
			values[(i*7)+4] = e.EventId
			values[(i*7)+5] = e.EventType
			values[(i*7)+6] = e.EventData
		}
		err := xsql.RunInTransaction(ctx, s.conn, func(tx *xsql.Tx) error {
			// Execute the query.
			_, err := tx.ExecContext(ctx, query, values...)
			if err != nil {
				return err
			}
			for _, agg := range aggregates {
				if _, err = tx.ExecContext(ctx, s.query.insertAggregate, agg.ID, agg.Type, agg.Timestamp); err != nil {
					return err
				}
			}

			return err
		})
		if err != nil {
			xlog.Debugf("Saving events failed: %v", err)
			return err
		}
		return nil
	}
}

// ListEvents gets the event stream for provided aggregate.
// Implements eventsource.Storage interface.
func (s *storage) ListEvents(ctx context.Context, aggId, aggType string) ([]*es.Event, error) {
	rows, err := s.conn.QueryContext(ctx, s.query.getEventStream, aggId, aggType)
	if err != nil {
		return nil, cgerrors.New("", err.Error(), s.d.ErrorCode(err))
	}
	defer rows.Close()

	var stream []*es.Event
	for rows.Next() {
		e := &es.Event{}
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
		return nil, err
	}
	return stream, nil
}

// SaveSnapshot stores the snapshot in the database.
// Implements eventsource.Storage interface.
func (s *storage) SaveSnapshot(ctx context.Context, snap *es.Snapshot) error {
	// aggregate_id, aggregate_type, aggregate_version, revision, timestamp, snapshot_data
	_, err := s.conn.ExecContext(ctx, s.query.saveSnapshot, snap.AggregateId, snap.AggregateType, snap.AggregateVersion, snap.Revision, snap.Timestamp, snap.SnapshotData)
	return err
}

// GetSnapshot gets the latest snapshot for given aggregate.
// Implements eventsource.Storage interface.
func (s *storage) GetSnapshot(ctx context.Context, aggId string, aggType string, aggVersion int64) (*es.Snapshot, error) {
	// aggregate_id = ? AND aggregate_type = ? AND aggregate_version = ?
	row := s.conn.QueryRowContext(ctx, s.query.getSnapshot, aggId, aggType, aggVersion)
	if err := row.Err(); err != nil {
		return nil, err
	}

	// aggregate_id, aggregate_type, aggregate_version, revision, timestamp, snapshot_data
	var snap es.Snapshot
	if err := row.Scan(&snap.AggregateId, &snap.AggregateType, &snap.AggregateVersion, &snap.Revision, &snap.Timestamp, &snap.SnapshotData); err != nil {
		if cgerrors.Is(err, sql.ErrNoRows) {
			return nil, cgerrors.ErrNotFound("snapshot not found")
		}
		return nil, cgerrors.New("", err.Error(), s.d.ErrorCode(err))
	}
	return &snap, nil
}

// ListEventsAfterRevision gets the event stream for given aggregate where the revision is subsequent from provided.
func (s *storage) ListEventsAfterRevision(ctx context.Context, aggId string, aggType string, after int64) ([]*es.Event, error) {
	rows, err := s.conn.QueryContext(ctx, s.query.getStreamAfterRevision, aggId, aggType, after)
	if err != nil {
		return nil, cgerrors.ErrInternal(err.Error())
	}
	defer rows.Close()

	var stream []*es.Event
	for rows.Next() {
		e := &es.Event{}
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

// StreamEvents opens the channel of the events stream that matches given request.
// Implements eventsource.Storage.
func (s *storage) StreamEvents(ctx context.Context, req *es.StreamEventsRequest) (<-chan *es.Event, error) {
	c := s.newStreamCursor(ctx, req)
	return c.openChannel()
}

func (s *storage) getTx(ctx context.Context) (tx *xsql.Tx, began bool, err error) {
	switch db := s.conn.(type) {
	case *xsql.Tx:
		tx = db
		began = true
	case *xsql.Conn:
		tx, err = db.BeginTx(ctx, nil)
		if err != nil {
			return nil, false, err
		}
	default:
		return nil, false, cgerrors.ErrInternalf("undefined connection type: %T", s.conn)
	}
	return tx, began, nil
}

type aggregate struct {
	ID        string
	Type      string
	Timestamp int64
}
