package xlog

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCtx(t *testing.T) {
	ctx := context.Background()
	ctx = CtxPut(ctx, WithField("test", "test"))
	entry := Ctx(ctx)
	assert.Equal(t, entry.Data["test"], "test")
	entry = Ctx(context.Background())
	assert.Equal(t, entry.Data["test"], nil)
}
