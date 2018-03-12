package controllers

import (
	"github.com/DanielRenne/goCoreAppTemplate/br"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	//"github.com/DanielRenne/goCoreAppTemplate/queries"
	"time"

	"github.com/DanielRenne/GoCore/core/utils"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
)

func (self *FileObjectListController) Root(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.FileObjectListViewModel
	//vm.SettingsBar = SetupVisibleButtons(context,ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_FILEOBJECT} )
	vm.WidgetList = viewModel.InitWidgetList()
	if !self.SearchCommon(context, respond, &vm, uriParams) {
		return
	}

	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)

}

func (self *FileObjectListController) SearchCommon(context session_functions.RequestContext, respond session_functions.ServerResponse, vm *viewModel.FileObjectListViewModel, uriParams map[string]string) bool {
	q := model.FileObjects.Query().RenderViews(session_functions.GetDataFormat(context())).InitAndOr().AddAnd()
	vm.WidgetList.DataKey = "FileObjects"
	vm.WidgetList.SearchFields = utils.Array(model.FIELD_FILEOBJECT_ID, model.FIELD_FILEOBJECT_ACCOUNTID, model.FIELD_FILEOBJECT_NAME)
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
		vm.WidgetList.ListTitle = "ShowingAllFileObjects"
	}
	err := q.All(&vm.FileObjects)

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_GENERIC_DB, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return false
	}

	return true
}

func (self *FileObjectListController) Search(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}

	var vm viewModel.FileObjectListViewModel
	vm.LoadDefaultState()
	//vm.SettingsBar = SetupVisibleButtons(context,ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_FILEOBJECT})
	vm.WidgetList = viewModel.InitWidgetListWithParams(uriParams)

	if !self.SearchCommon(context, respond, &vm, uriParams) {
		return
	}

	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)

}

func (self *FileObjectAddController) Root(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.FileObjectModifyViewModel
	vm.LoadDefaultState()
	//vm.SettingsBar = SetupVisibleButtons(context,ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_FILEOBJECT_ADD})
	message, err := br.FileObjects.GetVmDefaults(&vm)
	if err != nil {
		respond(constants.PARAM_REDIRECT_NONE, message, constants.PARAM_SNACKBAR_TYPE_ERROR, err, constants.PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	respond(constants.PARAM_REDIRECT_NONE, constants.PARAM_SNACKBAR_MESSAGE_NONE, constants.SNACKBAR_TYPE_SUCCESS, nil, "", vm)

}

func (self *FileObjectModifyController) Load(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {

	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.FileObjectModifyViewModel
	var err error
	vm.LoadDefaultState()
	//vm.SettingsBar = SetupVisibleButtons(context,ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_FILEOBJECT_MODIFY})
	id, ok := uriParams["Id"]

	if ok {
		err = model.FileObjects.Query().ById(id, &vm.FileObject)
		if err != nil {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_ACCOUNT_NOT_FOUND, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}

		message, err := br.FileObjects.GetVmDefaults(&vm)
		if err != nil {
			respond(constants.PARAM_REDIRECT_NONE, message, constants.PARAM_SNACKBAR_TYPE_ERROR, err, constants.PARAM_TRANSACTION_ID_NONE, vm)
			return
		}

	} else {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_ACCOUNT_NO_ID, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)
}
