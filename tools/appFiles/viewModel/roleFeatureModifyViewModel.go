package viewModel

import (
	"encoding/json"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
)

type RoleFeatureModifyViewModel struct {
	Constants struct {
	} `json:"constants"`
	RoleFeature model.RoleFeature `json:"RoleFeature"`
	Roles       []model.Role      `json:"Roles"`
	Features    []model.Feature   `json:"Features"`
	FileUpload  FileObject        `json:"FileUpload"`
	//AdditionalConstructs
}

func (this *RoleFeatureModifyViewModel) LoadDefaultState() {
	setConstants(this, "ROLEFEATUREMODIFY_CONST")
}

func (self *RoleFeatureModifyViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
