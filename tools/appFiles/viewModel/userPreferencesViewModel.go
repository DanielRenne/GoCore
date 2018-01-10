package viewModel

import (
	"encoding/json"

	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	_ "github.com/globalsign/mgo/bson"
)

type UserPreferences struct {
	UserPreferences []model.UsersPreference `json:"UserPreferences"`
}

func (self *UserPreferences) ParsePreferences(data string) {
	json.Unmarshal([]byte(data), &self)
}
