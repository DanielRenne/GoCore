package viewModel

import (
	"encoding/json"

	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
)

type AppErrorListViewModel struct {
	Constants struct {
	} `json:"constants"`
	AppErrors  []model.AppError                `json:"AppErrors"`
	WidgetList WidgetListUserControlsViewModel `json:"WidgetList"`

	DeletedAppErrors []model.AppError `json:"DeletedAppErrors"`

	//AdditionalConstructs
}

func (this *AppErrorListViewModel) LoadDefaultState() {
	setConstants(this, "APPERRORLIST_CONST")
}

func (self *AppErrorListViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
