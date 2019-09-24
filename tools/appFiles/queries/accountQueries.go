package queries

import (
	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/pkg/errors"
)

type queryAccounts struct{}

func (self queryAccounts) QueryByUserWithContext(context session_functions.RequestContext) (q *model.Query, err error) {
	user, err := getUser(context)

	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}

	account, err := session_functions.GetSessionAccount(context())

	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}

	return self.QueryByUser(user, account)
}

func (self queryAccounts) ByUserAllInclusive(user model.User, account model.Account) (accounts []model.Account, err error) {
	account.Id = ""
	q, err := self.QueryByUser(user, account)
	err = q.All(&accounts)
	return
}

func (self queryAccounts) ByUserWithContext(context session_functions.RequestContext) (accounts []model.Account, err error) {
	user, err := getUser(context)

	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}

	account, err := session_functions.GetSessionAccount(context())

	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}

	return self.ByUser(user, account)
}

func (self queryAccounts) QueryByUser2(user model.User, account model.Account) (q *model.Query, err error) {
	accountIds, err := self.GetAccountIds(user, account)
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}
	session_functions.Dump(accountIds)

	q = model.Accounts.Query().InitAndOr().AddAnd().AndIn(1, model.Q("Id", accountIds))
	return
}

func (self queryAccounts) QueryByUser(user model.User, account model.Account) (q *model.Query, err error) {
	accountIds, err := self.GetAccountIds(user, account)
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}
	q = model.Accounts.Query().In(model.Q("Id", accountIds))
	return
}

func (self queryAccounts) GetAccountIds(user model.User, account model.Account) (accountIds []string, err error) {
	var accountRoles []model.AccountRole

	if account.IsSystemAccount {
		var accounts []model.Account
		err = model.Accounts.Query().All(&accounts)
		if err != nil {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			return
		}
		for _, acct := range accounts {
			accountIds = append(accountIds, acct.Id.Hex())
		}
		return
	} else {
		// normal users can see their account and sub accounts
		err = model.AccountRoles.Query().In(model.Q("UserId", user.Id.Hex())).All(&accountRoles)

		if err != nil {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			return
		}

		for _, accR := range accountRoles {
			accountIds = append(accountIds, accR.AccountId)
		}
	}
	return
}

func (self queryAccounts) ByUser(user model.User, account model.Account) (accounts []model.Account, err error) {
	q, err := self.QueryByUser(user, account)
	err = q.All(&accounts)
	return
}

//Get an Account by Id
func (self queryAccounts) ById(id string) (account model.Account, err error) {
	err = model.Accounts.Query().ById(id, &account)
	return
}

//Get an Account by Id
func (self queryAccounts) ByIds(criteria map[string]interface{}) (accounts []model.Account, err error) {
	err = model.Accounts.Query().In(criteria).Sort(model.FIELD_ACCOUNT_ACCOUNTNAME).All(&accounts)
	return
}

//Get an Account by Id
func (self queryAccounts) Query() *model.Query {
	return model.Accounts.Query()
}

//Get all Accounts associated by a CompanyId
func (self queryAccounts) QueryByCompanyId(id string) (q *model.Query, err error) {
	q = model.Accounts.Query().Filter(model.Q("CompanyId", id))
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}
	return
}

func (self queryAccounts) ByCompanyIdDDD(id string) (accounts []model.Account, err error) {
	q, err := self.QueryByCompanyId(id)
	err = q.All(&accounts)
	return
}

func (self queryAccounts) Default() (account model.Account, err error) {
	err = model.Accounts.Query().Exclude(model.Q(model.FIELD_ACCOUNT_EMAIL, "root@root-company.com")).One(&account)
	return
}
