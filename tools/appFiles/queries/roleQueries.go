package queries

import (
	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/pkg/errors"
)

type queryRoles struct{}

func (self queryRoles) QueryByAccountWithContext(context session_functions.RequestContext) (q *model.Query) {
	account, err := getAccount(context)
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}

	q = self.QueryByAccount(account)
	return
}

func (self queryRoles) ByAccountWithContext(context session_functions.RequestContext) (roles []model.Role, err error) {
	account, err := getAccount(context)
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}
	roles, err = self.ByAccount(account)
	return
}

func (self queryRoles) QueryByAccount(account model.Account) (q *model.Query) {
	q = model.Roles.Query().In(model.Q("AccountId", []string{"", account.Id.Hex()}))
	return
}

func (self queryRoles) ByAccount(account model.Account) (roles []model.Role, err error) {
	q := self.QueryByAccount(account).Filter(model.Q(model.FIELD_ROLE_ACCOUNTTYPE, account.AccountTypeShort))
	q.All(&roles)
	return
}