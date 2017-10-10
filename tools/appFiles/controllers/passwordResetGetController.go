package controllers

import (
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
)

func (self *PasswordResetController) Root(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {

	var vm viewModel.PasswordResetViewModel
	vm.LoadDefaultState()
	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)

}

func (self *PasswordResetController) Load(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {

	var vm viewModel.PasswordResetViewModel
	vm.LoadDefaultState()

	var err error
	idParam, ok := uriParams["Id"]

	if ok == false {
		respond(PARAM_REDIRECT_NONE, "Password Reset Record Missing.", SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	vm.PasswordReset, err = queries.PasswordResets.ById(idParam)

	if err != nil {
		respond(PARAM_REDIRECT_NONE, "PasswordReset Not Found", SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)

}
