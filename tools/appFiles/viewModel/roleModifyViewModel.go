package viewModel

import (
	"encoding/json"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
)

type RoleModifyViewModel struct {
	Constants struct {
	} `json:"constants"`
	Role            model.Role                 `json:"Role"`
	FeatureGroups   []model.FeatureGroup       `json:"FeatureGroups"`
	FeaturesEnabled map[string]bool            `json:"FeaturesEnabled"`
	FileUpload      FileObject                 `json:"FileUpload"`
	SettingsBar     SettingsButtonBarViewModel `json:"SettingsBar"`
	//AdditionalConstructs
}

func (this *RoleModifyViewModel) LoadDefaultState() {
	setConstants(this, "ROLEMODIFY_CONST")
	this.FeaturesEnabled = make(map[string]bool, 0)
}

func (self *RoleModifyViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
