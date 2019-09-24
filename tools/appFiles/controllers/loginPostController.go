package controllers

import (
	"encoding/base64"
	"os"
	"runtime/trace"
	"strings"
	"time"

	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/ginServer"
	"github.com/DanielRenne/goCoreAppTemplate/br"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/settings"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
)

func (self *LoginController) ShutDownServer(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	if settings.AppSettings.DeveloperGoTrace {
		trace.Stop()
	}
	os.Exit(0)
}

func (self *LoginController) Authorize(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {

	var vm viewModel.LoginViewModel
	vm.Parse(state)

	var validationFailed bool

	if vm.Username == "" {
		validationFailed = true
		vm.UserNameError = model.VALIDATION_ERROR_SPECIFIC_REQUIRED
	}

	if vm.Password == "" {
		validationFailed = true
		vm.PasswordError = model.VALIDATION_ERROR_SPECIFIC_REQUIRED
	}

	if validationFailed {
		respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	maxLoginAttempts, err := queries.ServerSettings.LoginAttempts()
	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_GENERIC_DB, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	user, authorized := br.Users.Authorize(vm, extensions.StringToInt(maxLoginAttempts.Value))
	if !authorized {
		if user != nil {
			maxLoginlockLimit := extensions.StringToInt(maxLoginAttempts.Value)
			if maxLoginlockLimit > 0 && user.LoginAttempts > maxLoginlockLimit {
				vm.AuthMessage = "LoginPageLockLimitReached"
				if !user.Locked {
					t, _ := model.Transactions.New(constants.APP_CONSTANTS_USERS_ANONYMOUS_ID)
					user.Locked = true
					user.SaveWithTran(t)
					t.Commit()
				}
				respond(PARAM_REDIRECT_NONE, constants.APP_CONSTANTS_ACCOUNT_LOCKED, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
				return
			} else {
				t, _ := model.Transactions.New(constants.APP_CONSTANTS_USERS_ANONYMOUS_ID)
				user.LoginAttempts++
				user.SaveWithTran(t)
				t.Commit()

			}
		}

		vm.AuthMessage = "LoginPageAuthMessageInvalid"
		respond(PARAM_REDIRECT_NONE, constants.APP_CONSTANTS_INVALID_AUTH, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	vm.AuthMessage = ""
	self.SetSession(context, user)

	respond(PARAM_REDIRECT_REFRESH, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)

}

func (self *LoginController) SetSession(context session_functions.RequestContext, user *model.User) {
	c := context()
	ginServer.SetSessionKey(c, constants.COOKIE_AUTH_TOKEN, constants.COOKIE_AUTHED)
	ginServer.SetSessionKey(c, constants.COOKIE_AUTH_USER_ID, user.Id.Hex())
	ginServer.SetSessionKey(c, constants.COOKIE_DATE_CREATED, time.Now().String())

	//Pick the first related account or Default account
	accountRoleId := ""

	if user.DefaultAccountId != "" {

		var accountRole model.AccountRole
		err := model.AccountRoles.Query().In(model.Q("UserId", user.Id.Hex())).Filter(model.Q("AccountId", user.DefaultAccountId)).One(&accountRole)
		if err == nil {
			accountRoleId = accountRole.Id.Hex()
		} else { //Find the First Account they have access to.
			model.AccountRoles.Query().In(model.Q("UserId", user.Id.Hex())).Limit(1).One(&accountRole)
			accountRoleId = accountRole.Id.Hex()
		}
	}

	var accountRole model.AccountRole

	err := model.AccountRoles.Query().ById(accountRoleId, &accountRole)

	if err == nil {
		ginServer.SetSessionKey(c, constants.COOKIE_AUTH_ACCOUNT_ID, accountRole.AccountId)
		ginServer.SetSessionKey(c, constants.COOKIE_AUTH_ROLE_ID, accountRole.RoleId)
		ginServer.SetSessionKey(c, constants.COOKIE_AUTH_ACCOUNTROLE_ID, accountRole.Id.Hex())

		// Add this later when things are stable in production
		//
		//"github.com/DanielRenne/goCoreAppTemplate/settings"
		//	if settings.AppSettings.DeveloperMode && accountRole.RoleId == constants.ROLE_SUPER {
		//		go func() {
		//			_ = notifications.SMTP.Send([]string{user.Email}, "Login", "Please ensure it was you who just logged in with your username and password.  If this is not the case, please contact the administrator")
		//		}()
		//	}

	}

	ginServer.SaveSession(c)
}

func (self *LoginController) Logout(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {

	c := context()

	var vm viewModel.AppViewModel
	vm.Parse(state)
	ginServer.ClearSession(c)
	ginServer.SaveSession(c)
	respond(PARAM_REDIRECT_REFRESH, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)
}

func GetRequestorIp(context session_functions.RequestContext) (ip string) {
	ip = context().Request.RemoteAddr
	idx := strings.Index(ip, ":")
	if idx != -1 {
		ip = ip[0:idx]
	}
	return ip
}

func (self *LoginController) ForgotPassword(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm viewModel.LoginViewModel
	vm.Parse(state)
	c := context()

	var validationFailed bool

	if vm.Username == "" {
		validationFailed = true
		vm.UserNameError = model.VALIDATION_ERROR_SPECIFIC_REQUIRED
	}

	if validationFailed {
		respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	vm.Username = strings.ToLower(vm.Username)

	user, err := queries.Users.ByEmail(vm.Username)

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.APP_CONSTANTS_INVALID_AUTH, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	pr, err := queries.PasswordResets.ByUserId(user.Id.Hex())

	if err == nil {
		origin := c.Request.Header.Get("Origin")
		origin += "/#/passwordReset?action=Load&uriParams="
		uriParamsJSON := "{\"Id\":\"" + pr.Id.Hex() + "\"}"
		uriParamsValue := base64.StdEncoding.EncodeToString([]byte(uriParamsJSON))
		origin += uriParamsValue

		//localeLanguage := ginServer.GetLocaleLanguage(context())
		//title, body, errTranslation := queries.PasswordResets.GetPasswordResetTranslation(localeLanguage, vm.Username, origin)
		//
		//if errTranslation != nil {
		//	respond(PARAM_REDIRECT_NONE, constants.INVITATION_ERROR, PARAM_SNACKBAR_TYPE_ERROR, errTranslation, PARAM_TRANSACTION_ID_NONE, vm)
		//	return
		//}

		//errEmail := notifications.SMTP.Send([]string{vm.Username}, title, body)
		//if errEmail != nil {
		//	respond(PARAM_REDIRECT_NONE, constants.ERRORS_PASSWORD_RESET_EMAIL, PARAM_SNACKBAR_TYPE_ERROR, errEmail, PARAM_TRANSACTION_ID_NONE, vm)
		//	return
		//}

		respond(PARAM_REDIRECT_NONE, constants.PASSWORD_RESET_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	t, err := model.Transactions.New(constants.APP_CONSTANTS_USERS_ANONYMOUS_ID)

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	var passReset model.PasswordReset
	passReset.UserId = user.Id.Hex()
	err = passReset.SaveWithTran(t)

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_PASSWORD_RESET, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	origin := c.Request.Header.Get("Origin")
	origin += "/#/passwordReset?action=Load&uriParams="
	uriParamsJSON := "{\"Id\":\"" + passReset.Id.Hex() + "\"}"
	uriParamsValue := base64.StdEncoding.EncodeToString([]byte(uriParamsJSON))
	origin += uriParamsValue

	passReset.Url = origin

	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	//errEmail := notifications.SMTP.Send([]string{vm.Username}, title, body)
	//if errEmail != nil {
	//	respond(PARAM_REDIRECT_NONE, constants.ERRORS_PASSWORD_RESET_EMAIL, PARAM_SNACKBAR_TYPE_ERROR, errEmail, PARAM_TRANSACTION_ID_NONE, vm)
	//	return
	//}

	respond(PARAM_REDIRECT_NONE, constants.PASSWORD_RESET_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), vm)

}
