package controllers

import (
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
)

func (self *RoleListController) Root(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {
	if session_functions.BlockByRoleAccess(context(), constants.FEATURE_ROLE_VIEW) {
		return
	}
	var vm viewModel.RoleListViewModel
	vm.SettingsBar = SetupVisibleButtons(context, ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_ROLE})
	vm.WidgetList = viewModel.InitWidgetList()
	if !self.SearchCommon(context, respond, &vm, uriParams) {
		return
	}

	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)

}

func (self *RoleListController) SearchCommon(context session_functions.RequestContext, respond session_functions.ServerResponse, vm *viewModel.RoleListViewModel, uriParams map[string]string) bool {
	vm.WidgetList.ListTitle = "ShowingAllRoles"
	vm.WidgetList.DataKey = "Roles"
	var err error

	account, _ := session_functions.GetSessionAccount(context())
	q := model.Roles.Query().RenderViews(session_functions.GetDataFormat(context())).Join("RoleFeatures")
	acct, err := session_functions.GetSessionAccount(context())
	q = q.Filter(model.Q(model.FIELD_ROLE_ACCOUNTID, acct.Id.Hex())).Filter(model.Q(model.FIELD_ROLE_ACCOUNTTYPE, account.AccountTypeShort))
	viewModel.FilterWidgetList(vm.WidgetList, q)
	q = q.Join("LastUpdateUser")
	err = q.All(&vm.Roles)

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_GENERIC_DB, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return false
	}

	var or []model.Role
	q2 := model.Roles.Query().RenderViews(session_functions.GetDataFormat(context())).Join("RoleFeatures").Filter(model.Q(model.FIELD_ROLE_CANDELETE, false)).Filter(model.Q(model.FIELD_ROLE_ACCOUNTID, "")).Filter(model.Q(model.FIELD_ROLE_ACCOUNTTYPE, account.AccountTypeShort))

	viewModel.FilterWidgetList(vm.WidgetList, q2)
	q2 = q2.Join("LastUpdateUser")
	err = q2.All(&or)

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_GENERIC_DB, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return false
	}
	for _, row := range or {
		vm.Roles = append(vm.Roles, row)
	}

	return true
}

func (self *RoleListController) Search(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {
	if session_functions.BlockByRoleAccess(context(), constants.FEATURE_ROLE_VIEW) {
		return
	}

	var vm viewModel.RoleListViewModel
	vm.LoadDefaultState()
	vm.SettingsBar = SetupVisibleButtons(context, ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_ROLE})
	vm.WidgetList = viewModel.InitWidgetListWithParams(uriParams)

	if !self.SearchCommon(context, respond, &vm, uriParams) {
		return
	}

	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)

}

func (self *RoleAddController) Root(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {
	if session_functions.BlockByRoleAccess(context(), constants.FEATURE_ROLE_ADD) {
		return
	}
	var vm viewModel.RoleModifyViewModel
	vm.LoadDefaultState()
	vm.SettingsBar = SetupVisibleButtons(context, ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_ROLE_ADD})
	if !vmRoleAddEditGetCommon(context, respond, &vm) {
		return
	}
	respond("", "", SNACKBAR_TYPE_SUCCESS, nil, "", vm)

}

func (self *RoleModifyController) Load(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {
	if session_functions.BlockByRoleAccess(context(), constants.FEATURE_ROLE_MODIFY) {
		return
	}
	var vm viewModel.RoleModifyViewModel
	var err error
	vm.LoadDefaultState()
	vm.SettingsBar = SetupVisibleButtons(context, ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_ROLE_MODIFY})
	id, ok := uriParams["Id"]

	if ok {
		err = model.Roles.Query().ById(id, &vm.Role)
		if !vmRoleAddEditGetCommon(context, respond, &vm) {
			return
		}
		if err != nil {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_ACCOUNT_NOT_FOUND, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
	} else {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_ACCOUNT_NO_ID, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)
}
