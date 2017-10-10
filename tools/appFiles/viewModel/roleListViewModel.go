package viewModel

import (
	"encoding/json"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
)

type RoleListViewModel struct {
	Constants struct {
	} `json:"constants"`
	Roles        []model.Role                    `json:"Roles"`
	WidgetList   WidgetListUserControlsViewModel `json:"WidgetList"`
	FileUpload   FileObject                      `json:"FileUpload"`
	DeletedRoles []model.Role                    `json:"DeletedRoles"`
	SettingsBar  SettingsButtonBarViewModel      `json:"SettingsBar"`

	//AdditionalConstructs
}

func (this *RoleListViewModel) LoadDefaultState() {
	setConstants(this, "ROLELIST_CONST")
}

func (self *RoleListViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}

type RoleListFilterModel struct {
	FeatureKey string `json:"FeatureKey"`
}

func (self *RoleListFilterModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
