package appGen

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/logger"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/DanielRenne/GoCore/core/utils"
	"github.com/globalsign/mgo/bson"
)

func GenerateApp() {
	moveAppFiles()
}

// create if not exists
func createFile(path string, contents string) {
	_, err := os.Stat(serverSettings.APP_LOCATION + path)
	if err != nil {
		extensions.WriteToFile(contents, serverSettings.APP_LOCATION+path, 0644)
	}
}

func copyFolder(path string) (wasCopied bool) {
	_, err := os.Stat(serverSettings.APP_LOCATION + path)

	if err != nil {
		wasCopied = true
		os.MkdirAll(serverSettings.APP_LOCATION+path, 0777)
		extensions.CopyFolder(serverSettings.GOCORE_PATH+"/tools/appFiles"+path, serverSettings.APP_LOCATION+path)
		logger.Message("Created "+path+" in Application.", logger.GREEN)
	}
	return
}

func replacePath(path string, newpath string, newGithubUser string, newProject string) {
	filepath.Walk(serverSettings.APP_LOCATION+path, func(pth string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			utils.ReplaceTokenInFile(pth, "DanielRenne/goCoreAppTemplate", newpath)
			// we cant just globally replace DanielRenne with the new github username, so we use the special token DanielRenneFolder
			utils.ReplaceTokenInFile(pth, "DanielRenneFolder", newGithubUser)
			//Finally any straggler templates such as XXXX.go for main need to be replaced
			utils.ReplaceTokenInFile(pth, "goCoreAppTemplate", newProject)
			utils.ReplaceTokenInFile(pth, "goCoreUpperAppTemplate", strings.ToUpper(newProject[:1])+newProject[1:])
		}
		return err
	})
}

func replaceAnything(path string, find string, replace string) {
	filepath.Walk(serverSettings.APP_LOCATION+path, func(pth string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			utils.ReplaceTokenInFile(pth, find, replace)
		}
		return err
	})
}

func moveAppFiles() {
	humanTitle, err := extensions.ReadFile("/tmp/humanTitle")
	if err != nil {
		log.Println("error reading humanTitle")
	}
	parts := strings.Split(serverSettings.APP_LOCATION, "/")
	appName := parts[len(parts)-1]
	githubName := parts[len(parts)-2]
	project := githubName + "/" + appName
	//First check for the WebConfig.json file
	_, errNoWebConfig := os.Stat(serverSettings.APP_LOCATION + "/webConfig.json")
	if errNoWebConfig != nil {
		extensions.CopyFile(serverSettings.GOCORE_PATH+"/tools/appFiles/webConfig.json", serverSettings.APP_LOCATION+"/webConfig.json")
		logger.Message("Copied webConfig.json to Application.", logger.GREEN)
	}

	_, err = os.Stat(serverSettings.APP_LOCATION + "/webConfig.prod.json")
	if err != nil {
		extensions.CopyFile(serverSettings.GOCORE_PATH+"/tools/appFiles/webConfig.prod.json", serverSettings.APP_LOCATION+"/webConfig.prod.json")
		logger.Message("Copied webConfig.json to Application.", logger.GREEN)
	}

	_, err = os.Stat(serverSettings.APP_LOCATION + "/webConfig.dev.json")
	if err != nil {
		extensions.CopyFile(serverSettings.GOCORE_PATH+"/tools/appFiles/webConfig.dev.json", serverSettings.APP_LOCATION+"/webConfig.dev.json")
		logger.Message("Copied webConfig.json to Application.", logger.GREEN)
	}

	for _, v := range utils.Array("webConfig.prod.json", "webConfig.dev.json", "webConfig.json") {
		id1 := bson.NewObjectId()
		id2 := bson.NewObjectId()
		utils.ReplaceTokenInFile(serverSettings.APP_LOCATION+"/"+v, "goCoreProductName", appName+"BaseProduct")
		utils.ReplaceTokenInFile(serverSettings.APP_LOCATION+"/"+v, "goCoreCsrfSecret", id1.Hex())
		utils.ReplaceTokenInFile(serverSettings.APP_LOCATION+"/"+v, "goCoreSessionKey", id2.Hex())
	}

	_, err = os.Stat(serverSettings.APP_LOCATION + "/log")
	if err != nil {
		os.MkdirAll(serverSettings.APP_LOCATION+"/log/plugins", 0777)
	}
	var wasCopied bool
	wasCopied = copyFolder("/vendorPackages")
	wasCopied = copyFolder("/keys")
	wasCopied = copyFolder("/web")
	if wasCopied {
		utils.ReplaceTokenInFile(serverSettings.APP_LOCATION+"/web/app/watchFile.json", "DanielRenne/goCoreAppTemplate", project)
		utils.ReplaceTokenInFile(serverSettings.APP_LOCATION+"/web/app/javascript/build-css.sh", "DanielRenne/goCoreAppTemplate", project)
		replacePath("/web/app/javascript/pages/template", project, githubName, appName)
		for _, v := range utils.Array("/web/app/manifests", "/web/app/globalization/translations", "/web/app/javascript/pages/logs", "/web/app/javascript/globals", "/web/app/markup/app") {
			replaceAnything(v, "GoCoreAppHumanName", strings.TrimSpace(string(humanTitle)))
		}
	}
	wasCopied = copyFolder("/payloads")
	wasCopied = copyFolder("/constants")
	if wasCopied {
		replacePath("/constants", project, githubName, appName)
	}
	wasCopied = copyFolder("/controllers")
	if wasCopied {
		replacePath("/controllers", project, githubName, appName)
		utils.ReplaceTokenInFile(serverSettings.APP_LOCATION+"/controllers/homeGetController.go", "goCoreProductName", appName)
		utils.ReplaceTokenInFile(serverSettings.APP_LOCATION+"/controllers/homeGetController.go", "-APPNAME", "-"+appName)
	}
	wasCopied = copyFolder("/bin")
	if wasCopied {
		replacePath("/bin", project, githubName, appName)
	}
	wasCopied = copyFolder("/cron")
	if wasCopied {
		replacePath("/cron", project, githubName, appName)
	}
	wasCopied = copyFolder("/constants")
	if wasCopied {
		replacePath("/constants", project, githubName, appName)
	}

	// no web config only do this once when theres no webConfig.  In theory we can go back to the old buildApp if we ever merge to DanielRenne
	if errNoWebConfig != nil {
		copyFolder("/install")
		replacePath("/install", project, githubName, appName)
		err = os.Rename(serverSettings.APP_LOCATION+"/install/install.go", serverSettings.APP_LOCATION+"/install/install"+strings.Title(appName)+".go")
		if err != nil {
			log.Println("error renaming install")
			os.Exit(1)
		}
		err = os.Rename(serverSettings.APP_LOCATION+"/install", serverSettings.APP_LOCATION+"/install"+strings.Title(appName))
		if err != nil {
			log.Println("error renaming install folder")
			os.Exit(1)
		}
	}
	wasCopied = copyFolder("/br")
	if wasCopied {
		replacePath("/br", project, githubName, appName)
	}
	wasCopied = copyFolder("/scheduleEngine")
	if wasCopied {
		replacePath("/scheduleEngine", project, githubName, appName)
	}
	wasCopied = copyFolder("/password")
	if wasCopied {
		replacePath("/password", project, githubName, appName)
		secret := bson.NewObjectId()
		replaceAnything("/password", "GoCorePasswordSecret", secret.Hex())
	}
	wasCopied = copyFolder("/queries")
	if wasCopied {
		replacePath("/queries", project, githubName, appName)
	}
	wasCopied = copyFolder("/settings")
	if wasCopied {
		replacePath("/settings", project, githubName, appName)
	}
	wasCopied = copyFolder("/sessionFunctions")
	if wasCopied {
		replacePath("/sessionFunctions", project, githubName, appName)
	}
	wasCopied = copyFolder("/viewModel")
	if wasCopied {
		replacePath("/viewModel", project, githubName, appName)
	}
	wasCopied = copyFolder("/errors")
	if wasCopied {
		replacePath("/errors", project, githubName, appName)
	}
	wasCopied = copyFolder("/networks")
	if wasCopied {
		replacePath("/networks", project, githubName, appName)
	}

	_, err = os.Stat(serverSettings.APP_LOCATION + "/db")

	if err != nil {
		os.MkdirAll(serverSettings.APP_LOCATION+"/db/schemas/1.0.0", 0777)
		os.MkdirAll(serverSettings.APP_LOCATION+"/db/bootstrap", 0777)
		os.MkdirAll(serverSettings.APP_LOCATION+"/db/goFiles/v1", 0777)
		extensions.CopyFolder(serverSettings.GOCORE_PATH+"/tools/appFiles/db/schemas", serverSettings.APP_LOCATION+"/db/schemas/1.0.0")
		extensions.CopyFolder(serverSettings.GOCORE_PATH+"/tools/appFiles/db/bootstrap", serverSettings.APP_LOCATION+"/db/bootstrap")
		logger.Message("Created db/schemas/1.0.0 in Application.", logger.GREEN)
	}

	_, err = os.Stat(serverSettings.APP_LOCATION + "/models")

	if err != nil {
		os.MkdirAll(serverSettings.APP_LOCATION+"/models/v1/model", 0777)
		logger.Message("Created models/v1/model in Application.", logger.GREEN)
	}

	createFile("/releaseNotes.txt", `
`+strings.ToUpper(appName)+` Release Notes:

Legend:
				[+] new feature
				[-] removed function
				[*] bug fixed and improvement made

-`+appName+` 0.0.2 Firmware
				[*] Initial Changes By Developer

-`+appName+` 0.0.1 Firmware
				[*] GoCore Application Generated and Committed Base App

				-APPNAME`)

	createFile("/"+appName+".go", `
package main

import (
	"log"
	"os"

	"runtime/debug"
	"runtime/trace"
	"time"

	"fmt"
	"net/http"

	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/GoCore/core/app"
	"github.com/DanielRenne/GoCore/core/ginServer"
	"github.com/DanielRenne/GoCore/core/logger"
	_ "`+strings.Replace(serverSettings.APP_LOCATION, "src/", "", -1)+`/controllerRegistry"
	"`+strings.Replace(serverSettings.APP_LOCATION, "src/", "", -1)+`/br"
	"`+strings.Replace(serverSettings.APP_LOCATION, "src/", "", -1)+`/controllers"
	"`+strings.Replace(serverSettings.APP_LOCATION, "src/", "", -1)+`/cron"
	_ "`+strings.Replace(serverSettings.APP_LOCATION, "src/", "", -1)+`/models/v1/model"
	"`+strings.Replace(serverSettings.APP_LOCATION, "src/", "", -1)+`/sessionFunctions"
	"`+strings.Replace(serverSettings.APP_LOCATION, "src/", "", -1)+`/settings"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			session_functions.Print("\n\nPanic Stack: " + string(debug.Stack()))
			session_functions.Log("studio.go", "Panic Recovered at main():"+fmt.Sprintf("%+v", r))
			time.Sleep(time.Millisecond * 3000)
			main()
			return
		}
	}()

	err := app.Initialize("`+serverSettings.APP_LOCATION+`", "webConfig.json")
	settings.Initialize()
	br.Schedules.UpdateLinuxToGMT()

	if err != nil {
		//lastError := err.Error()
		ginServer.Router.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "%v", "An error occurred and the `+appName+` app cannot run (most likely due to mongo database services being down).\n\nError description: "+err.Error())
		})
		app.Run()
	} else {
		if settings.AppSettings.DeveloperGoTrace {
			f, err := os.Create("`+serverSettings.APP_LOCATION+`/log/trace.log")
			if err != nil {
				panic(err)
			}
			defer f.Close()

			err = trace.Start(f)
			if err != nil {
				panic(err)
			}
			mgo.SetDebug(true)

			file, _ := os.OpenFile("`+serverSettings.APP_LOCATION+`/log/studioMongo.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)

			var aLogger *log.Logger
			aLogger = log.New(file, "", log.LstdFlags)

			mgo.SetLogger(aLogger)
			mgo.SetStats(true)
		}

		controllers.Initialize()

		core.CronJobs.Start()
		cron.Start()

		go logger.GoRoutineLogger(func() {
			time.Sleep(time.Millisecond * 5000)
			br.Schedules.LoadDay(time.Now())
		}, "main->Loading Schedules")

		app.Run()
	}
}`)

	createFile("/.gitignore", `*.idea
*.pyc
db/bootstrap/*/mongoDump
localWebConfig.json
nohup.out
debug
.happypack
web/swagger/dist/swagger.*
/models/
/webAPIs/
/log/*
webConfig.json
/dist
*.upx
docker/dist
/updates/latest
.DS_Store
.history
*.vscode
`+appName)
}
