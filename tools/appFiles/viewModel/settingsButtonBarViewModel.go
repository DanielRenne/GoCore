package viewModel

import (
	"encoding/json"
)

var SETTINGS_CONST map[string]string

//MoreSettingsVars
var SETTINGS_CONST_ROLE string
var SETTINGS_CONST_ROLE_ADD string
var SETTINGS_CONST_ROLE_MODIFY string
var SETTINGS_CONST_ACCOUNT_INSTALLS string
var SETTINGS_CONST_ACCOUNT_INSTALL_ADD string
var SETTINGS_CONST_ACCOUNT_INSTALL_MODIFY string
var SETTINGS_CONST_USERS string
var SETTINGS_CONST_USER_MODIFY string
var SETTINGS_CONST_USER_ADD string
var SETTINGS_CONST_USER_PROFILE string
var SETTINGS_CONST_SERVER_SETTINGS string

func init() {
	// The following settings make up a TITLE which is visible in the ButtonBar and an Action which is passed to the getvars
	SETTINGS_CONST = make(map[string]string)
	//MoreSettingsInit
	SETTINGS_CONST_ROLE = "RoleList"
	SETTINGS_CONST[SETTINGS_CONST_ROLE] = SETTINGS_CONST_ROLE
	SETTINGS_CONST_ROLE_ADD = "RoleAdd"
	SETTINGS_CONST[SETTINGS_CONST_ROLE_ADD] = SETTINGS_CONST_ROLE_ADD
	SETTINGS_CONST_ROLE_MODIFY = "RoleModify"
	SETTINGS_CONST[SETTINGS_CONST_ROLE_MODIFY] = SETTINGS_CONST_ROLE_MODIFY
	SETTINGS_CONST_ACCOUNT_INSTALLS = "AccountList"
	SETTINGS_CONST[SETTINGS_CONST_ACCOUNT_INSTALLS] = SETTINGS_CONST_ACCOUNT_INSTALLS
	SETTINGS_CONST_ACCOUNT_INSTALL_ADD = "AccountAdd"
	SETTINGS_CONST[SETTINGS_CONST_ACCOUNT_INSTALL_ADD] = SETTINGS_CONST_ACCOUNT_INSTALL_ADD
	SETTINGS_CONST_ACCOUNT_INSTALL_MODIFY = "AccountModify"
	SETTINGS_CONST[SETTINGS_CONST_ACCOUNT_INSTALL_MODIFY] = SETTINGS_CONST_ACCOUNT_INSTALL_MODIFY
	SETTINGS_CONST_USERS = "UserList"
	SETTINGS_CONST[SETTINGS_CONST_USERS] = SETTINGS_CONST_USERS
	SETTINGS_CONST_USER_MODIFY = "UserModify"
	SETTINGS_CONST[SETTINGS_CONST_USER_MODIFY] = SETTINGS_CONST_USER_MODIFY
	SETTINGS_CONST_USER_ADD = "UserAdd"
	SETTINGS_CONST[SETTINGS_CONST_USER_ADD] = SETTINGS_CONST_USER_ADD
	SETTINGS_CONST_USER_PROFILE = "UserProfile"
	SETTINGS_CONST[SETTINGS_CONST_USER_PROFILE] = SETTINGS_CONST_USER_PROFILE
	SETTINGS_CONST_SERVER_SETTINGS = "ModifyServerSettings"
	SETTINGS_CONST[SETTINGS_CONST_SERVER_SETTINGS] = SETTINGS_CONST_SERVER_SETTINGS
}

type SettingsButtonBarViewModel struct {
	Constants struct {
		//MoreSettings
		RoleList             string `json:"RoleList"`
		RoleModify           string `json:"RoleModify"`
		RoleAdd              string `json:"RoleAdd"`
		AccountList          string `json:"AccountList"`
		AccountModify        string `json:"AccountModify"`
		AccountAdd           string `json:"AccountAdd"`
		UserList             string `json:"UserList"`
		UserModify           string `json:"UserModify"`
		UserProfile          string `json:"UserProfile"`
		UserAdd              string `json:"UserAdd"`
		ModifyServerSettings string `json:"ModifyServerSettings"`
	} `json:"Constants"`
	ButtonBar ButtonBar `json:"ButtonBar"`
}

func (this *SettingsButtonBarViewModel) LoadDefaultState() {
	setConstants(this, "SETTINGS_CONST")
}

func (self *SettingsButtonBarViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
