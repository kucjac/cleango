package xlog_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kucjac/cleango/pkg/xlog"
	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	buf := new(bytes.Buffer)
	lx := xlog.New()
	lx.SetFormatter(xlog.NewTextFormatter(false))
	lx.SetOutput(buf)
	lx.WithField("msg", "helloooo").Print("hello")
	str := buf.String()
	assert.True(t, strings.Contains(str, "INFO  hello \"msg\"=\"helloooo\""))
}
