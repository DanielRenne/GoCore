package br

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/DanielRenne/GoCore/core/app"
	"github.com/DanielRenne/GoCore/core/store"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/GoCore/core/utils"

)

const (
	storeKey = "Store"
)

type storeBr struct{}

type storePubPayload struct {
	Collection string      `json:"Collection"`
	ID         string      `json:"Id"`
	Path       string      `json:"Path"`
	Value      interface{} `json:"Value"`
	Error      string      `json:"Error"`
}

func init() {
	store.OnChange = Store.onStoreChange
}

func (storeBr) onStoreChange(key string, id string, path string, x interface{}, err error) {
	var vm storePubPayload
	vm.Collection = key
	vm.ID = id
	vm.Path = path
	vm.Value = x
	if err != nil {
		vm.Error = err.Error()
	}

	app.PublishWebSocketJSON(storeKey, vm)
}

//BroadcastError will broadcast the error for the path to display to the client.
func (storeBr) BroadcastError(key string, id string, path string, x interface{}, err error) {
	var vm storePubPayload
	vm.Collection = key
	vm.ID = id
	vm.Path = "Errors." + path
	vm.Value = x
	if err != nil {
		vm.Error = err.Error()
	}

	app.PublishWebSocketJSON(storeKey, vm)
}

//Broadcast will broadcast the store value
func (storeBr) Broadcast(key string, id string, path string, x interface{}) {
	var vm storePubPayload
	vm.Collection = key
	vm.ID = id
	vm.Path = path
	vm.Value = x

	app.PublishWebSocketJSON(storeKey, vm)
}

func (storeBr) Execute(path string, x interface{}, raw []byte) (y interface{}) {
	defer func() {
		if r := recover(); r != nil {
			session_functions.Log("Error", "Panic at br->store.go->Execute:  "+fmt.Sprintf("%+v", r)+"  \nPath:  "+path+"\nValue:  "+fmt.Sprintf("%+v", x))
			return
		}
	}()

	// session_functions.VelocityLog("br.Store.Exectute", fmt.Sprintf("path: %s, x: %+v, raw: % x", path, x, raw))

	pathItems := strings.Split(path, ".")
	brType := pathItems[0]
	functionCall := pathItems[1]

	brObj, ok := GetBr(brType)
	if !ok {
		utils.TalkDirtyToMe("Unable to locate BR type: " + brType)
		return
	}
	method := brObj.MethodByName(functionCall)
	if !method.IsValid() {
		utils.TalkDirtyToMe("BR method not registered.  Check init exists.")
		return
	}

	in := []reflect.Value{}

	methodType := method.Type()
	paramCnt := methodType.NumIn()
	if paramCnt != 1 {
		utils.TalkDirtyToMe(fmt.Sprintf("BR method %s expects %d parameters only 1 is supported.", functionCall, paramCnt))
		return
	}
	paramType := methodType.In(0)

	genericType := reflect.TypeOf((*interface{})(nil))
	if paramType == genericType {
		if x != nil {
			in = append(in, reflect.ValueOf(x))
		}
	} else {
		param := reflect.New(paramType)
		err1 := json.Unmarshal(raw, param.Interface())
		if err1 != nil {
			session_functions.Log("br.Store.Execute", "Error parsing param: "+err1.Error())
			return
		}

		in = append(in, param.Elem())
	}

	value := method.Call(in)
	if len(value) > 0 {
		y = value[0].Interface()
	}

	return
}
