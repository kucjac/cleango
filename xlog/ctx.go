package xlog

import (
	"context"

	"github.com/kucjac/cleango/meta"
	"github.com/sirupsen/logrus"
)

// CtxKey represents custom type for context
type CtxKey string

// CtxLogger represents key in context
const CtxLogger CtxKey = "xlog:entry"

// Ctx will extract xlog from context
func Ctx(ctx context.Context) *logrus.Entry {
	if ctx == nil {
		return logrus.NewEntry(std.Logger)
	}
	l, ok := ctx.Value(CtxLogger).(*logrus.Entry)
	if !ok || l == nil {
		return logrus.NewEntry(std.Logger)
	}

	return l
}

// CtxPut will put log into context
func CtxPut(ctx context.Context, entry *logrus.Entry) context.Context {
	return context.WithValue(ctx, CtxLogger, entry)
}

// CtxFields will add xlog fields
func CtxFields(ctx context.Context, fields logrus.Fields) context.Context {
	return CtxPut(ctx, Ctx(ctx).WithFields(fields))
}

var _ logrus.Hook = RequestIDHook{}

// RequestIDHook is a logging hook used to log the request id stored in the context metadata.
type RequestIDHook struct{}

func (m RequestIDHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (m RequestIDHook) Fire(entry *logrus.Entry) error {
	rID, ok := meta.RequestID(entry.Context)
	if !ok {
		return nil
	}
	entry.WithField("request-id", rID)
	return nil
}

var _ logrus.Hook = UserIDHook{}

// UserIDHook is a hook that gets the user id from the context metadata and logs it on each log level.
type UserIDHook struct{}

func (u UserIDHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (u UserIDHook) Fire(entry *logrus.Entry) error {
	uID, ok := meta.UserID(entry.Context)
	if !ok {
		return nil
	}
	entry.WithField("user-id", uID)
	return nil
}
