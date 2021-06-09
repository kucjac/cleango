package es

import (
	"errors"

	"github.com/kucjac/cleango/codec"
)

// SnapshotCodec is the type wrapper over the codec.Codec used for wire injection.
type SnapshotCodec codec.Codec

// EventCodec is the type wrapper over the codec.Codec used for event encoding in wire injection.
type EventCodec codec.Codec

// Config is the configuration for the eventsource storage.
type Config struct {
	BufferSize int
}

// DefaultConfig sets up the default config for the event store.
func DefaultConfig() *Config {
	return &Config{
		BufferSize: 100,
	}
}

func (c *Config) Validate() error {
	if c == nil {
		return nil
	}
	if c.BufferSize < 0 {
		return errors.New("event store buffer size is lower than 0")
	}
	return nil
}
