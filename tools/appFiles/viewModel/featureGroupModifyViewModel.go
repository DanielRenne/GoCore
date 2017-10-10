package viewModel

import (
	"encoding/json"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
)

type FeatureGroupModifyViewModel struct {
	Constants struct {
	} `json:"constants"`
	FeatureGroup model.FeatureGroup `json:"FeatureGroup"`
	FileUpload   FileObject         `json:"FileUpload"`
	//AdditionalConstructs
}

func (this *FeatureGroupModifyViewModel) LoadDefaultState() {
	setConstants(this, "FEATUREGROUPMODIFY_CONST")
}

func (self *FeatureGroupModifyViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
