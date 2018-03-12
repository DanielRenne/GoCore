package viewModel

import (
	"encoding/json"

	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
)

type AppErrorModifyViewModel struct {
	Constants struct {
	} `json:"constants"`
	AppError   model.AppError `json:"AppError"`
	FileUpload FileObject     `json:"FileUpload"`
	//AdditionalConstructs
}

func (this *AppErrorModifyViewModel) LoadDefaultState() {
	setConstants(this, "APPERRORMODIFY_CONST")
}

func (self *AppErrorModifyViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
