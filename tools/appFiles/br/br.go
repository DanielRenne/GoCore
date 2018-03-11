package br

import (
	"reflect"
	"sync"

	"github.com/pkg/errors"
)

var registry sync.Map

func ValidationError(msg string, errInfo error) (message string, err error) {
	err = errors.Wrap(errInfo, msg)
	message = msg
	return
}

var FileObjects fileObjectsBr
var Users usersBr
var Passwords passwordsBr
var AccountRoles accountRolesBr
var Server Server_Br
var Schedules schedulesBr
var Store storeBr

func GetBr(key string) (brObject reflect.Value, ok bool) {

	brObject, ok = getRegistry(key)
	return
}

func RegisterBr(brObj interface{}) {
	registry.Store(getType(brObj), reflect.ValueOf(brObj))
}

func RegisterBrByKey(key string, brObj interface{}) {
	registry.Store(key, reflect.ValueOf(brObj))
}

func getRegistry(key string) (brObj reflect.Value, ok bool) {

	obj, ok := registry.Load(key)
	if ok {
		brObj = obj.(reflect.Value)
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

//NewBrVarsDontDeleteMe
