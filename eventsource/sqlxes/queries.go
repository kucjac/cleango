package sqlxes

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

const (
	insertEventQuery           = `INSERT INTO %s (aggregate_id, aggregate_type, revision, timestamp, event_id, event_type, event_data) VALUES (?,?,?,?,?,?,?)`
	getEventStreamQuery        = `SELECT aggregate_id, aggregate_type, revision, timestamp, event_id, event_type, event_data FROM %s WHERE aggregate_id = ? AND aggregate_type = ?`
	saveSnapshotQuery          = `INSERT INTO %s (aggregate_id, aggregate_type, aggregate_version, revision, timestamp, snapshot_data) VALUES (?,?,?,?,?,?)`
	getSnapshotQuery           = `SELECT aggregate_id, aggregate_type, aggregate_version, revision, timestamp, snapshot_data FROM %s WHERE aggregate_id = ? AND aggregate_type = ? AND aggregate_version = ? ORDER BY timestamp DESC LIMIT 1`
	getStreamFromRevisionQuery = `SELECT aggregate_id, aggregate_type, revision, timestamp, event_id, event_type, event_data FROM %s WHERE aggregate_id = ? AND aggregate_type = ? AND revision > ?`
	batchInsertQueryBase       = `INSERT INTO %s (aggregate_id, aggregate_type, revision, timestamp, event_id, event_type, event_data) VALUES `
	insertAggregate            = `INSERT INTO %s (aggregate_id, inserted_at) VALUES (?, ?)`
	// listAggregates             = `SELECT id, aggregate_id FROM %s WHERE aggregate_type = ? LIMIT ? ORDER BY id`
	listNextAggregates   = `SELECT id, aggregate_id FROM %s WHERE aggregate_type = ? AND id > ? LIMIT ? ORDER BY id`
	listEventStreamQuery = `SELECT id, aggregate_id, aggregate_type, revision, timestamp, event_id, event_type, event_data FROM %s `
)

type queries struct {
	getEventStream        string
	saveSnapshot          string
	getSnapshot           string
	getStreamFromRevision string
	insertEvent           string
	insertAggregate       string
	listNextAggregates    string
	listEventStreamQuery  string
}

func (q queries) batchInsertEvent(length int) string {
	sb := strings.Builder{}
	sb.WriteString(batchInsertQueryBase)
	for i := 0; i < length; i++ {
		sb.WriteString("(?,?,?,?,?,?,?)")
		sb.WriteRune(',')
	}
	return sb.String()
}

func newQueries(conn *sqlx.DB, c *Config) queries {
	return queries{
		getEventStream:        conn.Rebind(fmt.Sprintf(getEventStreamQuery, c.eventTableName())),
		saveSnapshot:          conn.Rebind(fmt.Sprintf(saveSnapshotQuery, c.snapshotTableName())),
		getSnapshot:           conn.Rebind(fmt.Sprintf(getSnapshotQuery, c.snapshotTableName())),
		getStreamFromRevision: conn.Rebind(fmt.Sprintf(getStreamFromRevisionQuery, c.eventTableName())),
		insertEvent:           conn.Rebind(fmt.Sprintf(insertEventQuery, c.eventTableName())),
		insertAggregate:       conn.Rebind(fmt.Sprintf(insertAggregate, c.aggregateTableName())),
		listNextAggregates:    conn.Rebind(fmt.Sprintf(listNextAggregates, c.aggregateTableName())),
		listEventStreamQuery:  conn.Rebind(fmt.Sprintf(listEventStreamQuery, c.eventTableName())),
	}
}
