package queries

import (
	"encoding/json"
	"strings"

	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/ginServer"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/settings"
)

type PasswordResetTranslation struct {
	PasswordResetTitle   string `json:"PasswordResetTitle"`
	PasswordResetMessage string `json:"PasswordResetMessage"`
}

type queryPasswordResets struct{}

//Read the Invitation Language file and parse the json into an InvitaitonTranslation Struct
func (self *queryPasswordResets) GetPasswordResetTranslation(localeLanguage ginServer.LocaleLanguage, email string, link string) (title string, body string, err error) {

	contentData, err := extensions.ReadFile(settings.WebRoot + "/globalization/translations/passwordReset/" + localeLanguage.Language + "/" + localeLanguage.Locale + ".json")

	if err != nil {
		contentData, err = extensions.ReadFile(settings.WebRoot + "/globalization/translations/passwordReset/en/US.json")
	}

	var prt PasswordResetTranslation
	err = json.Unmarshal(contentData, &prt)

	if err != nil {
		return
	}

	title = prt.PasswordResetTitle

	body = strings.Replace(prt.PasswordResetMessage, "{Email}", email, -1)
	body = strings.Replace(body, "{Link}", "<a href=\""+link+"\">here</a>", -1)

	return
}

//Get an PasswordReset by Id
func (self *queryPasswordResets) ById(id string) (passwordReset model.PasswordReset, err error) {
	err = model.PasswordResets.Query().ById(id, &passwordReset)
	return
}

func (self *queryPasswordResets) QueryByUserId(id string) *model.Query {
	return model.PasswordResets.Query().Filter(model.Q("UserId", id)).Filter(model.Q("Complete", false))
}

//Get an PasswordReset by UserId
func (self *queryPasswordResets) ByUserId(id string) (passwordReset model.PasswordReset, err error) {
	q := self.QueryByUserId(id)
	err = q.One(&passwordReset)
	return
}
