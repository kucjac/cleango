package esxsql

import (
	"fmt"
	"strings"

	"github.com/kucjac/cleango/database/xsql"
)

const (
	insertEventQuery           = `INSERT INTO %s (aggregate_id, aggregate_type, revision, timestamp, event_id, event_type, event_data) VALUES (?,?,?,?,?,?,?)`
	getEventStreamQuery        = `SELECT aggregate_id, aggregate_type, revision, timestamp, event_id, event_type, event_data FROM %s WHERE aggregate_id = ? AND aggregate_type = ? ORDER BY timestamp`
	saveSnapshotQuery          = `INSERT INTO %s (aggregate_id, aggregate_type, aggregate_version, revision, timestamp, snapshot_data) VALUES (?,?,?,?,?,?)`
	getSnapshotQuery           = `SELECT aggregate_id, aggregate_type, aggregate_version, revision, timestamp, snapshot_data FROM %s WHERE aggregate_id = ? AND aggregate_type = ? AND aggregate_version = ? ORDER BY timestamp DESC LIMIT 1`
	getStreamFromRevisionQuery = `SELECT aggregate_id, aggregate_type, revision, timestamp, event_id, event_type, event_data FROM %s WHERE aggregate_id = ? AND aggregate_type = ? AND revision > ? ORDER BY timestamp`
	batchInsertQueryBase       = `INSERT INTO %s (aggregate_id, aggregate_type, revision, timestamp, event_id, event_type, event_data) VALUES `
	insertAggregate            = `INSERT INTO %s (aggregate_id, aggregate_type, inserted_at) VALUES (?,?,?)`
	// listAggregates             = `SELECT id, aggregate_id FROM %s WHERE aggregate_type = ? LIMIT ? ORDER BY id`
	listNextAggregates   = `SELECT id, aggregate_id FROM %s WHERE aggregate_type = ? AND id > ? LIMIT ? ORDER BY id`
	listEventStreamQuery = `SELECT id, aggregate_id, aggregate_type, revision, timestamp, event_id, event_type, event_data FROM %s `
	registerHandler      = `INSERT INTO %s (handler_name, event_type) VALUES (?,?)`
	listHandlers         = `
SELECT handler_name, array_agg(event_type) AS event_types 
FROM %s 
WHERE event_type = ?
GROUP BY handler_name`
	updateEventState = `UPDATE %s SET state = ?, timestamp = ? WHERE event_id = ? AND handler_name = ?`
	insertEventState = `INSERT INTO %s (event_id, state, handler_name, timestamp) 
SELECT ?,?,handler_name,?
FROM %s AS h
WHERE h.event_type = ?`
	insertHandlingFailure = `INSERT INTO %s (event_id, handler_name, timestamp, error_message, error_code, retry_no) VALUES (?,?,?,?,?,?)`
	findHandlerEvents     = `SELECT es.event_id, es.handler_name FROM %s AS es`
	findHandlingFailures  = `SELECT ef.event_id, ef.handler_name, ef.timestamp, ef.error_message, ef.error_code, ef.retry_no 
FROM %s AS ef`
)

type queries struct {
	batchInsertQueryBase   string
	getEventStream         string
	saveSnapshot           string
	getSnapshot            string
	getStreamAfterRevision string
	insertEvent            string
	insertAggregate        string
	listNextAggregates     string
	listEventStreamQuery   string
	registerHandler        string
	listHandlers           string
	updateEventState       string
	insertEventState       string
	insertHandlingFailure  string
	findHandlerEvents      string
	findHandlingFailures   string
}

func (q queries) batchInsertEvent(length int) string {
	sb := strings.Builder{}
	sb.WriteString(q.batchInsertQueryBase)
	for i := 0; i < length; i++ {
		sb.WriteString("(?,?,?,?,?,?,?)")
		if i != length-1 {
			sb.WriteRune(',')
		}
	}
	return sb.String()
}

func newQueries(conn xsql.DB, c *Config) queries {
	return queries{
		batchInsertQueryBase:   fmt.Sprintf(batchInsertQueryBase, c.eventTableName()),
		getEventStream:         conn.Rebind(fmt.Sprintf(getEventStreamQuery, c.eventTableName())),
		saveSnapshot:           conn.Rebind(fmt.Sprintf(saveSnapshotQuery, c.snapshotTableName())),
		getSnapshot:            conn.Rebind(fmt.Sprintf(getSnapshotQuery, c.snapshotTableName())),
		getStreamAfterRevision: conn.Rebind(fmt.Sprintf(getStreamFromRevisionQuery, c.eventTableName())),
		insertEvent:            conn.Rebind(fmt.Sprintf(insertEventQuery, c.eventTableName())),
		insertAggregate:        conn.Rebind(fmt.Sprintf(insertAggregate, c.aggregateTableName())),
		listNextAggregates:     conn.Rebind(fmt.Sprintf(listNextAggregates, c.aggregateTableName())),
		listEventStreamQuery:   conn.Rebind(fmt.Sprintf(listEventStreamQuery, c.eventTableName())),
		registerHandler:        conn.Rebind(fmt.Sprintf(registerHandler, c.handlerTableName())),
		listHandlers:           conn.Rebind(fmt.Sprintf(listHandlers, c.handlerTableName())),
		updateEventState:       conn.Rebind(fmt.Sprintf(updateEventState, c.eventStateTableName())),
		insertEventState:       conn.Rebind(fmt.Sprintf(insertEventState, c.eventStateTableName(), c.handlerTableName())),
		insertHandlingFailure:  conn.Rebind(fmt.Sprintf(insertHandlingFailure, c.eventHandleFailureTableName())),
		findHandlerEvents:      conn.Rebind(fmt.Sprintf(findHandlerEvents, c.eventStateTableName())),
		findHandlingFailures:   conn.Rebind(fmt.Sprintf(findHandlingFailures, c.eventHandleFailureTableName())),
	}
}
