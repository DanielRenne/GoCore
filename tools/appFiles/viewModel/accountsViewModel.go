package viewModel

import (
	"encoding/json"
)

type AccountsViewModel struct {
	ButtonBar ButtonBar `json:"ButtonBar"`
	Constants struct {
	} `json:"constants"`
}

func (this *AccountsViewModel) LoadDefaultState() {
	setConstants(this, "ACCOUNTS_CONST")
}

func (self *AccountsViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
