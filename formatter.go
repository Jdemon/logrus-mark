package logm

import (
	"fmt"
	m "github.com/ggwhite/go-masker"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"slices"
	"strings"
)

type JSONFormatter struct {
	defaultField    logrus.Fields
	callerSkip      int
	maskingEnabled  bool
	sensitiveFields []string
	jsonFormatter   *logrus.JSONFormatter
}

func (f *JSONFormatter) initJSONFormatter() {
	f.jsonFormatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "@timestamp",
			logrus.FieldKeyMsg:  "message",
		},
	}
}

func (f *JSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if f.jsonFormatter == nil {
		f.initJSONFormatter()
	}
	newEntry := entry.WithFields(f.defaultField)
	if f.maskingEnabled {
		if maskData, ok := f.maskFields(newEntry.Data).(logrus.Fields); ok {
			newEntry.Data = maskData
		}
	}
	entry.Data = newEntry.Data
	return f.jsonFormatter.Format(entry)
}

func (f *JSONFormatter) maskFields(fields interface{}) interface{} {
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

func (f *JSONFormatter) valueMasking(newData logrus.Fields, key string, fieldValue interface{}) logrus.Fields {
	snakeKey := strings.ToLower(strcase.ToSnake(key))
	//If the field is sensitive
	if slices.Contains(f.sensitiveFields, snakeKey) {
		if valueStr, ok := fieldValue.(string); ok {
			newData[snakeKey] = maskKeyValue(snakeKey, valueStr)
		} else {
			newData[snakeKey] = "******"
		}
		return newData
	}
	switch subFieldValue := fieldValue.(type) {
	case logrus.Fields:
		newData[snakeKey] = f.maskFields(subFieldValue)
	case map[string]interface{}:
		newData[snakeKey] = f.maskFields(subFieldValue)
	case []interface{}:
		newData[snakeKey] = f.maskArrayFields(subFieldValue)
	default:
		newData[snakeKey] = subFieldValue
	}
	return newData
}

func maskKeyValue(key string, value string) string {
	switch key {
	case "name", "firstname", "first_name", "lastname", "last_name":
		return m.Name(value)
	case "addr", "address":
		return m.Address(value)
	case "email", "mail":
		return m.Email(value)
	case "mobile_no", "mobile_phone", "mobile_number", "mobile":
		return m.Mobile(value)
	case "id", "passport", "passport_id", "passport_no", "passport_number":
		return m.ID(value)
	case "phone_number", "phone_no", "tel", "telephone_no", "telephone", "phone":
		return m.Mobile(value)
	case "card_no", "credit_card", "debit_card", "credit_card_no", "debit_card_no":
		return m.CreditCard(value)
	case "national_id", "cid", "citizen_id":
		return cidMasker(value)
	default:
		return m.Password(value)
	}
}

// The maskArrayFields method is added to handle array type fieldValue separately.
func (f *JSONFormatter) maskArrayFields(arrayFieldValue []interface{}) []interface{} {
	//Looping over array elements to mask them
	for index, value := range arrayFieldValue {
		arrayFieldValue[index] = f.maskFields(value)
	}
	return arrayFieldValue
}

func cidMasker(input string) string {
	// Check if input is empty
	if len(input) == 0 {
		return ""
	}

	// Remove hyphens from the input
	input = strings.ReplaceAll(input, "-", "")

	// Check if the length of the cleaned input is not 13
	if len(input) != 13 {
		return "*-****-*****-**-*"
	}

	// Construct the masked CID
	maskedCID := fmt.Sprintf(
		"%s-%s**-*****-%s-%s",
		input[:1],
		input[1:3],
		input[10:12],
		input[12:],
	)

	return maskedCID
}
