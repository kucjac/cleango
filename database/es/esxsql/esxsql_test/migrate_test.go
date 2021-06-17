package esxsql_tst

import (
	"os"
	"testing"

	"github.com/kucjac/cleango/database/es/esxsql"
)

func TestMigrate(t *testing.T) {
	t.Run("Postgres", func(t *testing.T) {
		if os.Getenv("MIGRATE") != "1" {
			t.Skip("MIGRATE environment variable not set")
		}
		conn := testPostgresConn(t)
		config := esxsql.DefaultConfig()
		err := esxsql.Migrate(conn, config, "TestAggregateType", "AnotherAggregate")
		if err != nil {
			t.Fatalf("migrating failed: %v", err)
		}

		err = esxsql.MigratePartitions(conn, config, "Orders", "TestAggregateType")
		if err != nil {
			t.Fatalf("migrating failed: %v", err)
		}

	})
}
