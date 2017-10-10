package viewModel

import (
	"encoding/json"
)

type UsersViewModel struct {
	Username           string `json:"Username"`
	Password           string `json:"Password"`
	PasswordValidation string `json:"PasswordValidation"`
}

func (this *UsersViewModel) LoadDefaultState() {

}

func (self *UsersViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
