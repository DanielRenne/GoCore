package viewModel

import (
	"encoding/json"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
)

type RoleFeatureListViewModel struct {
	Constants struct {
	} `json:"constants"`
	RoleFeatures []model.RoleFeature             `json:"RoleFeatures"`
	WidgetList   WidgetListUserControlsViewModel `json:"WidgetList"`
	FileUpload   FileObject                      `json:"FileUpload"`

	DeletedRoleFeatures []model.RoleFeature `json:"DeletedRoleFeatures"`

	//AdditionalConstructs
}

func (this *RoleFeatureListViewModel) LoadDefaultState() {
	setConstants(this, "ROLEFEATURELIST_CONST")
}

func (self *RoleFeatureListViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
