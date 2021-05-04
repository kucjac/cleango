package eventsource

import (
	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/codec"
)

type Config struct {
	BufferSize    int
	EventCodec    codec.Codec
	SnapshotCodec codec.Codec
}

// DefaultConfig sets up the default config for the event store.
func DefaultConfig() *Config {
	return &Config{
		BufferSize:    100,
		EventCodec:    codec.Proto(),
		SnapshotCodec: codec.Proto(),
	}
}

func (c *Config) Validate() error {
	if c.EventCodec == nil {
		return cgerrors.ErrInternal("event codec not defined")
	}
	if c.SnapshotCodec == nil {
		return cgerrors.ErrInternal("snapshot codec not defined")
	}
	return nil
}
