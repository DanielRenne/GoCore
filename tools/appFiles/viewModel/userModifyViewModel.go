package viewModel

import (
	"encoding/json"

	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
)

type UserModifyViewModel struct {
	SettingsBar           SettingsButtonBarViewModel `json:"SettingsBar"`
	User                  model.User                 `json:"User"`
	Accounts              []model.Account            `json:"Accounts"`
	Roles                 []model.Role               `json:"Roles"`
	AccountRole           model.AccountRole          `json:"AccountRole"`
	Password              string                     `json:"Password"`
	PasswordErrors        string                     `json:"PasswordErrors"`
	ConfirmPassword       string                     `json:"ConfirmPassword"`
	ConfirmPasswordErrors string                     `json:"ConfirmPasswordErrors"`
	EmailChanged          bool                       `json:"EmailChanged"`
	TimeZones             []model.Timezone           `json:"TimeZones"`
	Locales               []model.Locale             `json:"Locales"`
	UserLocale            string                     `json:"UserLocale"`
	CurrentPage           string                     `json:"CurrentPage"`
}

func (this *UserModifyViewModel) LoadDefaultState() {
	setConstants(this, "USERMODIFY_CONST")
	this.Accounts = make([]model.Account, 0)
	this.EmailChanged = false
}

func (self *UserModifyViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
