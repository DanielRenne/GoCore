package viewModel

import (
	"encoding/json"

	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	_ "gopkg.in/mgo.v2/bson"
)

type UserPreferences struct {
	UserPreferences []model.UsersPreference `json:"UserPreferences"`
}

func (self *UserPreferences) ParsePreferences(data string) {
	json.Unmarshal([]byte(data), &self)
}
