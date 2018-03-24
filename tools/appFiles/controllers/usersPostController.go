package controllers

import (
	"strings"

	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/password"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

func (self *UserModifyController) CommonUserValidation(vm *viewModel.UserModifyViewModel) int {
	valid := constants.BITWISE_TRUE
	if vm.User.Email == "" {
		valid &= constants.BITWISE_FALSE
		vm.User.Errors.Email = constants.ERROR_REQUIRED_FIELD
	}
	vm.User.Email = strings.ToLower(vm.User.Email)

	isNew := !bson.IsObjectIdHex(vm.User.Id.Hex())
	if isNew || vm.EmailChanged {
		qs := model.Users.Query().Filter(model.Q("Email", strings.ToLower(vm.User.Email)))
		if vm.EmailChanged {
			qs = qs.Exclude(model.Q("Id", vm.User.Id))
		}
		var users []model.User
		count, err := qs.Count(&users)
		if err != nil {
			valid &= constants.BITWISE_FALSE
		}
		if count > 0 {
			valid &= constants.BITWISE_FALSE
			vm.User.Errors.Email = constants.USER_ERROR_EMAIL_EXISTS
		}
	}

	if vm.User.First == "" {
		valid &= constants.BITWISE_FALSE
		vm.User.Errors.First = constants.ERROR_REQUIRED_FIELD
	}

	if vm.User.Last == "" {
		valid &= constants.BITWISE_FALSE
		vm.User.Errors.Last = constants.ERROR_REQUIRED_FIELD
	}

	// if vm.User.CompanyName == "" {
	// 	valid &= constants.BITWISE_FALSE
	// 	vm.User.Errors.CompanyName = constants.ERROR_REQUIRED_FIELD
	// }
	return valid
}

func (self *UserModifyController) UpdateUserDetails(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm viewModel.UserModifyViewModel
	vm.Parse(state)
	user, _ := session_functions.GetSessionUser(context())

	if user.Id.Hex() != vm.User.Id.Hex() {
		if !session_functions.CheckRoleAccess(context(), constants.FEATURE_USER_MODIFY) {
			respond(PARAM_REDIRECT_NONE, "NoAccess", PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
	}

	t, err := session_functions.StartTransaction(context())
	valid := constants.BITWISE_TRUE
	valid &= self.CommonUserValidation(&vm)
	if valid != constants.BITWISE_TRUE {
		// log.Printf("Entered Valid conditional of UpdateUserDetails.")
		respond(PARAM_REDIRECT_NONE, PARAM_REDIRECT_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	err = vm.User.SaveWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Validation Error exists", vm.User.Errors)
			respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Golang error", err)
			respond(PARAM_REDIRECT_NONE, queries.AppContent.GetTranslation(context, constants.ERRORS_USER_SAVE)+core.Debug.HandleError(err), PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		}
		return
	}

	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	session_functions.StoreDataFormat(context(), vm.User.Language, vm.User.TimeZone, vm.User.DateFormat)

	if vm.CurrentPage == "userModify" {
		respond(CONTROLLER_USERLIST, constants.USER_SAVE_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), vm)
		return
	}
	respond(PARAM_REDIRECT_NONE, constants.USER_SAVE_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), vm)
}

func (self *UserModifyController) UpdateAccountRole(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm viewModel.UserModifyViewModel
	vm.Parse(state)
	if !session_functions.CheckRoleAccess(context(), constants.FEATURE_USER_CHANGE_ROLE) {
		respond(PARAM_REDIRECT_NONE, "NoAccess", PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	t, err := session_functions.StartTransaction(context())

	err = vm.AccountRole.SaveWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Validation Error exists", vm.AccountRole.Errors)
			respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Golang error", err)
			respond(PARAM_REDIRECT_NONE, queries.AppContent.GetTranslation(context, constants.ERRORS_USER_SAVE)+core.Debug.HandleError(err), PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		}
		return
	}

	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	if vm.CurrentPage == "userModify" {
		respond(CONTROLLER_USERLIST, constants.USER_SAVE_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), vm)
		return
	}
	respond(PARAM_REDIRECT_NONE, constants.USER_SAVE_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), vm)
}

func (self *UserModifyController) ChangeUserPassword(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm viewModel.UserModifyViewModel
	vm.Parse(state)

	user, err := session_functions.GetSessionUser(context())
	if vm.User.Id.Hex() != user.Id.Hex() {
		respond(PARAM_REDIRECT_NONE, "Unauthorized", PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
	}

	t, err := session_functions.StartTransaction(context())
	valid := constants.BITWISE_TRUE
	valid &= self.CommonUserValidation(&vm)

	if vm.Password != vm.ConfirmPassword {
		valid &= constants.BITWISE_FALSE
		vm.PasswordErrors = constants.USER_ERROR_PASSWORD_MISMATCH
		vm.ConfirmPasswordErrors = constants.USER_ERROR_PASSWORD_MISMATCH
		respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	validPassword := verifyPasswordCheck(vm.Password)
	if validPassword != true {
		valid &= constants.BITWISE_FALSE
		vm.PasswordErrors = constants.ERROR_PASSWORD_NOTVALID
		respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	if vm.Password == "" || vm.ConfirmPassword == "" {
		valid &= constants.BITWISE_FALSE
		vm.PasswordErrors = constants.ERROR_REQUIRED_FIELD
		vm.ConfirmPasswordErrors = constants.ERROR_REQUIRED_FIELD
		respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	bcryptData, err := bcrypt.GenerateFromPassword([]byte(vm.Password), constants.BCRYPT_COST)
	if err != nil {
		respond(PARAM_REDIRECT_NONE, "Failed to generate a hash for the password.", PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	var password model.Password
	password.Value = string(bcryptData[:])

	err = password.SaveWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Golang error", err)
			respond(PARAM_REDIRECT_NONE, queries.AppContent.GetTranslation(context, constants.ERRORS_USER_SAVE)+core.Debug.HandleError(err), PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		}
		return
	}

	vm.User.PasswordId = password.Id.Hex()
	err = vm.User.SaveWithTran(t)
	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Validation Error exists", vm.User.Errors)
			respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Golang error", err)
			respond(PARAM_REDIRECT_NONE, queries.AppContent.GetTranslation(context, constants.ERRORS_USER_SAVE)+core.Debug.HandleError(err), PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		}
		return
	}
	t.Commit()
	if err != nil {
		respond(PARAM_REDIRECT_NONE, "Failed to commit", PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	vm.Password = ""
	vm.ConfirmPassword = ""
	respond(PARAM_REDIRECT_NONE, constants.USER_SAVE_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)
}

func (self *UserModifyController) CreateNewUser(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm viewModel.UserModifyViewModel
	vm.Parse(state)
	if !session_functions.CheckRoleAccess(context(), constants.FEATURE_USER_ADD) {
		respond(PARAM_REDIRECT_NONE, "NoAccess", PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	t, err := session_functions.StartTransaction(context())
	valid := constants.BITWISE_TRUE
	valid &= self.CommonUserValidation(&vm)

	vm.User.EnforcePasswordChange = true
	if vm.Password == "" {
		valid &= constants.BITWISE_FALSE
		vm.PasswordErrors = constants.ERROR_REQUIRED_FIELD
		vm.ConfirmPasswordErrors = constants.ERROR_REQUIRED_FIELD
	}

	if vm.AccountRole.RoleId == "" {
		valid &= constants.BITWISE_FALSE
		vm.AccountRole.Errors.RoleId = constants.ERROR_REQUIRED_FIELD
	}

	//log.Printf("Valid = %v", valid)
	if valid != constants.BITWISE_TRUE {
		//log.Printf("Entered Valid conditional of Create New User", valid)
		respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	bcryptData, err := bcrypt.GenerateFromPassword([]byte(vm.Password), constants.BCRYPT_COST)
	if err != nil {
		respond(PARAM_REDIRECT_NONE, "Failed to generate a hash for the password.", PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	var password model.Password
	password.Value = string(bcryptData[:])

	err = password.SaveWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Golang error", err)
			respond(PARAM_REDIRECT_NONE, queries.AppContent.GetTranslation(context, constants.ERRORS_USER_SAVE)+core.Debug.HandleError(err), PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		}
		return
	}

	vm.User.PasswordId = password.Id.Hex()
	vm.User.DefaultAccountId = vm.AccountRole.AccountId

	err = vm.User.SaveWithTran(t)
	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Validation Error exists", vm.User.Errors)
			respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Golang error", err)
			respond(PARAM_REDIRECT_NONE, queries.AppContent.GetTranslation(context, constants.ERRORS_USER_SAVE)+core.Debug.HandleError(err), PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		}
		return
	}

	vm.AccountRole.UserId = vm.User.Id.Hex()

	err = vm.AccountRole.SaveWithTran(t)

	if err != nil {
		respond(PARAM_REDIRECT_NONE, queries.AppContent.GetTranslation(context, constants.ERRORS_USER_SAVE)+core.Debug.HandleError(err), PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	//go func() {
	//	defer func() {
	//		if r := recover(); r != nil {
	//			session_functions.Print("\n\nPanic Stack: " + string(debug.Stack()))
	//			session_functions.Log("Panic Recovered at usersPostController email1", fmt.Sprintf("%+v", r))
	//			return
	//		}
	//	}()
	//	c := context()
	//	origin := session_functions.GetProtocol(c)
	//	origin += c.Request.Host
	//	acct, _ := queries.Accounts.ById(vm.AccountRole.AccountId)
	//
	//	replacements := queries.TagReplacements{
	//		Tag1: queries.Q("password", vm.Password),
	//		Tag2: queries.Q("account_name", acct.AccountName),
	//		Tag3: queries.Q("origin", origin),
	//	}
	//
	//	notifications.SMTP.Send([]string{vm.User.Email}, queries.AppContent.GetTranslationFromUser(vm.User, "NewUserEmailSubject"), queries.AppContent.GetTranslationWithReplacementsFromUser(vm.User, "NewUserEmailBody", &replacements))
	//}()

	respond(CONTROLLER_USERLIST, constants.USER_SAVE_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), vm)
}

func (self *UserModifyController) UserEnforcePasswordChange(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm viewModel.UserModifyViewModel
	vm.Parse(state)

	err := model.Users.Query().ById(vm.User.Id, &vm.User)
	if err != nil {
		respond(PARAM_REDIRECT_NONE, "Failed to get user", PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	valid := constants.BITWISE_TRUE
	if vm.Password != vm.ConfirmPassword {
		valid &= constants.BITWISE_FALSE
		vm.PasswordErrors = constants.USER_ERROR_PASSWORD_MISMATCH
		vm.ConfirmPasswordErrors = constants.USER_ERROR_PASSWORD_MISMATCH
	}
	validPassword := verifyPasswordCheck(vm.Password)
	if validPassword != true {
		valid &= constants.BITWISE_FALSE
		vm.PasswordErrors = constants.ERROR_PASSWORD_NOTVALID
	}
	if vm.Password == "" || vm.ConfirmPassword == "" {
		valid &= constants.BITWISE_FALSE
		vm.PasswordErrors = constants.ERROR_REQUIRED_FIELD
		vm.ConfirmPasswordErrors = constants.ERROR_REQUIRED_FIELD
	}

	if valid != constants.BITWISE_TRUE {
		respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	bcryptData, err := bcrypt.GenerateFromPassword([]byte(vm.Password), constants.BCRYPT_COST)
	if err != nil {
		respond(PARAM_REDIRECT_NONE, "Failed to generate a hash for the password.", PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	var password model.Password
	password.Value = string(bcryptData[:])

	err = password.Save()
	if err != nil {
		respond(PARAM_REDIRECT_NONE, queries.AppContent.GetTranslation(context, constants.ERRORS_USER_SAVE)+core.Debug.HandleError(err), PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	vm.User.PasswordId = password.Id.Hex()

	err = vm.User.Save()

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Validation Error exists", vm.User.Errors)
			respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Golang error", err)
			respond(PARAM_REDIRECT_NONE, queries.AppContent.GetTranslation(context, constants.ERRORS_USER_SAVE)+core.Debug.HandleError(err), PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		}
		return
	}

	vm.User.EnforcePasswordChange = false
	err = vm.User.Save()

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Validation Error exists", vm.User.Errors)
			respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Dump("Desc->Golang error", err)
			respond(PARAM_REDIRECT_NONE, queries.AppContent.GetTranslation(context, constants.ERRORS_USER_SAVE)+core.Debug.HandleError(err), PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		}
		return
	}

	respond(PARAM_REDIRECT_NONE, "PasswordChangeSuccess", PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)
}

// Verifies that the password has 8 characters, and at least 1 number, upper, lower, and special char within the password.
func verifyPasswordCheck(pass string) (isValid bool) {
	isValid = false
	number, upper, lower, special, eightOrMore := password.VerifyPassword(pass)
	if number != true {
		isValid = false
	} else if upper != true {
		isValid = false
	} else if lower != true {
		isValid = false
	} else if special != true {
		isValid = false
	} else if eightOrMore != true {
		isValid = false
	} else {
		isValid = true
	}
	// log.Printf("number: %v, \nupper: %v, \nlower: %v, \nspecial: %v, \n8orMore: %v, \nisValid: %v", number, upper, lower, special, eightOrMore, isValid)
	return isValid
}

func (self *UsersController) UpdatePreferences(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm viewModel.UserPreferences
	vm.ParsePreferences(state)
	user, err := session_functions.GetSessionUser(context())
	if err != nil {
		respond("", "", PARAM_SNACKBAR_TYPE_SUCCESS, nil, "", viewModel.EmptyViewModel{})
		return
	}
	t, err := session_functions.StartTransaction(context())
	if err != nil {
		respond("", "", PARAM_SNACKBAR_TYPE_SUCCESS, nil, "", viewModel.EmptyViewModel{})
		return
	}
	user.Preferences = vm.UserPreferences
	err = user.SaveWithTran(t)

	if err != nil {
		respond("", "", PARAM_SNACKBAR_TYPE_SUCCESS, nil, "", viewModel.EmptyViewModel{})
		return
	}

	err = t.Commit()
	if err != nil {
		respond("", "", PARAM_SNACKBAR_TYPE_SUCCESS, nil, "", viewModel.EmptyViewModel{})
		return
	}
	respond("", "", PARAM_SNACKBAR_TYPE_SUCCESS, nil, "", viewModel.EmptyViewModel{})
	return
}
