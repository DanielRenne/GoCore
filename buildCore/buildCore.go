package buildCore

import (
	"github.com/DanielRenne/GoCore/core/appGen"
	"github.com/DanielRenne/GoCore/core/dbServices"
	"github.com/DanielRenne/GoCore/core/serverSettings"
)

func Initialize(path string, fileName string) {

	serverSettings.Initialize(path, fileName)
	dbServices.Initialize()
	appGen.GenerateApp()
	dbServices.RunDBCreate()
}
