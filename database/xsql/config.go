package xsql

import (
	"time"
)

// Config is the configuration for the database connection.
type Config struct {
	// LongQueriesTime is the duration at which the query is marked to be long-running.
	LongQueriesTime time.Duration
	// WarnLongQueries is a flag which set to true warns on long-running queries.
	WarnLongQueries bool
}
