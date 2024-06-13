package mlog

import (
	"github.com/sirupsen/logrus"
)

type contextKey int

const (
	TraceID contextKey = iota
	CorrelationID
	RequestID
)

type ContextHook struct {
}

func (hook *ContextHook) Fire(entry *logrus.Entry) error {
	if entry.Context == nil {
		return nil
	}

	setEntryData := func(key string, value interface{}) {
		if value != nil {
			entry.Data[key] = value
		}
	}
	setEntryData("traceID", entry.Context.Value(TraceID))
	setEntryData("requestID", entry.Context.Value(RequestID))
	setEntryData("correlationID", entry.Context.Value(CorrelationID))
	return nil
}

func (hook *ContextHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.TraceLevel,
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}
