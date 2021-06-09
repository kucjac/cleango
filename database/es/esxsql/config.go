package esxsql

import (
	"strings"

	"github.com/kucjac/cleango/cgerrors"
)

// Config is the configuration for the event storage.
type Config struct {
	EventTable     string
	SnapshotTable  string
	SchemaName     string // Optional
	AggregateTable string
	WorkersCount   int
}

// DefaultConfig creates a new default config.
func DefaultConfig() *Config {
	return &Config{
		EventTable:     "event",
		SnapshotTable:  "snapshot",
		AggregateTable: "aggregate",
		WorkersCount:   10,
	}
}

// Validate checks if the config is valid to use.
func (c *Config) Validate() error {
	if c.EventTable == "" {
		return cgerrors.ErrInternal("no event table name provided")
	}
	if c.SnapshotTable == "" {
		return cgerrors.ErrInternal("no snapshot table name provided")
	}
	if c.AggregateTable == "" {
		return cgerrors.ErrInternalf("no aggregate table name provided")
	}
	return nil
}

func (c *Config) eventTableName() string {
	sb := strings.Builder{}
	if c.SchemaName != "" {
		sb.WriteString(c.SchemaName)
		sb.WriteRune('.')
	}
	sb.WriteString(c.EventTable)
	return sb.String()
}

func (c *Config) snapshotTableName() string {
	sb := strings.Builder{}
	if c.SchemaName != "" {
		sb.WriteString(c.SchemaName)
		sb.WriteRune('.')
	}
	sb.WriteString(c.SnapshotTable)
	return sb.String()
}

func (c *Config) aggregateTableName() string {
	sb := strings.Builder{}
	if c.SchemaName != "" {
		sb.WriteString(c.SchemaName)
		sb.WriteRune('.')
	}
	sb.WriteString(c.AggregateTable)
	return sb.String()
}
