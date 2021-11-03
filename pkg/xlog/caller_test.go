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
	lx.SetFormatter(xlog.NewTextFormatter(false))
	lx.SetOutput(buf)
	lx.WithField("msg", "helloooo").Print("hello")
	str := buf.String()
	assert.True(t, strings.Contains(str, "INFO  hello \"msg\"=\"helloooo\""))
}
