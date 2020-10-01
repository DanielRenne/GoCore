package commonStubs

var Locales string

func init() {

	Locales = `
package model

import "strings"

type Locale struct {
	Language string ` + "`" + `json:"Language"` + "`" + `
	Value    string ` + "`" + `json:"Value"` + "`" + `
}

var Locales = []Locale{
	{"Arabic", "ar"},
	{"Chinese", "zh"},
	{"Croatian", "hr"},
	{"Czech", "cs"},
	{"Danish", "da"},
	{"Dutch", "nl"},
	{"English", "en"},
	{"Finnish", "fi"},
	{"French", "fr"},
	{"German", "de"},
	{"Greek", "el"},
	{"Hebrew", "he"},
	{"Hindi", "hi"},
	{"Hungarian", "hu"},
	{"Indonesian", "id"},
	{"Italian", "it"},
	{"Japanese", "ja"},
	{"Korean", "ko"},
	{"Norwegian", "no"},
	{"Polish", "pl"},
	{"Portuguese", "pt"},
	{"Romanian", "ro"},
	{"Russian", "ru"},
	{"Turkish", "tr"},
	{"Thai", "th"},
	{"Spanish", "es"},
	{"Swedish", "sv"},
	{"Debug", "dev"},
}

func GetDefaultLocale(language string) string {
	if strings.Contains(language, "en") {
		return "en"
	}
	if strings.Contains(language, "es") {
		return "es"
	}
	if strings.Contains(language, "fr") {
		return "fr"
	}
	if strings.Contains(language, "ru") {
		return "ru"
	}
	if strings.Contains(language, "de") {
		return "de"
	}
	if strings.Contains(language, "it") {
		return "it"
	}
	if strings.Contains(language, "sv") {
		return "sv"
	}
	if strings.Contains(language, "ro") {
		return "ro"
	}
	if strings.Contains(language, "pt") {
		return "pt"
	}
	if strings.Contains(language, "hu") {
		return "hu"
	}
	if strings.Contains(language, "nl") {
		return "nl"
	}
	if strings.Contains(language, "ar") {
		return "ar"
	}
	if strings.Contains(language, "ko") {
		return "ko"
	}
	if strings.Contains(language, "ja") {
		return "ja"
	}
	if strings.Contains(language, "zh") {
		return "zh"
	}
	if strings.Contains(language, "he") {
		return "he"
	}
	if strings.Contains(language, "tr") {
		return "tr"
	}
	if strings.Contains(language, "th") {
		return "th"
	}
	if strings.Contains(language, "pl") {
		return "pl"
	}
	if strings.Contains(language, "hi") {
		return "hi"
	}
	if strings.Contains(language, "el") {
		return "el"
	}
	if strings.Contains(language, "hr") {
		return "hr"
	}
	if strings.Contains(language, "cs") {
		return "cs"
	}
	if strings.Contains(language, "id") {
		return "id"
	}
	if strings.Contains(language, "cs") {
		return "cs"
	}
	if strings.Contains(language, "da") {
		return "da"
	}
	if strings.Contains(language, "fi") {
		return "fi"
	}
	if strings.Contains(language, "no") {
		return "no"
	}
	return "en"
}
`
}
