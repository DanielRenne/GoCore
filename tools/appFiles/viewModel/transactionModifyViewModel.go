package viewModel

import (
	"encoding/json"

	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
)

type TransactionModifyViewModel struct {
	Constants struct {
	} `json:"constants"`
	Transaction model.Transaction `json:"Transaction"`

	//AdditionalConstructs
}

func (this *TransactionModifyViewModel) LoadDefaultState() {
	setConstants(this, "TRANSACTIONMODIFY_CONST")
}

func (self *TransactionModifyViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
