package xlog

import (
	"github.com/sirupsen/logrus"
)

// Logger is interface which satisfies Entry and Log
type Logger interface {
	logrus.FieldLogger
}

var _ Logger = (*logrus.Entry)(nil)
var _ Logger = (*Log)(nil)

// Log is wrapper around logrus
type Log struct {
	*logrus.Logger
}

// NewEntry creates
func (logger *Log) NewEntry() *logrus.Entry {
	return logrus.NewEntry(logger.Logger)
}

// New will create new instance if logrus
func New() *Log {
	l := logrus.New()
	// l.SetFormatter(NewTextFormatter(false))
	l.SetLevel(logrus.DebugLevel)
	return &Log{Logger: l}
}
