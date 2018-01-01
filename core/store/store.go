//Package store provides a registry and interface to interact with a store repository against model entities.
package store

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/DanielRenne/GoCore/core/extensions"
)

//OnChange provides inserts, updates, and deletes to the store key.
var OnChange func(key string, id string, path string, x interface{}, err error)

//Get gets a collection entity by id, path.
func Get(key string, id string) (x interface{}, err error) {
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

	obj, err := collection.ById(id)
	if err != nil {
		return
	}

	x = obj.Elem().Interface()
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

	obj, err := collection.ById(id)
	if err != nil {
		log.Printf("%s%s", "Error Getting Collection Object by id.  ", err.Error())
		return
	}

	objElem := obj.Elem()

	// log.Printf("%s%+v\n", "Original Entity ", objElem)

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
				log.Printf("%s%+v", "Set the Object ", x)
				logger("Setting Field", fmt.Sprintf("%+v", x))

				propType := reflect.TypeOf(properties[i].Interface()).String()
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

				properties[i].Set(reflect.ValueOf(x))
				logger("Done Setting Field", fmt.Sprintf("%+v", x))
			}
		}
	}

	method := obj.MethodByName("Save")

	in := []reflect.Value{}
	logger("Calling Call", "")
	values := method.Call(in)
	logger("Saved", "")
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
