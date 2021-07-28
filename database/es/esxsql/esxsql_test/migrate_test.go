package esxsql_tst

import (
	"testing"

	"github.com/kucjac/cleango/database/es/esxsql"
	"github.com/kucjac/cleango/ddd/events/eventstate"
)

func TestMigrate(t *testing.T) {
	t.Run("Postgres", func(t *testing.T) {
		t.Run("NoPartitions", func(t *testing.T) {
			conn := testPostgresConn(t)
			tx, err := conn.Begin()
			if err != nil {
				t.Fatalf("new transaction failed: %v", err)
			}
			defer tx.Rollback()
			config := esxsql.DefaultConfig("Orders", "TestAggregateType").
				WithEventState(esxsql.DefaultEventStateConfig(
					eventstate.Handler{Name: "TestHandler1", EventTypes: []string{"event_type_1", "event_type_2"}},
					eventstate.Handler{Name: "TestHandler2", EventTypes: []string{eventType}},
				))

			err = esxsql.Migrate(tx, config)
			if err != nil {
				t.Fatalf("migrating failed: %v", err)
			}
		})
		t.Run("Partitioned", func(t *testing.T) {
			conn := testPostgresConn(t)
			tx, err := conn.Begin()
			if err != nil {
				t.Fatalf("new transaction failed: %v", err)
			}
			defer tx.Rollback()
			config := esxsql.DefaultConfig("Orders", "TestAggregateType").
				WithEventState(esxsql.DefaultEventStateConfig(
					eventstate.Handler{Name: "TestHandler1", EventTypes: []string{"event_type_1", "event_type_2"}},
					eventstate.Handler{Name: "TestHandler2", EventTypes: []string{eventType}},
				))

			config.PartitionEventTable = true
			config.EventState.PartitionState = true

			err = esxsql.Migrate(tx, config)
			if err != nil {
				t.Fatalf("migrating failed: %v", err)
			}
		})
	})

}
