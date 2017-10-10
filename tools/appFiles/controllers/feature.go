package controllers

import (
	"encoding/json"
	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/DanielRenne/GoCore/core/utils"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
)

func addEditFeatureValidateCommon(context session_functions.RequestContext, vm *viewModel.FeatureModifyViewModel, t *model.Transaction) (bool, session_functions.ServerResponseStruct) {
	//var err error

	// Add custom logic here and return false when theres an error
	//	Like this:
	//return false, session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, constants..ERRORS_FEATURE_SAVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)

	return true, session_functions.ServerResponseStruct{}
}

func vmFeatureAddEditGetCommon(context session_functions.RequestContext, respond session_functions.ServerResponse, vm *viewModel.FeatureModifyViewModel) bool {
	err := model.FeatureGroups.Query().All(&vm.FeatureGroups)

	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_ACCOUNT_NOT_FOUND, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return false
	}

	return true
}

func WrapFeatureViewModel(FeatureInstance model.Feature) viewModel.FeatureModifyViewModel {
	vm := viewModel.FeatureModifyViewModel{}
	vm.Feature = FeatureInstance
	return vm
}

func FeaturePostCommitHook(actionPerformed string, context session_functions.RequestContext, id string) (ret bool) {
	if utils.InArray(actionPerformed, utils.Array("CreateFeature", "CopyFeature")) && id != "" {
		var row model.Feature
		err := model.Features.Query().ById(id, &row)
		if err != nil {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			return false
		}
		if err == nil {
			//row.Key
			systemRoleFile := serverSettings.APP_LOCATION + "/constants/systemRoles.go"
			row.LastUpdateId = "57d9b383dcba0f51172f1f57"
			row.Key = strings.ToUpper(row.Key)
			utils.ReplaceTokenInFile(systemRoleFile, "//MoreConstants", "\tFEATURE_"+row.Key+"= \""+row.Id.Hex()+"\"\n//MoreConstants")
			err = row.Save()
			if err != nil {
				err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
				return false
			}

		} else {
			session_functions.Dump("Desc->Error in LocationEntityPostSaveHook", err)
		}
	}
	var all []model.Feature
	err := model.Features.Query().All(&all)
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
	err = ioutil.WriteFile(serverSettings.APP_LOCATION+"/db/bootstrap/features/features.json", []byte(strjson), 0644)
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}
	ret = true
	return ret
}

func CopyFeatureRow(context session_functions.RequestContext, copyRowVm *viewModel.FeatureModifyViewModel, t *model.Transaction) session_functions.ServerResponseStruct {
	var err error
	err = model.Features.Query().ById(copyRowVm.Feature.Id, &copyRowVm.Feature)
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return session_functions.ServerResponseToStruct(PARAM_FAILED, context, queries.AppContent.GetTranslation(context, constants.ERRORS_FEATURE_COPY)+core.Debug.HandleError(err), constants.ERRORS_FEATURE_COPY, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, copyRowVm)
	}
	copyRowVm.Feature.Id = ""

	// Copy Name
	replacements := queries.TagReplacements{
		Tag1: queries.Q("old_name_of_row", copyRowVm.Feature.Name),
	}
	copyRowVm.Feature.Name = queries.AppContent.GetTranslationWithReplacements(context, "CopyRow", &replacements)

	r := CreateFeatureRow(context, copyRowVm, t)

	if !r.CompletedSuccessfully {
		return session_functions.ServerResponseToStruct(r.CompletedSuccessfully, r.Context, r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
	}

	return session_functions.ServerResponseToStruct(PARAM_SUCCESS, context, PARAM_REDIRECT_RERENDER, queries.AppContent.GetTranslation(context, constants.FEATURE_COPY_SUCCESSFUL), PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), copyRowVm)
}

func UpdateFeatureRow(context session_functions.RequestContext, vm *viewModel.FeatureModifyViewModel, t *model.Transaction) session_functions.ServerResponseStruct {
	var err error

	success, r := addEditFeatureValidateCommon(context, vm, t)

	if !success {
		return r
	}

	err = vm.Feature.SaveWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Validation Error exists", vm.Feature.Errors)
			return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Golang error", err)
			return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, queries.AppContent.GetTranslation(context, constants.ERRORS_FEATURE_SAVE)+core.Debug.HandleError(err), PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		}
	}

	return session_functions.ServerResponseToStruct(PARAM_SUCCESS, context, PARAM_REDIRECT_BACK, constants.FEATURE_SAVE_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, "", vm)
}

func CreateFeatureRow(context session_functions.RequestContext, vm *viewModel.FeatureModifyViewModel, t *model.Transaction) session_functions.ServerResponseStruct {
	var err error

	success, r := addEditFeatureValidateCommon(context, vm, t)

	if !success {
		return r
	}

	err = vm.Feature.SaveWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Validation Error exists", vm.Feature.Errors)
			return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Golang error", err)
			return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, queries.AppContent.GetTranslation(context, constants.ERRORS_FEATURE_SAVE)+core.Debug.HandleError(err), PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		}
	}

	return session_functions.ServerResponseToStruct(PARAM_SUCCESS, context, PARAM_REDIRECT_BACK, constants.FEATURE_SAVE_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, "", vm)

}

func DeleteFeatureRow(context session_functions.RequestContext, vm *viewModel.FeatureModifyViewModel, t *model.Transaction) session_functions.ServerResponseStruct {
	var err error

	err = vm.Feature.DeleteWithTran(t)

	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, constants.ERRORS_FEATURE_DELETE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
	}

	return session_functions.ServerResponseToStruct(PARAM_SUCCESS, context, PARAM_REDIRECT_RERENDER, constants.FEATURE_DELETE_SUCCESSFUL, PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), vm)
}
