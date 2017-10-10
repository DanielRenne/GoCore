package queries

import (
	"encoding/json"
	"strings"

	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/settings"
	"github.com/davidrenne/reflections"
)

type AppContentJson struct {
	//AdditionalConstructs
}

type queryAppContent struct{}

type TagReplacements struct {
	Tag1  map[string]string
	Tag2  map[string]string
	Tag3  map[string]string
	Tag4  map[string]string
	Tag5  map[string]string
	Tag6  map[string]string
	Tag7  map[string]string
	Tag8  map[string]string
	Tag9  map[string]string
	Tag10 map[string]string
}

func Q(k string, v string) map[string]string {
	var ret map[string]string
	ret = make(map[string]string)
	ret[k] = v
	return ret
}

func (self *queryAppContent) GetTranslationWithReplacementsFromUser(user model.User, key string, tags *TagReplacements) (originalReplacement string) {
	originalReplacement = self.GetTranslationFromUser(user, key)
	if originalReplacement != "" && tags != nil {
		for k, v := range tags.Tag1 {
			originalReplacement = strings.Replace(originalReplacement, "{"+k+"}", v, -1)
		}
		for k, v := range tags.Tag2 {
			originalReplacement = strings.Replace(originalReplacement, "{"+k+"}", v, -1)
		}
		for k, v := range tags.Tag3 {
			originalReplacement = strings.Replace(originalReplacement, "{"+k+"}", v, -1)
		}
		for k, v := range tags.Tag4 {
			originalReplacement = strings.Replace(originalReplacement, "{"+k+"}", v, -1)
		}
		for k, v := range tags.Tag5 {
			originalReplacement = strings.Replace(originalReplacement, "{"+k+"}", v, -1)
		}
		for k, v := range tags.Tag6 {
			originalReplacement = strings.Replace(originalReplacement, "{"+k+"}", v, -1)
		}
		for k, v := range tags.Tag7 {
			originalReplacement = strings.Replace(originalReplacement, "{"+k+"}", v, -1)
		}
		for k, v := range tags.Tag8 {
			originalReplacement = strings.Replace(originalReplacement, "{"+k+"}", v, -1)
		}
		for k, v := range tags.Tag9 {
			originalReplacement = strings.Replace(originalReplacement, "{"+k+"}", v, -1)
		}
		for k, v := range tags.Tag10 {
			originalReplacement = strings.Replace(originalReplacement, "{"+k+"}", v, -1)
		}
	}
	return originalReplacement
}

func (self *queryAppContent) GetTranslationWithReplacements(context session_functions.RequestContext, key string, tags *TagReplacements) (originalReplacement string) {
	user, err := getUser(context)
	if err != nil {
		return
	}
	return self.GetTranslationWithReplacementsFromUser(user, key, tags)
}

func (self *queryAppContent) GetTranslationFromUser(user model.User, key string) (translatedText string) {
	contentData, err := extensions.ReadFile(settings.WebRoot + "/globalization/translations/app/" + user.Language + "/" + user.Language + ".json")
	if err != nil {
		contentData, err = extensions.ReadFile(settings.WebRoot + "/globalization/translations/app/en/US.json")
	}
	var appContent AppContentJson
	err = json.Unmarshal(contentData, &appContent)
	if err != nil {
		return ""
	}
	translatedText, _ = reflections.GetFieldAsString(appContent, key)
	return translatedText
}

func (self *queryAppContent) GetTranslation(context session_functions.RequestContext, key string) (translatedText string) {

	user, err := getUser(context)
	if err != nil { //Just get the English Language for no context
		var usr model.User
		usr.Language = "en"
		translatedText = self.GetTranslationFromUser(usr, key)
		return
	}
	return self.GetTranslationFromUser(user, key)
}
