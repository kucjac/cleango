package eventsource

import (
	"sync"

	"github.com/kucjac/cleango/codec"
	"github.com/kucjac/cleango/xlog"
)

type loader struct {
	aggType    string
	aggVersion int64
	c          Cursor
	snapCodec  codec.Codec
	isClosed   bool
	numWorkers int
}

type aggregateLoader struct {
	loader
	factory AggregateFactory
	channel chan Aggregate
}

func newAggregateLoader(c Cursor, aggType string, aggVersion int64, factory AggregateFactory, size int) *aggregateLoader {
	return &aggregateLoader{
		loader: loader{
			c:          c,
			aggVersion: aggVersion,
			aggType:    aggType,
		},
		factory: factory,
		channel: make(chan Aggregate, size),
	}
}

func (l *aggregateLoader) readAggregateChannel() (<-chan Aggregate, error) {
	r, err := l.c.GetAggregateStream(true)
	if err != nil {
		return nil, err
	}
	go l.disposeWorkers(r, l.readAggregate, l.closeChannel)

	return l.channel, nil
}

func (l *aggregateLoader) readAggregate(r <-chan *CursorAggregate, wg *sync.WaitGroup) {
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

func (l *aggregateLoader) closeChannel() {
	close(l.channel)
}

// Scan scans then next aggregate returned by the loader and applies all events it
func (l *aggregateLoader) unmarshalAggregate(ca *CursorAggregate) (Aggregate, error) {
	agg := l.factory.New(l.aggType, l.aggVersion)
	b := agg.AggBase()
	b.SetID(ca.AggregateID)

	if ca.Snapshot != nil {
		if err := l.snapCodec.Unmarshal(ca.Snapshot.SnapshotData, agg); err != nil {
			return nil, err
		}
	}
	for _, e := range ca.Events {
		if err := agg.Apply(e); err != nil {
			return nil, err
		}
	}
	return agg, nil
}

func (l *loader) disposeWorkers(r <-chan *CursorAggregate, fn func(<-chan *CursorAggregate, *sync.WaitGroup), after func()) {
	wg := &sync.WaitGroup{}
	wg.Add(l.numWorkers)
	for i := 1; i < l.numWorkers; i++ {
		go fn(r, wg)
	}
	wg.Wait()

	after()
}
