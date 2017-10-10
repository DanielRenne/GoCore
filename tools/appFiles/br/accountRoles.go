package br

import (
	"fmt"

	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
	"github.com/pkg/errors"
)

type accountRolesBr struct{}

func (self accountRolesBr) Create(context session_functions.RequestContext, vm *viewModel.AccountRoleModifyViewModel, t *model.Transaction) (message string, err error) {

	message, err = self.Validate(vm)
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}

	err = vm.AccountRole.SaveWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Log("Error->Validation_Error->accountRolesBr->Create", fmt.Sprintf("%+v", vm.AccountRole.Errors))
			message = constants.PARAM_SNACKBAR_MESSAGE_NONE
			return
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Log("Error->accountRolesBr->Create", err.Error())
			message = queries.AppContent.GetTranslation(context, constants.ERRORS_ACCOUNTROLE_SAVE) + core.Debug.HandleError(err)
			return
		}
	}

	message = constants.ACCOUNTROLE_SAVE_SUCCESS
	return
}

func (self accountRolesBr) Update(context session_functions.RequestContext, vm *viewModel.AccountRoleModifyViewModel, t *model.Transaction) (message string, err error) {

	message, err = self.Validate(vm)
	if err != nil {
		return
	}

	err = vm.AccountRole.SaveWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Log("Error->Validation_Error->accountRolesBr->Update", fmt.Sprintf("%+v", vm.AccountRole.Errors))
			message = constants.PARAM_SNACKBAR_MESSAGE_NONE
			return
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Log("Error->accountRolesBr->Update", err.Error())
			message = queries.AppContent.GetTranslation(context, constants.ERRORS_ACCOUNTROLE_SAVE) + core.Debug.HandleError(err)
			return
		}
	}

	message = constants.ACCOUNTROLE_SAVE_SUCCESS
	return
}

func (self accountRolesBr) Delete(context session_functions.RequestContext, vm *viewModel.AccountRoleModifyViewModel, t *model.Transaction) (message string, err error) {

	err = vm.AccountRole.DeleteWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Log("Error->Validation_Error->accountRolesBr->Delete", fmt.Sprintf("%+v", vm.AccountRole.Errors))
			message = constants.ERRORS_ACCOUNTROLE_DELETE
			return
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Log("Error->accountRolesBr->Delete", err.Error())
			message = constants.ERRORS_ACCOUNTROLE_DELETE
			return
		}
	}

	message = constants.ACCOUNTROLE_DELETE_SUCCESSFUL
	return
}

func (self accountRolesBr) WrapViewModel(AccountRoleInstance model.AccountRole) viewModel.AccountRoleModifyViewModel {
	vm := viewModel.AccountRoleModifyViewModel{}
	vm.AccountRole = AccountRoleInstance
	return vm
}

func (self accountRolesBr) Validate(vm *viewModel.AccountRoleModifyViewModel) (message string, err error) {
	message = ""
	return
}
