package modelBuild 

import (
	"github.com/cloud-ignite/GoCore/core/dbServices"
	"github.com/cloud-ignite/GoCore/core/serverSettings"
)

func Initialize(path string, fileName string) {

	serverSettings.Initialize(path, fileName)
	dbServices.Initialize()
	dbServices.RunDBCreate()
}
