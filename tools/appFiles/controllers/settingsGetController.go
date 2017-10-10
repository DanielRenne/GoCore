package controllers

import (
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
)

type ButtonBarMetaData struct {
	CurrentTab        string
	IsServerSettings  bool
	IsAccountSettings bool
}

func (self *SettingsController) Root(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {
	var vm viewModel.SettingsViewModel
	vm.SettingsBar = SetupVisibleButtons(context, ButtonBarMetaData{IsAccountSettings: true, CurrentTab: ""})
	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)

}

func SetupVisibleButtons(context session_functions.RequestContext, meta ButtonBarMetaData) viewModel.SettingsButtonBarViewModel {
	// build this based off what the user can see in the business rules

	var vm viewModel.SettingsButtonBarViewModel
	vm.LoadDefaultState()
	vm.ButtonBar.Config.VisibleTabs = make(map[string]string)

	c := ""

	if meta.IsAccountSettings {
		if session_functions.CheckRoleAccess(context(), constants.FEATURE_ACCOUNT_VIEW) {
			c = viewModel.SETTINGS_CONST_ACCOUNT_INSTALLS
			vm.ButtonBar.Config.VisibleTabs[c] = c //vm.Constants.AccountList
			vm.ButtonBar.Config.TabActions = append(vm.ButtonBar.Config.TabActions, "Root")
			vm.ButtonBar.Config.TabControllers = append(vm.ButtonBar.Config.TabControllers, CONTROLLER_ACCOUNTLIST)
			vm.ButtonBar.Config.TabOrder = append(vm.ButtonBar.Config.TabOrder, c)
			vm.ButtonBar.Config.TabIsVisible = append(vm.ButtonBar.Config.TabIsVisible, true)
			vm.ButtonBar.Config.OtherTabSelected = append(vm.ButtonBar.Config.OtherTabSelected, "")
		}

		c = viewModel.SETTINGS_CONST_ACCOUNT_INSTALL_ADD
		vm.ButtonBar.Config.VisibleTabs[c] = c
		vm.ButtonBar.Config.TabActions = append(vm.ButtonBar.Config.TabActions, "Root")
		vm.ButtonBar.Config.TabControllers = append(vm.ButtonBar.Config.TabControllers, CONTROLLER_ACCOUNTLIST)
		vm.ButtonBar.Config.TabOrder = append(vm.ButtonBar.Config.TabOrder, c)
		vm.ButtonBar.Config.TabIsVisible = append(vm.ButtonBar.Config.TabIsVisible, false)
		vm.ButtonBar.Config.OtherTabSelected = append(vm.ButtonBar.Config.OtherTabSelected, viewModel.SETTINGS_CONST_ACCOUNT_INSTALLS)

		c = viewModel.SETTINGS_CONST_ACCOUNT_INSTALL_MODIFY
		vm.ButtonBar.Config.VisibleTabs[c] = c
		vm.ButtonBar.Config.TabActions = append(vm.ButtonBar.Config.TabActions, "Root")
		vm.ButtonBar.Config.TabControllers = append(vm.ButtonBar.Config.TabControllers, CONTROLLER_ACCOUNTLIST)
		vm.ButtonBar.Config.TabOrder = append(vm.ButtonBar.Config.TabOrder, c)
		vm.ButtonBar.Config.TabIsVisible = append(vm.ButtonBar.Config.TabIsVisible, false)
		vm.ButtonBar.Config.OtherTabSelected = append(vm.ButtonBar.Config.OtherTabSelected, viewModel.SETTINGS_CONST_ACCOUNT_INSTALLS)
	}

	if meta.IsAccountSettings {
		c = viewModel.SETTINGS_CONST_USERS
		vm.ButtonBar.Config.VisibleTabs[c] = c
		vm.ButtonBar.Config.TabActions = append(vm.ButtonBar.Config.TabActions, "Root")
		vm.ButtonBar.Config.TabControllers = append(vm.ButtonBar.Config.TabControllers, CONTROLLER_USERLIST)
		vm.ButtonBar.Config.TabOrder = append(vm.ButtonBar.Config.TabOrder, c)
		vm.ButtonBar.Config.TabIsVisible = append(vm.ButtonBar.Config.TabIsVisible, true)
		vm.ButtonBar.Config.OtherTabSelected = append(vm.ButtonBar.Config.OtherTabSelected, "")
		c = viewModel.SETTINGS_CONST_USER_ADD
		vm.ButtonBar.Config.VisibleTabs[c] = c
		vm.ButtonBar.Config.TabActions = append(vm.ButtonBar.Config.TabActions, "Root")
		vm.ButtonBar.Config.TabControllers = append(vm.ButtonBar.Config.TabControllers, CONTROLLER_USERLIST)
		vm.ButtonBar.Config.TabOrder = append(vm.ButtonBar.Config.TabOrder, c)
		vm.ButtonBar.Config.TabIsVisible = append(vm.ButtonBar.Config.TabIsVisible, false)
		vm.ButtonBar.Config.OtherTabSelected = append(vm.ButtonBar.Config.OtherTabSelected, viewModel.SETTINGS_CONST_USERS)

		c = viewModel.SETTINGS_CONST_USER_MODIFY
		vm.ButtonBar.Config.VisibleTabs[c] = c
		vm.ButtonBar.Config.TabActions = append(vm.ButtonBar.Config.TabActions, "Root")
		vm.ButtonBar.Config.TabControllers = append(vm.ButtonBar.Config.TabControllers, CONTROLLER_USERLIST)
		vm.ButtonBar.Config.TabOrder = append(vm.ButtonBar.Config.TabOrder, c)
		vm.ButtonBar.Config.TabIsVisible = append(vm.ButtonBar.Config.TabIsVisible, false)
		vm.ButtonBar.Config.OtherTabSelected = append(vm.ButtonBar.Config.OtherTabSelected, viewModel.SETTINGS_CONST_USERS)

		c = viewModel.SETTINGS_CONST_USER_PROFILE
		vm.ButtonBar.Config.VisibleTabs[c] = c
		vm.ButtonBar.Config.TabActions = append(vm.ButtonBar.Config.TabActions, "Root")
		vm.ButtonBar.Config.TabControllers = append(vm.ButtonBar.Config.TabControllers, CONTROLLER_USERPROFILE)
		vm.ButtonBar.Config.TabOrder = append(vm.ButtonBar.Config.TabOrder, c)
		vm.ButtonBar.Config.TabIsVisible = append(vm.ButtonBar.Config.TabIsVisible, true)
		vm.ButtonBar.Config.OtherTabSelected = append(vm.ButtonBar.Config.OtherTabSelected, "")

		if session_functions.CheckRoleAccess(context(), constants.FEATURE_ROLE_VIEW) {
			c = viewModel.SETTINGS_CONST_ROLE
			vm.ButtonBar.Config.VisibleTabs[c] = c
			vm.ButtonBar.Config.TabActions = append(vm.ButtonBar.Config.TabActions, "Root")
			vm.ButtonBar.Config.TabControllers = append(vm.ButtonBar.Config.TabControllers, CONTROLLER_ROLELIST)
			vm.ButtonBar.Config.TabOrder = append(vm.ButtonBar.Config.TabOrder, c)
			vm.ButtonBar.Config.TabIsVisible = append(vm.ButtonBar.Config.TabIsVisible, true)
			vm.ButtonBar.Config.OtherTabSelected = append(vm.ButtonBar.Config.OtherTabSelected, "")
		}
		c = viewModel.SETTINGS_CONST_ROLE_ADD
		vm.ButtonBar.Config.VisibleTabs[c] = c
		vm.ButtonBar.Config.TabActions = append(vm.ButtonBar.Config.TabActions, "Root")
		vm.ButtonBar.Config.TabControllers = append(vm.ButtonBar.Config.TabControllers, CONTROLLER_ROLELIST)
		vm.ButtonBar.Config.TabOrder = append(vm.ButtonBar.Config.TabOrder, c)
		vm.ButtonBar.Config.TabIsVisible = append(vm.ButtonBar.Config.TabIsVisible, false)
		vm.ButtonBar.Config.OtherTabSelected = append(vm.ButtonBar.Config.OtherTabSelected, viewModel.SETTINGS_CONST_ROLE)

		c = viewModel.SETTINGS_CONST_ROLE_MODIFY
		vm.ButtonBar.Config.VisibleTabs[c] = c
		vm.ButtonBar.Config.TabActions = append(vm.ButtonBar.Config.TabActions, "Root")
		vm.ButtonBar.Config.TabControllers = append(vm.ButtonBar.Config.TabControllers, CONTROLLER_ROLELIST)
		vm.ButtonBar.Config.TabOrder = append(vm.ButtonBar.Config.TabOrder, c)
		vm.ButtonBar.Config.TabIsVisible = append(vm.ButtonBar.Config.TabIsVisible, false)
		vm.ButtonBar.Config.OtherTabSelected = append(vm.ButtonBar.Config.OtherTabSelected, viewModel.SETTINGS_CONST_ROLE)

	}

	if meta.IsServerSettings {
		if session_functions.CheckRoleAccess(context(), constants.FEATURE_SERVER_SETTING_MODIFY) {
			c = viewModel.SETTINGS_CONST_SERVER_SETTINGS
			vm.ButtonBar.Config.VisibleTabs[c] = c
			vm.ButtonBar.Config.TabActions = append(vm.ButtonBar.Config.TabActions, "Root")
			vm.ButtonBar.Config.TabControllers = append(vm.ButtonBar.Config.TabControllers, CONTROLLER_SERVERSETTINGSMODIFY)
			vm.ButtonBar.Config.TabOrder = append(vm.ButtonBar.Config.TabOrder, c)
			vm.ButtonBar.Config.TabIsVisible = append(vm.ButtonBar.Config.TabIsVisible, true)
			vm.ButtonBar.Config.OtherTabSelected = append(vm.ButtonBar.Config.OtherTabSelected, "")
		}
	}

	if meta.CurrentTab == "" {
		// Add logic based on role here to initialize the first landing page and button
		if meta.IsServerSettings {
			meta.CurrentTab = viewModel.SETTINGS_CONST_SERVER_SETTINGS
		}
		if meta.IsAccountSettings {
			meta.CurrentTab = viewModel.SETTINGS_CONST_ACCOUNT_INSTALLS
		}

	}

	// set this based off the tab on the far left for the user
	vm.ButtonBar.Config.CurrentTab = meta.CurrentTab
	return vm
}
