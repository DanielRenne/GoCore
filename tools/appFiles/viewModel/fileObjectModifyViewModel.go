package viewModel

import (
	"encoding/json"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
)

type FileObjectModifyViewModel struct {
	Constants struct {
	} `json:"constants"`
	FileObject model.FileObject `json:"FileObject"`
	FileUpload FileObject       `json:"FileUpload"`
	//AdditionalConstructs
}

func (this *FileObjectModifyViewModel) LoadDefaultState() {
	setConstants(this, "FILEOBJECTMODIFY_CONST")
}

func (self *FileObjectModifyViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
