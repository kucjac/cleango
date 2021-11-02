package xlog

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	buf := new(bytes.Buffer)
	lx := New()
	lx.SetFormatter(NewTextFormatter(false))
	lx.SetOutput(buf)
	lx.WithField("msg", "helloooo").Print("hello")
	str := buf.String()
	assert.True(t, strings.Contains(str, "INFO  hello \"msg\"=\"helloooo\""))
}
