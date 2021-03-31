package xlog

import (
	"github.com/sirupsen/logrus"
)

// WithErr will format error. If errdef.ErrSet then maps
// of both are merged
func (logger *Log) WithErr(err error) *logrus.Entry {
		return logger.Logger.WithError(err)
}

