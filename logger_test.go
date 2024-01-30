package logm

import (
	"errors"
	"github.com/magiconair/properties/assert"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestGetLogger(t *testing.T) {
	log := GetLogger()
	log.Infoln("test")
}

func TestGetLevel(t *testing.T) {
	level := GetLevel()
	assert.Equal(t, level, logrus.DebugLevel)
}

func TestIsLevelEnabled(t *testing.T) {
	result := IsLevelEnabled(logrus.DebugLevel)
	assert.Equal(t, result, true)

	result = IsLevelEnabled(logrus.TraceLevel)
	assert.Equal(t, result, false)
}

func TestLogger(t *testing.T) {
	Info("test")
	Infof("%s", "test")
	Debug("test")
	Debugf("%s", "test")
	Error("test")
	Errorf("%s", "test")
	Warn("test")
	Warnf("%s", "test")
}

func TestWithField(t *testing.T) {
	WithField("data", "new_data").Info("test")
	WithField("data", "new_data").Infof("%s", "test")
	WithField("data", "new_data").Debug("test")
	WithField("data", "new_data").Debugf("%s", "test")
	WithField("data", "new_data").Error("test")
	WithField("data", "new_data").Errorf("%s", "test")
	WithField("data", "new_data").Warn("test")
	WithField("data", "new_data").Warnf("%s", "test")
}

func TestWithFields(t *testing.T) {
	WithFields(logrus.Fields{"data": "new_data"}).Info("test")
	WithFields(logrus.Fields{"data": "new_data"}).Infof("%s", "test")
	WithFields(logrus.Fields{"data": "new_data"}).Debug("test")
	WithFields(logrus.Fields{"data": "new_data"}).Debugf("%s", "test")
	WithFields(logrus.Fields{"data": "new_data"}).Error("test")
	WithFields(logrus.Fields{"data": "new_data"}).Errorf("%s", "test")
	WithFields(logrus.Fields{"data": "new_data"}).Warn("test")
	WithFields(logrus.Fields{"data": "new_data"}).Warnf("%s", "test")
}

func TestWithError(t *testing.T) {
	WithError(errors.New("error")).Info("test")
}

func TestWithTime(t *testing.T) {
	WithTime(time.Now()).Info("test")
}

func TestWithContext(t *testing.T) {
	WithContext(context.Background()).Info("test")
}

func TestPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	Panic("test")
}

func TestPanicF(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	Panicf("%s", "test")
}
