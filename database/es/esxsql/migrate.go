package esxsql

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/database/xsql"

	"github.com/kucjac/cleango/xlog"
)

//go:embed mysql.tmpl
var mysqlMigrateQuery string

//go:embed postgres.tmpl
var pgMigrateQuery string

//go:embed postgres_partitions.tmpl
var pgPartitionsQuery string

var (
	migrateMySQL        *template.Template
	migratePostgres     *template.Template
	migratePgPartitions *template.Template
)

func init() {
	migrateMySQL = template.Must(template.New("").Parse(mysqlMigrateQuery))
	migratePostgres = template.Must(template.New("").Parse(pgMigrateQuery))
	migratePgPartitions = template.Must(template.New("").Parse(pgPartitionsQuery))
}

// Migrate executes table and types migration for the event store and snapshot.
// The table names are taken from the config.
func Migrate(conn xsql.DB, config *Config, aggregateTypes ...string) error {
	var buf bytes.Buffer
	if err := config.Validate(); err != nil {
		return err
	}

	switch conn.DriverName() {
	case "pg", "postgres", "postgresql", "gopg", "pgx":
		type aggregate struct {
			Type  string
			Value string
		}
		type postgresInput struct {
			Config
			Aggregates []aggregate
		}
		cfg := *config
		if cfg.SchemaName != "" {
			cfg.SchemaName += "."
		}
		pi := postgresInput{Config: cfg}
		for _, agg := range aggregateTypes {
			pi.Aggregates = append(pi.Aggregates, aggregate{Type: ToSnakeCase(agg), Value: agg})
		}
		xlog.Infoln("Migrating esxsql with postgres driver")
		if err := migratePostgres.Execute(&buf, pi); err != nil {
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

	commands := strings.Split(buf.String(), ";")
	var resCmds []string
	for _, cmd := range commands {
		cmd = strings.TrimSpace(cmd)
		if cmd != "" {
			resCmds = append(resCmds, cmd)
		}
	}

	var err error
	for _, cmd := range resCmds {
		if _, err = conn.Exec(cmd); err != nil {
			if conn.ErrorCode(err) == cgerrors.ErrorCode_AlreadyExists {
				xlog.Debugf("%v", err)
				continue
			}
			return err
		}
	}

	return nil
}

// MigrateEventPartitions migrates event partitions
func MigrateEventPartitions(conn xsql.DB, config *Config, aggregateTypes ...string) error {
	var buf bytes.Buffer
	if err := config.Validate(); err != nil {
		return err
	}

	switch conn.DriverName() {
	case "pg", "postgres", "postgresql", "gopg", "pgx":
		type aggregate struct {
			Type  string
			Value string
		}
		type postgresInput struct {
			Config
			Aggregates []aggregate
		}
		cfg := *config
		if cfg.SchemaName != "" {
			cfg.SchemaName += "."
		}
		pi := postgresInput{Config: cfg}
		for _, agg := range aggregateTypes {
			pi.Aggregates = append(pi.Aggregates, aggregate{Type: ToSnakeCase(agg), Value: agg})
		}
		xlog.Infoln("Migrating esxsql partitions with postgres driver")
		if err := migratePgPartitions.Execute(&buf, pi); err != nil {
			return err
		}
	case "mysql":
		return errors.New("partitions not implemented for the mysql yet")
	default:
		return errors.New("driver not supported by the esxsql migration tool")
	}

	commands := strings.Split(buf.String(), ";")
	var resCmds []string
	for _, cmd := range commands {
		cmd = strings.TrimSpace(cmd)
		if cmd != "" {
			resCmds = append(resCmds, cmd)
		}
	}

	var err error
	for _, cmd := range resCmds {
		if _, err = conn.Exec(cmd); err != nil {
			if conn.ErrorCode(err) == cgerrors.ErrorCode_AlreadyExists {
				xlog.Debugf("%v", err)
			} else {
				xlog.Errorf("Code: %d - %v", conn.ErrorCode(err), err)
			}

		}
	}

	return nil
}

// MigrateEventStatePartitions create partitions on event state by its type.
func MigrateEventStatePartitions(conn xsql.DB, cfg *Config, handlerNames ...string) error {
	xlog.Infof("Migrating esxsql event state partitions - handlers: (%s)", strings.Join(handlerNames, ","))
	for _, handlerName := range handlerNames {
		r := strings.NewReplacer(":", "_", ".", "_", ",", "_")
		nt := r.Replace(handlerName)
		queries := []string{
			conn.Rebind(fmt.Sprintf(`CREATE TABLE %s PARTITION OF %s FOR VALUES IN('%s')`, nt, cfg.eventStateTableName(), handlerName)),
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
