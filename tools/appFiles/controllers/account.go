package controllers

import (
	"strings"

	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
	"github.com/pkg/errors"
)

func addEditAccountValidateCommon(context session_functions.RequestContext, vm *viewModel.AccountModifyViewModel, t *model.Transaction) (bool, session_functions.ServerResponseStruct) {
	return true, session_functions.ServerResponseStruct{}
}

func constructStatesMap(states []model.State, countries []model.Country) (result map[string][]model.State) {
	result = map[string][]model.State{}
	mapCountry := map[string]string{}

	for _, country := range countries {
		mapCountry[country.Iso] = country.Id.Hex()
	}

	for _, state := range states {
		country := mapCountry[strings.ToLower(state.Country)]
		result[country] = append(result[country], state)
	}
	return
}

func vmAccountAddEditGetCommon(context session_functions.RequestContext, respond session_functions.ServerResponse, vm *viewModel.AccountModifyViewModel) bool {

	var countries []model.Country
	err := model.Countries.Query().All(&countries)
	vm.Countries = countries

	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_COUNTRY, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return false
	}

	return true
}

func UpdateAccountRow(context session_functions.RequestContext, vm *viewModel.AccountModifyViewModel, t *model.Transaction) session_functions.ServerResponseStruct {
	var err error

	success, r := addEditAccountValidateCommon(context, vm, t)

	if !success {
		return r
	}

	err = vm.Account.SaveWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump(vm.Account.Errors)
			return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, queries.AppContent.GetTranslation(context, constants.ERRORS_ACCOUNT_SAVE)+core.Debug.HandleError(err), PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		}
	}

	return session_functions.ServerResponseToStruct(PARAM_SUCCESS, context, PARAM_REDIRECT_BACK, constants.ACCOUNT_SAVE_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, "", viewModel.EmptyViewModel{})
}

func CreateAccountRow(context session_functions.RequestContext, vm *viewModel.AccountModifyViewModel, t *model.Transaction) session_functions.ServerResponseStruct {
	success, r := addEditAccountValidateCommon(context, vm, t)
	if !success {
		return r
	}
	acct, err := session_functions.GetSessionAccount(context())
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, constants.ERRORS_ACCOUNT_NOT_FOUND, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
	}

	vm.Account.RelatedAcctId = acct.Id.Hex()

	err = vm.Account.SaveWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump(vm.Account.Errors)
			return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, queries.AppContent.GetTranslation(context, constants.ERRORS_ACCOUNT_SAVE)+core.Debug.HandleError(err), PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		}
	}

	return session_functions.ServerResponseToStruct(PARAM_SUCCESS, context, PARAM_REDIRECT_BACK, constants.ACCOUNT_SAVE_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, "", vm)
}

func DeleteAccountRow(context session_functions.RequestContext, vm *viewModel.AccountModifyViewModel, t *model.Transaction) session_functions.ServerResponseStruct {
	var err error

	err = vm.Account.DeleteWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, constants.ERRORS_ACCOUNT_DELETE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, constants.ERRORS_ACCOUNT_DELETE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		}
	}

	return session_functions.ServerResponseToStruct(PARAM_SUCCESS, context, PARAM_REDIRECT_BACK, constants.ACCOUNT_DELETE_SUCCESSFUL, PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), vm)
}
