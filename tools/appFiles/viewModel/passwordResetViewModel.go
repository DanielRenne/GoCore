package viewModel

import (
	"encoding/json"

	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
)

type PasswordResetViewModel struct {
	Constants struct {
	} `json:"constants"`
	PasswordReset         model.PasswordReset `json:"PasswordReset"`
	Password              string              `json:"Password"`
	PasswordErrors        string              `json:"PasswordErrors"`
	ConfirmPassword       string              `json:"ConfirmPassword"`
	ConfirmPasswordErrors string              `json:"ConfirmPasswordErrors"`
}

func (this *PasswordResetViewModel) LoadDefaultState() {
	setConstants(this, "PASSWORDRESET_CONST")
}

func (self *PasswordResetViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
