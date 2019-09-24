//Package store provides an API into the DB GoCore Store operations.
package store

import (
	"encoding/json"
	"fmt"

	coreStore "github.com/DanielRenne/GoCore/core/store"
	"github.com/DanielRenne/goCoreAppTemplate/br"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/controllers"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
)

//StoreController is an api controller for interacting with the GoCore Data Store.
type StoreController struct{}

func init() {
	controllers.RegisterController(&StoreController{})
}

type storePostPayload struct {
	Collection    string                 `json:"Collection"`
	ID            string                 `json:"Id"`
	Filter        map[string]interface{} `json:"Filter"`
	InFilter      map[string]interface{} `json:"InFilter"`
	ExcludeFilter map[string]interface{} `json:"ExcludeFilter"`
	Joins         []string               `json:"Joins"`
	Path          string                 `json:"Path"`
	ValueRaw      json.RawMessage        `json:"Value"`
	Value         interface{}
}

//Parse will parse the json datastring for the GetRoomDevices.
func (obj *storePostPayload) Parse(data string) {
	json.Unmarshal([]byte(data), &obj)
	json.Unmarshal(obj.ValueRaw, &obj.Value)
}

//Get returns an entity from a collection store.
func (sc *StoreController) Get(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm storePostPayload
	vm.Parse(state)

	x, err := coreStore.Get(vm.Collection, vm.ID, vm.Joins)
	if err != nil {
		message := "Failed to retrieve entity.  Get->Collection:  " + vm.Collection + "\r\n  Id:  " + vm.ID + "\r\n" + fmt.Sprintf("%+v", vm)
		respond(constants.PARAM_REDIRECT_NONE, message, constants.PARAM_SNACKBAR_TYPE_ERROR, err, constants.PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
		return
	}

	respond(constants.PARAM_REDIRECT_NONE, constants.PARAM_SNACKBAR_MESSAGE_NONE, constants.PARAM_SNACKBAR_TYPE_SUCCESS, nil, constants.PARAM_TRANSACTION_ID_NONE, x)
}

//GetByPath returns an entity field value from a collection store.
func (sc *StoreController) GetByPath(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm storePostPayload
	vm.Parse(state)

	if vm.Collection == "GoCore" {
		y := sc.processGoCore(vm)
		respond(constants.PARAM_REDIRECT_NONE, constants.PARAM_SNACKBAR_MESSAGE_NONE, constants.PARAM_SNACKBAR_TYPE_ERROR, nil, constants.PARAM_TRANSACTION_ID_NONE, y)
		return
	}

	x, err := coreStore.GetByPath(vm.Collection, vm.ID, vm.Joins, vm.Path)
	if err != nil {
		message := "Failed to retrieve entity.  GetByPath->Collection:  " + vm.Collection + "\r\n  Id:  " + vm.ID + "\r\n" + fmt.Sprintf("%+v", vm)
		respond(constants.PARAM_REDIRECT_NONE, message, constants.PARAM_SNACKBAR_TYPE_ERROR, err, constants.PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
		return
	}

	respond(constants.PARAM_REDIRECT_NONE, constants.PARAM_SNACKBAR_MESSAGE_NONE, constants.PARAM_SNACKBAR_TYPE_SUCCESS, nil, constants.PARAM_TRANSACTION_ID_NONE, x)
}

//GetByFilter returns an array of entities from a collection store.
func (*StoreController) GetByFilter(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm storePostPayload
	vm.Parse(state)

	x, err := coreStore.GetByFilter(vm.Collection, vm.Filter, vm.InFilter, vm.ExcludeFilter, vm.Joins)
	if err != nil {

		message := "Failed to retrieve entity.  GetByFilter->Collection:  " + vm.Collection + "\r\n" + fmt.Sprintf("%+v", vm)
		respond(constants.PARAM_REDIRECT_NONE, message, constants.PARAM_SNACKBAR_TYPE_ERROR, err, constants.PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
		return
	}

	respond(constants.PARAM_REDIRECT_NONE, constants.PARAM_SNACKBAR_MESSAGE_NONE, constants.PARAM_SNACKBAR_TYPE_SUCCESS, nil, constants.PARAM_TRANSACTION_ID_NONE, x)
}

//Remove removes the document from the collection.
func (*StoreController) Remove(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm storePostPayload
	vm.Parse(state)

	err := coreStore.Remove(vm.Collection, vm.ID)
	if err != nil {
		respond(constants.PARAM_REDIRECT_NONE, "Failed to remove entity.", constants.PARAM_SNACKBAR_TYPE_ERROR, err, constants.PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
		return
	}

	respond(constants.PARAM_REDIRECT_NONE, constants.PARAM_SNACKBAR_MESSAGE_NONE, constants.PARAM_SNACKBAR_TYPE_SUCCESS, nil, constants.PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
}

//Set sets an entity field in the collection store.
func (sc *StoreController) Set(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm storePostPayload
	vm.Parse(state)

	if vm.Collection == "GoCore" {
		y := sc.processGoCore(vm)
		respond(constants.PARAM_REDIRECT_NONE, constants.PARAM_SNACKBAR_MESSAGE_NONE, constants.PARAM_SNACKBAR_TYPE_ERROR, nil, constants.PARAM_TRANSACTION_ID_NONE, y)
		return
	}

	coreStore.Set(vm.Collection, vm.ID, vm.Path, vm.Value, session_functions.Log)

	respond(constants.PARAM_REDIRECT_NONE, constants.PARAM_SNACKBAR_MESSAGE_NONE, constants.PARAM_SNACKBAR_TYPE_ERROR, nil, constants.PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
}

//Publish will fetch the store record and publish to all subscribers.
func (sc *StoreController) Publish(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm storePostPayload
	vm.Parse(state)

	coreStore.Publish(vm.Collection, vm.ID, vm.Path, session_functions.Log)
	respond(constants.PARAM_REDIRECT_NONE, constants.PARAM_SNACKBAR_MESSAGE_NONE, constants.PARAM_SNACKBAR_TYPE_ERROR, nil, constants.PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
}

//Add creates a new collection object to the collection
func (*StoreController) Add(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm storePostPayload
	vm.Parse(state)

	x, _ := coreStore.Add(vm.Collection, vm.Value, session_functions.Log)

	respond(constants.PARAM_REDIRECT_NONE, constants.PARAM_SNACKBAR_MESSAGE_NONE, constants.PARAM_SNACKBAR_TYPE_ERROR, nil, constants.PARAM_TRANSACTION_ID_NONE, x)
}

//Append adds a new row object to a collection path
func (*StoreController) Append(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm storePostPayload
	vm.Parse(state)

	x, _ := coreStore.Append(vm.Collection, vm.ID, vm.Path, vm.Value, session_functions.Log)

	respond(constants.PARAM_REDIRECT_NONE, constants.PARAM_SNACKBAR_MESSAGE_NONE, constants.PARAM_SNACKBAR_TYPE_ERROR, nil, constants.PARAM_TRANSACTION_ID_NONE, x)
}

func (*StoreController) processGoCore(vm storePostPayload) (y interface{}) {

	defer func() {
		if r := recover(); r != nil {
			session_functions.Log("Error", "Panic at store.go processGoCore:  "+fmt.Sprintf("%+v", r)+"  \nPath:  "+vm.Path+"\nValue:  "+fmt.Sprintf("%+v", vm.Value))
			return
		}
	}()

	y = br.Store.Execute(vm.Path, vm.Value, vm.ValueRaw)
	return
}
