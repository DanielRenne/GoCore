package store

import (
	"reflect"
	"strings"
	"sync"

	"github.com/globalsign/mgo"
)

var registry sync.Map
var registryHistory sync.Map
var Version string

type collectionStore interface {
	ById(objectID interface{}, joins []string) (value reflect.Value, err error)
	ByFilter(filter map[string]interface{}, inFilter map[string]interface{}, excludeFilter map[string]interface{}, joins []string) (value reflect.Value, err error)
	CountByFilter(filter map[string]interface{}, inFilter map[string]interface{}, excludeFilter map[string]interface{}, joins []string) (count int, err error)
	ReflectByFieldName(fieldName string, x interface{}) (value reflect.Value, err error)
	ReflectBaseTypeByFieldName(fieldName string, x interface{}) (value reflect.Value, err error)
	NewByReflection() (value reflect.Value)

	SetCollection(mdb *mgo.Database)
	Bootstrap() error
	BootStrapComplete()
	Index() error
}

type collectionHistoryStore interface {
	SetCollection(mdb *mgo.Database)
	Index() error
}

//RegisterStore will register a new store to the store registry.
func RegisterStore(x interface{}) {
	key := strings.Replace(getType(x), "model", "", -1)
	registry.Store(key, x)
}

//GetCollection will return the collection by key string.
func GetCollection(key string) (x collectionStore, ok bool) {
	obj, ok := registry.Load(key)
	if ok {
		x = obj.(collectionStore)
	}
	return
}

//RegisterStore will register a new store to the store registry.
func RegisterHistoryStore(x interface{}) {
	key := strings.Replace(getType(x), "model", "", -1)
	registryHistory.Store(key, x)
}

//GetCollection will return the collection by key string.
func GetCollectionHistory(key string) (x collectionHistoryStore, ok bool) {
	obj, ok := registryHistory.Load(key)
	if ok {
		x = obj.(collectionHistoryStore)
	}
	return
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
