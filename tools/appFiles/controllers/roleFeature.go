package controllers

import (
	"encoding/json"
	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/settings"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
	"github.com/pkg/errors"
	"io/ioutil"
)

func addEditRoleFeatureValidateCommon(context session_functions.RequestContext, vm *viewModel.RoleFeatureModifyViewModel, t *model.Transaction) (bool, session_functions.ServerResponseStruct) {
	//var err error

	// Add custom logic here and return false when theres an error
	//	Like this:
	//return false, session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, constants..ERRORS_ROLEFEATURE_SAVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)

	return true, session_functions.ServerResponseStruct{}
}

func vmRoleFeatureAddEditGetCommon(context session_functions.RequestContext, respond session_functions.ServerResponse, vm *viewModel.RoleFeatureModifyViewModel) bool {
	// Use this method as a common way to load up any VM defaults that both the add and edit pages need!

	//companies, err := queries.Companies.ByAccountWithContext(context)

	//if err != nil {
	//	respond(PARAM_REDIRECT_NONE, constants.ERRORS_ACCOUNT_NOT_FOUND, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
	//	return false
	//}

	//if !vm.RoleFeature.XXXXFIELDXXXX {
	//	vm.RoleFeature.Errors.XXXXFIELDXXXX = constants.ERROR_REQUIRED_FIELD
	//	return false, session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
	//}

	return true
}

func WrapRoleFeatureViewModel(RoleFeatureInstance model.RoleFeature) viewModel.RoleFeatureModifyViewModel {
	vm := viewModel.RoleFeatureModifyViewModel{}
	vm.RoleFeature = RoleFeatureInstance
	return vm
}

func RoleFeaturePostCommitHook(actionPerformed string, context session_functions.RequestContext, id string) (ret bool) {
	if settings.AppSettings.DeveloperMode {
		var all []model.RoleFeature
		err := model.RoleFeatures.Query().All(&all)
		if err != nil {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			return
		}
		for i, _ := range all {
			all[i].BootstrapMeta = &model.BootstrapMeta{
				AlwaysUpdate: true,
			}
			all[i].LastUpdateId = constants.APP_CONSTANTS_USERS_ANONYMOUS_ID
		}
		strjson, err := json.MarshalIndent(all, "", "\t")
		err = ioutil.WriteFile(serverSettings.APP_LOCATION+"/db/bootstrap/roleFeatures/roleFeatures.json", []byte(strjson), 0644)
		if err != nil {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			return
		}
		ret = true
		return ret
	}
	return ret
}

func RoleFfffffeaturePostCommitHook(actionPerformed string, context session_functions.RequestContext, id string) (ret bool) {
	//if settings.AppSettings.DeveloperMode {
	//	var all []model.RoleFeature
	//	err := model.RoleFeatures.Query().All(&all)
	//	if err != nil {
	//		return
	//	}
	//	strjson, err := json.MarshalIndent(all, "", "\t")
	//	err = ioutil.WriteFile(serverSettings.APP_LOCATION+"/db/bootstrap/roleFeatures/roleFeatures.json", []byte(strjson), 0644)
	//	if err != nil {
	//		return
	//	}
	//	ret = true
	//	return ret
	//}
	ret = true
	return ret
}

func CopyRoleFeatureRow(context session_functions.RequestContext, copyRowVm *viewModel.RoleFeatureModifyViewModel, t *model.Transaction) session_functions.ServerResponseStruct {
	var err error
	err = model.RoleFeatures.Query().ById(copyRowVm.RoleFeature.Id, &copyRowVm.RoleFeature)
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return session_functions.ServerResponseToStruct(PARAM_FAILED, context, queries.AppContent.GetTranslation(context, constants.ERRORS_ROLEFEATURE_COPY)+core.Debug.HandleError(err), constants.ERRORS_ROLEFEATURE_COPY, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, copyRowVm)
	}
	copyRowVm.RoleFeature.Id = ""

	r := CreateRoleFeatureRow(context, copyRowVm, t)

	if !r.CompletedSuccessfully {
		return session_functions.ServerResponseToStruct(r.CompletedSuccessfully, r.Context, r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
	}

	return session_functions.ServerResponseToStruct(PARAM_SUCCESS, context, PARAM_REDIRECT_RERENDER, queries.AppContent.GetTranslation(context, constants.ROLEFEATURE_COPY_SUCCESSFUL), PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), copyRowVm)
}

func UpdateRoleFeatureRow(context session_functions.RequestContext, vm *viewModel.RoleFeatureModifyViewModel, t *model.Transaction) session_functions.ServerResponseStruct {
	var err error

	success, r := addEditRoleFeatureValidateCommon(context, vm, t)

	if !success {
		return r
	}

	err = vm.RoleFeature.SaveWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Validation Error exists", vm.RoleFeature.Errors)
			return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Golang error", err)
			return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, queries.AppContent.GetTranslation(context, constants.ERRORS_ROLEFEATURE_SAVE)+core.Debug.HandleError(err), PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		}
	}

	return session_functions.ServerResponseToStruct(PARAM_SUCCESS, context, PARAM_REDIRECT_BACK, constants.ROLEFEATURE_SAVE_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, "", vm)
}

func CreateRoleFeatureRow(context session_functions.RequestContext, vm *viewModel.RoleFeatureModifyViewModel, t *model.Transaction) session_functions.ServerResponseStruct {
	var err error

	success, r := addEditRoleFeatureValidateCommon(context, vm, t)

	if !success {
		return r
	}

	err = vm.RoleFeature.SaveWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Validation Error exists", vm.RoleFeature.Errors)
			return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Golang error", err)
			return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, queries.AppContent.GetTranslation(context, constants.ERRORS_ROLEFEATURE_SAVE)+core.Debug.HandleError(err), PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		}
	}

	return session_functions.ServerResponseToStruct(PARAM_SUCCESS, context, PARAM_REDIRECT_BACK, constants.ROLEFEATURE_SAVE_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, "", vm)

}

func DeleteRoleFeatureRow(context session_functions.RequestContext, vm *viewModel.RoleFeatureModifyViewModel, t *model.Transaction) session_functions.ServerResponseStruct {
	var err error

	err = vm.RoleFeature.DeleteWithTran(t)

	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, constants.ERRORS_ROLEFEATURE_DELETE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
	}

	return session_functions.ServerResponseToStruct(PARAM_SUCCESS, context, PARAM_REDIRECT_RERENDER, constants.ROLEFEATURE_DELETE_SUCCESSFUL, PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), vm)
}
