package logm

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/mcuadros/go-defaults"
	"github.com/sirupsen/logrus"
)

type (
	Logger struct {
		*logrus.Logger `yaml:"log"`
	}
	Config struct {
		Level   string        `yaml:"level" env:"LOG_LEVEL" env-default:"debug" mapstructure:"level" default:"debug"`
		Masking ConfigMasking `yaml:"masking" mapstructure:"masking"`
	}

	ConfigMasking struct {
		Enabled    bool     `yaml:"enabled" env:"LOG_MASKING_ENABLED" env-default:"true" mapstructure:"enabled" default:"true"`
		FieldNames []string `yaml:"field-names" env:"LOG_MASKING_ENABLED" mapstructure:"field-names" default:""`
	}
)

var defaultSensitiveFields = []string{
	"name", "firstname", "first_name", "lastname", "last_name", "addr",
	"address", "email", "mail", "mobile_no", "mobile_phone", "mobile_number", "mobile",
	"phone_number", "phone_no", "tel", "telephone_no", "telephone", "phone", "card_no",
	"credit_card", "debit_card", "credit_card_no", "debit_card_no",
	"id", "passport", "passport_id", "passport_no", "passport_number",
	"national_id", "cid", "citizen_id", "cvc", "password",
	"x-api-key", "authorization", "x-authorization", "field-names",
}

const appNameKey = "app_name"

func newStandardLogger() *Logger {
	configDefault := &Config{}
	defaults.SetDefaults(configDefault)
	return newLogger(configDefault, "", nil)
}

func newLogger(config *Config, appName string, formatter logrus.Formatter) *Logger {
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
	logger.SetFormatter(&LoggerFormatter{
		defaultField:    defaultField,
		maskingEnabled:  config.Masking.Enabled,
		sensitiveFields: removeDuplicates(sensitiveFields),
		formatter:       formatter,
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
