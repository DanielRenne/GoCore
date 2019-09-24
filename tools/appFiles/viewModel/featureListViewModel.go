package viewModel

import (
	"encoding/json"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
)

type FeatureListViewModel struct {
	Constants struct {
	} `json:"constants"`
	Features   []model.Feature                 `json:"Features"`
	WidgetList WidgetListUserControlsViewModel `json:"WidgetList"`
	FileUpload FileObject                      `json:"FileUpload"`

	DeletedFeatures []model.Feature `json:"DeletedFeatures"`

	//AdditionalConstructs
}

func (this *FeatureListViewModel) LoadDefaultState() {
	setConstants(this, "FEATURELIST_CONST")
}

func (self *FeatureListViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
