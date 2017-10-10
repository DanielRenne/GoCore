package viewModel

import (
	"encoding/json"
	"time"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
)

type ServerSettingsModifyViewModel struct {
	// SettingsBar            SettingsButtonBarViewModel `json:"SettingsBar"`
	SelectedTab            string                  `json:"SelectedTab"`
	TimeZones   []model.Timezone    `json:"TimeZones"`
	LockoutSettings        LockoutSettingsModel    `json:"LockoutSettings"`
	TimeZone    model.ServerSetting `json:"TimeZone"`
	CurrentDate string              `json:"CurrentDate"`
	CurrentTime string              `json:"CurrentTime"`
	DateToSet   time.Time           `json:"DateToSet"`
	TimeToSet   time.Time           `json:"TimeToSet"`
}

type LockoutSettingsModel struct {
	Lockout model.ServerSetting `json:"Lockout"`
}

//type TestSettingsModel struct {
//	XXXX        model.ServerSetting `json:"XXXX"`
//	TTTT model.ServerSetting `json:"TTTT"`
//}

func (this *ServerSettingsModifyViewModel) LoadDefaultState() {
	this.SelectedTab = "settings"
	setConstants(this, "SERVERSETTINGSMODIFY_CONST")
}

func (self *ServerSettingsModifyViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
