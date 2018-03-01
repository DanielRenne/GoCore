package constants

const (
	HTTP_NOT_AUTHORIZED = "NOT AUTHORIZED"

	GET_INDEX_CONTROLLER_METHOD = "Root"

	BITWISE_TRUE  = 0x01
	BITWISE_FALSE = 0x00

	BANNER_COLOR_DEFAULT = "#607d8b"
	BANNER_COLOR_OTHER   = "#404347"

	MAX_COOKIE_AGE              = 315360000
	COOKIE_AUTH_TOKEN           = "AuthToken"
	COOKIE_AUTHED               = "Authorized"
	COOKIE_AUTH_USER_ID         = "AuthUserId"
	COOKIE_AUTH_ROLE_ID         = "AuthRoleId"
	COOKIE_AUTH_ACCOUNTROLE_ID  = "AuthAccountRoleId"
	COOKIE_AUTH_ACCOUNT_ID      = "AuthAccountId"
	COOKIE_AUTH_REDIRECTION_URL = "AuthRedirectionURL"
	COOKIE_DATA_FORMAT          = "DataFormat"
	COOKIE_DATE_CREATED         = "DateCreated"
	MGO_RECORD_NOT_FOUND        = "not found"

	BCRYPT_COST = 12

	APP_CONSTANTS_INVALID_AUTH       = "InvalidAuth"
	APP_CONSTANTS_ACCOUNT_LOCKED     = "LoginPageLockLimitReached"
	APP_CONSTANTS_AUTHORIZED         = "Authorized"
	APP_CONSTANTS_LOGGED_OUT         = "LoggedOut"
	APP_CONSTANTS_USERS_ANONYMOUS_ID = "57d9b383dcba0f51172f1f57"
	APP_CONSTANTS_CRONJOB_ID         = "5835eb61e9f1283d495114c1"

	VIEWMODEL_SNACKBAR_TRANSACTION = "SnackBarUndoTransactionId"
	VIEWMODEL_SNACKBAR_MESSAGE     = "SnackbarMessage"
	VIEWMODEL_SNACKBAR_OPEN        = "SnackbarOpen"
	VIEWMODEL_SNACKBAR_TYPE        = "SnackbarType"

	PARAM_SNACKBAR_TYPE_SUCCESS = ""
	PARAM_SNACKBAR_TYPE_WARNING = "Warning"
	PARAM_SNACKBAR_TYPE_ERROR   = "Error"

	PARAM_SUCCESS                        = true
	PARAM_FAILED                         = false
	PARAM_REDIRECT_REFRESH               = "refresh"
	PARAM_REDIRECT_RERENDER              = "rerender"
	PARAM_REDIRECT_REFRESH_HOME          = "homeRefresh"
	PARAM_REDIRECT_BACK                  = "back"
	PARAM_REDIRECT_NONE                  = ""
	SNACKBAR_TYPE_SUCCESS                = ""
	PARAM_SNACKBAR_MESSAGE_NONE          = ""
	PARAM_SNACKBAR_TYPE_DOWNLOAD_CONTENT = "DownloadContent"
	PARAM_SNACKBAR_TYPE_DOWNLOAD_FILE    = "DownloadFile"
	PARAM_SNACKBAR_TYPE_ALERT            = "Alert"
	PARAM_SNACKBAR_TYPE_POPUP            = "Popup"
	PARAM_TRANSACTION_ID_NONE            = ""

	VIEWMODEL_DIALOG_MESSAGE = "DialogMessage"
	VIEWMODEL_DIALOG_OPEN    = "DialogOpen"
	VIEWMODEL_DIALOG_TITLE   = "DialogTitle"

	VIEWMODEL_DIALOG_SUCCESS_TITLE = "DialogSuccess"
	VIEWMODEL_DIALOG_WARNING_TITLE = "DialogWarning"
	VIEWMODEL_DIALOG_ERROR_TITLE   = "DialogError"

	ERRORS_USER_RELATED_ACCOUNTS         = "ErrorUserRelatedAccounts"
	ERRORS_GENERIC_DB                    = "GenericDatabaseErrorOccurred"
	SUCCESS_USER_UPDATED_DEFAULT_ACCOUNT = "SuccessUserSaveDefaultAccount"

	ERRORS_ACCOUNT_COULD_NOT_CREATE_COMPANY = "ErrorAccountCompanyFailed"
	ERRORS_COUNTRY                          = "GlobalNoCountriesFound"

	ERRORS_TRANSACTION_FAILED_TO_RETRIEVE = "ErrorTranFailedRetrieve"
	ERRORS_TRANSACTION_FAILED_TO_COMMIT   = "ErrorTranCommitFailed"
	ERRORS_TRANSACTION_FAILED_POST_HOOK   = "PostHookFailure"
	ERRORS_TRANSACTION_FAILED_TO_ROLLBACK = "ErrorTranRollbackFailed"

	ERRORS_ROLES_FAILED_TO_RETRIEVE       = "ErrorRolesFailedRetrieve"
	ERRORS_ACCOUNTROLE_FAILED_TO_RETRIEVE = "ErrorAccountRoleFailedRetrieve"
	ERRORS_ACCOUNT_FAILED_TO_RETRIEVE     = "ErrorAccountFailedRetrieve"
	ERRORS_COMPANY_FAILED_TO_RETRIEVE     = "ErrorCompanyFailedRetrieve"
	ERRORS_DELETING_CURRENT_ACCOUNT       = "ErrorsDeletingCurrentAccount"
	ERRORS_UNAUTHORIZED_DELETE            = "ErrorUnauthorizedDelete"
	ERRORS_ACCOUNT_ROLE_DELETE            = "ErrorDeleteAccountRole"

	ERRORS_ACCOUNT_INVITATION   = "ErrorAccountInvitation"
	ERRORS_PASSWORD_RESET       = "ErrorPasswordReset"
	ERRORS_PASSWORD_RESET_EMAIL = "ErrorPasswordResetEmail"

	ERROR_REQUIRED_FIELD = "ErrorRequiredField"

	USER_CONST_SAVECHANGES = "SaveChanges"
	PASSWORD_RESET_SUCCESS = "PasswordResetSuccess"
)
