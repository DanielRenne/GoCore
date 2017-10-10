package viewModel

import (
	"encoding/json"
)

type SettingsViewModel struct {
	SettingsBar SettingsButtonBarViewModel `json:"SettingsBar"`
}

func (this *SettingsViewModel) LoadDefaultState() {
	setConstants(this, "SETTINGS_CONST")
}

func (self *SettingsViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
