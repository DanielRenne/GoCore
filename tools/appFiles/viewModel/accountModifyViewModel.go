package viewModel

import (
	"encoding/json"

	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
)

type AccountModifyViewModel struct {
	SettingsBar         SettingsButtonBarViewModel `json:"SettingsBar"`
	Countries           []model.Country            `json:"Countries"`
	Account             model.Account              `json:"Account"`
	AccountRole         model.AccountRole          `json:"AccountRole"`
	ImageFileName       string                     `json:"ImageFileName"`
	ImageFileNameErrors string                     `json:"ImageFileNameErrors"`
	States              map[string][]model.State   `json:"States"`
}

func (this *AccountModifyViewModel) LoadDefaultState() {
	setConstants(this, "ACCOUNTMODIFY_CONST")
}

func (self *AccountModifyViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
