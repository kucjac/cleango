package xlog_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kucjac/cleango/pkg/xlog"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCallerHook(t *testing.T) {
	buf := new(bytes.Buffer)
	lx := xlog.New()
	lx.AddHook(xlog.NewCallerHook([]logrus.Level{logrus.InfoLevel}))
	lx.SetReportCaller(true)
	tf := xlog.NewTextFormatter(false)
	tf.DisableTimestamp = true
	lx.SetFormatter(tf)
	lx.SetOutput(buf)
	lx.WithField("msg", "helloooo").Print("hello")
	str := buf.String()
	assert.True(t, strings.Contains(str, "INFO  hello"))
	assert.True(t, strings.Contains(str, "\"msg\"=\"helloooo\""))
	assert.True(t, strings.Contains(str, "\"func\"=\"github.com/kucjac/cleango/pkg/xlog_test.TestCallerHook\""))
	assert.True(t, strings.Contains(str, "github.com/kucjac/cleango/pkg/xlog/caller_test.go:22"), str)
}
