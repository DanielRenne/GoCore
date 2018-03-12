package viewModel

import (
	"encoding/json"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
)

type FeatureGroupListViewModel struct {
	Constants struct {
	} `json:"constants"`
	FeatureGroups []model.FeatureGroup            `json:"FeatureGroups"`
	WidgetList    WidgetListUserControlsViewModel `json:"WidgetList"`
	FileUpload    FileObject                      `json:"FileUpload"`

	DeletedFeatureGroups []model.FeatureGroup `json:"DeletedFeatureGroups"`

	//AdditionalConstructs
}

func (this *FeatureGroupListViewModel) LoadDefaultState() {
	setConstants(this, "FEATUREGROUPLIST_CONST")
}

func (self *FeatureGroupListViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
