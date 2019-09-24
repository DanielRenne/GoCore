package controllers

import (
	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"

	"github.com/DanielRenne/GoCore/core/utils"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
	"github.com/pkg/errors"
)

func addEditAppErrorValidateCommon(context session_functions.RequestContext, vm *viewModel.AppErrorModifyViewModel, t *model.Transaction) (bool, session_functions.ServerResponseStruct) {
	//var err error

	// Add custom logic here and return false when theres an error
	//	Like this:
	//return false, session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, constants..ERRORS_APPERROR_SAVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)

	return true, session_functions.ServerResponseStruct{}
}

func vmAppErrorAddEditGetCommon(context session_functions.RequestContext, respond session_functions.ServerResponse, vm *viewModel.AppErrorModifyViewModel) bool {
	// Use this method as a common way to load up any VM defaults that both the add and edit pages need!

	//companies, err := queries.Companies.ByAccountWithContext(context)

	//if err != nil {
	//	respond(PARAM_REDIRECT_NONE, constants.ERRORS_ACCOUNT_NOT_FOUND, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
	//	return false
	//}

	return true
}

func WrapAppErrorViewModel(AppErrorInstance model.AppError) viewModel.AppErrorModifyViewModel {
	vm := viewModel.AppErrorModifyViewModel{}
	vm.AppError = AppErrorInstance
	return vm
}

func AppErrorPostCommitHook(actionPerformed string, context session_functions.RequestContext, id string) bool {
	if utils.InArray(actionPerformed, utils.Array("CreateAppError", "CopyAppError", "UpdateAppError")) && id != "" {
		var row model.AppError
		err := model.AppErrors.Query().ById(id, &row)
		if err == nil {
			t, err := session_functions.StartTransaction(context())
			if err == nil {
				// Add custom business rules to clean up data based on relationships or changes here

				//if row.Joins.Buildings.Count > 1 {
				//	row.IsCampus = true
				//} else {
				//	row.IsCampus = false
				//}

				err = row.SaveWithTran(t)
				if err == nil {
					err = t.Commit()
					if err == nil {
						return true
					}
				}
			} else {
				session_functions.Dump("Desc->Error in LocationEntityPostSaveHook", err)
			}
		} else {
			session_functions.Dump("Desc->Error in LocationEntityPostSaveHook 29", err)
		}
		return false
	} else {
		return true
	}
}

func CopyAppErrorRow(context session_functions.RequestContext, copyRowVm viewModel.AppErrorModifyViewModel, t *model.Transaction) session_functions.ServerResponseStruct {
	var err error
	copyVm := WrapAppErrorViewModel(copyRowVm.AppError)
	err = model.AppErrors.Query().ById(copyVm.AppError.Id, &copyVm.AppError)
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, constants.ERRORS_APPERROR_COPY, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, copyVm)
	}
	copyVm.AppError.Id = ""

	r := CreateAppErrorRow(context, &copyVm, t)

	if !r.CompletedSuccessfully {
		return session_functions.ServerResponseToStruct(r.CompletedSuccessfully, r.Context, r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
	}

	return session_functions.ServerResponseToStruct(PARAM_SUCCESS, context, PARAM_REDIRECT_RERENDER, constants.APPERROR_COPY_SUCCESSFUL, PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), copyVm)
}

func UpdateAppErrorRow(context session_functions.RequestContext, vm *viewModel.AppErrorModifyViewModel, t *model.Transaction) session_functions.ServerResponseStruct {
	var err error

	success, r := addEditAppErrorValidateCommon(context, vm, t)

	if !success {
		return r
	}

	err = vm.AppError.SaveWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Validation Error exists", vm.AppError.Errors)
			return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, queries.AppContent.GetTranslation(context, constants.ERRORS_APPERROR_SAVE)+core.Debug.HandleError(err), PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Golang error", err)
			return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, queries.AppContent.GetTranslation(context, constants.ERRORS_APPERROR_SAVE)+core.Debug.HandleError(err), PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		}
	}

	return session_functions.ServerResponseToStruct(PARAM_SUCCESS, context, PARAM_REDIRECT_BACK, constants.APPERROR_SAVE_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, "", viewModel.EmptyViewModel{})
}

func CreateAppErrorRow(context session_functions.RequestContext, vm *viewModel.AppErrorModifyViewModel, t *model.Transaction) session_functions.ServerResponseStruct {
	var err error

	success, r := addEditAppErrorValidateCommon(context, vm, t)

	if !success {
		return r
	}

	err = vm.AppError.SaveWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Validation Error exists", vm.AppError.Errors)
			return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, queries.AppContent.GetTranslation(context, constants.ERRORS_APPERROR_SAVE)+core.Debug.HandleError(err), PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Golang error", err)
			return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, queries.AppContent.GetTranslation(context, constants.ERRORS_APPERROR_SAVE)+core.Debug.HandleError(err), PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		}
	}

	return session_functions.ServerResponseToStruct(PARAM_SUCCESS, context, PARAM_REDIRECT_BACK, constants.APPERROR_SAVE_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, "", viewModel.EmptyViewModel{})

}

func DeleteAppErrorRow(context session_functions.RequestContext, vm *viewModel.AppErrorModifyViewModel, t *model.Transaction) session_functions.ServerResponseStruct {
	var err error

	err = vm.AppError.DeleteWithTran(t)

	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, constants.ERRORS_APPERROR_DELETE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
	}

	return session_functions.ServerResponseToStruct(PARAM_SUCCESS, context, PARAM_REDIRECT_RERENDER, constants.APPERROR_DELETE_SUCCESSFUL, PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), vm)
}
