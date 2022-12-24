package main

import (
	// Example module package name is github.com/DanielRenne/GoCore/core/dbServices/example
	// When a webConfig.json exists like the file in this directory
	// models/v1/model will generate structs and methods for your ORM schemas
	// you should never checkin this folder and .gitignore it and require developers to run the buildCore command you have in a binary somewhere to generate the models
	"fmt"
	"time"

	"github.com/DanielRenne/GoCore/core/dbServices/example/models/v1/model"

	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/GoCore/core/dbServices"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/common-nighthawk/go-figure"
)

func addLineEndings(heading string) {
	fmt.Println("------------------------------------------------------------------------------------------")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	myFigure := figure.NewFigure(heading, "", true)
	myFigure.Print()

	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("------------------------------------------------------------------------------------------")

	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
}

func exampleInsertions() {
	var allBootstrappedSites []model.Site
	err := model.Sites.Query().Join("Buildings").Join("LastUpdateUser").Join("Country").Join("Account").Join("FileObjects").All(&allBootstrappedSites)
	if err != nil {
		core.Dump("Error: " + err.Error())
		return
	}
	var allBootstrappedBuildings []model.Building
	err = model.Sites.Query().Join("Site").Join("LastUpdateUser").Join("Floors").Join("Account").Join("FileObjects").All(&allBootstrappedBuildings)
	if err != nil {
		core.Dump("Error: " + err.Error())
		return
	}
	core.Dump("exampleInsertions output", "This is the data bootstrapped for sites schema (with joins to related colletions)", allBootstrappedSites, "This is the data bootstrapped for sites schema (with joins to related colletions)", allBootstrappedBuildings, "Yay bootstrapped data exists in your new database even though you didnt write any code to insert it!")
	var country model.Country
	err = model.Countries.Query().Where("Iso", "us").One(&country)
	if err != nil {
		core.Dump("Error: " + err.Error())
		return
	}
	user := model.User{}
	user.Email = "testing12356789@test.com"
	user.First = "Go"
	user.Last = "Core"
	err = user.Save()
	if err != nil {
		core.Dump("Error: " + err.Error())
		return
	}
	account := model.Account{}
	account.AccountName = "GoCore App"
	account.SecondaryPhone.CountryISO = "us"
	account.SecondaryPhone.DialCode = "1"
	account.SecondaryPhone.Numeric = "2483339223"
	account.SecondaryPhone.Value = "1 2483339223"
	err = account.Save()
	if err != nil {
		core.Dump("Error: " + err.Error())
		return
	}
	site := model.Site{}
	site.LastUpdateId = user.Id.Hex()
	site.AccountId = account.Id.Hex()
	site.Name = "GoCore Location"
	site.CountryId = country.Id.Hex()
	site.ImageCustom = "585f966a1d41c87c55ce450f"
	err = site.Save()
	if err != nil {
		core.Dump("Error: " + err.Error())
		return
	}
	building1 := model.Building{}
	building1.SiteId = site.Id.Hex()
	building1.Name = "WebSocket Building"
	building1.AccountId = account.Id.Hex()
	building1.LastUpdateId = user.Id.Hex()
	building1.ImageCustom = "585f966a1d41c87c55ce450f"
	err = building1.Save()
	if err != nil {
		core.Dump("Error: " + err.Error())
		return
	}
	building2 := model.Building{}
	building2.SiteId = site.Id.Hex()
	building2.AccountId = account.Id.Hex()
	building2.Name = "Channels and GoRoutine Building"
	building2.LastUpdateId = user.Id.Hex()
	building2.ImageCustom = "585f966a1d41c87c55ce450f"
	err = building2.Save()
	if err != nil {
		core.Dump("Error: " + err.Error())
		return
	}
	floor := model.Floor{}
	floor.SiteId = site.Id.Hex()
	floor.BuildingId = building1.Id.Hex()
	floor.LastUpdateId = user.Id.Hex()
	floor.AccountId = account.Id.Hex()
	floor.Name = "Floor 1"
	err = floor.Save()

	if err != nil {
		core.Dump("Error: " + err.Error())
		return
	}
	floor = model.Floor{}
	floor.SiteId = site.Id.Hex()
	floor.BuildingId = building1.Id.Hex()
	floor.LastUpdateId = user.Id.Hex()
	floor.AccountId = account.Id.Hex()
	floor.Name = "Floor 2"
	err = floor.Save()
	if err != nil {
		core.Dump("Error: " + err.Error())
		return
	}
	floor = model.Floor{}
	floor.SiteId = site.Id.Hex()
	floor.BuildingId = building2.Id.Hex()
	floor.LastUpdateId = user.Id.Hex()
	floor.AccountId = account.Id.Hex()
	floor.Name = "Floor 1"
	err = floor.Save()
	if err != nil {
		core.Dump("Error: " + err.Error())
		return
	}
	floor = model.Floor{}
	floor.SiteId = site.Id.Hex()
	floor.BuildingId = building2.Id.Hex()
	floor.LastUpdateId = user.Id.Hex()
	floor.AccountId = account.Id.Hex()
	floor.Name = "Floor 2"
	err = floor.Save()
	if err != nil {
		core.Dump("Error: " + err.Error())
		return
	}
}

func exampleViews() {
	var site model.Site
	err := model.Sites.Query().RenderViews(model.DataFormat{Language: "en", LocalTimeZone: "US/Eastern", DateFormat: "mm/dd/yyyy"}).ById("633e21412a1b49f431ee6f4d", &site)
	if err != nil {
		core.Dump("Error: " + err.Error())
		return
	}
	core.Dump("This is a calculated builtin view off a field based on the time elapsed since the last update: " + site.Views.UpdateFromNow)
}

func exampleWhereWithJoins() {

	var allSites []model.Site
	err := model.Sites.Query().Join("Buildings").Join("LastUpdateUser").Join("Country").Join("Account").Join("FileObjects").All(&allSites)
	if err != nil {
		core.Dump("Error: " + err.Error())
		return
	}
	var allBuildings []model.Building
	err = model.Buildings.Query().Join("Site").Join("LastUpdateUser").Join("Floors").Join("Account").Join("FileObjects").All(&allBuildings)
	if err != nil {
		core.Dump("Error: " + err.Error())
		return
	}
	core.Dump("exampleWhereWithJoins output", allSites, allBuildings)
}

func exampleGetFieldsOfCollectionViaReflection() {
	core.Dump("exampleGetFieldsOfCollectionViaReflection output", model.Sites.New().Reflect())
}

func exampleGreaterThan() {
	var fileObjects []model.FileObject
	err := model.FileObjects.Query().GreaterThanEqualTo(model.MaxQ("CreateDate", time.Date(2022, 9, 20, 0, 0, 0, 0, time.Now().UTC().Location()))).All(&fileObjects)
	if err != nil {
		core.Dump("Error: " + err.Error())
		return
	}
	core.Dump("exampleGreaterThan output only pulls one record that was bootstrapped", fileObjects)
}

func exampleDistinct() {
	var buildingsImageCustom []string
	err := model.Buildings.Query().Distinct("ImageCustom", &buildingsImageCustom)
	if err != nil {
		core.Dump("Error: " + err.Error())
		return
	}
	core.Dump("exampleDistinct output should have a length of 3 even though 4 rows exist with data were getting the distinct list of ImageCustom", buildingsImageCustom)
}

func main() {
	serverSettings.Init()
	dbServices.Initialize()
	model.ConnectDB()
	addLineEndings("Insertions")
	exampleInsertions()
	addLineEndings("Views")
	exampleViews()
	addLineEndings("Where / Joins")
	exampleWhereWithJoins()
	addLineEndings("Reflection")
	exampleGetFieldsOfCollectionViaReflection()
	addLineEndings("Greater Than")
	exampleGreaterThan()
	addLineEndings("Distinct")
	exampleDistinct()
	addLineEndings("Thanks for playing!")
	addLineEndings("With GoCore!")
}
