package controllers

import (
	"time"

	"github.com/DanielRenne/GoCore/core/utils"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
	"strings"
)

func (self *RoleFeatureListController) Root(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	//
	//_, ok := uriParams["DumpAllFeaturesToAllRoles"]
	//
	//if ok {
	//	var features []model.Feature
	//	_ = model.Features.Query().All(&features)
	//	var roles []model.Role
	//	_ = model.Roles.Query().All(&roles)
	//	var rfs []model.RoleFeature
	//	_ = model.RoleFeatures.Query().All(&rfs)
	//	for _, rf := range rfs {
	//		rf.Delete()
	//	}
	//
	//	t, _ := session_functions.StartTransaction(context())
	//	for _, feature := range features {
	//		for _, role := range roles {
	//			rf := model.RoleFeature{
	//				FeatureId: feature.Id.Hex(),
	//				RoleId:    role.Id.Hex(),
	//			}
	//			rf.SaveWithTran(t)
	//		}
	//	}
	//	t.Commit()
	//}

	var vm viewModel.RoleFeatureListViewModel
	//vm.SettingsBar = SetupVisibleButtons(context,ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_ROLEFEATURE} )
	vm.WidgetList = viewModel.InitWidgetList()
	if !self.SearchCommon(context, respond, &vm, uriParams) {
		return
	}

	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)

}

func (self *RoleFeatureListController) SearchCommon(context session_functions.RequestContext, respond session_functions.ServerResponse, vm *viewModel.RoleFeatureListViewModel, uriParams map[string]string) bool {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return false
	}
	q := model.RoleFeatures.Query().RenderViews(session_functions.GetDataFormat(context())).Join("Feature").Join("Role").InitAndOr().AddAnd()
	vm.WidgetList.DataKey = "RoleFeatures"
	vm.WidgetList.SearchFields = utils.Array(model.FIELD_ROLEFEATURE_FEATUREID)
	q = q.Join("LastUpdateUser")

	customCriteria, ok := uriParams["CustomCriteria"]

	if ok && len(customCriteria) > 0 && customCriteria != "last_hour" { // last_hour is busted everywhere for some reason.  Dont have time to fix mongo issues.
		if customCriteria == "last_hour" {
			vm.WidgetList.ListTitle = "ShowingModifiedLast15Minutes"
			vm.WidgetList.IsDefaultFilter = false
			q = q.AndRange(1, model.RangeQ("UpdateDate", time.Now().Add(-15*time.Minute).UTC(), time.Now().UTC()))
			viewModel.FilterWidgetList(vm.WidgetList, q)
		} else if strings.Index(customCriteria, "{") != -1 {
			viewModel.FilterWidgetListNoLimit(vm.WidgetList, q)
		}
	} else {
		viewModel.FilterWidgetList(vm.WidgetList, q)
		vm.WidgetList.ListTitle = "ShowingAllRoleFeatures"
	}
	err := q.All(&vm.RoleFeatures)

	if ok {
		if strings.Index(customCriteria, "{") != -1 {
			var filters viewModel.RoleListFilterModel
			filters.Parse(customCriteria)
			if filters.FeatureKey != "" {
				//vm.RoleFeatures
				for i := 0; i < len(vm.RoleFeatures); i++ {
					r := vm.RoleFeatures[i]
					if r.Joins.Feature == nil || r.Joins.Feature.Key != filters.FeatureKey {
						vm.RoleFeatures = append(vm.RoleFeatures[:i], vm.RoleFeatures[i+1:]...)
						i--
					}
				}
				vm.WidgetList.IsDefaultFilter = false
			}
		} else if customCriteria == "last_hour" {
			vm.WidgetList.ListTitle = "ShowingModifiedLast15Minutes"
			vm.WidgetList.IsDefaultFilter = false
			q.AndRange(1, model.RangeQ("UpdateDate", time.Now().Add(-15*time.Minute).UTC(), time.Now().UTC()))
		}
	}

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_GENERIC_DB, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return false
	}

	return true
}

func (self *RoleFeatureListController) Search(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {

	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.RoleFeatureListViewModel
	vm.LoadDefaultState()
	//vm.SettingsBar = SetupVisibleButtons(context,ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_ROLEFEATURE})
	vm.WidgetList = viewModel.InitWidgetListWithParams(uriParams)

	if !self.SearchCommon(context, respond, &vm, uriParams) {
		return
	}

	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)

}

func (self *RoleFeatureAddController) Root(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.RoleFeatureModifyViewModel
	//vm.SettingsBar = SetupVisibleButtons(context,ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_ROLEFEATURE_ADD})
	if !vmRoleFeatureAddEditGetCommon(context, respond, &vm) {
		return
	}

	// vm.Roles, err := queries.Roles.ByAccountWithContext(context)
	err := model.Roles.Query().All(&vm.Roles)
	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_ROLES_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	err = model.Features.Query().All(&vm.Features)
	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_ROLES_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	respond("", "", SNACKBAR_TYPE_SUCCESS, nil, "", vm)

}

func (self *RoleFeatureModifyController) Load(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {

	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.RoleFeatureModifyViewModel
	var err error
	vm.LoadDefaultState()

	//vm.SettingsBar = SetupVisibleButtons(context,ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_ROLEFEATURE_MODIFY})
	id, ok := uriParams["Id"]

	if ok {
		err = model.RoleFeatures.Query().ById(id, &vm.RoleFeature)
		if !vmRoleFeatureAddEditGetCommon(context, respond, &vm) {
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

	vm.Roles, err = queries.Roles.ByAccountWithContext(context)
	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_ROLES_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	err = model.Features.Query().All(&vm.Features)
	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_ROLES_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)
}
