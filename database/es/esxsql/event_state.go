package esxsql

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/database/es/esstate"
	"github.com/kucjac/cleango/database/xsql"
	"github.com/kucjac/cleango/ddd/events/eventstate"
	uuid "github.com/satori/go.uuid"
)

// Compile time check if StateStorage implements eventstate.Storage interface.
var _ esstate.Storage = (*StateStorage)(nil)

// StateStorage is the implementation of the eventstate.Storage interface.
// It also implements es.StorageBase.
type StateStorage struct {
	storage
}

// BeginTx starts a new transaction.
func (s *StateStorage) BeginTx(ctx context.Context) (esstate.TxStorage, error) {
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
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if cfg.EventState == nil {
		return nil, cgerrors.ErrInternal("no event state config found")
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
	args := []interface{}{esstate.StateUnhandled}
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
		var eventID, handlerName string
		// aggregate_id, aggregate_type, revision, timestamp, event_id, event_type, event_data
		err = rows.Scan(
			&eventID, &handlerName,
		)
		if err != nil {
			return nil, cgerrors.ErrInternalf("scanning  event row failed: %v", err.Error())
		}
		result = append(result, eventstate.Unhandled{EventID: eventID, HandlerName: handlerName})
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
		sb.WriteString(" WHERE ")
		sb.WriteString("ef.handler_name IN (")
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
		var (
			eventID     string
			handlerName string
			timestamp   int64
			errMessage  string
			errCode     int
			retryNo     int
		)

		err = rows.Scan(&eventID, &handlerName, &timestamp, &errMessage, &errCode, &retryNo)
		if err != nil {
			return nil, cgerrors.ErrInternalf("scanning  event row failed: %v", err.Error())
		}
		result = append(result, eventstate.HandleFailure{
			EventID:     eventID,
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
func (s *storage) MarkUnhandled(ctx context.Context, eventID, eventType string, timestamp int64) error {
	q := s.query.insertEventState
	_, err := s.conn.ExecContext(ctx, q, eventID, esstate.StateUnhandled, timestamp, eventType)
	if err != nil {
		return cgerrors.New("", "marking event unhandled failed", s.conn.ErrorCode(err)).
			WithMeta("err", err.Error())
	}
	return nil
}

// StartHandling implements eventstate.StorageBase.
func (s *storage) StartHandling(ctx context.Context, eventID string, handlerName string, timestamp int64) error {
	q := s.query.updateEventState
	_, err := s.conn.ExecContext(ctx, q, esstate.StateStarted, timestamp, eventID, handlerName)
	if err != nil {
		return s.Err(err)
	}
	return nil
}

// FinishHandling implements eventstate.StorageBase.
func (s *storage) FinishHandling(ctx context.Context, eventID string, handlerName string, timestamp int64) error {
	q := s.query.updateEventState
	_, err := s.conn.ExecContext(ctx, q, esstate.StateFinished, timestamp, eventID, handlerName)
	if err != nil {
		return s.Err(err)
	}
	return nil
}

// HandlingFailed implements eventstate.StorageBase.
func (s *storage) HandlingFailed(ctx context.Context, failure *eventstate.HandleFailure) error {
	q := s.query.updateEventState
	_, err := s.conn.ExecContext(ctx, q, esstate.StateFailed, failure.Timestamp.UnixNano(), failure.EventID, failure.HandlerName)
	if err != nil {
		return s.Err(err)
	}

	q = s.query.insertHandlingFailure
	_, err = s.conn.ExecContext(ctx, q, failure.EventID, failure.HandlerName, failure.Timestamp.UnixNano(), failure.Err, failure.ErrCode, failure.RetryNo)
	if err != nil {
		return s.Err(err)
	}
	return nil
}
