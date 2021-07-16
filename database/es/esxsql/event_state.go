package esxsql

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/database/es"
	"github.com/kucjac/cleango/database/es/eventstate"
	"github.com/kucjac/cleango/database/xsql"
	uuid "github.com/satori/go.uuid"
)

// Compile time check if StateStorage implements eventstate.Storage interface.
var _ eventstate.Storage = (*StateStorage)(nil)

// StateStorage is the implementation of the eventstate.Storage interface.
// It also implements es.StorageBase.
type StateStorage struct {
	storage
}

// BeginTx starts a new transaction.
func (s *StateStorage) BeginTx(ctx context.Context) (eventstate.TxStorage, error) {
	tx, _, err := s.getTx(ctx)
	if err != nil {
		return nil, err
	}
	st := s.storage
	st.conn = tx
	return &Transaction{id: uuid.NewV4().String(), storage: st}, nil
}

// As exposes driver specific implementation.
func (s *StateStorage) As(dst interface{}) error {
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

// Config gets the storage config.
func (s *StateStorage) Config() Config {
	return *s.cfg
}

// NewStateStorage creates a new event storage based on provided sqlx connection.
func NewStateStorage(conn *xsql.Conn, cfg *Config) (*StateStorage, error) {
	if cfg == nil {
		return nil, cgerrors.ErrInternal("no storage config provided")
	}
	cfg.useEventState = true
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if cfg.WorkersCount == 0 {
		cfg.WorkersCount = 10
	}
	return &StateStorage{storage: storage{conn: conn, cfg: cfg, query: newQueries(conn, cfg)}}, nil
}

// FindUnhandled implements eventstate.StorageBase interface.
// Finds all unhandled event state matching given query.
func (s *storage) FindUnhandled(ctx context.Context, query eventstate.FindUnhandledQuery) ([]eventstate.Unhandled, error) {
	q := s.query.findHandlerEvents
	args := []interface{}{eventstate.StateUnhandled}
	sb := strings.Builder{}
	sb.WriteString(q)
	sb.WriteString(" WHERE state = ?")
	if len(query.HandlerNames) != 0 {
		sb.WriteString(" AND es.handler_name IN (")
		for i, hn := range query.HandlerNames {
			sb.WriteRune('?')
			if i < len(query.HandlerNames)-1 {
				sb.WriteRune(',')
			}
			args = append(args, hn)
		}
		sb.WriteRune(')')
	}
	q = s.conn.Rebind(sb.String())

	rows, err := s.conn.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, cgerrors.New("", "finding unhandled event state failed", s.conn.ErrorCode(err)).
			WithMeta("err", err.Error())
	}
	defer rows.Close()

	var result []eventstate.Unhandled
	for rows.Next() {
		e := &es.Event{}
		var handlerName string
		// aggregate_id, aggregate_type, revision, timestamp, event_id, event_type, event_data
		err = rows.Scan(
			&e.AggregateId, &e.AggregateType, &e.Revision, &e.Timestamp, &e.EventId, &e.EventType, &e.EventData,
			&handlerName,
		)
		if err != nil {
			return nil, cgerrors.ErrInternalf("scanning  event row failed: %v", err.Error())
		}
		result = append(result, eventstate.Unhandled{Event: e, HandlerName: handlerName})
	}
	if err = rows.Err(); err != nil {
		if cgerrors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, cgerrors.New("", "finding unhandled event state - rows error", s.conn.ErrorCode(err)).
			WithMeta("err", err.Error())
	}
	return result, nil
}

// FindFailures implements eventstate.StorageBase.
func (s *storage) FindFailures(ctx context.Context, query eventstate.FindFailureQuery) ([]eventstate.HandleFailure, error) {
	q := s.query.findHandlingFailures
	var args []interface{}
	if len(query.HandlerNames) != 0 {
		sb := strings.Builder{}
		sb.WriteString(q)
		sb.WriteString(" AND es.handler_name IN (")
		for i, hn := range query.HandlerNames {
			sb.WriteRune('?')
			if i < len(query.HandlerNames)-1 {
				sb.WriteRune(',')
			}
			args = append(args, hn)
		}
		sb.WriteRune(')')
		q = s.conn.Rebind(sb.String())
	}

	rows, err := s.conn.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, cgerrors.New("", "finding unhandled event state failed", s.conn.ErrorCode(err)).
			WithMeta("err", err.Error())
	}
	defer rows.Close()

	var result []eventstate.HandleFailure
	for rows.Next() {
		e := &es.Event{}
		var (
			handlerName string
			timestamp   int64
			errMessage  string
			errCode     int
			retryNo     int
		)
		// aggregate_id, aggregate_type, revision, timestamp, event_id, event_type, event_data
		err = rows.Scan(
			&e.AggregateId, &e.AggregateType, &e.Revision, &e.Timestamp, &e.EventId, &e.EventType, &e.EventData,
			&handlerName, &timestamp, &errMessage, &errCode, &retryNo,
		)
		if err != nil {
			return nil, cgerrors.ErrInternalf("scanning  event row failed: %v", err.Error())
		}
		result = append(result, eventstate.HandleFailure{
			Event:       e,
			HandlerName: handlerName,
			Err:         errMessage,
			ErrCode:     cgerrors.ErrorCode(errCode),
			RetryNo:     retryNo,
			Timestamp:   time.Unix(0, timestamp),
		})
	}
	if err = rows.Err(); err != nil {
		if cgerrors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, cgerrors.New("", "finding unhandled event state - rows error", s.conn.ErrorCode(err)).
			WithMeta("err", err.Error())
	}
	return result, nil
}

// MarkUnhandled implements eventstate.StorageBase.
func (s *storage) MarkUnhandled(ctx context.Context, events ...*es.Event) error {
	for _, e := range events {
		q := s.query.insertEventState
		_, err := s.conn.ExecContext(ctx, q, e.EventId, eventstate.StateUnhandled, e.Timestamp, e.EventType)
		if err != nil {
			return cgerrors.New("", "marking event unhandled failed", s.conn.ErrorCode(err)).
				WithMeta("err", err.Error())
		}
	}
	return nil
}

// StartHandling implements eventstate.StorageBase.
func (s *storage) StartHandling(ctx context.Context, e *es.Event, handlerName string, timestamp int64) error {
	q := s.query.updateEventState
	_, err := s.conn.ExecContext(ctx, q, eventstate.StateStarted, timestamp, e.EventId, handlerName)
	if err != nil {
		return s.Err(err)
	}
	return nil
}

// FinishHandling implements eventstate.StorageBase.
func (s *storage) FinishHandling(ctx context.Context, e *es.Event, handlerName string, timestamp int64) error {
	q := s.query.updateEventState
	_, err := s.conn.ExecContext(ctx, q, eventstate.StateFinished, timestamp, e.EventId, handlerName)
	if err != nil {
		return s.Err(err)
	}
	return nil
}

// HandlingFailed implements eventstate.StorageBase.
func (s *storage) HandlingFailed(ctx context.Context, failure *eventstate.HandleFailure) error {
	q := s.query.updateEventState
	_, err := s.conn.ExecContext(ctx, q, eventstate.StateFailed, failure.Timestamp.UnixNano(), failure.Event.EventId, failure.HandlerName)
	if err != nil {
		return s.Err(err)
	}

	q = s.query.insertHandlingFailure
	_, err = s.conn.ExecContext(ctx, q, failure.Event.EventId, failure.HandlerName, failure.Timestamp.UnixNano(), failure.Err, failure.ErrCode, failure.RetryNo)
	if err != nil {
		return s.Err(err)
	}
	return nil
}
