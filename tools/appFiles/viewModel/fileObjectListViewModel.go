package viewModel

import (
	"encoding/json"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
)

type FileObjectListViewModel struct {
	Constants struct {
	} `json:"constants"`
	FileObjects []model.FileObject              `json:"FileObjects"`
	WidgetList  WidgetListUserControlsViewModel `json:"WidgetList"`
	FileUpload  FileObject                      `json:"FileUpload"`

	DeletedFileObjects []model.FileObject `json:"DeletedFileObjects"`

	//AdditionalConstructs
}

func (this *FileObjectListViewModel) LoadDefaultState() {
	setConstants(this, "FILEOBJECTLIST_CONST")
}

func (self *FileObjectListViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
