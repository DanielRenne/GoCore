package queries

import (
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
)

type queryPasswords struct{}

//Get an PasswordReset by Id
func (self *queryPasswords) ById(id string) (password model.Password, err error) {
	err = model.Passwords.Query().ById(id, &password)
	return
}
