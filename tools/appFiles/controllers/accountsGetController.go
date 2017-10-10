package controllers

import (
	"time"

	"github.com/DanielRenne/GoCore/core/utils"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
)

func (self *AccountListController) Root(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {
	if session_functions.BlockByRoleAccess(context(), constants.FEATURE_ACCOUNT_VIEW) {
		return
	}
	var vm viewModel.AccountListViewModel
	vm.SettingsBar = SetupVisibleButtons(context, ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_ACCOUNT_INSTALLS})
	vm.WidgetList = viewModel.InitWidgetList()
	if !self.SearchCommon(context, respond, &vm, uriParams) {
		return
	}

	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)

}

func (self *AccountListController) SearchCommon(context session_functions.RequestContext, respond session_functions.ServerResponse, vm *viewModel.AccountListViewModel, uriParams map[string]string) bool {
	if session_functions.BlockByRoleAccess(context(), constants.FEATURE_ACCOUNT_VIEW) {
		return false
	}
	user, err := session_functions.GetSessionUser(context())
	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_USER_NOT_FOUND, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return false
	}

	account, err := session_functions.GetSessionAccount(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_ACCOUNT_NOT_FOUND, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return false
	}

	if account.AccountTypeShort == "cust" {
		respond(CONTROLLER_USERLIST, "CustomersNoAccounts", PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return false
	}
	q, _ := queries.Accounts.QueryByUser2(user, account)
	vm.WidgetList.DataKey = "Accounts"
	vm.WidgetList.SearchFields = utils.Array(model.FIELD_ACCOUNT_ACCOUNTNAME, model.FIELD_ACCOUNT_ADDRESS1, model.FIELD_ACCOUNT_EMAIL, model.FIELD_ACCOUNT_REGION, model.FIELD_ACCOUNT_PRIMARYPHONE, model.FIELD_ACCOUNT_ID, model.FIELD_ACCOUNT_ACCOUNTTYPESHORT)

	viewModel.FilterWidgetList(vm.WidgetList, q.RenderViews(session_functions.GetDataFormat(context())).Join("LastUpdateUser").Join("RelatedAccount"))
	customCriteria, ok := uriParams["CustomCriteria"]
	if ok && customCriteria != "last_hour" { // last_hour is busted everywhere for some reason.  Dont have time to fix mongo issues.
		if customCriteria == "last_hour" {
			vm.WidgetList.ListTitle = "ShowingModifiedLast15Minutes"
			vm.WidgetList.IsDefaultFilter = false
			q = q.AndRange(1, model.RangeQ("UpdateDate", time.Now().Add(-15*time.Minute).UTC(), time.Now().UTC()))
		}
	} else {
		vm.WidgetList.ListTitle = "ShowingAllAccounts"
	}

	err = q.All(&vm.Accounts)

	vm.Roles, err = queries.Roles.ByAccount(account)

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_ROLES_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return false
	}

	vm.User = user

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_ROLES_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return false
	}

	return true
}

func (self *AccountListController) Search(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {
	if session_functions.BlockByRoleAccess(context(), constants.FEATURE_ACCOUNT_VIEW) {
		return
	}
	var vm viewModel.AccountListViewModel
	vm.LoadDefaultState()
	vm.SettingsBar = SetupVisibleButtons(context, ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_ACCOUNT_INSTALLS})
	vm.WidgetList = viewModel.InitWidgetListWithParams(uriParams)

	if !self.SearchCommon(context, respond, &vm, uriParams) {
		return
	}

	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)

}

func (self *AccountsController) Root(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {
	var vm viewModel.AccountsViewModel
	vm.LoadDefaultState()
	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)

}

func (self *AccountAddController) Root(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {
	if session_functions.BlockByRoleAccess(context(), constants.FEATURE_ACCOUNT_ADD) {
		return
	}
	var vm viewModel.AccountModifyViewModel
	vm.LoadDefaultState()
	vm.SettingsBar = SetupVisibleButtons(context, ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_ACCOUNT_INSTALL_ADD})
	if !vmAccountAddEditGetCommon(context, respond, &vm) {
		return
	}

	var countries []model.Country
	_ = model.Countries.Query().All(&countries)
	vm.Countries = countries

	var states []model.State
	_ = model.States.Query().All(&states)
	vm.States = constructStatesMap(states, countries)

	// For new records, lets set the country to your current
	account, _ := session_functions.GetSessionAccount(context())
	vm.Account.CountryId = account.CountryId

	respond("", "", SNACKBAR_TYPE_SUCCESS, nil, "", vm)

}

func (self *AccountModifyController) Load(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {

	if session_functions.BlockByRoleAccess(context(), constants.FEATURE_ACCOUNT_MODIFY) {
		return
	}
	var vm viewModel.AccountModifyViewModel
	var err error
	vm.LoadDefaultState()

	var countries []model.Country
	_ = model.Countries.Query().All(&countries)
	vm.Countries = countries

	var states []model.State
	_ = model.States.Query().All(&states)
	vm.States = constructStatesMap(states, countries)
	vm.SettingsBar = SetupVisibleButtons(context, ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_ACCOUNT_INSTALL_MODIFY})
	id, ok := uriParams["Id"]

	if ok {
		err = model.Accounts.Query().LeftJoin("Company").ById(id, &vm.Account)
		if !vmAccountAddEditGetCommon(context, respond, &vm) {
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
