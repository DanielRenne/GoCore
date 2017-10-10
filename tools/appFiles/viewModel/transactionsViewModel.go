package viewModel

import (
	"encoding/json"
)

type TransactionsViewModel struct {
	Constants struct {
	} `json:"constants"`
	//AdditionalConstructs
}

func (this *TransactionsViewModel) LoadDefaultState() {
	setConstants(this, "TRANSACTIONS_CONST")
}

func (self *TransactionsViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
