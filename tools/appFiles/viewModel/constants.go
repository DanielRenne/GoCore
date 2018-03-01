package viewModel

import (
	"reflect"
)

const (
	//-DONT-REMOVE-NEW-CONST
	VIEWMODEL_FILEOBJECTADD        = "fileObjectAdd"
	VIEWMODEL_FILEOBJECTLIST       = "fileObjectList"
	VIEWMODEL_FILEOBJECTMODIFY     = "fileObjectModify"
	VIEWMODEL_FILEOBJECTS          = "fileObjects"
	VIEWMODEL_ROLEADD              = "roleAdd"
	VIEWMODEL_ROLELIST             = "roleList"
	VIEWMODEL_ROLEMODIFY           = "roleModify"
	VIEWMODEL_ROLES                = "roles"
	VIEWMODEL_SETTINGS             = "settings"
	VIEWMODEL_FEATUREGROUPADD      = "featureGroupAdd"
	VIEWMODEL_FEATUREGROUPLIST     = "featureGroupList"
	VIEWMODEL_FEATUREGROUPMODIFY   = "featureGroupModify"
	VIEWMODEL_FEATUREGROUPS        = "featureGroups"
	VIEWMODEL_ROLEFEATUREADD       = "roleFeatureAdd"
	VIEWMODEL_ROLEFEATURELIST      = "roleFeatureList"
	VIEWMODEL_ROLEFEATUREMODIFY    = "roleFeatureModify"
	VIEWMODEL_ROLEFEATURES         = "roleFeatures"
	VIEWMODEL_FEATUREADD           = "featureAdd"
	VIEWMODEL_FEATURELIST          = "featureList"
	VIEWMODEL_FEATUREMODIFY        = "featureModify"
	VIEWMODEL_FEATURES             = "features"
	VIEWMODEL_APPERRORADD          = "appErrorAdd"
	VIEWMODEL_APPERRORLIST         = "appErrorList"
	VIEWMODEL_APPERRORMODIFY       = "appErrorModify"
	VIEWMODEL_APPERRORS            = "appErrors"
	IMPORT_REQUIRED                = true
	IMPORT_NOT_REQUIRED            = false
	VIEWMODEL_TRANSACTIONADD       = "transactionAdd"
	VIEWMODEL_TRANSACTIONLIST      = "transactionList"
	VIEWMODEL_TRANSACTIONMODIFY    = "transactionModify"
	VIEWMODEL_TRANSACTIONS         = "transactions"
	VIEWMODEL_LOGIN                = "login"
	VIEWMODEL_HOME                 = "home"
	VIEWMODEL_APP                  = "app"
	VIEWMODEL_USERS                = "users"
	VIEWMODEL_FILEOBJECT           = "fileObject"
	VIEWMODEL_PASSWORDRESET        = "passwordReset"
	VIEWMODEL_USERPROFILE          = "userProfile"
	VIEWMODEL_USERADD              = "userAdd"
	VIEWMODEL_ACCOUNTADD           = "accountAdd"
	VIEWMODEL_USERLIST             = "userList"
	VIEWMODEL_SERVERSETTINGSMODIFY = "serverSettingsModify"
	VIEWMODEL_USERMODIFY           = "userModify"
	VIEWMODEL_ACCOUNTMODIFY        = "accountModify"
	VIEWMODEL_ACCOUNTLIST          = "accountList"
	VIEWMODEL_ACCOUNTS             = "accounts"
	VIEWMODEL_SNACKBAR_MESSAGE     = "SnackbarMessage"
	VIEWMODEL_SNACKBAR_OPEN        = "SnackbarOpen"
	VIEWMODEL_SNACKBAR_TYPE        = "SnackbarType"
)

func setConstants(v interface{}, identifier string) {
	vm := reflect.ValueOf(v).Elem()
	viewConstants := vm.FieldByName("Constants")

	var constants map[string]string
	switch identifier {
	case "SETTINGS_CONST":
		constants = SETTINGS_CONST
	}

	for key, value := range constants {
		field := viewConstants.FieldByName(key)
		if field.String() == "<invalid Value>" {
			continue
		}
		viewConstants.FieldByName(key).Set(reflect.ValueOf(value))
	}
}
