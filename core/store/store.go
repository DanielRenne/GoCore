//Package store provides a registry and interface to interact with a store repository against model entities.
package store

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/DanielRenne/GoCore/core/extensions"
)

const (
	WebSocketStoreKey = "WebSocket"
	PathAdd           = "Add"
	PathRemove        = "Remove"
)

//OnChange provides inserts, updates, and deletes to the store key.
var OnChange func(key string, id string, path string, x interface{}, err error)

//Get gets a collection entity by id.
func Get(key string, id string, joins []string) (x interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%+v", r)
			return
		}
	}()

	collection, ok := getRegistry(key)
	if !ok {
		return
	}

	obj, err := collection.ById(id, joins)
	if err != nil {
		return
	}

	x = obj.Elem().Interface()
	return
}

//GetByFilter gets a collection entity by filter.
func GetByFilter(key string, filter map[string]interface{}, inFilter map[string]interface{}, excludeFilter map[string]interface{}, joins []string) (x interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%+v", r)
			return
		}
	}()

	collection, ok := getRegistry(key)
	if !ok {
		return
	}

	obj, err := collection.ByFilter(filter, inFilter, excludeFilter, joins)
	if err != nil {
		return
	}

	x = obj.Elem().Interface()
	return
}

//GetByPath gets a collection entity-property value by id & path.
func GetByPath(key string, id string, joins []string, path string) (x interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%+v", r)
			return
		}
	}()

	collection, ok := getRegistry(key)
	if !ok {
		return
	}

	obj, err := collection.ById(id, joins)
	if err != nil {
		return
	}

	objElem := obj.Elem()

	fields := strings.Split(path, ".")
	depth := len(fields)

	properties := []reflect.Value{}

	for i := range fields {
		fieldName := fields[i]
		arrayIndex := -1

		if strings.Contains(fieldName, "[") {
			arraySplit := strings.Split(fieldName, "[")
			fieldName = arraySplit[0]
			arrayIndex = extensions.StringToInt(strings.Replace(arraySplit[1], "]", "", -1))
		}

		var fieldValue reflect.Value

		if i == 0 {
			fieldValue = objElem.FieldByName(fieldName)
		} else {
			fieldValue = properties[i-1].FieldByName(fieldName)
		}

		if arrayIndex == -1 {
			properties = append(properties, fieldValue)
		} else {
			properties = append(properties, fieldValue.Index(arrayIndex))
		}

		if i+1 == depth {
			x = properties[i].Interface()
		}
	}

	return
}

//Set updates a collection by id, path.
func Set(key string, id string, path string, x interface{}, logger func(string, string)) {

	defer func() {
		if r := recover(); r != nil {
			logger("Recover", fmt.Sprintf("%+v", r))
			if OnChange != nil {
				OnChange(key, id, path, x, fmt.Errorf("%+v", r))
			}
		}
	}()

	collection, ok := getRegistry(key)
	if !ok {
		return
	}

	obj, err := collection.ById(id, []string{})
	if err != nil {
		log.Printf("%s%s", "Error Getting Collection Object by id.  ", err.Error())
		return
	}

	objElem := obj.Elem()

	if path == "" {

		method := obj.MethodByName("ParseInterface")

		in := []reflect.Value{}
		in = append(in, reflect.ValueOf(x))
		values := method.Call(in)

		if values[0].Interface() != nil {
			err, ok := values[0].Interface().(error)
			if ok {
				logger("Error", err.Error())
				log.Printf("%s%+v\n", "Error Parsing Object.", err.Error())
				if OnChange != nil {
					OnChange(key, id, path, x, err)
				}
				return
			}
		}

	} else {
		fields := strings.Split(path, ".")
		depth := len(fields)

		properties := []reflect.Value{}

		for i := range fields {
			fieldName := fields[i]
			arrayIndex := -1

			if strings.Contains(fieldName, "[") {
				arraySplit := strings.Split(fieldName, "[")
				fieldName = arraySplit[0]
				arrayIndex = extensions.StringToInt(strings.Replace(arraySplit[1], "]", "", -1))
			}

			var fieldValue reflect.Value

			if i == 0 {
				fieldValue = objElem.FieldByName(fieldName)
			} else {
				fieldValue = properties[i-1].FieldByName(fieldName)
			}

			if arrayIndex == -1 {
				properties = append(properties, fieldValue)
			} else {
				properties = append(properties, fieldValue.Index(arrayIndex))
			}

			if i+1 == depth {
				if properties[i].CanSet() {
					propInterface := properties[i].Interface()
					if propInterface != nil {
						propType := reflect.TypeOf(propInterface).String()
						logger("PropType", propType)
						if propType == "int" {
							floatVal, ok := x.(float64)
							if ok {
								x = int(floatVal)
							}
						} else if propType == "float64" {
							intVal, ok := x.(int)
							if ok {
								x = float64(intVal)
							}
						}
					}

					// logger("Trying to set", fmt.Sprintf("%+v", x))

					valueToSet, err := collection.ReflectByFieldName(fieldName, x)
					if err != nil {
						logger("Error Setting Value to Store", fmt.Sprintf("%+v", valueToSet)+"\nError:  "+err.Error())
						OnChange(key, id, path, x, err)
					}

					// properties[i].Set(reflect.ValueOf(x))
					properties[i].Set(valueToSet)
				}
			}
		}
	}

	method := obj.MethodByName("Save")

	in := []reflect.Value{}
	values := method.Call(in)

	if values[0].Interface() != nil {
		err, ok := values[0].Interface().(error)
		if ok {
			logger("Error", err.Error())
			log.Printf("%s%+v\n", "Error Saving Object.", err.Error())
			if OnChange != nil {
				OnChange(key, id, path, x, err)
			}
			return
		}
	}

	if OnChange != nil {
		OnChange(key, id, path, x, nil)
	}

	// log.Printf("%s%+v\n", "UPdated Entity ", objElem)

}

//Append adds to an array field.
func Append(key string, id string, path string, x interface{}, logger func(string, string)) (y interface{}, err error) {

	defer func() {
		if r := recover(); r != nil {
			logger("Recover", fmt.Sprintf("%+v", r))
			if OnChange != nil {
				OnChange(key, id, path, x, fmt.Errorf("%+v", r))
			}
		}
	}()

	collection, ok := getRegistry(key)
	if !ok {
		return
	}

	obj, err := collection.ById(id, []string{})
	if err != nil {
		log.Printf("%s%s", "Error Getting Collection Object by id.  ", err.Error())
		return
	}

	objElem := obj.Elem()

	fields := strings.Split(path, ".")
	depth := len(fields)

	properties := []reflect.Value{}

	var updatedArray reflect.Value

	for i := range fields {
		fieldName := fields[i]
		arrayIndex := -1

		if strings.Contains(fieldName, "[") {
			arraySplit := strings.Split(fieldName, "[")
			fieldName = arraySplit[0]
			arrayIndex = extensions.StringToInt(strings.Replace(arraySplit[1], "]", "", -1))
		}

		var fieldValue reflect.Value

		if i == 0 {
			fieldValue = objElem.FieldByName(fieldName)
		} else {
			fieldValue = properties[i-1].FieldByName(fieldName)
		}

		if arrayIndex == -1 {
			properties = append(properties, fieldValue)
		} else {
			properties = append(properties, fieldValue.Index(arrayIndex))
		}

		if i+1 == depth {
			if properties[i].CanSet() {

				propType := reflect.TypeOf(properties[i].Interface()).String()

				if propType == "int" {
					floatVal, ok := x.(float64)
					if ok {
						x = int(floatVal)
					}
				} else if propType == "float64" {
					intVal, ok := x.(int)
					if ok {
						x = float64(intVal)
					}
				}

				valueToSet, err := collection.ReflectBaseTypeByFieldName(fieldName, x)
				if err != nil {
					logger("Error Setting Value to Store", fmt.Sprintf("%+v", valueToSet)+"\nError:  "+err.Error())
					OnChange(key, id, path, x, err)
				}

				properties[i].Set(reflect.Append(properties[i], valueToSet))
				updatedArray = properties[i]
			}
		}
	}

	method := obj.MethodByName("Save")

	in := []reflect.Value{}
	values := method.Call(in)

	if values[0].Interface() != nil {
		errInterface, ok := values[0].Interface().(error)
		if ok {
			logger("Error", errInterface.Error())
			log.Printf("%s%+v\n", "Error Saving Object.", errInterface.Error())
			if OnChange != nil {
				OnChange(key, id, path, x, errInterface)
			}
			err = errInterface
			return
		}
	}

	if OnChange != nil {
		OnChange(key, id, path, updatedArray.Interface(), nil)
	}

	y = updatedArray.Interface()
	return
}

//Add adds an entity to the collection and returns it.
func Add(key string, x interface{}, logger func(string, string)) (y interface{}, err error) {

	defer func() {
		if r := recover(); r != nil {
			logger("Recover", fmt.Sprintf("%+v", r))
			if OnChange != nil {
				OnChange(key, "", "", x, fmt.Errorf("%+v", r))
			}
		}
	}()

	collection, ok := getRegistry(key)
	if !ok {
		return
	}

	obj := collection.NewByReflection()

	if x != nil {
		data, errMarshal := json.Marshal(x)
		if errMarshal != nil {
			if OnChange != nil {
				OnChange(key, "", "", x, errMarshal)
			}
			return
		}
		errUnMarshal := json.Unmarshal(data, obj.Interface())
		if errUnMarshal != nil {
			if OnChange != nil {
				OnChange(key, "", "", x, errUnMarshal)
			}
			return
		}
	}

	method := obj.MethodByName("Save")

	in := []reflect.Value{}
	values := method.Call(in)

	if values[0].Interface() != nil {
		errSave, ok := values[0].Interface().(error)
		if ok {
			logger("Error", err.Error())
			log.Printf("%s%+v\n", "Error Saving Object.", err.Error())
			if OnChange != nil {
				OnChange(key, "", "", x, err)
			}
			err = errSave
			return
		}
	}

	methodGetID := obj.MethodByName("GetId")

	inID := []reflect.Value{}
	idValues := methodGetID.Call(inID)

	if OnChange != nil {
		OnChange(key, idValues[0].Interface().(string), "", obj.Interface(), nil)
	}
	y = obj.Interface()
	return
}

//Remove removes an entity from the collection and returns the Id Removed.
func Remove(key string, id string) (err error) {

	defer func() {
		if r := recover(); r != nil {
			if OnChange != nil {
				OnChange(key, "", "", nil, fmt.Errorf("%+v", r))
			}
		}
	}()

	collection, ok := getRegistry(key)
	if !ok {
		return
	}

	obj, err := collection.ById(id, []string{})
	if err != nil {
		return
	}

	methodGetID := obj.MethodByName("GetId")

	inID := []reflect.Value{}
	idValues := methodGetID.Call(inID)

	method := obj.MethodByName("Delete")
	in := []reflect.Value{}
	values := method.Call(in)

	if values[0].Interface() != nil {
		errSave, ok := values[0].Interface().(error)
		if ok {
			log.Printf("%s%+v\n", "Error Deleting Object.", err.Error())
			if OnChange != nil {
				OnChange(key, "", "", nil, err)
			}
			err = errSave
			return
		}
	}

	if OnChange != nil {
		OnChange(key, idValues[0].Interface().(string), "", nil, nil)
	}
	return
}
