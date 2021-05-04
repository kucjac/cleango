package eventsource

import (
	codec2 "github.com/kucjac/cleango/codec"
)

// newAggregateBaseSetter creates new aggregate setter.
func newAggregateBaseSetter(eventCodec, snapCodec codec2.Codec, idGen IdGenerator) *aggBaseSetter {
	return &aggBaseSetter{
		eventCodec: eventCodec,
		snapCodec:  snapCodec,
		idGen:      idGen,
	}
}

type aggBaseSetter struct {
	eventCodec, snapCodec codec2.Codec
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
