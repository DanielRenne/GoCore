package controllers

import (
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
)

func (self *AppErrorsController) CreateAppError(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm viewModel.AppErrorModifyViewModel
	vm.Parse(state)
	account, err := session_functions.GetSessionAccount(context())
	if err != nil {
		respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
		return
	}

	user, err := session_functions.GetSessionUser(context())
	if err != nil {
		respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
		return
	}

	t, err := session_functions.StartTransaction(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
		return
	}
	vm.AppError.AccountId = account.Id.Hex()
	vm.AppError.UserId = user.Id.Hex()
	r := CreateAppErrorRow(context, &vm, t)
	if !r.CompletedSuccessfully {
		respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
		return
	}

	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
		return
	}

	hookSuccess := AppErrorPostCommitHook("CreateAppError", context, vm.AppError.Id.Hex())
	if !hookSuccess {
		err = queries.Transactions.Rollback(context, "AppErrorPostCommitHook failed", t.Id.Hex())
		if err != nil {
			respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
			return
		}
		respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
		return
	}
	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
}

func (self *AppErrorsController) DeleteManyAppErrors(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm viewModel.AppErrorListViewModel
	vm.Parse(state)
	t, err := session_functions.StartTransaction(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
		return
	}
	var AppErrorId string
	for i := 0; i < len(vm.DeletedAppErrors); i++ {
		var vmModify viewModel.AppErrorModifyViewModel
		vmModify.AppError = vm.DeletedAppErrors[i]
		AppErrorId = vmModify.AppError.Id.Hex()
		r := DeleteAppErrorRow(context, &vmModify, t)
		if !r.CompletedSuccessfully {
			respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
			return
		}
	}
	err = t.Commit()

	if err != nil {
		session_functions.Dump(err)
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	hookSuccess := AppErrorPostCommitHook("DeleteManyAppErrors", context, AppErrorId)
	if !hookSuccess {
		err = queries.Transactions.Rollback(context, "AppErrorPostCommitHook failed", t.Id.Hex())
		if err != nil {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_POST_HOOK, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	respond(PARAM_REDIRECT_RERENDER, constants.APPERROR_MANY_DELETE, PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), vm)
}

func (self *AppErrorsController) DeleteAppError(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm viewModel.AppErrorModifyViewModel
	vm.Parse(state)

	t, err := session_functions.StartTransaction(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
		return
	}

	r := DeleteAppErrorRow(context, &vm, t)
	if !r.CompletedSuccessfully {
		respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
		return
	}

	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	hookSuccess := AppErrorPostCommitHook("DeleteAppError", context, vm.AppError.Id.Hex())
	if !hookSuccess {
		err = queries.Transactions.Rollback(context, "AppErrorPostCommitHook failed", t.Id.Hex())
		if err != nil {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_POST_HOOK, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
}
