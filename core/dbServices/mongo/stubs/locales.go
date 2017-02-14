package model

import "strings"

type Locale struct {
	Language string `json:"Language"`
	Value    string `json:"Value"`
}

var Locales = []Locale{
	{"English", "en"},
	{"Spanish", "es"},
	{"French", "fr"},
	{"German", "de"},
	{"Italian", "it"},
	{"Swedish", "sv"},
	{"Portuguese", "pt"},
	{"Hungarian", "hu"},
	{"Dutch", "nl"},
}

//
//{"Romanian", "ro"},
//{"Arabic", "ar"},
//{"Korean", "ko"},
//{"Japanese", "ja"},
//{"Chinese", "zh"},

// we have this translated but have some latin1 conversion work to do before it would be ready
//	{"Russian", "ru"},

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
	return "en"
}
