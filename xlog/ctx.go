package xlog

import (
	"context"

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

// CtxPut will put xlog into context
func CtxPut(ctx context.Context, entry *logrus.Entry) context.Context {
	return context.WithValue(ctx, CtxLogger, entry)
}

// CtxFields will add xlog fields
func CtxFields(ctx context.Context, fields logrus.Fields) context.Context {
	return CtxPut(ctx, Ctx(ctx).WithFields(fields))
}
