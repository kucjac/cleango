package eventsource

import (
	"github.com/kucjac/cleango/codec"
)

// NewAggregateBaseSetter creates new aggregate setter.
func NewAggregateBaseSetter(eventCodec, snapCodec codec.Codec, idGen IdGenerator) *AggregateBaseSetter {
	return &AggregateBaseSetter{
		eventCodec: eventCodec,
		snapCodec:  snapCodec,
		idGen:      idGen,
	}
}

// AggregateBaseSetter is a structure responsible for setting the aggregate base.
type AggregateBaseSetter struct {
	eventCodec, snapCodec codec.Codec
	idGen                 IdGenerator
}

// SetAggregateBase implements AggregateBaseSetter interface.
func (a *AggregateBaseSetter) SetAggregateBase(agg Aggregate, aggId, aggType string, version int64) {
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
