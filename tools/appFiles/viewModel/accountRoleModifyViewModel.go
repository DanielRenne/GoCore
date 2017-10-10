package viewModel

import (
	"encoding/json"

	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
)

type AccountRoleModifyViewModel struct {
	Constants struct {
	} `json:"constants"`
	AccountRole model.AccountRole `json:"AccountRole"`
	FileUpload  FileObject        `json:"FileUpload"`
	//AdditionalConstructs
}

func (this *AccountRoleModifyViewModel) LoadDefaultState() {
	setConstants(this, "ACCOUNTROLEMODIFY_CONST")
}

func (self *AccountRoleModifyViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
