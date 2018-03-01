package controllers

import (
	_ "fmt"

	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
)

const (
	ACCOUNTS_URI_PARAM_ACCOUNT_ID = "accountId"
	ACCOUNT_TYPE_LONG_YOURGOCORE  = "Long Type"
	ACCOUNT_TYPE_SHORT_YOURGOCORE = "short..."
)

func (self *AccountsController) UpdateAccountDetails(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm viewModel.AccountModifyViewModel
	t, err := session_functions.StartTransaction(context())
	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
		return
	}
	vm.Parse(state)
	if !session_functions.CheckRoleAccess(context(), constants.FEATURE_ACCOUNT_MODIFY) {
		respond(PARAM_REDIRECT_NONE, "NoAccess", PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	valid := constants.BITWISE_TRUE
	if vm.Account.StateName == "" && vm.Account.StateId == "" {
		valid &= constants.BITWISE_FALSE
		vm.Account.Errors.StateName = constants.ERROR_REQUIRED_FIELD
		vm.Account.Errors.StateId = constants.ERROR_REQUIRED_FIELD
	} else if vm.Account.StateName == "" {
		valid &= constants.BITWISE_FALSE
		vm.Account.Errors.StateName = constants.ERROR_REQUIRED_FIELD
	}

	if valid != constants.BITWISE_TRUE {
		respond(PARAM_REDIRECT_NONE, PARAM_REDIRECT_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	r := UpdateAccountRow(context, &vm, t)
	if !r.CompletedSuccessfully {
		respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
		return
	}

	err = t.Commit()
	if err != nil {
		r = session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
	}

	ActivityLog.UpsertActivityByContext(context, queries.ACTIVITY_ACCOUNT_ENTITY, vm.Account.Id.Hex(), queries.ACTIVITY_ACTION_ACCESS, "")

	if !vmAccountAddEditGetCommon(context, respond, &vm) {
		return
	}

	respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
}

func (self *AccountsController) CreateAccount(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm viewModel.AccountModifyViewModel
	vm.Parse(state)
	if !session_functions.CheckRoleAccess(context(), constants.FEATURE_ACCOUNT_ADD) {
		respond(PARAM_REDIRECT_NONE, "NoAccess", PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	valid := constants.BITWISE_TRUE
	if vm.Account.StateName == "" && vm.Account.StateId == "" {
		valid &= constants.BITWISE_FALSE
		vm.Account.Errors.StateName = constants.ERROR_REQUIRED_FIELD
		vm.Account.Errors.StateId = constants.ERROR_REQUIRED_FIELD
	} else if vm.Account.StateName == "" {
		valid &= constants.BITWISE_FALSE
		vm.Account.Errors.StateName = constants.ERROR_REQUIRED_FIELD
	}

	if valid != constants.BITWISE_TRUE {
		respond(PARAM_REDIRECT_NONE, PARAM_REDIRECT_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	t, err := session_functions.StartTransaction(context())

	if !vmAccountAddEditGetCommon(context, respond, &vm) {
		respond(PARAM_REDIRECT_NONE, "Failed common", PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
		return
	}

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
		return
	}

	r := CreateAccountRow(context, &vm, t)

	if !r.CompletedSuccessfully {
		respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
		return
	}

	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	ActivityLog.UpsertActivityByContext(context, queries.ACTIVITY_ACCOUNT_ENTITY, vm.Account.Id.Hex(), queries.ACTIVITY_ACTION_ACCESS, "")
	respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)

}

func (self *AccountsController) DeleteAccount(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm viewModel.AccountModifyViewModel
	vm.Parse(state)
	if !session_functions.CheckRoleAccess(context(), constants.FEATURE_ACCOUNT_DELETE) {
		respond(PARAM_REDIRECT_NONE, "NoAccess", PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	t, err := session_functions.StartTransaction(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
		return
	}

	act, err := session_functions.GetSessionAccount(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_ACCOUNTROLE_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	acct, err := queries.Accounts.ById(vm.Account.Id.Hex())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_DELETING_CURRENT_ACCOUNT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	var users []model.User
	err = model.Users.Query().Filter(model.Q(model.FIELD_USER_DEFAULTACCOUNTID, vm.Account.Id.Hex())).All(&users)

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_ACCOUNT_DELETE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	if len(users) > 0 {

		for i, _ := range users {
			user := users[i]
			//For these users if they have a 2nd account we can update their default account.  If they have no other accounts they will simply be orphan users.
			accountRoles, err := queries.AccountRoles.ByUser(user)
			if err == nil {
				if len(accountRoles) == 1 {
					user.DefaultAccountId = ""
					user.SaveWithTran(t)
				} else {
					for _, ar := range accountRoles {
						if ar.AccountId != acct.Id.Hex() {
							user.DefaultAccountId = ar.AccountId
							user.SaveWithTran(t)
							break
						}
					}
				}
			}
		}
	}

	if vm.Account.Id == act.Id {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_DELETING_CURRENT_ACCOUNT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	acct, err = queries.Accounts.ById(vm.Account.Id.Hex())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_DELETING_CURRENT_ACCOUNT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	err = acct.DeleteWithTran(t)

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_ACCOUNT_DELETE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	actRole, err := queries.AccountRoles.ByAccount(vm.Account)

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_ACCOUNTROLE_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	for i, _ := range actRole {
		actRole[i].DeleteWithTran(t)
		if err != nil {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_ACCOUNT_ROLE_DELETE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
	}

	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	respond(PARAM_REDIRECT_BACK, constants.ACCOUNT_DELETE_SUCCESSFUL, PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), vm)
}
