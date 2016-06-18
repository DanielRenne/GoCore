package appGen

import (

	// "fmt"
	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/log"
	"os"
)

// const APP_LOCATION = "src/bitbucket.org/DRenne/ohHell"
const APP_LOCATION = "src/github.com/thomsonreuters/aumentum-Web"

func GenerateApp() {
	moveAppFiles()

}

func moveAppFiles() {

	//First check for the WebConfig.json file
	var err error
	_, err = os.Stat(APP_LOCATION + "/webConfig.json")

	if err != nil {
		extensions.CopyFile(extensions.GOCORE_PATH+"/tools/appFiles/webConfig.json", APP_LOCATION+"/webConfig.json")
		log.Message("Copied webConfig.json to Application.", log.GREEN)
	}

	_, err = os.Stat(APP_LOCATION + "/keys")

	if err != nil {
		extensions.CopyFolder(extensions.GOCORE_PATH+"/tools/appFiles/keys", APP_LOCATION+"/keys")
		log.Message("Copied keys to Application.", log.GREEN)
	}

	extensions.RemoveDirectory(APP_LOCATION + "/web/core")

	// _, err = os.Stat(APP_LOCATION + "/web/core")
	// if err != nil {
	extensions.CopyFolder(extensions.GOCORE_PATH+"/tools/appFiles/web/core", APP_LOCATION+"/web/core")
	log.Message("Copied web/core to Application.", log.GREEN)
	// }

	_, err = os.Stat(APP_LOCATION + "/web/swagger")
	if err != nil {
		extensions.CopyFolder(extensions.GOCORE_PATH+"/tools/appFiles/web/swagger", APP_LOCATION+"/web/swagger")
		log.Message("Copied web/swagger to Application.", log.GREEN)
	}

	_, err = os.Stat(APP_LOCATION + "/db")

	if err != nil {
		os.MkdirAll(APP_LOCATION+"/db/schemas/1.0.0", 0777)
		log.Message("Created db/schemas/1.0.0 to Application.", log.GREEN)
	}
}
