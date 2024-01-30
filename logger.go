package logm

import (
	"fmt"
	"github.com/mcuadros/go-defaults"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

type (
	Logger struct {
		*logrus.Logger
	}
	Config struct {
		Level   string        `mapstructure:"level" default:"debug"`
		Masking ConfigMasking `mapstructure:"masking"`
	}

	ConfigMasking struct {
		Enabled    bool     `mapstructure:"enabled" default:"true"`
		FieldNames []string `mapstructure:"field-names" default:"[password]"`
	}
)

var defaultSensitiveFields = []string{
	"name", "firstname", "first_name", "lastname", "last_name", "addr",
	"address", "email", "mail", "mobile_no", "mobile_phone", "mobile_number", "mobile",
	"phone_number", "phone_no", "tel", "telephone_no", "telephone", "phone", "card_no",
	"credit_card", "debit_card", "credit_card_no", "debit_card_no",
	"id", "passport", "passport_id", "passport_no", "passport_number",
	"national_id", "cid", "citizen_id", "cvc", "password",
	"x-api-key", "authorization", "x-authorization",
}

const appNameKey = "app_name"

func newStandardLogger() *Logger {
	configDefault := &Config{}
	defaults.SetDefaults(configDefault)
	return newLogger(configDefault, "")
}

func newLogger(config *Config, appName string) *Logger {
	logger := logrus.New()
	logger.SetLevel(parseLogLevel(config.Level))
	logger.SetReportCaller(false)

	defaultField := logrus.Fields{}
	if appName != "" {
		defaultField = logrus.Fields{
			appNameKey: appName,
		}
	}

	sensitiveFields := config.Masking.FieldNames
	sensitiveFields = append(sensitiveFields, defaultSensitiveFields...)
	logger.SetFormatter(&JSONFormatter{
		defaultField:    defaultField,
		maskingEnabled:  config.Masking.Enabled,
		sensitiveFields: removeDuplicates(sensitiveFields),
	})
	log := Logger{
		Logger: logger,
	}

	entry := log.WithFields(ConvertStructToDataFields(config))
	entry.Data["file"] = fileInfo(1)
	entry.Info("initial logger")
	return &log
}

var (
	logLevel = map[string]logrus.Level{
		"info":  logrus.InfoLevel,
		"panic": logrus.PanicLevel,
		"fatal": logrus.FatalLevel,
		"error": logrus.ErrorLevel,
		"warn":  logrus.WarnLevel,
		"debug": logrus.DebugLevel,
		"trace": logrus.TraceLevel,
	}
)

func parseLogLevel(str string) logrus.Level {
	c, ok := logLevel[strings.ToLower(str)]
	if !ok {
		return logrus.InfoLevel
	}
	return c
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
