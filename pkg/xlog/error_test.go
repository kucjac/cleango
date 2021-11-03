package xlog_test

import (
	"errors"
	"testing"

	"github.com/kucjac/cleango/pkg/xlog"
	"github.com/stretchr/testify/assert"
)

func TestWithErr(t *testing.T) {
	entry := xlog.StdLogger().WithErr(errors.New("errr"))
	assert.Equal(t, entry.Data["error"], errors.New("errr"))
}
