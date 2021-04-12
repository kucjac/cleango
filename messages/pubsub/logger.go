package pubsub

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/sirupsen/logrus"

	"github.com/kucjac/cleango/xlog"
)

var _ watermill.LoggerAdapter = (*loggerAdapter)(nil)

// loggerAdapter is the logger adapter for the watermill.LoggerAdapter messages.
type loggerAdapter struct {
	xlogger xlog.Logger
}

// NewLoggerAdapter creates new logger adapter for the watermill LoggerAdapter.
func NewLoggerAdapter(logger xlog.Logger) watermill.LoggerAdapter {
	return &loggerAdapter{xlogger: logger}
}

func (l *loggerAdapter) Error(msg string, err error, fields watermill.LogFields) {
	l.xlogger.WithFields(logrus.Fields(fields)).Errorf("%s - %v", msg, err)
}

func (l *loggerAdapter) Info(msg string, fields watermill.LogFields) {
	l.xlogger.WithFields(logrus.Fields(fields)).Info(msg)
}

func (l *loggerAdapter) Debug(msg string, fields watermill.LogFields) {
	l.xlogger.WithFields(logrus.Fields(fields)).Debug(msg)
}

func (l *loggerAdapter) Trace(msg string, fields watermill.LogFields) {
	l.xlogger.WithFields(logrus.Fields(fields)).Trace(msg)
}

func (l *loggerAdapter) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return &entryAdapter{entry: l.xlogger.WithFields(logrus.Fields(fields))}
}

var _ watermill.LoggerAdapter = (*entryAdapter)(nil)

type entryAdapter struct {
	entry *logrus.Entry
}

func (e *entryAdapter) Error(msg string, err error, fields watermill.LogFields) {
	e.entry.WithFields(logrus.Fields(fields)).Errorf("%s - %v")
}

func (e *entryAdapter) Info(msg string, fields watermill.LogFields) {
	e.entry.WithFields(logrus.Fields(fields)).Info(msg)
}

func (e *entryAdapter) Debug(msg string, fields watermill.LogFields) {
	e.entry.WithFields(logrus.Fields(fields)).Debug(msg)
}

func (e *entryAdapter) Trace(msg string, fields watermill.LogFields) {
	e.entry.WithFields(logrus.Fields(fields)).Trace(msg)
}

func (e *entryAdapter) With(fields watermill.LogFields) watermill.LoggerAdapter {
	e.entry = e.entry.WithFields(logrus.Fields(fields))
	return e
}
