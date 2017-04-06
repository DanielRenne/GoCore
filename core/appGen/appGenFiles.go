package appGen

import (
	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/logger"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"os"
)

func GenerateApp() {
	moveAppFiles()
}

func moveAppFiles() {

	//First check for the WebConfig.json file
	var err error
	_, err = os.Stat(serverSettings.APP_LOCATION + "/webConfig.json")

	if err != nil {
		extensions.CopyFile(serverSettings.GOCORE_PATH+"/tools/appFiles/webConfig.json", serverSettings.APP_LOCATION+"/webConfig.json")
		logger.Message("Copied webConfig.json to Application.", logger.GREEN)
	}

	_, err = os.Stat(serverSettings.APP_LOCATION + "/keys")

	if err != nil {
		extensions.CopyFolder(serverSettings.GOCORE_PATH+"/tools/appFiles/keys", serverSettings.APP_LOCATION+"/keys")
		logger.Message("Copied keys to Application.", logger.GREEN)
	}

	_, err = os.Stat(serverSettings.APP_LOCATION + "/web/core")
	if err != nil {
		extensions.CopyFolder(serverSettings.GOCORE_PATH+"/tools/appFiles/web/core", serverSettings.APP_LOCATION+"/web/core")
		logger.Message("Copied web/core to Application.", logger.GREEN)
	}

	_, err = os.Stat(serverSettings.APP_LOCATION + "/web/swagger")
	if err != nil {
		extensions.CopyFolder(serverSettings.GOCORE_PATH+"/tools/appFiles/web/swagger", serverSettings.APP_LOCATION+"/web/swagger")
		logger.Message("Copied web/swagger to Application.", logger.GREEN)
	}

	_, err = os.Stat(serverSettings.APP_LOCATION + "/db")

	if err != nil {
		os.MkdirAll(serverSettings.APP_LOCATION+"/db/schemas/1.0.0", 0777)
		logger.Message("Created db/schemas/1.0.0 to Application.", logger.GREEN)
	}
}
