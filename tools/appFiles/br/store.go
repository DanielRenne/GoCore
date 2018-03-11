package br

import (
	"github.com/DanielRenne/GoCore/core/app"
	"github.com/DanielRenne/GoCore/core/store"
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
