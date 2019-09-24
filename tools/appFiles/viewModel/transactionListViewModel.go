package viewModel

import (
	"encoding/json"

	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
)

type TransactionListViewModel struct {
	Constants struct {
	} `json:"constants"`
	Transactions []model.Transaction             `json:"Transactions"`
	WidgetList   WidgetListUserControlsViewModel `json:"WidgetList"`

	//AdditionalConstructs
}

func (this *TransactionListViewModel) LoadDefaultState() {
	setConstants(this, "TRANSACTIONLIST_CONST")
}

func (self *TransactionListViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
