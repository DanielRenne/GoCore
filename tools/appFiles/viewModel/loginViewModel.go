package viewModel

import (
	"encoding/json"
)

type LoginViewModel struct {
	AuthMessage   string `json:"authMessage" js:"authMessage"`
	Username      string `json:"username" js:"username"`
	Password      string `json:"password" js:"password"`
	UserNameError string `json:"UserNameError"`
	PasswordError string `json:"PasswordError"`
}

func (self *LoginViewModel) LoadDefaultState() {
}

func (self *LoginViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
