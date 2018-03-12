package controllers

import (
	"github.com/DanielRenne/GoCore/core/ginServer"
	"github.com/DanielRenne/GoCore/core/utils"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
	"gopkg.in/mgo.v2/bson"
	"regexp"
)

func (self *UserListController) Root(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {

	if session_functions.BlockByRoleAccess(context(), constants.FEATURE_USER_VIEW) {
		return
	}
	var vm viewModel.UserListViewModel
	vm.SettingsBar = SetupVisibleButtons(context, ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_USERS})
	vm.WidgetList = viewModel.InitWidgetList()
	if !self.SearchCommon(context, respond, &vm, true) {
		return
	}

	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)

}

func (self *UserListController) SearchCommon(context session_functions.RequestContext, respond session_functions.ServerResponse, vm *viewModel.UserListViewModel, applyLimit bool) bool {

	q, err := queries.AccountRoles.QueryByAccountWithContext(context)
	vm.WidgetList.DataKey = "Users"
	vm.WidgetList.SearchFields = utils.Array(model.FIELD_USER_FIRST, model.FIELD_USER_LAST, model.FIELD_USER_EMAIL)

	q.Sort("UpdateDate")

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_ACCOUNT_NOT_FOUND, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return false
	}

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_GENERIC_DB, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return false
	}

	vm.Roles, err = queries.Roles.ByAccountWithContext(context)

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_ROLES_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return false
	}
	q.RenderViews(session_functions.GetDataFormat(context())).Join("User.LastUpdateUser").Join("Account").Join("Role")

	if applyLimit {
		q.Limit(vm.WidgetList.PerPage)
		if vm.WidgetList.Page > 1 {
			q.Skip((vm.WidgetList.Page - 1) * vm.WidgetList.PerPage)
		}
		if vm.WidgetList.SortBy != "" {
			q.Sort(vm.WidgetList.SortDirection + vm.WidgetList.SortBy)
		}
	}

	if vm.WidgetList.Criteria != "" && len(vm.WidgetList.SearchFields) > 0 {
		if vm.WidgetList.Criteria != "" && len(vm.WidgetList.SearchFields) > 0 {
			var users []model.User
			qUsers := model.Users.Query()
			for _, field := range vm.WidgetList.SearchFields {
				if field == "Id" && len(vm.WidgetList.Criteria) == 24 {
					field = "_id"
					qUsers.Or(model.Q(field, bson.ObjectIdHex(vm.WidgetList.Criteria)))
				} else {
					qUsers.Or(model.Q(field, bson.M{"$regex": bson.RegEx{`.*` + regexp.QuoteMeta(vm.WidgetList.Criteria) + `.*`, "i"}}))
				}
			}
			err := qUsers.All(&users)
			if err != nil || len(users) == 0 {
				session_functions.Dump(users, err)
				q.Filter(model.Q("AccountId", "neverwillmatch"))
			} else {
				var userIds []string
				for _, user := range users {
					userIds = append(userIds, user.Id.Hex())
				}

				session_functions.Dump(userIds)
				q.In(model.Q(model.FIELD_ACCOUNTROLE_USERID, userIds))
			}
		}
	}
	err = q.All(&vm.Users)

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_GENERIC_DB, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return false
	}

	for idx := range vm.Users {
		vm.Users[idx].Errors.Id = "false"
	}

	return true
}

func (self *UserListController) Search(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {

	if session_functions.BlockByRoleAccess(context(), constants.FEATURE_USER_VIEW) {
		return
	}

	var vm viewModel.UserListViewModel
	vm.LoadDefaultState()
	vm.SettingsBar = SetupVisibleButtons(context, ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_USERS})
	vm.WidgetList = viewModel.InitWidgetListWithParams(uriParams)

	if !self.SearchCommon(context, respond, &vm, true) {
		return
	}

	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)

}

func (self *UserAddController) Root(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {
	if session_functions.BlockByRoleAccess(context(), constants.FEATURE_USER_ADD) {
		return
	}
	controller := UserModifyController{}
	uriParams["Id"] = "New"
	controller.Load(context, uriParams, respond)

}

func (self *UserProfileController) Root(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {

	controller := UserModifyController{}
	controller.Load(context, uriParams, respond)
}

func (self *UserModifyController) Load(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {
	userId, ok := uriParams["Id"]
	user, _ := session_functions.GetSessionUser(context())

	if ok && session_functions.BlockByRoleAccess(context(), constants.FEATURE_USER_MODIFY) && user.Id.Hex() != userId {
		return
	}

	var vm viewModel.UserModifyViewModel
	vm.LoadDefaultState()
	vm.Locales = model.Locales
	vm.TimeZones = model.TimeZoneLocations
	localeLanguage := ginServer.GetLocaleLanguage(context())
	vm.UserLocale = localeLanguage.Language

	var err error
	acct, errAcct := session_functions.GetSessionAccount(context())
	if errAcct != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_ACCOUNT_NOT_FOUND, PARAM_SNACKBAR_TYPE_ERROR, errAcct, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	if ok {
		var controller string
		if userId == "New" {
			controller = viewModel.SETTINGS_CONST_USER_ADD
			vm.AccountRole.AccountId = acct.Id.Hex()
		} else {
			controller = viewModel.SETTINGS_CONST_USER_MODIFY
			err = model.Users.Query().ById(userId, &vm.User)
			if err != nil {
				respond(PARAM_REDIRECT_NONE, constants.ERRORS_USER_NOT_FOUND, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
				return
			}
			err = model.AccountRoles.Query().Filter(model.Q(model.FIELD_ACCOUNTROLE_ACCOUNTID, acct.Id)).Filter(model.Q(model.FIELD_ACCOUNTROLE_USERID, vm.User.Id)).One(&vm.AccountRole)
			if err != nil {
				respond(PARAM_REDIRECT_NONE, err.Error(), PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
				return
			}
		}
		vm.SettingsBar = SetupVisibleButtons(context, ButtonBarMetaData{IsAccountSettings: true, CurrentTab: controller})
	} else {
		// User Profile editing myself here
		vm.SettingsBar = SetupVisibleButtons(context, ButtonBarMetaData{IsAccountSettings: true, CurrentTab: viewModel.SETTINGS_CONST_USER_PROFILE})
		vm.User = user
		err = model.AccountRoles.Query().Filter(model.Q(model.FIELD_ACCOUNTROLE_ACCOUNTID, acct.Id)).Filter(model.Q(model.FIELD_ACCOUNTROLE_USERID, vm.User.Id)).One(&vm.AccountRole)
		if err != nil {
			respond(PARAM_REDIRECT_NONE, err.Error(), PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
	}

	vm.Accounts, err = queries.Accounts.ByUserAllInclusive(vm.User, acct)
	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_USER_RELATED_ACCOUNTS, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	vm.Roles, err = queries.Roles.ByAccountWithContext(context)
	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_ACCOUNT_NOT_FOUND, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_ROLES_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)
}
