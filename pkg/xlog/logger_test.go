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
	tf := xlog.NewTextFormatter(false)
	tf.DisableTimestamp = true
	lx.SetFormatter(tf)
	lx.SetOutput(buf)
	lx.WithField("msg", "helloooo").Print("hello")
	str := buf.String()
	assert.True(t, strings.Contains(str, "INFO  hello \"msg\"=\"helloooo\""))
}
