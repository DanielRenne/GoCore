package dbServices

import (
	"github.com/fatih/color"
)

func createTieDotCollections(collections []NOSQLCollection) {

	for _, collection := range collections {
		createTieDotCollection(collection)
	}
}

func createTieDotCollection(collection NOSQLCollection) {

	if err := TiedotDB.Create(collection.Name); err != nil {
		color.Red("Failed to create tiedot Collection for " + collection.Name + ":  " + err.Error())
	} else {
		color.Green("Created tiedot collection " + collection.Name)
	}

	//Create any indexes
	if len(collection.Indexes) > 0 {

		myCollection := TiedotDB.Use(collection.Name)

		if err := myCollection.Index(collection.Indexes); err != nil {
			color.Red("Failed to create tiedot Indexes for " + collection.Name + " collection:  " + err.Error())
		} else {
			color.Green("Created tiedot indexs for " + collection.Name)
		}

	}

}
