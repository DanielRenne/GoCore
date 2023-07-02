package extensions

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Title(str string) string {
	caser := cases.Title(language.English)
	return caser.String(str)
}
