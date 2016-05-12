package helloWorld

import (
	"fmt"
	_ "helloWorld/httpServices"
	"helloWorld/models/v1/model"
	_ "helloWorld/webAPIs/v1/webAPI"
	_ "helloWorld/webAPIs/v2/webAPI"
)

func init() {

	var person model.Person

	person.Worth = 23.5

	err := person.Save()

	if err != nil {
		fmt.Println(err.Error())
	}

	// var persons model.Persons

}
