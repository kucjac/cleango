package xlog

import (
	"github.com/sirupsen/logrus"
)

// WithErr format error.
func (logger *Log) WithErr(err error) *logrus.Entry {
	return logger.Logger.WithError(err)
}
