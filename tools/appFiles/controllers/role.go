package controllers

import (
	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/GoCore/core/utils"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
	"github.com/pkg/errors"
)

func addEditRoleValidateCommon(context session_functions.RequestContext, vm *viewModel.RoleModifyViewModel, t *model.Transaction) (bool, session_functions.ServerResponseStruct) {
	//var err error
	acct, err := session_functions.GetSessionAccount(context())
	if err == nil {
		vm.Role.AccountId = acct.Id.Hex()
	}
	vm.Role.CanDelete = true
	return true, session_functions.ServerResponseStruct{}
}

func (self *RolesController) MapRoleFeatures(vm viewModel.RoleModifyViewModel, t *model.Transaction) (ret bool) {
	var roleFeatures []model.RoleFeature
	err := model.RoleFeatures.Query().Filter(model.Q(model.FIELD_ROLEFEATURE_ROLEID, vm.Role.Id.Hex())).All(&roleFeatures)
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}
	var ok bool
	ok = true
	for _, rf := range roleFeatures {
		err := rf.Delete()
		if err != nil {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			ok = false
		}
	}
	if !ok {
		return
	}

	for featureId, isEnabled := range vm.FeaturesEnabled {
		if isEnabled {
			rf := model.RoleFeature{
				FeatureId: featureId,
				RoleId:    vm.Role.Id.Hex(),
			}
			rf.SaveWithTran(t)
		}
	}
	ret = true
	return ret
}

func vmRoleAddEditGetCommon(context session_functions.RequestContext, respond session_functions.ServerResponse, vm *viewModel.RoleModifyViewModel) bool {
	// Use this method as a common way to load up any VM defaults that both the add and edit pages need!
	q := model.FeatureGroups.Query().Join("Features")
	account, _ := session_functions.GetSessionAccount(context())
	q = q.Filter(model.Q(model.FIELD_FEATUREGROUP_ACCOUNTTYPE, account.AccountTypeShort))

	err := q.All(&vm.FeatureGroups)
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		respond(PARAM_REDIRECT_NONE, "Feature groups not found", PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return false
	}

	var orBlanks []model.FeatureGroup
	q2 := model.FeatureGroups.Query().Filter(model.Q(model.FIELD_FEATUREGROUP_ACCOUNTTYPE, "")).Join("Features")
	err = q2.All(&orBlanks)
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return false
	}
	for _, row := range orBlanks {
		vm.FeatureGroups = append(vm.FeatureGroups, row)
	}
	//for _, feature := range vm.FeatureGroups.Joins.
	for _, group := range vm.FeatureGroups {
		for _, feature := range *group.Joins.Features.Items {
			vm.FeaturesEnabled[feature.Id.Hex()] = session_functions.CheckRoleAccessByRole("12345", vm.Role.Id.Hex(), feature.Id.Hex())
		}
	}
	return true

}

func WrapRoleViewModel(RoleInstance model.Role) viewModel.RoleModifyViewModel {
	vm := viewModel.RoleModifyViewModel{}
	vm.Role = RoleInstance
	return vm
}

func RolePostCommitHook(actionPerformed string, context session_functions.RequestContext, id string) (ret bool) {
	if utils.InArray(actionPerformed, utils.Array("DeleteRole")) && id != "" {
		var roleFeatures []model.RoleFeature
		err := model.RoleFeatures.Query().Filter(model.Q(model.FIELD_ROLEFEATURE_ROLEID, id)).All(&roleFeatures)

		if err != nil {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			return false
		}
		t, err := session_functions.StartTransaction(context())
		if err != nil {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			return false
		}

		for _, f := range roleFeatures {
			feature := f
			vmRoleFeature := WrapRoleFeatureViewModel(feature)
			r := DeleteRoleFeatureRow(context, &vmRoleFeature, t)
			if !r.CompletedSuccessfully {
				err := queries.Transactions.Rollback(context, "Failed to complete all floor copies", t.Id.Hex())
				if err != nil {
					err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
					return false
				}
			}
		}
		err = t.Commit()
		if err != nil {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			return false
		}
	}
	return true
}

func CopyRoleRow(context session_functions.RequestContext, copyRowVm *viewModel.RoleModifyViewModel, t *model.Transaction) session_functions.ServerResponseStruct {
	var err error
	originalId := copyRowVm.Role.Id
	err = model.Roles.Query().ById(copyRowVm.Role.Id, &copyRowVm.Role)
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return session_functions.ServerResponseToStruct(PARAM_FAILED, context, queries.AppContent.GetTranslation(context, constants.ERRORS_ROLE_COPY)+core.Debug.HandleError(err), constants.ERRORS_ROLE_COPY, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, copyRowVm)
	}
	copyRowVm.Role.Id = ""

	// Copy Name
	replacements := queries.TagReplacements{
		Tag1: queries.Q("old_name_of_row", copyRowVm.Role.Name),
	}
	copyRowVm.Role.Name = queries.AppContent.GetTranslationWithReplacements(context, "CopyRow", &replacements)

	r := CreateRoleRow(context, copyRowVm, t)

	if !r.CompletedSuccessfully {
		return session_functions.ServerResponseToStruct(r.CompletedSuccessfully, r.Context, r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
	}

	var roleFeatures []model.RoleFeature
	err = model.RoleFeatures.Query().Filter(model.Q(model.FIELD_ROLEFEATURE_ROLEID, originalId)).All(&roleFeatures)

	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, "An error occurred", PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
	}

	for _, f := range roleFeatures {
		feature := f
		vmRoleFeature := WrapRoleFeatureViewModel(feature)
		vmRoleFeature.RoleFeature.Id = ""
		vmRoleFeature.RoleFeature.RoleId = copyRowVm.Role.Id.Hex()
		r := CreateRoleFeatureRow(context, &vmRoleFeature, t)
		if !r.CompletedSuccessfully {
			err := queries.Transactions.Rollback(context, "Failed to complete all floor copies", t.Id.Hex())
			if err != nil {
				err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
				return session_functions.ServerResponseToStruct(r.CompletedSuccessfully, r.Context, r.Redirect, "Failed to rollback on roleFeature copies", r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
			}
		}
	}

	return session_functions.ServerResponseToStruct(PARAM_SUCCESS, context, PARAM_REDIRECT_RERENDER, queries.AppContent.GetTranslation(context, constants.ROLE_COPY_SUCCESSFUL), PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), copyRowVm)
}

func UpdateRoleRow(context session_functions.RequestContext, vm *viewModel.RoleModifyViewModel, t *model.Transaction) session_functions.ServerResponseStruct {
	var err error

	success, r := addEditRoleValidateCommon(context, vm, t)

	if !success {
		return r
	}

	err = vm.Role.SaveWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Validation Error exists", vm.Role.Errors)
			return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Golang error", err)
			return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, queries.AppContent.GetTranslation(context, constants.ERRORS_ROLE_SAVE)+core.Debug.HandleError(err), PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		}
	}

	return session_functions.ServerResponseToStruct(PARAM_SUCCESS, context, PARAM_REDIRECT_BACK, constants.ROLE_SAVE_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, "", vm)
}

func CreateRoleRow(context session_functions.RequestContext, vm *viewModel.RoleModifyViewModel, t *model.Transaction) session_functions.ServerResponseStruct {
	var err error

	success, r := addEditRoleValidateCommon(context, vm, t)

	if !success {
		return r
	}

	err = vm.Role.SaveWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Validation Error exists", vm.Role.Errors)
			return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Golang error", err)
			return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, queries.AppContent.GetTranslation(context, constants.ERRORS_ROLE_SAVE)+core.Debug.HandleError(err), PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		}
	}

	return session_functions.ServerResponseToStruct(PARAM_SUCCESS, context, PARAM_REDIRECT_BACK, constants.ROLE_SAVE_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, "", vm)

}

func IsSystemRole(vm *viewModel.RoleModifyViewModel) (isDefaultRole bool) {
	return true
}

func DeleteRoleRow(context session_functions.RequestContext, vm *viewModel.RoleModifyViewModel, t *model.Transaction) session_functions.ServerResponseStruct {
	var err error

	acct, err := session_functions.GetSessionAccount(context())
	if err == nil {
		vm.Role.AccountId = acct.Id.Hex()
	}
	vm.Role.CanDelete = true
	var total []model.AccountRole
	if model.AccountRoles.Query().Filter(model.Q(model.FIELD_ACCOUNTROLE_ROLEID, vm.Role.Id.Hex())).TotalRows(&total) > 0 {
		return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, "CannotDelete", PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
	}

	if IsSystemRole(vm) {
		return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, "CannotDelete", PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
	}

	if !vm.Role.CanDelete {
		return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, "CannotDeleteDueToField", PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
	}

	err = vm.Role.DeleteWithTran(t)

	if err != nil {
		return session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, constants.ERRORS_ROLE_DELETE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
	}

	return session_functions.ServerResponseToStruct(PARAM_SUCCESS, context, PARAM_REDIRECT_RERENDER, constants.ROLE_DELETE_SUCCESSFUL, PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), vm)
}
