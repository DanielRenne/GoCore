package br

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
	"github.com/pkg/errors"
)

type usersBr struct{}

func (self usersBr) Create(context session_functions.RequestContext, vm *viewModel.UserModifyViewModel, t *model.Transaction) (message string, err error) {

	message, err = self.Validate(vm)
	if err != nil {
		return
	}

	err = vm.User.SaveWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Log("Error->Validation_Error->usersBr->Create", fmt.Sprintf("%+v", vm.User.Errors))
			message = constants.PARAM_SNACKBAR_MESSAGE_NONE
			return
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Log("Error->usersBr->Create", err.Error())
			message = queries.AppContent.GetTranslation(context, constants.ERRORS_USER_SAVE) + core.Debug.HandleError(err)
			return
		}
	}

	message = constants.USER_SAVE_SUCCESS
	return
}

func (self usersBr) Update(context session_functions.RequestContext, vm *viewModel.UserModifyViewModel, t *model.Transaction) (message string, err error) {

	message, err = self.Validate(vm)
	if err != nil {
		return
	}

	err = vm.User.SaveWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Log("Error->Validation_Error->usersBr->Update", fmt.Sprintf("%+v", vm.User.Errors))
			message = constants.PARAM_SNACKBAR_MESSAGE_NONE
			return
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Log("Error->usersBr->Update", err.Error())
			message = queries.AppContent.GetTranslation(context, constants.ERRORS_USER_SAVE) + core.Debug.HandleError(err)
			return
		}
	}

	message = constants.USER_SAVE_SUCCESS
	return
}

func (self usersBr) Delete(context session_functions.RequestContext, vm *viewModel.UserModifyViewModel, t *model.Transaction) (message string, err error) {

	err = vm.User.DeleteWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Log("Error->Validation_Error->usersBr->Delete", fmt.Sprintf("%+v", vm.User.Errors))
			message = constants.ERRORS_USER_DELETE
			return
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Log("Error->usersBr->Delete", err.Error())
			message = constants.ERRORS_USER_DELETE
			return
		}
	}

	message = constants.USER_DELETE_SUCCESSFUL
	return
}

func (self usersBr) WrapViewModel(UserInstance model.User) viewModel.UserModifyViewModel {
	vm := viewModel.UserModifyViewModel{}
	vm.User = UserInstance
	return vm
}

func (self usersBr) Validate(vm *viewModel.UserModifyViewModel) (message string, err error) {
	message = ""
	return
}

func (self usersBr) Authorize(vm viewModel.LoginViewModel, maxLoginAttempts int) (*model.User, bool) {

	var user model.User
	var err error

	filter := make(map[string]interface{}, 1)
	filter[model.FIELD_USER_EMAIL] = strings.TrimSpace(strings.ToLower(vm.Username))
	session_functions.Log("Login Authorize filter", fmt.Sprintf("%+v", filter))
	err = model.Users.Query().Filter(filter).One(&user)

	if err != nil {
		session_functions.Log("Login->authorize", err.Error())
		return nil, false
	}

	if maxLoginAttempts != 0 && user.LoginAttempts > maxLoginAttempts {
		return &user, false
	}

	//https://github.com/DanielRenne/goCoreAppTemplate/commit/055d02ab94f4ea452dbd7f18df7288c28efe2b22

	var password model.Password

	err = model.Passwords.Query().ById(user.PasswordId, &password)
	session_functions.Dump(password)

	if err != nil {
		session_functions.Log("Login->authorize err 2", err.Error())
		return nil, false
	}
	session_functions.Log("Login Authorize Check Password", fmt.Sprintf("%+v", password))

	//Now check the password.
	comparePassword := vm.Password

	err = bcrypt.CompareHashAndPassword([]byte(password.Value), []byte(comparePassword))
	if err == nil {
		session_functions.Log("Login Authorized", user.Email)
		t, _ := model.Transactions.New(constants.APP_CONSTANTS_USERS_ANONYMOUS_ID)
		user.LoginAttempts = 0
		user.LastLoginDate = time.Now()
		user.SaveWithTran(t)
		t.Commit()

		return &user, true
	} else {
		session_functions.Log("Error->Login->authorize.", "Failed to compare user password with hash:  "+err.Error())
	}

	return &user, false
}
