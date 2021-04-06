package eventsource

import (
	"github.com/kucjac/cleango/messages/codec"
)

// AggregateBaseSetter is an interface that allows to set the aggregate base in the given aggregate.
type AggregateBaseSetter interface {
	SetAggregateBase(agg Aggregate, aggId, aggType string, version int64)
}

// NewAggregateBaseSetter creates new aggregate setter.
func NewAggregateBaseSetter(eventCodec, snapCodec codec.Codec, idGen IdGenerator) AggregateBaseSetter {
	return &aggBaseSetter{
		eventCodec: eventCodec,
		snapCodec:  snapCodec,
		idGen:      idGen,
	}
}

type aggBaseSetter struct {
	eventCodec, snapCodec codec.Codec
	idGen                 IdGenerator
}

// SetAggregateBase implements AggregateBaseSetter interface.
func (a *aggBaseSetter) SetAggregateBase(agg Aggregate, aggId, aggType string, version int64) {
	base := &AggregateBase{
		id:         aggId,
		aggType:    aggType,
		agg:        agg,
		eventCodec: a.eventCodec,
		snapCodec:  a.snapCodec,
		idGen:      a.idGen,
		version:    version,
	}
	agg.SetBase(base)
}
