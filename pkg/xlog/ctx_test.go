package xlog_test

import (
	"context"
	"testing"

	"github.com/kucjac/cleango/pkg/xlog"
	"github.com/stretchr/testify/assert"
)

func TestCtx(t *testing.T) {
	ctx := context.Background()
	ctx = xlog.CtxPut(ctx, xlog.WithField("test", "test"))
	entry := xlog.Ctx(ctx)
	assert.Equal(t, entry.Data["test"], "test")
	entry = xlog.Ctx(context.Background())
	assert.Equal(t, entry.Data["test"], nil)
}
