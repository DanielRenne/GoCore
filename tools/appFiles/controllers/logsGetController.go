package controllers

import (
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
)

func (self *LogsController) Root(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {

	var vm viewModel.LogsViewModel
	vm.LoadDefaultState()
	if !session_functions.CheckRoleAccess(context(), constants.FEATURE_LOGS_VIEW) {
		respond(PARAM_REDIRECT_NONE, "NoAccess", PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	val, ok := uriParams["Id"]
	if ok {
		vm.Id = val
	} else {
		respond(PARAM_REDIRECT_NONE, "Id invalid", PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)

}
