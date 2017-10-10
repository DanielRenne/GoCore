package controllers

import (
	"encoding/json"
	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
	"github.com/pkg/errors"
	"io/ioutil"
)

func addEditFeatureGroupValidateCommon(context session_functions.RequestContext, vm *viewModel.FeatureGroupModifyViewModel, t *model.Transaction) (bool, session_functions.ServerResponseStruct) {
	//var err error

	// Add custom logic here and return false when theres an error
	//	Like this:
	//return false, session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, constants..ERRORS_FEATUREGROUP_SAVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)

	return true, session_functions.ServerResponseStruct{}
}

func vmFeatureGroupAddEditGetCommon(context session_functions.RequestContext, respond session_functions.ServerResponse, vm *viewModel.FeatureGroupModifyViewModel) bool {
	// Use this method as a common way to load up any VM defaults that both the add and edit pages need!

	//companies, err := queries.Companies.ByAccountWithContext(context)

	//if err != nil {
	//	respond(PARAM_REDIRECT_NONE, constants.ERRORS_ACCOUNT_NOT_FOUND, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
	//	return false
	//}

	//if !vm.FeatureGroup.XXXXFIELDXXXX {
	//	vm.FeatureGroup.Errors.XXXXFIELDXXXX = constants.ERROR_REQUIRED_FIELD
	//	return false, session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
	//}

	return true
}

func WrapFeatureGroupViewModel(FeatureGroupInstance model.FeatureGroup) viewModel.FeatureGroupModifyViewModel {
	vm := viewModel.FeatureGroupModifyViewModel{}
	vm.FeatureGroup = FeatureGroupInstance
	return vm
}

func FeatureGroupPostCommitHook(actionPerformed string, context session_functions.RequestContext, id string) (ret bool) {
	var all []model.FeatureGroup
	err := model.FeatureGroups.Query().All(&all)
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
	err = ioutil.WriteFile(serverSettings.APP_LOCATION+"/db/bootstrap/featureGroups/featureGroups.json", []byte(strjson), 0644)
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}
	ret = true
	return ret
}

func CopyFeatureGroupRow(context session_functions.RequestContext, copyRowVm *viewModel.FeatureGroupModifyViewModel, t *model.Transaction) session_functions.ServerResponseStruct {
	var err error
	err = model.FeatureGroups.Query().ById(copyRowVm.FeatureGroup.Id, &copyRowVm.FeatureGroup)
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return session_functions.ServerResponseToStruct(PARAM_FAILED, context, queries.AppContent.GetTranslation(context, constants.ERRORS_FEATUREGROUP_COPY)+core.Debug.HandleError(err), constants.ERRORS_FEATUREGROUP_COPY, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, copyRowVm)
	}
	copyRowVm.FeatureGroup.Id = ""

	// Copy Name
	replacements := queries.TagReplacements{
		Tag1: queries.Q("old_name_of_row", copyRowVm.FeatureGroup.Name),
	}
	copyRowVm.FeatureGroup.Name = queries.AppContent.GetTranslationWithReplacements(context, "CopyRow", &replacements)

	r := CreateFeatureGroupRow(context, copyRowVm, t)

	if !r.CompletedSuccessfully {
		return session_functions.ServerResponseToStruct(r.CompletedSuccessfully, r.Context, r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
	}

	return session_functions.ServerResponseToStruct(PARAM_SUCCESS, context, PARAM_REDIRECT_RERENDER, queries.AppContent.GetTranslation(context, constants.FEATUREGROUP_COPY_SUCCESSFUL), PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), copyRowVm)
}

func UpdateFeatureGroupRow(context session_functions.RequestContext, vm *viewModel.FeatureGroupModifyViewModel, t *model.Transaction) session_functions.ServerResponseStruct {
	var err error

	success, r := addEditFeatureGroupValidateCommon(context, vm, t)

	if !success {
		return r
	}

	err = vm.FeatureGroup.SaveWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Validation Error exists", vm.FeatureGroup.Errors)
			return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Golang error", err)
			return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, queries.AppContent.GetTranslation(context, constants.ERRORS_FEATUREGROUP_SAVE)+core.Debug.HandleError(err), PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		}
	}

	return session_functions.ServerResponseToStruct(PARAM_SUCCESS, context, PARAM_REDIRECT_BACK, constants.FEATUREGROUP_SAVE_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, "", vm)
}

func CreateFeatureGroupRow(context session_functions.RequestContext, vm *viewModel.FeatureGroupModifyViewModel, t *model.Transaction) session_functions.ServerResponseStruct {
	var err error

	success, r := addEditFeatureGroupValidateCommon(context, vm, t)

	if !success {
		return r
	}

	err = vm.FeatureGroup.SaveWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Validation Error exists", vm.FeatureGroup.Errors)
			return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Golang error", err)
			return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, queries.AppContent.GetTranslation(context, constants.ERRORS_FEATUREGROUP_SAVE)+core.Debug.HandleError(err), PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		}
	}

	return session_functions.ServerResponseToStruct(PARAM_SUCCESS, context, PARAM_REDIRECT_BACK, constants.FEATUREGROUP_SAVE_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, "", vm)

}

func DeleteFeatureGroupRow(context session_functions.RequestContext, vm *viewModel.FeatureGroupModifyViewModel, t *model.Transaction) session_functions.ServerResponseStruct {
	var err error

	err = vm.FeatureGroup.DeleteWithTran(t)

	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, constants.ERRORS_FEATUREGROUP_DELETE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
	}

	return session_functions.ServerResponseToStruct(PARAM_SUCCESS, context, PARAM_REDIRECT_RERENDER, constants.FEATUREGROUP_DELETE_SUCCESSFUL, PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), vm)
}
