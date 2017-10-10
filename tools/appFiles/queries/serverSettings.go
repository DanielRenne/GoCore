package queries

import "github.com/DanielRenne/goCoreAppTemplate/models/v1/model"

const (
	SERVERSETTINGS_CATEGORY_USERS      = "users"
	SERVERSETTINGS_KEY_LOCKOUTATTEMPTS = "lockoutAttempts"
)

type queryServerSettings struct{}
func (self queryServerSettings) QueryUserSettings() (q *model.Query, err error) {
	q = model.ServerSettings.Query().In(model.Q(model.FIELD_SERVERSETTING_CATEGORY, SERVERSETTINGS_CATEGORY_USERS))
	return
}

func (self queryServerSettings) ById(id string) (setting model.ServerSetting, err error) {
	err = model.ServerSettings.Query().ById(id, &setting)
	return
}

func (self queryServerSettings) LoginAttempts() (setting model.ServerSetting, err error) {
	q, err := self.QueryUserSettings()
	q.Filter(model.Q(model.FIELD_SERVERSETTING_KEY, SERVERSETTINGS_KEY_LOCKOUTATTEMPTS))
	err = q.One(&setting)
	return
}
