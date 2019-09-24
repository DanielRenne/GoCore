package viewModel

import (
	"encoding/json"
)

type TransactionAddViewModel struct {
	Constants struct {
	} `json:"constants"`
	//AdditionalConstructs
}

func (this *TransactionAddViewModel) LoadDefaultState() {
	setConstants(this, "TRANSACTIONADD_CONST")
}

func (self *TransactionAddViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
