package queries

import (
	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/pkg/errors"
)

var Accounts queryAccounts
var AccountRoles queryAccountRoles
var Roles queryRoles
var ServerSettings queryServerSettings
var Users queryUsers
var PasswordResets queryPasswordResets
var Passwords queryPasswords
var Transactions queryTransactions
var AppContent queryAppContent
var ActivityLogs queryActivityLogs
var Common queryCommon

func getUser(context session_functions.RequestContext) (user model.User, err error) {
	if context == nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(errors.New("nil gin context")))
		return
	}
	user, err = session_functions.GetSessionUser(context())

	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}

	return user, err
}

func getAccount(context session_functions.RequestContext) (account model.Account, err error) {
	if context == nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(errors.New("nil gin context")))
		return
	}
	account, err = session_functions.GetSessionAccount(context())

	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}

	return account, err
}

type queryCommon struct {
}

func (self queryCommon) GetUser(context session_functions.RequestContext) (user model.User, err error) {
	user, err = getUser(context)
	return
}

func (self queryCommon) GetAccount(context session_functions.RequestContext) (acct model.Account, err error) {
	acct, err = getAccount(context)
	return
}
