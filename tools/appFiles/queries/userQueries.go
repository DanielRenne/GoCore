package queries

import (
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
)

//import "github.com/DanielRenne/goCoreAppTemplate/models/v1/model"

type queryUsers struct{}

func (self queryUsers) ById(id string) (user model.User, err error) {
	err = model.Users.Query().ById(id, &user)
	return
}

func (self queryUsers) QueryByEmail(email string) *model.Query {
	return model.Users.Query().Filter(model.Q("Email", email))
}

func (self queryUsers) ByEmail(email string) (user model.User, err error) {
	q := self.QueryByEmail(email)
	err = q.One(&user)
	return
}

func (self queryUsers) RoomUserById(id string) (user model.User, err error) {
	err = model.Users.Query().Filter(model.Q(model.FIELD_USER_ID, id)).Join("Password").One(&user)
	return
}

func (self queryUsers) TestSingleUserQuery(context session_functions.RequestContext) (user model.User, err error) {
	err = model.Users.Query().Filter(model.Criteria("Id", "57dda019929d244bb4c11706")).Join("DefaultAccount").RenderViews(session_functions.GetDataFormat(context())).One(&user)
	return
}

func (self queryUsers) TestAccountRoleQuery(context session_functions.RequestContext) (accountRoles []model.AccountRole, err error) {
	filter := make(map[string]interface{}, 1)
	filter["Id"] = "57c09d20dcba0f7a0be3480f"
	// filter["Id"] = "57dda057929d244bb4c11715"
	err = model.AccountRoles.Query().In(filter).Join("User.LastUpdateUser").Join("User.DefaultAccount").RenderViews(session_functions.GetDataFormat(context())).All(&accountRoles)
	return
}