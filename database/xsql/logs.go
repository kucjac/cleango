package xsql

import (
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kucjac/cleango/xlog"
	"github.com/sirupsen/logrus"
)

func logQuery(tx, query string, ts time.Time, config *Config, res sql.Result, args ...interface{}) {
	execDuration := time.Since(ts)
	isLongRunning := config.WarnLongQueries && execDuration > config.LongQueriesTime
	if !xlog.IsLevelEnabled(logrus.DebugLevel) && !isLongRunning {
		return
	}
	sb := strings.Builder{}
	sb.WriteString("query: ")
	sb.WriteString(query)
	sb.WriteRune(';')
	if len(args) > 0 {
		sb.WriteString(" args: (")
		for i, arg := range args {
			switch at := arg.(type) {
			case driver.Valuer:
				v, err := at.Value()
				if err != nil {
					sb.WriteString("{ERR-INVALID-VALUE}")
				} else {
					switch vt := v.(type) {
					case string:
						sb.WriteRune('\'')
						sb.WriteString(vt)
						sb.WriteRune('\'')
					case []byte:
						if len(vt) == 0 {
							sb.WriteString("NULL")
						} else {
							sb.WriteString("'x")
							sb.WriteString(hex.EncodeToString(vt))
							sb.WriteRune('\'')
						}
					default:
						sb.WriteString(fmt.Sprintf("%v", vt))
					}
				}
			case string:
				sb.WriteRune('\'')
				sb.WriteString(at)
				sb.WriteRune('\'')
			case []byte:
				if len(at) == 0 {
					sb.WriteString("NULL")
				} else {
					sb.WriteString("'x")
					sb.WriteString(hex.EncodeToString(at))
					sb.WriteRune('\'')
				}
			default:
				sb.WriteString(fmt.Sprintf("%v", arg))
			}

			if i != len(args)-1 {
				sb.WriteRune(',')
			}
		}
		sb.WriteRune(')')
	}
	if tx != "" {
		sb.WriteString(" tx: ")
		sb.WriteString(tx)
		sb.WriteRune(';')
	}
	if res != nil {
		sb.WriteString(" affected: ")
		affected, _ := res.RowsAffected()
		sb.WriteString(strconv.Itoa(int(affected)))
		if affected == 1 {
			sb.WriteString(" row")
		} else {
			sb.WriteString(" rows")
		}
		sb.WriteRune(';')
	}
	sb.WriteString(" taken: ")
	sb.WriteString(execDuration.String())

	// Check if the query was marked to be long-running.
	if isLongRunning {
		xlog.Warn(sb.String())
	} else {
		xlog.Debug(sb.String())
	}
}
