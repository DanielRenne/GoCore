package store

import (
	"reflect"
	"strings"
	"sync"
)

var registry sync.Map

type collectionStore interface {
	ById(objectID interface{}) (value reflect.Value, err error)
}

//RegisterStore will register a new store to the store registry.
func RegisterStore(x interface{}) {
	key := strings.Replace(getType(x), "model", "", -1)
	registry.Store(key, x)
}

func getRegistry(key string) (x collectionStore, ok bool) {

	obj, ok := registry.Load(key)
	if ok {
		x = obj.(collectionStore)
	}
	return
}

func getType(myvar interface{}) string {

	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}
