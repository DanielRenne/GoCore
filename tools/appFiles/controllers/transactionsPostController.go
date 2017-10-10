package controllers

import (
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
)

func (self *TransactionsController) UpdateTransactionDetails(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.TransactionModifyViewModel
	var isValid = false

	vm, isValid = checkTransactionRequiredFields(state)
	if isValid != true {
		respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	err := vm.Transaction.Save()

	if err != nil {
		if model.IsValidationError(err) {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_SAVE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		} else {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_SAVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		}
		return
	}

	respond(PARAM_REDIRECT_BACK, constants.TRANSACTION_SAVE_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, "", vm)
}

func (self *TransactionsController) CreateTransaction(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.TransactionModifyViewModel

	vm, isValid := checkTransactionRequiredFields(state)
	if isValid != true {
		respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	err := vm.Transaction.Save()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_SAVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	respond(PARAM_REDIRECT_BACK, constants.TRANSACTION_SAVE_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, "", vm)
}

func (self *TransactionsController) DeleteTransaction(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.TransactionModifyViewModel
	vm.Parse(state)

	err := vm.Transaction.Delete()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_DELETE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	respond(PARAM_REDIRECT_BACK, constants.TRANSACTION_DELETE_SUCCESSFUL, PARAM_SNACKBAR_TYPE_SUCCESS, nil, "", vm)
}

func checkTransactionRequiredFields(state string) (vm viewModel.TransactionModifyViewModel, isValid bool) {
	vm.Parse(state)

	valid := constants.BITWISE_TRUE
	isValid = false

	if valid == constants.BITWISE_TRUE {
		isValid = true
	}
	return vm, isValid
}
