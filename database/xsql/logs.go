package xsql

import (
	"fmt"
	"strings"
	"time"

	"github.com/kucjac/cleango/xlog"
	"github.com/sirupsen/logrus"
)

func logQuery(tx, query string, ts time.Time, args ...interface{}) {
	if !xlog.IsLevelEnabled(logrus.DebugLevel) {
		return
	}
	sb := strings.Builder{}
	sb.WriteString("query: ")
	sb.WriteString(query)
	if len(args) > 0 {
		sb.WriteString(" args: (")
		for i, arg := range args {
			sb.WriteString(fmt.Sprintf("%v", arg))
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
	sb.WriteString(time.Since(ts).String())
	xlog.Debug(sb.String())
}
