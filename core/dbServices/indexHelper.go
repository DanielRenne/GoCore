package dbServices

import (
	"github.com/DanielRenne/GoCore/core/extensions"
	"reflect"
	"strconv"
	"strings"
)

//GetIndexes provides a way to reflect on your structure to get structs tagged with `dbIndex`.
//This function is used to generate Indexes for MongoDB and other databases.
func GetDBIndexes(x interface{}) map[string]string {
	keys := make(map[string]string)
	getDBIndexesRecursive(reflect.ValueOf(x), keys, "")
	return keys
}

func getDBIndexesRecursive(val reflect.Value, keys map[string]string, key string) {

	kind := val.Kind()

	if kind == reflect.Slice {
		return
	}

	for i := 0; i < val.NumField(); i++ {

		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		index := typeField.Tag.Get("dbIndex")

		appendKey := strings.Title(typeField.Name)

		if !valueField.CanInterface() {
			continue
		}

		field := valueField.Interface()
		fieldval := reflect.ValueOf(field)

		switch fieldval.Kind() {

		case reflect.Array, reflect.Slice, reflect.Struct:
			getDBIndexesRecursive(fieldval, keys, key+appendKey+".")

		default:
			if index != "" {
				keys[key+appendKey] = index
			}
		}
	}
}

//GetValidationTags provides a way to reflect on your structure to get structs tagged with `dbIndex`.
//This function is used to generate Indexes for MongoDB and other databases.
func GetValidationTags(x interface{}) map[string]string {
	keys := make(map[string]string)
	getValidationTagsRecursive(reflect.ValueOf(x), keys, "")
	return keys
}

func getValidationTagsRecursive(val reflect.Value, keys map[string]string, key string) {

	kind := val.Kind()

	if kind == reflect.Slice {
		return
	}

	for i := 0; i < val.NumField(); i++ {

		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		index := typeField.Tag.Get("validate")

		appendKey := strings.Title(typeField.Name)

		if !valueField.CanInterface() {
			continue
		}

		field := valueField.Interface()
		fieldval := reflect.ValueOf(field)

		switch fieldval.Kind() {

		case reflect.Array, reflect.Slice, reflect.Struct:
			getValidationTagsRecursive(fieldval, keys, key+appendKey+".")

		default:
			if index != "" {
				keys[key+appendKey] = index
			}
		}
	}
}

func GetReflectionFieldValue(key string, x interface{}) string {

	splitKey := strings.Split(key, ".")

	propertyName := splitKey[0]
	originalPropertyName := propertyName

	arrayIndex := strings.Index(propertyName, "[")
	if arrayIndex != -1 {
		propertyName = propertyName[:arrayIndex]
	}

	val := reflect.ValueOf(x).Elem()

	for i := 0; i < val.NumField(); i++ {

		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		name := typeField.Name

		if name != propertyName {
			continue
		}

		if !valueField.CanInterface() {
			continue
		}

		f := valueField.Interface()
		val := reflect.ValueOf(f)

		switch val.Kind() {

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return strconv.FormatInt(val.Int(), 10)
		case reflect.Float32:
			return strconv.FormatFloat(val.Float(), 'E', -1, 32)
		case reflect.Float64:
			return strconv.FormatFloat(val.Float(), 'E', -1, 64)
		case reflect.String:
			return val.String()
		case reflect.Array, reflect.Slice:

			lastIndex := strings.Index(originalPropertyName, "]")
			index := extensions.StringToInt(originalPropertyName[arrayIndex+1 : lastIndex])
			arrayItem := val.Index(index)

			switch arrayItem.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return strconv.FormatInt(arrayItem.Int(), 10)
			case reflect.Float32:
				return strconv.FormatFloat(arrayItem.Float(), 'E', -1, 32)
			case reflect.Float64:
				return strconv.FormatFloat(arrayItem.Float(), 'E', -1, 64)
			case reflect.String:
				return arrayItem.String()

			case reflect.Struct:
				if len(splitKey) > 1 {
					return getStructReflectionValue(strings.Replace(key, originalPropertyName+".", "", 1), arrayItem)
				}
			}

		case reflect.Struct:
			if len(splitKey) > 1 {
				return getStructReflectionValue(strings.Replace(key, propertyName+".", "", 1), val)
			}
		}
	}
	return ""
}

func getStructReflectionValue(key string, val reflect.Value) string {

	splitKey := strings.Split(key, ".")

	propertyName := splitKey[0]
	originalPropertyName := propertyName

	arrayIndex := strings.Index(propertyName, "[")
	if arrayIndex != -1 {
		propertyName = propertyName[:arrayIndex]
	}

	for i := 0; i < val.NumField(); i++ {

		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		name := typeField.Name

		if name != propertyName {
			continue
		}

		if !valueField.CanInterface() {
			continue
		}

		f := valueField.Interface()
		val := reflect.ValueOf(f)

		switch val.Kind() {

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return strconv.FormatInt(val.Int(), 10)
		case reflect.Float32:
			return strconv.FormatFloat(val.Float(), 'E', -1, 32)
		case reflect.Float64:
			return strconv.FormatFloat(val.Float(), 'E', -1, 64)
		case reflect.String:
			return val.String()
		case reflect.Array, reflect.Slice:

			lastIndex := strings.Index(originalPropertyName, "]")
			index := extensions.StringToInt(originalPropertyName[arrayIndex+1 : lastIndex])
			arrayItem := val.Index(index)

			switch arrayItem.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return strconv.FormatInt(arrayItem.Int(), 10)
			case reflect.Float32:
				return strconv.FormatFloat(arrayItem.Float(), 'E', -1, 32)
			case reflect.Float64:
				return strconv.FormatFloat(arrayItem.Float(), 'E', -1, 64)
			case reflect.String:
				return arrayItem.String()

			case reflect.Struct:
				if len(splitKey) > 1 {
					return getStructReflectionValue(strings.Replace(key, originalPropertyName+".", "", 1), arrayItem)
				}
			}
		case reflect.Struct:
			if len(splitKey) > 1 {
				return getStructReflectionValue(strings.Replace(key, propertyName+".", "", 1), val)
			}

		}
	}
	return ""
}

func SetFieldValue(key string, val reflect.Value, value interface{}) {

	splitKey := strings.Split(key, ".")

	propertyName := splitKey[0]
	originalPropertyName := propertyName

	arrayIndex := strings.Index(propertyName, "[")
	if arrayIndex != -1 {
		propertyName = propertyName[:arrayIndex]
	}

	for i := 0; i < val.NumField(); i++ {

		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		name := typeField.Name

		if name != propertyName {
			continue
		}

		if !valueField.CanInterface() {
			continue
		}

		fieldToSet := val.Field(i)

		f := valueField.Interface()
		valType := reflect.ValueOf(f)

		switch valType.Kind() {

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fieldToSet.SetInt(value.(int64))
		case reflect.Float32:
			fieldToSet.SetFloat(value.(float64))
		case reflect.Float64:
			fieldToSet.SetFloat(value.(float64))
		case reflect.String:
			fieldToSet.SetString(value.(string))
		case reflect.Array, reflect.Slice:

			lastIndex := strings.Index(originalPropertyName, "]")
			index := extensions.StringToInt(originalPropertyName[arrayIndex+1 : lastIndex])
			arrayItem := fieldToSet.Index(index)

			switch arrayItem.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				arrayItem.SetInt(value.(int64))
			case reflect.Float32:
				arrayItem.SetFloat(value.(float64))
			case reflect.Float64:
				arrayItem.SetFloat(value.(float64))
			case reflect.String:
				arrayItem.SetString(value.(string))

			case reflect.Struct:
				if len(splitKey) > 1 {
					setStructReflectionValue(strings.Replace(key, originalPropertyName+".", "", 1), arrayItem, value)
				}
			}

		case reflect.Struct:
			if len(splitKey) > 1 {
				setStructReflectionValue(strings.Replace(key, propertyName+".", "", 1), fieldToSet, value)
			}
		}
	}

}

func setStructReflectionValue(key string, val reflect.Value, value interface{}) {

	splitKey := strings.Split(key, ".")

	propertyName := splitKey[0]
	originalPropertyName := propertyName

	arrayIndex := strings.Index(propertyName, "[")
	if arrayIndex != -1 {
		propertyName = propertyName[:arrayIndex]
	}

	for i := 0; i < val.NumField(); i++ {

		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		name := typeField.Name

		if name != propertyName {
			continue
		}

		if !valueField.CanInterface() {
			continue
		}

		fieldToSet := val.Field(i)
		f := valueField.Interface()
		valType := reflect.ValueOf(f)

		switch valType.Kind() {

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fieldToSet.SetInt(value.(int64))
		case reflect.Float32:
			fieldToSet.SetFloat(value.(float64))
		case reflect.Float64:
			fieldToSet.SetFloat(value.(float64))
		case reflect.String:
			fieldToSet.SetString(value.(string))
		case reflect.Array, reflect.Slice:

			lastIndex := strings.Index(originalPropertyName, "]")
			index := extensions.StringToInt(originalPropertyName[arrayIndex+1 : lastIndex])
			arrayItem := fieldToSet.Index(index)

			switch arrayItem.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				arrayItem.SetInt(value.(int64))
			case reflect.Float32:
				arrayItem.SetFloat(value.(float64))
			case reflect.Float64:
				arrayItem.SetFloat(value.(float64))
			case reflect.String:
				arrayItem.SetString(value.(string))

			case reflect.Struct:
				if len(splitKey) > 1 {
					setStructReflectionValue(strings.Replace(key, originalPropertyName+".", "", 1), arrayItem, value)
				}
			}
		case reflect.Struct:
			if len(splitKey) > 1 {
				setStructReflectionValue(strings.Replace(key, propertyName+".", "", 1), fieldToSet, value)
			}

		}
	}
}
