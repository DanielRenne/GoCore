package extensions

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Title(str string) string {
	return cases.Title(language.English, cases.NoLower).String(str)
}
