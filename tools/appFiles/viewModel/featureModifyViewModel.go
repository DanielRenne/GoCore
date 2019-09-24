package viewModel

import (
	"encoding/json"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
)

type FeatureModifyViewModel struct {
	Constants struct {
	} `json:"constants"`
	Feature       model.Feature        `json:"Feature"`
	FeatureGroups []model.FeatureGroup `json:"FeatureGroups"`
	FileUpload    FileObject           `json:"FileUpload"`
	//AdditionalConstructs
}

func (this *FeatureModifyViewModel) LoadDefaultState() {
	setConstants(this, "FEATUREMODIFY_CONST")
}

func (self *FeatureModifyViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
