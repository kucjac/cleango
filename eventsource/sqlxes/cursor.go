package sqlxes

import (
	"context"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/xlog"

	"github.com/kucjac/cleango/eventsource"
	"github.com/kucjac/cleango/xservice"
)

var _ eventsource.Cursor = (*cursor)(nil)

type cursor struct {
	ctx          context.Context
	cancelFunc   context.CancelFunc
	l            sync.Mutex
	conn         *sqlx.DB
	s            *sqlStorage
	driver       xservice.Driver
	aggType      string
	aggVersion   int64
	lastTakenID  uint64
	query        queries
	limit        int64
	workersCount int
	workers      chan struct{}
}

func (s *sqlStorage) newCursor(ctx context.Context, aggType string, aggVersion int64) eventsource.Cursor {
	ctx, cancelFunc := context.WithCancel(ctx)
	return &cursor{
		ctx:          ctx,
		cancelFunc:   cancelFunc,
		aggType:      aggType,
		aggVersion:   aggVersion,
		conn:         s.conn,
		query:        s.query,
		driver:       s.d,
		limit:        1000,
		s:            s,
		workersCount: s.cfg.WorkersCount,
		workers:      make(chan struct{}, s.cfg.WorkersCount),
	}
}

func (c *cursor) OpenChannel(withSnapshots bool) (<-chan *eventsource.CursorAggregate, error) {
	ch := make(chan *eventsource.CursorAggregate, c.limit)
	go c.readAggregates(ch, withSnapshots)
	return ch, nil
}

func (c *cursor) readAggregates(ca chan *eventsource.CursorAggregate, withSnapshots bool) {
	var err error
	errChan := make(chan error, 1)

	c.initializeWorkers()

	for {
		select {
		case e := <-errChan:
			err = e
			break
		case <-c.ctx.Done():
			err = c.ctx.Err()
			break
		default:
		}
		var rows *sqlx.Rows
		rows, err = c.conn.QueryxContext(c.ctx, c.query.listNextAggregates, c.aggType, c.lastTakenID, c.limit)
		if err != nil {
			if c.driver.CanRetry(err) {
				continue
			}
			break
		}

		var rowsCount int

	rowsLoop:
		for rows.Next() {
			var (
				id          uint64
				aggregateId string
			)
			if err = rows.Scan(&id, &aggregateId); err != nil {
				break
			}
			c.lastTakenID = id

			// Either error channel context or workers are finished.
			select {
			case e := <-errChan:
				err = e
				break
			case <-c.ctx.Done():
				err = c.ctx.Err()
				break
			case _, ok := <-c.workers:
				// The channel of workers is already closed.
				if !ok {
					break rowsLoop
				}
			}
			go c.readAggregate(aggregateId, ca, errChan, withSnapshots)
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

	// If an error occurred log it's content, no matter what we need to close the channels.
	if err != nil {
		xlog.Errorf("reading aggregates failed: %v", err)
	}

	// Release all the workers so that the worker channel is empty.
	c.releaseWorkers()

	// Now we're sure there is no worker left that would put anything in the channel.
	close(ca)
	close(c.workers)
}

func (c *cursor) releaseWorkers() {
	c.cancelFunc()
	// Release all workers.
	for i := 0; i < c.workersCount; i++ {
		<-c.workers
	}
}

func (c *cursor) initializeWorkers() {
	// Initialise workers.
	for i := 0; i < c.workersCount; i++ {
		c.workers <- struct{}{}
	}
}

func (c *cursor) readAggregate(aggregateId string, ac chan<- *eventsource.CursorAggregate, ec chan<- error, withSnapshots bool) {
	defer func() {
		// Return the worker to the pool of workers.
		c.workers <- struct{}{}
	}()

	var err error
	agg := &eventsource.CursorAggregate{AggregateID: aggregateId}
	if withSnapshots {
		agg.Snapshot, err = c.s.GetSnapshot(c.ctx, aggregateId, c.aggType, c.aggVersion)
		if err != nil {
			if !cgerrors.IsNotFound(err) {
				ec <- err
				return
			}
		}
	}

	var revision int64
	if agg.Snapshot != nil {
		revision = agg.Snapshot.Revision
	}

	agg.Events, err = c.s.GetStreamFromRevision(c.ctx, aggregateId, c.aggType, revision)
	if err != nil {
		ec <- err
		return
	}

	ac <- agg
}
