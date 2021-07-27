package esxsql

import (
	"strings"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/database/es/eventstate"
)

// Config is the configuration for the event storage.
type Config struct {
	SchemaName          string // Optional
	EventTable          string
	PartitionEventTable bool
	SnapshotTable       string
	AggregateTable      string
	AggregateTypes      []string
	EventState          *EventStateConfig
	WorkersCount        int
}

// DefaultConfig creates a new default config.
func DefaultConfig(aggregateTypes ...string) *Config {
	return &Config{
		EventTable:     "event",
		SnapshotTable:  "snapshot",
		AggregateTable: "aggregate",
		WorkersCount:   10,
		AggregateTypes: aggregateTypes,
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
	// Validate event state inputs.
	if c.EventState != nil {
		if err := c.EventState.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// WithEventState sets up the event state for given config.
func (c *Config) WithEventState(config EventStateConfig) *Config {
	c.EventState = &config
	return c
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
	if c.EventState == nil {
		return ""
	}
	sb := strings.Builder{}
	if c.SchemaName != "" {
		sb.WriteString(c.SchemaName)
		sb.WriteRune('.')
	}
	sb.WriteString(c.EventState.HandleFailureTable)
	return sb.String()
}

func (c *Config) eventStateTableName() string {
	if c.EventState == nil {
		return ""
	}
	sb := strings.Builder{}
	if c.SchemaName != "" {
		sb.WriteString(c.SchemaName)
		sb.WriteRune('.')
	}
	sb.WriteString(c.EventState.EventStateTable)
	return sb.String()
}

func (c *Config) handlerTableName() string {
	if c.EventState == nil {
		return ""
	}
	sb := strings.Builder{}
	if c.SchemaName != "" {
		sb.WriteString(c.SchemaName)
		sb.WriteRune('.')
	}
	sb.WriteString(c.EventState.HandlerTable)
	return sb.String()
}

// EventStateConfig is a configuration for the event state part.
type EventStateConfig struct {
	EventStateTable    string
	PartitionState     bool
	Handlers           []eventstate.Handler
	HandleFailureTable string
	HandlerTable       string
}

// DefaultEventStateConfig is a default configuration for the event state tracking.
func DefaultEventStateConfig(handlers ...eventstate.Handler) EventStateConfig {
	return EventStateConfig{
		EventStateTable:    "event_state",
		PartitionState:     false,
		Handlers:           handlers,
		HandleFailureTable: "event_handle_failure",
		HandlerTable:       "event_handler",
	}
}

func (c *EventStateConfig) Validate() error {
	if c.HandlerTable == "" {
		return cgerrors.ErrInternalf("no handler registry table name provided")
	}
	if c.EventStateTable == "" {
		return cgerrors.ErrInternalf("no event state table name provided")
	}
	if c.HandleFailureTable == "" {
		return cgerrors.ErrInternalf("no event handle failure table name provided")
	}
	return nil
}
