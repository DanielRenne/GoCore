package dbServices

import (
	"github.com/DanielRenne/GoCore/core/extensions"
	"reflect"
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

		appendKey := extensions.MakeFirstLowerCase(typeField.Name)

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
