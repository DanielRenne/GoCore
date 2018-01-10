package viewModel

import (
	"encoding/json"

	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	_ "github.com/globalsign/mgo/bson"
)

type AccountListViewModel struct {
	SettingsBar        SettingsButtonBarViewModel      `json:"SettingsBar"`
	Accounts           []model.Account                 `json:"Accounts"`
	ComponentAccounts  []model.Account                 `json:"ComponentAccounts"`
	User               model.User                      `json:"User"`
	Roles              []model.Role                    `json:"Roles"`
	ImportAccountRoles []model.AccountRole             `json:"ImportAccountRoles"`
	WidgetList         WidgetListUserControlsViewModel `json:"WidgetList"`
	FileUpload         FileObject                      `json:"FileUpload"`
}

func (this *AccountListViewModel) LoadDefaultState() {
	setConstants(this, "ACCOUNTLIST_CONST")
	this.Accounts = make([]model.Account, 0)
	this.Roles = make([]model.Role, 0)
}

func (self *AccountListViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
