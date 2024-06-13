package mlog

import (
	"os"
	"slices"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
)

const logTimeFormat = "2006-01-02T15:04:05-0700"

type LoggerFormatter struct {
	defaultField    logrus.Fields
	callerSkip      int
	maskingEnabled  bool
	sensitiveFields []string
	formatter       logrus.Formatter
}

func (f *LoggerFormatter) initFormatter() {
	fieldMap := logrus.FieldMap{
		logrus.FieldKeyTime: "@timestamp",
		logrus.FieldKeyMsg:  "message",
	}

	env, _ := os.LookupEnv("ENV")
	if env == "dev" {
		textFormatter := &logrus.TextFormatter{
			FieldMap: fieldMap,
		}
		textFormatter.TimestampFormat = logTimeFormat
		f.formatter = textFormatter
	} else {
		jsonFormatter := &logrus.JSONFormatter{
			FieldMap: fieldMap,
		}
		jsonFormatter.TimestampFormat = logTimeFormat
		f.formatter = jsonFormatter
	}
}

func (f *LoggerFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if f.formatter == nil {
		f.initFormatter()
	}
	newEntry := entry.WithFields(f.defaultField)
	if f.maskingEnabled {
		if maskData, ok := f.maskFields(newEntry.Data).(logrus.Fields); ok {
			newEntry.Data = maskData
		}
	}
	entry.Data = newEntry.Data
	return f.formatter.Format(entry)
}

func (f *LoggerFormatter) maskFields(fields interface{}) interface{} {
	newData := make(logrus.Fields)

	switch value := fields.(type) {
	case logrus.Fields:
		for key, fieldValue := range value {
			newData = f.valueMasking(newData, key, fieldValue)
		}
	case map[string]interface{}:
		for key, fieldValue := range value {
			newData = f.valueMasking(newData, key, fieldValue)
		}
	default:
		return fields
	}
	return newData
}

func (f *LoggerFormatter) valueMasking(newData logrus.Fields, key string, fieldValue interface{}) logrus.Fields {
	snakeKey := strings.ToLower(strcase.ToSnake(key))
	//If the field is sensitive
	if slices.Contains(f.sensitiveFields, snakeKey) {
		if valueStr, ok := fieldValue.(string); ok {
			newData[key] = maskKeyValue(snakeKey, valueStr)
		} else {
			newData[key] = "<***mask***>"
		}
		return newData
	}
	switch subFieldValue := fieldValue.(type) {
	case logrus.Fields:
		newData[key] = f.maskFields(subFieldValue)
	case map[string]interface{}:
		newData[key] = f.maskFields(subFieldValue)
	case []interface{}:
		newData[key] = f.maskArrayFields(subFieldValue)
	default:
		newData[key] = subFieldValue
	}
	return newData
}

func maskKeyValue(key string, value string) string {
	switch key {
	case "name", "firstname", "firstName", "lastname", "lastName":
		return Name(value)
	case "addr", "address":
		return Address(value)
	case "email", "mail":
		return Email(value)
	case "mobile_no", "mobilePhone", "mobile_number", "mobile":
		return Mobile(value)
	case "id", "passport", "passport_id", "passport_no", "passport_number":
		return ID(value)
	case "phone_number", "phone_no", "tel", "telephone_no", "telephone", "phone":
		return Mobile(value)
	case "card_no", "credit_card", "debit_card", "credit_card_no", "debit_card_no":
		return CreditCard(value)
	case "national_id", "cid", "citizen_id":
		return CidMasker(value)
	default:
		return Password(value)
	}
}

// The maskArrayFields method is added to handle array type fieldValue separately.
func (f *LoggerFormatter) maskArrayFields(arrayFieldValue []interface{}) []interface{} {
	//Looping over array elements to mask them
	for index, value := range arrayFieldValue {
		arrayFieldValue[index] = f.maskFields(value)
	}
	return arrayFieldValue
}
