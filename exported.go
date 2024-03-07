package logm

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"time"
)

var std = newStandardLogger()

func NewLogger(config *Config, appName string) {
	std = newLogger(config, appName, nil)
}

func GetLogger() *Logger {
	return std
}

func GetLevel() logrus.Level {
	return std.GetLevel()
}

func IsLevelEnabled(level logrus.Level) bool {
	return std.IsLevelEnabled(level)
}

// AddHook adds a hook to the standard logger hooks.
func AddHook(hook logrus.Hook) {
	std.AddHook(hook)
}

// WithError creates an entry from the standard logger and adds an error to it, using the value defined in ErrorKey as key.
func WithError(err error) *logrus.Entry {
	entry := std.WithField(logrus.ErrorKey, err)
	callerEntry(entry, 3)
	return entry
}

// WithContext creates an entry from the standard logger and adds a context to it.
func WithContext(ctx context.Context) *logrus.Entry {
	entry := std.WithContext(ctx)
	callerEntry(entry, 3)
	return entry
}

// WithField creates an entry from the standard logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithField(key string, value interface{}) *logrus.Entry {
	entry := std.WithField(key, value)
	callerEntry(entry, 3)
	return entry
}

// WithFields creates an entry from the standard logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithFields(fields logrus.Fields) *logrus.Entry {
	entry := std.WithFields(fields)
	callerEntry(entry, 3)
	return entry
}

// WithTime creates an entry from the standard logger and overrides the time of
// logs generated with it.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithTime(t time.Time) *logrus.Entry {
	entry := std.WithTime(t)
	callerEntry(entry, 3)
	return std.WithTime(t)
}

func callerEntry(entry *logrus.Entry, skip int) *logrus.Entry {
	entry.Data["file"] = fileInfo(skip)
	return entry
}

func caller() *logrus.Entry {
	entry := std.WithFields(logrus.Fields{})
	entry.Data["file"] = fileInfo(3)
	return entry
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	caller().Debugln(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	caller().Infoln(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	caller().Warnln(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	caller().Errorln(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	caller().Panicln(args...)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	caller().Fatalln(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	caller().Debugf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	caller().Infof(format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	caller().Warnf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	caller().Errorf(format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	caller().Panicf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalf(format string, args ...interface{}) {
	caller().Fatalf(format, args...)
}

func ConvertStructToDataFields(v any) logrus.Fields {
	marshal, _ := json.Marshal(v)
	fields := make(logrus.Fields)
	_ = json.Unmarshal(marshal, &fields)
	return logrus.Fields{
		"data": fields,
	}
}

func removeDuplicates(input []string) []string {
	uniqueValues := make(map[string]bool)
	var result []string

	for _, val := range input {
		if !uniqueValues[val] {
			uniqueValues[val] = true
			result = append(result, val)
		}
	}

	return result
}
