package tieDotTest

import (
	"fmt"
	"github.com/fatih/color"
	"tieDotTest/model"
)

func init() {
	// fmt.Printf("%+v\n", model)
	var collection model.TestCollections
	tc := model.TestCollection{Field1: 123}

	docID, err := collection.Insert(tc)
	if err != nil {
		color.Red("Error inserting new TestCollection:  " + err.Error())
		return
	}

	fmt.Println(string(docID))

	readBack, err := collection.ReadById(docID)
	if err != nil {
		panic(err)
	}
	fmt.Println("Document", docID)
	fmt.Printf("%+v\n", readBack)
}
