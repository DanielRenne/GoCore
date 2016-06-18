package appGen

import (
	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/log"
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
		log.Message("Copied webConfig.json to Application.", log.GREEN)
	}

	_, err = os.Stat(serverSettings.APP_LOCATION + "/keys")

	if err != nil {
		extensions.CopyFolder(serverSettings.GOCORE_PATH+"/tools/appFiles/keys", serverSettings.APP_LOCATION+"/keys")
		log.Message("Copied keys to Application.", log.GREEN)
	}

	_, err = os.Stat(serverSettings.APP_LOCATION + "/web/core")
	if err != nil {
		extensions.CopyFolder(serverSettings.GOCORE_PATH+"/tools/appFiles/web/core", serverSettings.APP_LOCATION+"/web/core")
		log.Message("Copied web/core to Application.", log.GREEN)
	}

	_, err = os.Stat(serverSettings.APP_LOCATION + "/web/swagger")
	if err != nil {
		extensions.CopyFolder(serverSettings.GOCORE_PATH+"/tools/appFiles/web/swagger", serverSettings.APP_LOCATION+"/web/swagger")
		log.Message("Copied web/swagger to Application.", log.GREEN)
	}

	_, err = os.Stat(serverSettings.APP_LOCATION + "/db")

	if err != nil {
		os.MkdirAll(serverSettings.APP_LOCATION+"/db/schemas/1.0.0", 0777)
		log.Message("Created db/schemas/1.0.0 to Application.", log.GREEN)
	}
}
