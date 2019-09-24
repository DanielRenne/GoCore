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

func (self *FeatureListController) Root(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.FeatureListViewModel
	//vm.SettingsBar = SetupVisibleButtons(context,ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_FEATURE} )
	vm.WidgetList = viewModel.InitWidgetList()
	if !self.SearchCommon(context, respond, &vm, uriParams) {
		return
	}

	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)

}

func (self *FeatureListController) SearchCommon(context session_functions.RequestContext, respond session_functions.ServerResponse, vm *viewModel.FeatureListViewModel, uriParams map[string]string) bool {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return false
	}
	q := model.Features.Query().RenderViews(session_functions.GetDataFormat(context())).Join("FeatureGroup").InitAndOr().AddAnd()
	vm.WidgetList.DataKey = "Features"
	vm.WidgetList.SearchFields = utils.Array(model.FIELD_FEATURE_ID, model.FIELD_FEATURE_NAME, model.FIELD_FEATURE_DESCRIPTION)
	viewModel.FilterWidgetList(vm.WidgetList, q)
	q = q.Join("LastUpdateUser")

	customCriteria, ok := uriParams["CustomCriteria"]
	if ok && customCriteria != "last_hour" { // last_hour is busted everywhere for some reason.  Dont have time to fix mongo issues.
		if customCriteria == "last_hour" {
			vm.WidgetList.ListTitle = "ShowingModifiedLast15Minutes"
			vm.WidgetList.IsDefaultFilter = false
			q = q.AndRange(1, model.RangeQ("UpdateDate", time.Now().Add(-15*time.Minute).UTC(), time.Now().UTC()))
		}
	} else {
		vm.WidgetList.ListTitle = "ShowingAllFeatures"
	}
	err := q.All(&vm.Features)

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_GENERIC_DB, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return false
	}

	return true
}

func (self *FeatureListController) Search(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {

	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.FeatureListViewModel
	vm.LoadDefaultState()
	//vm.SettingsBar = SetupVisibleButtons(context,ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_FEATURE})
	vm.WidgetList = viewModel.InitWidgetListWithParams(uriParams)

	if !self.SearchCommon(context, respond, &vm, uriParams) {
		return
	}

	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)

}

func (self *FeatureAddController) Root(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.FeatureModifyViewModel
	//vm.SettingsBar = SetupVisibleButtons(context,ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_FEATURE_ADD})
	if !vmFeatureAddEditGetCommon(context, respond, &vm) {
		return
	}
	respond("", "", SNACKBAR_TYPE_SUCCESS, nil, "", vm)

}

func (self *FeatureModifyController) Load(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {

	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.FeatureModifyViewModel
	var err error
	vm.LoadDefaultState()
	//vm.SettingsBar = SetupVisibleButtons(context,ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_FEATURE_MODIFY})
	id, ok := uriParams["Id"]

	if ok {
		err = model.Features.Query().ById(id, &vm.Feature)
		if !vmFeatureAddEditGetCommon(context, respond, &vm) {
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
