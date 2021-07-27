package esxsql

import (
	"context"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/database/es/eventstate"
)

// RegisterHandlers implements eventstate.StorageBase.
func (s *storage) RegisterHandlers(ctx context.Context, eventHandlers ...eventstate.Handler) error {
	if s.cfg.EventState == nil {
		return cgerrors.ErrInternal("undefined event state for the storage")
	}
	q := s.query.registerHandler
	for _, eventHandler := range eventHandlers {
		for _, et := range eventHandler.EventTypes {
			_, err := s.conn.ExecContext(ctx, q, eventHandler.Name, et)
			if err != nil {
				return s.Err(err)
			}
		}
	}

	// Check if event state needs partition on top of new handlers.
	if !s.cfg.EventState.PartitionState {
		return nil
	}

	switch s.conn.DriverName() {
	case "pg", "postgres", "postgresql", "gopg", "pgx":
		handlerNames := make([]string, len(eventHandlers))
		for i, h := range eventHandlers {
			handlerNames[i] = h.Name
		}
		if err := migratePostgresEventStatePartitions(ctx, s.conn, s.cfg, handlerNames); err != nil {
			return err
		}
	}
	return nil
}

// ListHandlers implements eventstate.StorageBase.
func (s *storage) ListHandlers(ctx context.Context) ([]eventstate.Handler, error) {
	q := s.query.listHandlers
	rows, err := s.conn.QueryContext(ctx, q)
	if err != nil {
		return nil, s.Err(err)
	}
	defer rows.Close()

	var handlers []eventstate.Handler
	for rows.Next() {
		var h eventstate.Handler
		if err = rows.Scan(h.Name, h.EventTypes); err != nil {
			return nil, s.Err(err)
		}
		handlers = append(handlers, h)
	}
	if err = rows.Err(); err != nil {
		return nil, s.Err(err)
	}
	return handlers, nil
}
