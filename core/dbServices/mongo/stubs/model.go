package model

import (
	"encoding/base64"
	"errors"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/DanielRenne/GoCore/core/dbServices"
	"github.com/asaskevich/govalidator"
)

const (
	TRANSACTION_DATATYPE_ORIGINAL = 1
	TRANSACTION_DATATYPE_NEW      = 2

	TRANSACTION_CHANGETYPE_INSERT = 1
	TRANSACTION_CHANGETYPE_UPDATE = 2
	TRANSACTION_CHANGETYPE_DELETE = 3

	MGO_RECORD_NOT_FOUND = "not found"

	VALIDATION_ERROR                   = "ValidationError"
	VALIDATION_ERROR_REQUIRED          = "ValidationErrorRequiredFieldMissing"
	VALIDATION_ERROR_EMAIL             = "ValidationErrorInvalidEmail"
	VALIDATION_ERROR_SPECIFIC_REQUIRED = "ValidationFieldSpecificRequired"
	VALIDATION_ERROR_SPECIFIC_EMAIL    = "ValidationFieldSpecificEmailRequired"
)

type modelEntity interface {
	Save() error
	Delete() error
}

type modelCollection interface {
	Rollback(transactionId string) error
}

type collection interface {
	Query() *Query
}

type tQueue struct {
	sync.RWMutex
	queue map[string]*transactionsToPersist
}

type transactionsToPersist struct {
	t             *Transaction
	newItems      []entityTransaction
	originalItems []entityTransaction
	startTime     time.Time
}

type entityTransaction struct {
	changeType int
	committed  bool
	entity     modelEntity
}

var transactionQueue tQueue

func init() {
	transactionQueue.queue = make(map[string]*transactionsToPersist)
	go clearTransactionQueue()
}

func Q(k string, v interface{}) map[string]interface{} {
	return map[string]interface{}{k: v}
}

//Every 12 hours check the transactionQueue and remove any outstanding stale transactions > 48 hours old
func clearTransactionQueue() {

	transactionQueue.Lock()

	for key, value := range transactionQueue.queue {

		if time.Since(value.startTime).Hours() > 48 {
			delete(transactionQueue.queue, key)
		}
	}

	transactionQueue.Unlock()

	time.Sleep(12 * time.Hour)
	clearTransactionQueue()
}

func getBase64(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func decodeBase64(value string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}

	return string(data[:]), nil
}

func getNow() time.Time {
	return time.Now()
}

func removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func IsValidationError(err error) bool {
	if err == nil {
		return false
	}
	if err.Error() == VALIDATION_ERROR || err.Error() == VALIDATION_ERROR_EMAIL {
		return true
	}
	return false
}

func validateFields(x interface{}, objectToUpdate interface{}, val reflect.Value) error {

	isError := false
	for key, value := range dbServices.GetValidationTags(x) {

		fieldValue := dbServices.GetReflectionFieldValue(key, objectToUpdate)
		validations := strings.Split(value, ",")

		if validations[0] != "" {
			if err := validateRequired(fieldValue, validations[0]); err != nil {
				dbServices.SetFieldValue("Errors."+key, val, VALIDATION_ERROR_SPECIFIC_REQUIRED)
				isError = true
			}
		}
		if validations[1] != "" {

			cleanup, err := validateType(fieldValue, validations[1])

			if err != nil {
				if err.Error() == VALIDATION_ERROR_EMAIL {
					dbServices.SetFieldValue("Errors."+key, val, VALIDATION_ERROR_SPECIFIC_EMAIL)
				}
				isError = true
			}

			if cleanup != "" {
				dbServices.SetFieldValue(key, val, cleanup)
			}

		}

	}

	if isError {
		return errors.New(VALIDATION_ERROR)
	}

	return nil
}

func validateRequired(value string, tagValue string) error {
	if tagValue == "true" {
		if value == "" {
			return errors.New(VALIDATION_ERROR_REQUIRED)
		}
		return nil
	}
	return nil
}

func validateType(value string, tagValue string) (string, error) {
	switch tagValue {
	case dbServices.VALIDATION_TYPE_EMAIL:
		return "", validateEmail(value)
	}
	return "", nil
}

func validateEmail(value string) error {
	if !govalidator.IsEmail(value) {
		return errors.New(VALIDATION_ERROR_EMAIL)
	}
	return nil
}

func getJoins(x reflect.Value, remainingRecursions string) (joins []join) {
	if remainingRecursions == "" {
		return
	}

	fields := strings.Split(remainingRecursions, ".")
	fieldName := fields[0]

	joinsField := x.FieldByName("Joins")

	if joinsField.Kind() != reflect.Struct {
		return
	}

	if fieldName == JOIN_ALL {
		for i := 0; i < joinsField.NumField(); i++ {

			typeField := joinsField.Type().Field(i)
			name := typeField.Name
			tagValue := typeField.Tag.Get("join")
			splitValue := strings.Split(tagValue, ",")
			var j join
			j.collectionName = splitValue[0]
			j.joinSchemaName = splitValue[1]
			j.joinFieldRefName = splitValue[2]
			j.joinFieldName = name
			j.joinSpecified = JOIN_ALL
			joins = append(joins, j)
		}
	} else {
		typeField, ok := joinsField.Type().FieldByName(fieldName)
		if ok == false {
			return
		}
		name := typeField.Name
		tagValue := typeField.Tag.Get("join")
		splitValue := strings.Split(tagValue, ",")
		var j join
		j.collectionName = splitValue[0]
		j.joinSchemaName = splitValue[1]
		j.joinFieldRefName = splitValue[2]
		j.joinFieldName = name
		j.joinSpecified = strings.Replace(remainingRecursions, fieldName+".", "", 1)
		joins = append(joins, j)
	}
	return
}

// Start of autogenerated code....
