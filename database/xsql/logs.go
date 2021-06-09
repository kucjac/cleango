package xsql

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/kucjac/cleango/xlog"
	"github.com/sirupsen/logrus"
)

func logQuery(tx, query string, ts time.Time, config *Config, args ...interface{}) {
	execDuration := time.Since(ts)
	isLongRunning := config.WarnLongQueries && execDuration > config.LongQueriesTime
	if !xlog.IsLevelEnabled(logrus.DebugLevel) && !isLongRunning {
		return
	}
	sb := strings.Builder{}
	sb.WriteString("query: ")
	sb.WriteString(query)
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
