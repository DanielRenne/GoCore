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
	// Reinitialize if webConfig.json created
	serverSettings.Initialize(path, fileName)
	dbServices.RunDBCreate()
}

//GenerateModels will build model files based on your json schema
func GenerateModels(path string, fileName string) {

	serverSettings.Initialize(path, fileName)
	dbServices.Initialize()
	dbServices.RunDBCreate()
}
