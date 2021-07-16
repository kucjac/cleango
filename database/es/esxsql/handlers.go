package esxsql

import (
	"context"

	"github.com/kucjac/cleango/database/es/eventstate"
)

// RegisterHandlers implements eventstate.StorageBase.
func (s *storage) RegisterHandlers(ctx context.Context, eventHandlers ...eventstate.Handler) error {
	q := s.query.registerHandler
	for _, eventHandler := range eventHandlers {
		for _, et := range eventHandler.EventTypes {
			_, err := s.conn.ExecContext(ctx, q, eventHandler.Name, et)
			if err != nil {
				return s.Err(err)
			}
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
