package main

import (
	"github.com/DanielRenne/GoCore/core/appGen"
	"github.com/DanielRenne/GoCore/core/dbServices"
)

func main() {
	appGen.GenerateApp()
	dbServices.RunDBCreate()
}
