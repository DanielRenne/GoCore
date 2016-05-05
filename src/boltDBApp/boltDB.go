package boltDBApp

import (
	"boltDBApp/model"
	"fmt"
	"github.com/fatih/color"
)

func init() {
	// fmt.Printf("%+v\n", model)
	p := model.Person{Worth: 30.6665444, First: "Dave"}
	err := p.Save()
	if err != nil {
		color.Red("Error Saving Person:  " + err.Error())
	}
	color.Green("Saved Person Successfully")

	var persons model.Persons
	ps := persons.Range("", "Dave", "First")
	// if err != nil {
	// 	color.Red("Error Getting User:  " + err.Error())
	// }
	fmt.Printf("%+v\n", ps)

	// val, _ := ps.JSONString()
	// fmt.Println(val)

	// bytes, _ := ps.JSONBytes()

	// fmt.Printf("%+v\n", bytes)
}
