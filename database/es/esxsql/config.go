package esxsql

import (
	"strings"

	"github.com/kucjac/cleango/cgerrors"
)

// Config is the configuration for the event storage.
type Config struct {
	EventTable              string
	SnapshotTable           string
	SchemaName              string // Optional
	AggregateTable          string
	HandlerRegistryTable    string
	EventStateTable         string
	EventHandleFailureTable string
	WorkersCount            int
	useEventState           bool
}

// DefaultConfig creates a new default config.
func DefaultConfig() *Config {
	return &Config{
		EventTable:              "event",
		SnapshotTable:           "snapshot",
		AggregateTable:          "aggregate",
		HandlerRegistryTable:    "handler_registry",
		EventStateTable:         "event_state",
		EventHandleFailureTable: "event_handle_failure",
		WorkersCount:            10,
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
	if c.useEventState {
		if c.HandlerRegistryTable == "" {
			return cgerrors.ErrInternalf("no handler registry table name provided")
		}
		if c.EventStateTable == "" {
			return cgerrors.ErrInternalf("no event state table name provided")
		}
		if c.EventHandleFailureTable == "" {
			return cgerrors.ErrInternalf("no event handle failure table name provided")
		}
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

func (c *Config) eventHandleFailureTableName() string {
	sb := strings.Builder{}
	if c.SchemaName != "" {
		sb.WriteString(c.SchemaName)
		sb.WriteRune('.')
	}
	sb.WriteString(c.EventHandleFailureTable)
	return sb.String()
}

func (c *Config) eventStateTableName() string {
	sb := strings.Builder{}
	if c.SchemaName != "" {
		sb.WriteString(c.SchemaName)
		sb.WriteRune('.')
	}
	sb.WriteString(c.EventStateTable)
	return sb.String()
}

func (c *Config) handlerRegistryTableName() string {
	sb := strings.Builder{}
	if c.SchemaName != "" {
		sb.WriteString(c.SchemaName)
		sb.WriteRune('.')
	}
	sb.WriteString(c.HandlerRegistryTable)
	return sb.String()
}
