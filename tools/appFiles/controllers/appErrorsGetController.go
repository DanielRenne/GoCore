package controllers

import (
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	//"github.com/DanielRenne/goCoreAppTemplate/queries"
	"time"

	"github.com/DanielRenne/GoCore/core/utils"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
)

func (self *AppErrorListController) Root(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {
	var vm viewModel.AppErrorListViewModel
	//vm.SettingsBar = SetupVisibleButtons(context,ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_APPERROR} )
	vm.WidgetList = viewModel.InitWidgetList()
	if !self.SearchCommon(context, respond, &vm, uriParams) {
		return
	}

	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)

}

func (self *AppErrorListController) SearchCommon(context session_functions.RequestContext, respond session_functions.ServerResponse, vm *viewModel.AppErrorListViewModel, uriParams map[string]string) bool {
	q := model.AppErrors.Query().RenderViews(session_functions.GetDataFormat(context())).InitAndOr().AddAnd()
	vm.WidgetList.DataKey = "AppErrors"
	vm.WidgetList.SearchFields = utils.Array(model.FIELD_APPERROR_ID, model.FIELD_APPERROR_MESSAGE, model.FIELD_APPERROR_URL, model.FIELD_APPERROR_ACCOUNTID, model.FIELD_APPERROR_USERID)
	viewModel.FilterWidgetList(vm.WidgetList, q)
	q = q.Join("LastUpdateUser").Join("User").Join("Account")

	customCriteria, ok := uriParams["CustomCriteria"]
	if ok && customCriteria != "last_hour" { // last_hour is busted everywhere for some reason.  Dont have time to fix mongo issues.
		if customCriteria == "last_hour" {
			vm.WidgetList.ListTitle = "ShowingModifiedLast15Minutes"
			vm.WidgetList.IsDefaultFilter = false
			q = q.AndRange(1, model.RangeQ("UpdateDate", time.Now().Add(-15*time.Minute).UTC(), time.Now().UTC()))
		}
	} else {
		vm.WidgetList.ListTitle = "ShowingAllAppErrors"
	}
	err := q.All(&vm.AppErrors)

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_GENERIC_DB, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return false
	}
	return true
}

func (self *AppErrorListController) Search(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {

	var vm viewModel.AppErrorListViewModel
	vm.LoadDefaultState()
	//vm.SettingsBar = SetupVisibleButtons(context,ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_APPERROR})
	vm.WidgetList = viewModel.InitWidgetListWithParams(uriParams)

	if !self.SearchCommon(context, respond, &vm, uriParams) {
		return
	}

	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)

}

func (self *AppErrorAddController) Root(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {
	var vm viewModel.AppErrorModifyViewModel
	//vm.SettingsBar = SetupVisibleButtons(context,ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_APPERROR_ADD})
	if !vmAppErrorAddEditGetCommon(context, respond, &vm) {
		return
	}
	respond("", "", SNACKBAR_TYPE_SUCCESS, nil, "", vm)

}

func (self *AppErrorModifyController) Load(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {

	var vm viewModel.AppErrorModifyViewModel
	var err error
	vm.LoadDefaultState()
	//vm.SettingsBar = SetupVisibleButtons(context,ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_APPERROR_MODIFY})
	id, ok := uriParams["Id"]

	if ok {
		err = model.AppErrors.Query().ById(id, &vm.AppError)
		if !vmAppErrorAddEditGetCommon(context, respond, &vm) {
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
