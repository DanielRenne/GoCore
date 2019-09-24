package viewModel

import (
	"encoding/json"

	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
)

type PasswordModifyViewModel struct {
	Constants struct {
	} `json:"constants"`
	Password   model.Password `json:"Password"`
	FileUpload FileObject     `json:"FileUpload"`
	//AdditionalConstructs
}

func (this *PasswordModifyViewModel) LoadDefaultState() {
	setConstants(this, "PASSWORDMODIFY_CONST")
}

func (self *PasswordModifyViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
