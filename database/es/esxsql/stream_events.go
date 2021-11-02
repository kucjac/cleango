package esxsql

import (
	"context"
	"strings"
	"sync"

	"github.com/kucjac/cleango/database/es"
	"github.com/kucjac/cleango/database/xsql"
	"github.com/kucjac/cleango/pkg/xlog"
)

type streamEventsCursor struct {
	ctx         context.Context
	cancelFunc  context.CancelFunc
	l           sync.Mutex
	conn        xsql.DB
	s           *storage
	req         *es.StreamEventsRequest
	lastTakenID uint64
	query       queries
	limit       int64
}

func (s *storage) newStreamCursor(ctx context.Context, req *es.StreamEventsRequest) *streamEventsCursor {
	ctx, cancelFunc := context.WithCancel(ctx)
	return &streamEventsCursor{
		ctx:        ctx,
		cancelFunc: cancelFunc,
		conn:       s.conn,
		query:      s.query,
		limit:      int64(req.BuffSize),
		s:          s,
		req:        req,
	}
}

func (c *streamEventsCursor) openChannel() (<-chan *es.Event, error) {
	ch := make(chan *es.Event, c.limit)
	go c.startReadingEvents(ch)
	return ch, nil
}

func (c *streamEventsCursor) startReadingEvents(ca chan *es.Event) {
	var err error

	q := c.buildQuery()
	for {
		select {
		case <-c.ctx.Done():
			err = c.ctx.Err()
			break
		default:
		}

		var rows *xsql.Rows
		rows, err = c.conn.QueryContext(c.ctx, q.query, append(q.args, c.lastTakenID, c.limit)...)
		if err != nil {
			if c.conn.CanRetry(err) {
				continue
			}
			break
		}

		var rowsCount int
		for rows.Next() {
			var (
				id uint64
				e  es.Event
			)
			// id, aggregate_id, aggregate_type, revision, timestamp, event_id, event_type, event_data
			if err = rows.Scan(&id, &e.AggregateId, &e.AggregateType, &e.Revision, &e.Timestamp, &e.EventId, &e.EventType, &e.EventData); err != nil {
				break
			}

			c.lastTakenID = id
			ca <- &e
			rowsCount++
			// Either error channel context or workers are finished.
			select {
			case <-c.ctx.Done():
				err = c.ctx.Err()
				break
			default:
			}
		}
		rows.Close()
		if err != nil {
			break
		}

		// If there is no more rows to read break the loop and close the channel.
		if rowsCount == 0 {
			break
		}
	}
	if err != nil {
		xlog.Errorf("reading aggregates failed: %v", err)
	}

	close(ca)
}

type streamEventsQuery struct {
	query string
	args  []interface{}
}

func (c *streamEventsCursor) buildQuery() *streamEventsQuery {
	sb := strings.Builder{}
	sb.WriteString(c.query.listEventStreamQuery)
	sb.WriteString("WHERE ")

	q := streamEventsQuery{}
	if len(c.req.AggregateIDs) != 0 {
		sb.WriteString("aggregate_id IN (")
		for i, id := range c.req.AggregateIDs {
			sb.WriteRune('?')
			q.args = append(q.args, id)
			if i != len(c.req.AggregateIDs)-1 {
				sb.WriteRune(',')
			}
		}
		sb.WriteString(") ")

	}

	if len(c.req.AggregateTypes) != 0 {
		sb.WriteString("aggregate_type IN (")
		for i := range c.req.AggregateTypes {
			sb.WriteRune('?')
			if i != len(c.req.AggregateTypes)-1 {
				sb.WriteRune(',')
			}
		}
		sb.WriteString(") ")
	}

	if len(c.req.AggregateTypes) != 0 {
		sb.WriteString("aggregate_type IN (")
		for i, tp := range c.req.AggregateTypes {
			sb.WriteRune('?')
			q.args = append(q.args, tp)
			if i != len(c.req.AggregateTypes)-1 {
				sb.WriteRune(',')
			}
		}
		sb.WriteString(") ")
	}

	if len(c.req.EventTypes) != 0 {
		sb.WriteString("event_type IN (")
		for i, tp := range c.req.EventTypes {
			sb.WriteRune('?')
			q.args = append(q.args, tp)
			if i != len(c.req.EventTypes)-1 {
				sb.WriteRune(',')
			}
		}
		sb.WriteString(") ")
	}

	if len(c.req.ExcludeEventTypes) != 0 {
		sb.WriteString("event_type NOT IN (")
		for i, tp := range c.req.EventTypes {
			sb.WriteRune('?')
			q.args = append(q.args, tp)
			if i != len(c.req.EventTypes)-1 {
				sb.WriteRune(',')
			}
		}
		sb.WriteString(") ")
	}

	sb.WriteString("id > ? ")
	sb.WriteString("ORDER BY id ")
	sb.WriteString("LIMIT ?")

	q.query = c.conn.Rebind(sb.String())
	return &q
}
