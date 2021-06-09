package xpg

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/kucjac/cleango/xlog"
)

var spaceReg = regexp.MustCompile(`\s+`)

// Log is regular xlog
type Log struct {
	strip int
}

var _ pg.QueryHook = &Log{}

// NewLogger creates new xlog, specify strip if you want queries longer then
// x chars be striped
func NewLogger(strip int) Logger {
	return &Log{strip: strip}
}

const startAt = "start_at"

func (log *Log) truncate(query string) string {
	query = spaceReg.ReplaceAllString(query, " ")
	query = strings.TrimSpace(query)
	if log.strip == 0 {
		return query
	}
	if len(query) <= log.strip {
		return query
	}
	return fmt.Sprintf("%s ....", query[0:log.strip-1])
}

// BeforeQuery as before query is executed for log entry
func (log *Log) BeforeQuery(ctx context.Context, event *pg.QueryEvent) (context.Context, error) {
	if event.Stash == nil {
		event.Stash = make(map[interface{}]interface{})
	}
	event.Stash[startAt] = time.Now()
	return ctx, nil
}

// AfterQuery as after query went back
func (log *Log) AfterQuery(ctx context.Context, event *pg.QueryEvent) error {
	lgr := xlog.Ctx(ctx)
	query, err := event.FormattedQuery()
	if err != nil {
		lgr.Error(err)
		return err
	}
	q := log.truncate(fmt.Sprintf("%s", query))
	v := event.Stash[startAt].(time.Time)
	diff := time.Since(v)
	lgr = lgr.WithField("execution", diff.Seconds())
	switch {
	case strings.Contains(q, "gopg:ping"):
	case strings.Contains(q, "SELECT 'healthcheck'"):
	default:
		lgr.Debugln(q)
	}
	if diff < time.Second {
		return nil
	}
	lgr.Warnf("query %s took %s", query, diff)
	return nil
}
