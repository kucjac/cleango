package xlog

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestWithErr(t *testing.T) {
	entry := std.WithErr(errors.New("errr"))
	assert.Equal(t, entry.Data["error"], errors.New("errr"))
}

