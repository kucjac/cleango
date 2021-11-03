package xlog_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/kucjac/cleango/pkg/xlog"
	"github.com/stretchr/testify/require"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInfo(t *testing.T) {
	buf := new(bytes.Buffer)
	lx := xlog.New()
	f := xlog.NewTextFormatter(false)
	f.DisableTimestamp = true
	lx.SetFormatter(f)
	lx.SetOutput(buf)
	lx.Info("$TestT$")
	lx.Warnf("$TestT$ %s", "test!")
	lx.WithField("hello", "test!").Error("damm")
	lx.WithField("hello", "").Error("damm")
	str := buf.String()

	assert.True(t, strings.Contains(str, "INFO  $TestT$"))
	assert.True(t, strings.Contains(str, "WARN  $TestT$ test"))
	assert.True(t, strings.Contains(str, "ERROR damm \"hello\"=\"test!\""))
	assert.True(t, strings.Contains(str, "ERROR damm \"hello\"=\"\""))

	e1 := lx.WithField("test", "test")
	e1.Message = "!"
	e1.Level = logrus.FatalLevel
	data, err := f.Format(e1)
	require.NoError(t, err)
	assert.True(t, strings.Contains(string(data), "FATAL ! \"test\"=\"test\""))

	e2 := lx.WithField("test", "test")
	e2.Message = "!"
	e2.Level = logrus.PanicLevel
	data, err = f.Format(e2)
	require.NoError(t, err)
	assert.True(t, strings.Contains(string(data), "PANIC ! \"test\"=\"test\""))
}

func TestErrorFormat(t *testing.T) {
	buf := new(bytes.Buffer)
	lx := xlog.New()
	f := xlog.NewTextFormatter(false)
	f.DisableTimestamp = true
	lx.SetLevel(logrus.TraceLevel)
	lx.SetFormatter(f)
	lx.SetOutput(buf)
	req := xlog.HTTPRequest{
		RequestMethod: "GET",
		RequestURL:    "/wuf",
		Status:        200,
		Latency:       time.Second.String(),
	}
	lx.WithField(xlog.HTTPRequestKey, req).Debug(logrus.DebugLevel)
	str := buf.String()
	assert.Equal(t, "DEBUG [GET] 200 |       1s |            /wuf\n", str)
}
