package esxsql

import (
	"bytes"
	_ "embed"
	"errors"
	"strings"
	"text/template"

	"github.com/jmoiron/sqlx"

	"github.com/kucjac/cleango/xlog"
)

//go:embed mysql.tmpl
var mysqlMigrateQuery string

//go:embed mysql.tmpl
var postgresMysqlQuery string

var (
	migrateMySQL    *template.Template
	migratePostgres *template.Template
)

func init() {
	t := template.New("")
	migrateMySQL = template.Must(t.Parse(mysqlMigrateQuery))
	migratePostgres = template.Must(t.Parse(postgresMysqlQuery))
}

// Migrate executes table and types migration for the event store and snapshot.
// The table names are taken from the config.
func Migrate(config *Config, conn *sqlx.DB) error {
	var buf bytes.Buffer
	if err := config.Validate(); err != nil {
		return err
	}

	switch conn.DriverName() {
	case "pg", "postgres", "postgresql", "gopg", "pgx":
		if err := migratePostgres.Execute(&buf, config); err != nil {
			return err
		}
	case "mysql":
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

	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	for _, cmd := range resCmds {
		xlog.Debug(cmd)
		if _, err = tx.Exec(cmd); err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}
