package controllers

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
	"github.com/pkg/errors"
)

func (self *PasswordResetController) Reset(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {

	var vm viewModel.PasswordResetViewModel
	vm.Parse(state)

	if vm.PasswordReset.Complete == true {
		respond(PARAM_REDIRECT_REFRESH_HOME, "", SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	vm.ConfirmPasswordErrors = ""
	vm.PasswordErrors = ""

	if vm.Password == "" {
		vm.PasswordErrors = "FieldRequired"
	}
	if vm.ConfirmPassword == "" {
		vm.ConfirmPasswordErrors = "FieldRequired"
	}

	if vm.Password == "" || vm.ConfirmPassword == "" {
		respond(PARAM_REDIRECT_NONE, "Validation Error", SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	if vm.Password != vm.ConfirmPassword {
		vm.PasswordErrors = "PasswordsNoMatch"
		vm.ConfirmPasswordErrors = "PasswordsNoMatch"
		respond(PARAM_REDIRECT_NONE, "Validation Error", SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	if !verifyPasswordCheck(vm.Password) {
		vm.PasswordErrors = constants.ERROR_PASSWORD_NOTVALID
		vm.ConfirmPasswordErrors = constants.ERROR_PASSWORD_NOTVALID
		respond(PARAM_REDIRECT_NONE, "Validation Error", SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	user, err := queries.Users.ById(vm.PasswordReset.UserId)
	if err != nil {
		respond(PARAM_REDIRECT_NONE, "Failed to find User Profile.", SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	password, err := queries.Passwords.ById(user.PasswordId)
	if err != nil {
		respond(PARAM_REDIRECT_NONE, "Failed to find Password Record.", SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	bcryptData, err := bcrypt.GenerateFromPassword([]byte(vm.Password), constants.BCRYPT_COST)
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		message := "Failed to generate a hash for password."
		respond(PARAM_REDIRECT_NONE, message, SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	password.Value = string(bcryptData)
	vm.PasswordReset.Complete = true
	user.LoginAttempts = 0
	user.EnforcePasswordChange = false

	t, err := session_functions.StartTransaction(nil)
	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	err = user.SaveWithTran(t)

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_PASSWORD_SAVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	err = password.SaveWithTran(t)

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_PASSWORD_SAVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	err = vm.PasswordReset.SaveWithTran(t)

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_PASSWORD_SAVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	respond(PARAM_REDIRECT_REFRESH_HOME, "Password reset successfully.", SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)
}
