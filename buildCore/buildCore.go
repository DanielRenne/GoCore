package buildCore

import (
	"github.com/cloud-ignite/GoCore/core/appGen"
	"github.com/cloud-ignite/GoCore/core/dbServices"
	"github.com/cloud-ignite/GoCore/core/serverSettings"
)

func Initialize(path string, fileName string) {

	serverSettings.Initialize(path, fileName)
	dbServices.Initialize()
	appGen.GenerateApp()
	// Reinitialize if webConfig.json created
	serverSettings.Initialize(path, fileName)
	dbServices.RunDBCreate()
}
