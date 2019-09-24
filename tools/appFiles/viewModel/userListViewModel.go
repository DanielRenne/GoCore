package viewModel

import (
	"encoding/json"

	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
)

type UserListViewModel struct {
	SettingsBar SettingsButtonBarViewModel      `json:"SettingsBar"`
	WidgetList  WidgetListUserControlsViewModel `json:"WidgetList"`
	Users       []model.AccountRole             `json:"Users"`
	Roles       []model.Role                    `json:"Roles"`
}

func (this *UserListViewModel) LoadDefaultState() {
	setConstants(this, "USERLIST_CONST")
	this.Users = make([]model.AccountRole, 0)
}

func (self *UserListViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
