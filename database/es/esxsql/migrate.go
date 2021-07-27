package esxsql

import (
	"bytes"
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/database/es/eventstate"
	"github.com/kucjac/cleango/database/xsql"

	"github.com/kucjac/cleango/xlog"
)

//go:embed mysql.tmpl
var mysqlMigrateQuery string

var migrateMySQL *template.Template

func init() {
	migrateMySQL = template.Must(template.New("").Parse(mysqlMigrateQuery))
}

// Migrate executes table and types migration for the event store and snapshot.
// The table names are taken from the config.
func Migrate(conn xsql.DB, config *Config) error {
	var buf bytes.Buffer
	if err := config.Validate(); err != nil {
		return err
	}

	switch conn.DriverName() {
	case "pg", "postgres", "postgresql", "gopg", "pgx":
		if err := migratePostgresTables(context.Background(), conn, config); err != nil {
			return err
		}
	case "mysql":
		xlog.Infoln("Migrating esxsql with mysql driver")
		if err := migrateMySQL.Execute(&buf, config); err != nil {
			return err
		}
	default:
		return errors.New("driver not supported by the esxsql migration tool")
	}

	return nil
}

// MigrateEventPartitions migrates event partitions
func MigrateEventPartitions(conn xsql.DB, cfg *Config, aggregateTypes ...string) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	if !cfg.PartitionEventTable {
		return cgerrors.ErrInternal("partitioning of event table is not set in configuration")
	}

	switch conn.DriverName() {
	case "pg", "postgres", "postgresql", "gopg", "pgx":
		err := migratePostgresEventPartitions(context.Background(), conn, cfg, aggregateTypes...)
		if err != nil {
			return err
		}
	case "mysql":
		return errors.New("partitions not implemented for the mysql yet")
	default:
		return errors.New("driver not supported by the esxsql migration tool")
	}

	xlog.Infoln("Migrating esxsql partitions with postgres driver")
	return nil
}

// MigrateEventStatePartitions create partitions on event state by its type.
func MigrateEventStatePartitions(conn xsql.DB, cfg *Config, handlerNames ...string) error {
	xlog.Infof("Migrating esxsql event state partitions - handlers: (%s)", strings.Join(handlerNames, ","))

	switch conn.DriverName() {
	case "postgres":
		// Search for all partitions of the event state table.
		if err := migratePostgresEventStatePartitions(context.Background(), conn, cfg, handlerNames); err != nil {
			return err
		}
	default:
		return cgerrors.ErrUnimplemented("migration for driver not implemented").WithMeta("driverName", conn.DriverName())
	}

	return nil
}

func migratePostgresTables(ctx context.Context, conn xsql.DB, cfg *Config) error {
	// Migrate event table.
	err := migratePostgresEventTable(ctx, conn, cfg)
	if err != nil {
		return err
	}

	// Migrate snapshot table
	if err = migratePostgresSnapshotTable(ctx, conn, cfg); err != nil {
		return err
	}

	if err = migratePostgresAggregateTables(ctx, conn, cfg); err != nil {
		return err
	}

	// If the eventstate config is undefined, no tables should be migrated for the eventstate.
	if cfg.EventState == nil {
		return nil
	}

	if err = migratePostgresHandlersTables(ctx, conn, cfg); err != nil {
		return err
	}

	if err = migratePostgresEventStateTable(ctx, conn, cfg); err != nil {
		return err
	}

	if err = migratePostgresEventHandleFailureTable(ctx, conn, cfg); err != nil {
		return err
	}
	return nil
}

func migratePostgresSnapshotTable(ctx context.Context, conn xsql.DB, cfg *Config) error {
	schema := cfg.SchemaName
	if schema == "" {
		schema = "public"
	}

	exists, err := postgresTableExists(ctx, conn, schema, cfg.SnapshotTable)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	var sb strings.Builder
	sb.WriteString("CREATE TABLE ")
	sb.WriteString(schema)
	sb.WriteString(".")
	sb.WriteString(cfg.SnapshotTable)
	sb.WriteString(" (\n")
	sb.WriteString("\tid BIGSERIAL NOT NULL PRIMARY KEY,\n")
	sb.WriteString("\taggregate_id TEXT NOT NULL,\n")
	sb.WriteString("\taggregate_type TEXT NOT NULL,\n")
	sb.WriteString("\taggregate_version integer NOT NULL,\n")
	sb.WriteString("\trevision integer NOT NULL,\n")
	sb.WriteString("\ttimestamp bigint NOT NULL,\n")
	sb.WriteString("\tsnapshot_data bytea,\n")
	sb.WriteString("\tCONSTRAINT ")
	sb.WriteString(cfg.SnapshotTable)
	sb.WriteString("_aggregate_revision_uidx UNIQUE(aggregate_id, aggregate_type, revision)\n")
	sb.WriteString(")")

	_, err = conn.ExecContext(ctx, sb.String())
	return err
}

func migratePostgresAggregateTables(ctx context.Context, conn xsql.DB, cfg *Config) error {
	schema := cfg.SchemaName
	if schema == "" {
		schema = "public"
	}

	exists, err := postgresTableExists(ctx, conn, schema, cfg.AggregateTable)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	var sb strings.Builder
	sb.WriteString("CREATE TABLE ")
	sb.WriteString(schema)
	sb.WriteString(".")
	sb.WriteString(cfg.AggregateTable)
	sb.WriteString(" (\n")
	sb.WriteString("\tid BIGSERIAL NOT NULL PRIMARY KEY,\n")
	sb.WriteString("\taggregate_id TEXT NOT NULL,\n")
	sb.WriteString("\taggregate_type TEXT NOT NULL,\n")
	sb.WriteString("\tinserted_at bigint NOT NULL,\n")
	sb.WriteString("\tCONSTRAINT ")
	sb.WriteString(cfg.AggregateTable)
	sb.WriteString("_aggregate_id_aggregate_type_uidx UNIQUE(aggregate_id, aggregate_type)\n")
	sb.WriteString(")")

	_, err = conn.ExecContext(ctx, sb.String())
	return err
}

func migratePostgresHandlersTables(ctx context.Context, conn xsql.DB, cfg *Config) error {
	schema := cfg.SchemaName
	if schema == "" {
		schema = "public"
	}

	exists, err := postgresTableExists(ctx, conn, schema, cfg.EventState.HandlerTable)
	if err != nil {
		return err
	}
	if !exists {
		var sb strings.Builder
		sb.WriteString("CREATE TABLE ")
		sb.WriteString(schema)
		sb.WriteString(".")
		sb.WriteString(cfg.EventState.HandlerTable)
		sb.WriteString(" (\n")
		sb.WriteString("\tid BIGSERIAL NOT NULL PRIMARY KEY,\n")
		sb.WriteString("\thandler_name TEXT NOT NULL,\n")
		sb.WriteString("\tevent_type TEXT NOT NULL,\n")
		sb.WriteString("\tCONSTRAINT ")
		sb.WriteString(cfg.EventState.HandlerTable)
		sb.WriteString("_handler_name_event_type_uidx UNIQUE(handler_name, event_type)\n")
		sb.WriteString(")")

		_, err = conn.ExecContext(ctx, sb.String())
		if err != nil {
			return err
		}
	}

	idxName := fmt.Sprintf("%s_event_type_idx", cfg.EventState.HandlerTable)
	exists, err = postgresIndexExists(ctx, conn, schema, cfg.EventState.HandlerTable, idxName)
	if err != nil {
		return err
	}
	if !exists {
		_, err = conn.ExecContext(ctx, fmt.Sprintf("CREATE INDEX %s ON %s.%s (event_type)", idxName, schema, cfg.EventState.HandlerTable))
		if err != nil {
			return err
		}
	}
	handlers := cfg.EventState.Handlers

	err = insertHandlers(ctx, conn, cfg.handlerTableName(), handlers)
	if err != nil {
		return err
	}

	return nil
}

func insertHandlers(ctx context.Context, conn xsql.DB, handlersTable string, handlers []eventstate.Handler) (err error) {
	const q = `INSERT INTO %s (handler_name, event_type) VALUES ($1, $2);`
	for _, handler := range handlers {
		for _, eh := range handler.EventTypes {
			_, err = conn.ExecContext(ctx, fmt.Sprintf(q, handlersTable), handler.Name, eh)
			if err != nil {
				if cgerrors.IsAlreadyExists(err) {
					continue
				}
				return err
			}
		}
	}
	return nil
}

func migratePostgresEventStateTable(ctx context.Context, conn xsql.DB, cfg *Config) error {
	schema := cfg.SchemaName
	if schema == "" {
		schema = "public"
	}
	// language=PostgreSQL
	exists, err := postgresTableExists(ctx, conn, schema, cfg.EventState.EventStateTable)
	if err != nil {
		return err
	}

	sb := strings.Builder{}
	if !exists {
		sb.WriteString("CREATE TABLE ")
		sb.WriteString(schema)
		sb.WriteRune('.')
		sb.WriteString(cfg.EventState.EventStateTable)
		sb.WriteString(" (\n")
		sb.WriteString("\tid BIGSERIAL NOT NULL")
		if !cfg.EventState.PartitionState {
			sb.WriteString(" PRIMARY KEY")
		}
		sb.WriteString(",\n")
		sb.WriteString("\tevent_id TEXT NOT NULL,\n")
		sb.WriteString("\thandler_name TEXT NOT NULL,\n")
		sb.WriteString("\tstate smallint NOT NULL,\n")
		sb.WriteString("\ttimestamp bigint NOT NULL,\n")
		sb.WriteString("\tCONSTRAINT ")
		sb.WriteString(cfg.EventState.EventStateTable)
		sb.WriteString("_event_id_handler_name_uidx UNIQUE (event_id, handler_name)\n")
		if !cfg.EventState.PartitionState {
			sb.WriteString(")")
		} else {
			sb.WriteString(") PARTITION BY LIST(handler_name);")
		}

		q := sb.String()
		if _, err := conn.ExecContext(ctx, q); err != nil {
			return err
		}
		sb.Reset()
	}

	if cfg.EventState.PartitionState && len(cfg.EventState.Handlers) > 0 {
		handlerNames := make([]string, len(cfg.EventState.Handlers))
		for i, h := range cfg.EventState.Handlers {
			handlerNames[i] = h.Name
		}
		err = migratePostgresEventStatePartitions(ctx, conn, cfg, handlerNames)
		if err != nil {
			return err
		}
	}
	return nil
}

func migratePostgresEventHandleFailureTable(ctx context.Context, conn xsql.DB, cfg *Config) error {
	schema := cfg.SchemaName
	if schema == "" {
		schema = "public"
	}

	exists, err := postgresTableExists(ctx, conn, schema, cfg.EventState.HandleFailureTable)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	var sb strings.Builder
	sb.WriteString("CREATE TABLE ")
	sb.WriteString(schema)
	sb.WriteString(".")
	sb.WriteString(cfg.EventState.HandleFailureTable)
	sb.WriteString(" (\n")
	sb.WriteString("\tid BIGSERIAL NOT NULL PRIMARY KEY,\n")
	sb.WriteString("\tevent_id TEXT NOT NULL,\n")
	sb.WriteString("\thandler_name TEXT NOT NULL,\n")
	sb.WriteString("\ttimestamp bigint NOT NULL,\n")
	sb.WriteString("\terror_message TEXT NOT NULL,\n")
	sb.WriteString("\terror_code smallint NOT NULL,\n")
	sb.WriteString("\tretry_no smallint NOT NULL\n")
	sb.WriteString(")")

	_, err = conn.ExecContext(ctx, sb.String())
	return err
}

func migratePostgresEventTable(ctx context.Context, conn xsql.DB, cfg *Config) error {
	schema := cfg.SchemaName
	if schema == "" {
		schema = "public"
	}
	// language=PostgreSQL
	exists, err := postgresTableExists(ctx, conn, schema, cfg.EventTable)
	if err != nil {
		return err
	}

	sb := strings.Builder{}
	if !exists {
		sb.WriteString("CREATE TABLE ")
		sb.WriteString(schema)
		sb.WriteRune('.')
		sb.WriteString(cfg.EventTable)
		sb.WriteString(" (\n")
		sb.WriteString("\tid BIGSERIAL NOT NULL")
		if !cfg.PartitionEventTable {
			sb.WriteString(" PRIMARY KEY")
		}
		sb.WriteString(",\n")
		sb.WriteString("\tevent_id TEXT NOT NULL,\n")
		sb.WriteString("\taggregate_id TEXT NOT NULL,\n")
		sb.WriteString("\taggregate_type TEXT NOT NULL,\n")
		sb.WriteString("\trevision integer NOT NULL,\n")
		sb.WriteString("\ttimestamp bigint NOT NULL,\n")
		sb.WriteString("\tevent_type TEXT NOT NULL,\n")
		sb.WriteString("\tevent_data bytea,\n")
		sb.WriteString("\tCONSTRAINT ")
		sb.WriteString(cfg.EventTable)
		sb.WriteString("_aggregate_revision_uidx UNIQUE (aggregate_id, aggregate_type, revision)\n")
		if !cfg.PartitionEventTable {
			sb.WriteString(")")
		} else {
			sb.WriteString(") PARTITION BY LIST(aggregate_type);")
		}

		q := sb.String()
		if _, err := conn.ExecContext(ctx, q); err != nil {
			return err
		}
		sb.Reset()
	}

	// Check and create if not exists event_id index.
	sb.WriteString(cfg.EventTable)
	sb.WriteString("_event_id_")
	if !cfg.PartitionEventTable {
		sb.WriteRune('u')
	}
	sb.WriteString("idx")
	idxName := sb.String()
	sb.Reset()

	exists, err = postgresIndexExists(ctx, conn, schema, cfg.EventTable, idxName)
	if err != nil {
		return err
	}

	if !exists {
		sb.WriteString("CREATE ")
		if !cfg.PartitionEventTable {
			sb.WriteString("UNIQUE ")
		}
		sb.WriteString("INDEX ")
		sb.WriteString(idxName)
		sb.WriteString(" ON ")
		sb.WriteString(schema)
		sb.WriteRune('.')
		sb.WriteString(cfg.EventTable)
		sb.WriteString(" (event_id);")

		q := sb.String()
		sb.Reset()

		_, err = conn.ExecContext(ctx, q)
		if err != nil {
			return err
		}
	}

	// Check and create if not exists 'event_type' index.
	idxName = fmt.Sprintf("%s_event_type_idx", cfg.EventTable)
	exists, err = postgresIndexExists(ctx, conn, schema, cfg.EventTable, idxName)
	if err != nil {
		return err
	}
	if !exists {
		_, err = conn.ExecContext(ctx, fmt.Sprintf("CREATE INDEX %s ON %s.%s (event_type)", idxName, schema, cfg.EventTable))
		if err != nil {
			return err
		}
	}

	if cfg.PartitionEventTable && len(cfg.AggregateTypes) > 0 {
		err = migratePostgresEventPartitions(ctx, conn, cfg, cfg.AggregateTypes...)
		if err != nil {
			return err
		}
	}

	return nil
}

func postgresIndexExists(ctx context.Context, conn xsql.DB, schema string, eventTable string, idxName string) (bool, error) {
	// language=PostgreSQL
	q := `SELECT 1 FROM pg_indexes WHERE schemaname = $1 AND tablename = $2 and indexname = $3`
	row := conn.QueryRowContext(ctx, q, schema, eventTable, idxName)
	var exists int
	if err := row.Scan(&exists); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func postgresTableExists(ctx context.Context, conn xsql.DB, schema, table string) (bool, error) {
	q := `SELECT 1
	FROM information_schema.tables AS t
	WHERE t.table_schema = $1 AND t.table_name = $2`

	if schema == "" {
		schema = "public"
	}

	row := conn.QueryRowContext(ctx, q, schema, table)
	var exists int
	err := row.Scan(&exists)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func migratePostgresEventPartitions(ctx context.Context, conn xsql.DB, cfg *Config, aggregateTypes ...string) error {
	type partition struct {
		childSchema, child string
	}
	var partitions []partition
	// language=PostgreSQL
	const listPartitions = `SELECT        
    nmsp_child.nspname  AS child_schema,
    child.relname       AS child
FROM pg_inherits
         JOIN pg_class parent            ON pg_inherits.inhparent = parent.oid
         JOIN pg_class child             ON pg_inherits.inhrelid   = child.oid         
         JOIN pg_namespace nmsp_parent   ON nmsp_parent.oid  = parent.relnamespace
         JOIN pg_namespace nmsp_child    ON nmsp_child.oid   = child.relnamespace
WHERE parent.relname = $1 AND nmsp_parent.nspname = $2`

	schemaName := cfg.SchemaName
	if schemaName == "" {
		schemaName = "public"
	}

	rows, err := conn.QueryContext(ctx, listPartitions, cfg.EventTable, schemaName)
	if err != nil {
		return err
	}

	for rows.Next() {
		var p partition
		err = rows.Scan(&p.childSchema, &p.child)
		if err != nil {
			return err
		}
		partitions = append(partitions, p)
	}
	if err = rows.Err(); err != nil {
		xlog.Errorf("Listing event state partitions failed: %v", err)
		return err
	}

	mp := map[string]struct{}{}
	for _, p := range partitions {
		mp[p.child] = struct{}{}
	}

	r := strings.NewReplacer(":", "_", ".", "_", ",", "_")
	for _, aggregateName := range aggregateTypes {
		childTable := r.Replace(aggregateName)
		childTable = ToSnakeCase(childTable)
		childTable = cfg.EventTable + "_" + childTable

		if _, ok := mp[childTable]; ok {
			xlog.Debugf("Event partition table for aggregate: '%s' already exists", aggregateName)
			continue
		}
		queries := []string{
			conn.Rebind(fmt.Sprintf(`CREATE TABLE %s PARTITION OF %s FOR VALUES IN('%s')`, childTable, cfg.eventTableName(), aggregateName)),
			conn.Rebind(fmt.Sprintf("CREATE UNIQUE INDEX %s_event_id_uidx ON %s (event_id)", childTable, childTable)),
		}

		for _, q := range queries {
			_, err := conn.Exec(q)
			if err != nil {
				if conn.ErrorCode(err) == cgerrors.ErrorCode_AlreadyExists {
					xlog.Debugf("%v", err)
				} else {
					xlog.Errorf("Code: %d - %v", conn.ErrorCode(err), err)
					return err
				}
			}
		}
	}
	return nil
}

func migratePostgresEventStatePartitions(ctx context.Context, conn xsql.DB, cfg *Config, handlerNames []string) error {
	type partition struct {
		childSchema, child string
	}

	// Define all partitions for given parent.
	var partitions []partition

	// language=PostgreSQL
	const listPartitions = `SELECT        
    nmsp_child.nspname  AS child_schema,
    child.relname       AS child
FROM pg_inherits
         JOIN pg_class parent            ON pg_inherits.inhparent = parent.oid
         JOIN pg_class child             ON pg_inherits.inhrelid   = child.oid         
         JOIN pg_namespace nmsp_parent   ON nmsp_parent.oid  = parent.relnamespace
         JOIN pg_namespace nmsp_child    ON nmsp_child.oid   = child.relnamespace
WHERE parent.relname = ? AND nmsp_parent.nspname = ?`

	schemaName := cfg.SchemaName
	if schemaName == "" {
		schemaName = "public"
	}

	rows, err := conn.QueryContext(ctx, conn.Rebind(listPartitions), cfg.EventState.EventStateTable, schemaName)
	if err != nil {
		return err
	}

	for rows.Next() {
		var p partition
		err = rows.Scan(&p.childSchema, &p.child)
		if err != nil {
			return err
		}
		partitions = append(partitions, p)
	}
	if err = rows.Err(); err != nil {
		xlog.Errorf("Listing event state partitions failed: %v", err)
		return err
	}

	mp := map[string]struct{}{}
	for _, p := range partitions {
		mp[p.child] = struct{}{}
	}

	r := strings.NewReplacer(":", "_", ".", "_", ",", "_")
	for _, handlerName := range handlerNames {
		childTable := r.Replace(handlerName)
		childTable = ToSnakeCase(childTable)

		if _, ok := mp[childTable]; ok {
			xlog.Debugf("Event state partition table for handler: %s already exists", childTable)
			continue
		}

		queries := []string{
			conn.Rebind(fmt.Sprintf(`CREATE TABLE %s PARTITION OF %s FOR VALUES IN('%s')`, childTable, cfg.eventStateTableName(), handlerName)),
		}

		for _, q := range queries {
			_, err := conn.ExecContext(ctx, q)
			if err != nil {
				if conn.ErrorCode(err) == cgerrors.ErrorCode_AlreadyExists {
					xlog.Debugf("%v", err)
				} else {
					xlog.Errorf("Code: %d - %v", conn.ErrorCode(err), err)
					return err
				}
			}
		}
	}
	return nil
}
