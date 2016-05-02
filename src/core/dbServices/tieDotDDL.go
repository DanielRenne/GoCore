package dbServices

import (
	"github.com/fatih/color"
)

func createTieDotCollections(collections []NOSQLCollection) {

	for _, collection := range collections {
		createTieDotCollection(collection.Name)
	}
}

func createTieDotCollection(name string) {

	if err := TiedotDB.Create(name); err != nil {
		color.Red("Failed to create tiedot Collection for " + name + ":  " + err.Error())
		return
	}
	color.Green("Created tiedot collection " + name)

}
