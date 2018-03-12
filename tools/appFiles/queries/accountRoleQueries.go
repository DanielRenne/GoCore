package queries

import (
	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/pkg/errors"
)

type queryAccountRoles struct{}

func (self queryAccountRoles) QueryByAccountWithContext(context session_functions.RequestContext) (q *model.Query, err error) {
	account, err := getAccount(context)
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}

	q, err = self.QueryByAccount(account)

	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}
	return
}

func (self queryAccountRoles) ByAccountWithContext(context session_functions.RequestContext) (accountRoles []model.AccountRole, err error) {
	account, err := getAccount(context)
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}
	accountRoles, err = self.ByAccount(account)
	return
}

func (self queryAccountRoles) QueryByAccount(account model.Account) (q *model.Query, err error) {
	q = model.AccountRoles.Query().In(model.Q("AccountId", account.Id.Hex()))

	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}
	return
}

func (self queryAccountRoles) ByAccount(account model.Account) (accountRoles []model.AccountRole, err error) {
	q, err := self.QueryByAccount(account)
	err = q.All(&accountRoles)
	return
}

func (self queryAccountRoles) ByAccounts(criteria map[string]interface{}) (accountRoles []model.AccountRole, err error) {
	err = model.AccountRoles.Query().In(criteria).Join("Account").All(&accountRoles)
	return
}

func (self queryAccountRoles) QueryByUserWithContext(context session_functions.RequestContext) (q *model.Query, err error) {
	user, err := getUser(context)
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}

	q, err = self.QueryByUser(user)

	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}
	return
}

func (self queryAccountRoles) ByUserWithContext(context session_functions.RequestContext) (accountRoles []model.AccountRole, err error) {
	user, err := getUser(context)
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}
	accountRoles, err = self.ByUser(user)
	return
}

func (self queryAccountRoles) QueryByUser(user model.User) (q *model.Query, err error) {
	q = model.AccountRoles.Query().Filter(model.Q(model.FIELD_ACCOUNTROLE_USERID, user.Id.Hex()))

	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}
	return
}

func (self queryAccountRoles) ByUser(user model.User) (accountRoles []model.AccountRole, err error) {
	q, err := self.QueryByUser(user)
	err = q.All(&accountRoles)
	return
}

func (self queryAccountRoles) ByUserId(userId string) (accountRole model.AccountRole, err error) {
	err = model.AccountRoles.Query().Filter(model.Q(model.FIELD_ACCOUNTROLE_USERID, userId)).One(&accountRole)
	return
}
