package es

import (
	"sync"

	"github.com/kucjac/cleango/codec"
	"github.com/kucjac/cleango/xlog"
)

// Projection is an interface used to represent the query projeciton.
type Projection interface {
	Apply(c codec.Codec, e *Event) error
}

// ProjectionFactory is an interface used to create new projections.
type ProjectionFactory interface {
	NewProjection(id string) Projection
}

type projectionLoader struct {
	loader
	factory    ProjectionFactory
	channel    chan Projection
	eventCodec codec.Codec
}

func newProjectionLoader(cd codec.Codec, c Cursor, aggType string, aggVersion int64, factory ProjectionFactory, size int) *projectionLoader {
	return &projectionLoader{
		loader: loader{
			c:          c,
			aggVersion: aggVersion,
			aggType:    aggType,
		},
		factory:    factory,
		channel:    make(chan Projection, size),
		eventCodec: cd,
	}
}

func (l *projectionLoader) readProjectionChannel() (<-chan Projection, error) {
	r, err := l.c.GetAggregateStream(false)
	if err != nil {
		return nil, err
	}
	go l.disposeWorkers(r, l.readProjections, l.closeChannel)

	return l.channel, nil
}

func (l *projectionLoader) readProjections(r <-chan *CursorAggregate, wg *sync.WaitGroup) {
	for ca := range r {
		agg, err := l.unmarshalAggregate(ca)
		if err != nil {
			xlog.Errorf("unmarshalling aggregate failed: %v", err)
			continue
		}
		l.channel <- agg
	}
	wg.Done()
}

func (l *projectionLoader) closeChannel() {
	close(l.channel)
}

// Scan scans then next aggregate returned by the loader and applies all events it
func (l *projectionLoader) unmarshalAggregate(ca *CursorAggregate) (Projection, error) {
	p := l.factory.NewProjection(ca.AggregateID)
	for _, e := range ca.Events {
		if err := p.Apply(l.eventCodec, e); err != nil {
			return nil, err
		}
	}
	return p, nil
}
