// Package appGen internal only to scaffold applications running GoCore
package appGen

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/logger"
	"github.com/DanielRenne/GoCore/core/path"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/DanielRenne/GoCore/core/utils"
	"github.com/globalsign/mgo/bson"
)

func cdGoPath() {
	os.Chdir(os.Getenv("GOPATH"))
}

// GenerateApp internal only
func GenerateApp() {
	cdGoPath()
	moveAppFiles()
}

// GenerateServerApp internal only
func GenerateServerApp() {
	moveServerOnlyAppFiles()
}

func moveServerOnlyAppFiles() {

	//First check for the WebConfig.json file
	var err error
	var isInitializingApp bool
	var totalSuccesses int
	totalExpectedSuccesses := 12
	_, err = os.Stat(serverSettings.APP_LOCATION + path.PathSeparator + "webConfig.json")

	if err != nil {
		isInitializingApp = true
		var mainAppName string
		for {
			reader := bufio.NewReader(os.Stdin)
			log.Println("Please enter the camelCase name of your app:")
			mainAppName, _ = reader.ReadString('\n')
			mainAppName = strings.Trim(mainAppName, "\n")
			if mainAppName != "" {
				break
			}
		}

		webConfig := `
{
	"application":{
		"logGophers": true,
		"logGopherInterval": 15,
		"domain": "0.0.0.0",
		"serverFQDN": "0.0.0.0",
		"httpPort": 80,
		"httpsPort": 443,
		"releaseMode":"development",
		"versionNumeric": 1,
		"versionDot": "0.0.1",
		"productName": "` + mainAppName + `",
		"disableRootIndex": true,
		"disableWebSockets": false,
		"sessionKey":"` + mainAppName + `SessionKey",
		"sessionName":"` + mainAppName + `ProductName",
		"sessionExpirationDays":3650,
		"sessionSecureCookie":false,
		"bootstrapData":true, 
		"htmlTemplates":{
			"enabled":false,
			"directory":"templates",
			"directoryLevels": 1
		}
	},
	"dbConnections":[
		{
			"driver" : "mongoDB",
			"connectionString": "mongodb://127.0.0.1:27017/` + mainAppName + `",
			"database": "` + mainAppName + `"
		}
	]
}
`
		err := extensions.Write(webConfig, serverSettings.APP_LOCATION+path.PathSeparator+"webConfig.json")
		if err != nil {
			logger.Message("failed to create webConfig.json"+err.Error(), logger.RED)
		} else {
			totalSuccesses++
			logger.Message("Copied webConfig.json to Application.", logger.GREEN)
		}
	}

	if isInitializingApp {
		var mainCNKeys string
		for {
			reader := bufio.NewReader(os.Stdin)
			log.Println("We are now attempting to generate SSL self signed certificates.  Add your full cert information like this: \"/CN=www.mydom.com/O=My Company Name LTD./C=US\" (defaults to this if you just press enter)")
			mainCNKeys, _ = reader.ReadString('\n')
			mainCNKeys = strings.Trim(mainCNKeys, "\n")
			if mainCNKeys == "" {
				mainCNKeys = "/CN=www.mydom.com/O=My Company Name LTD./C=US"
			}
			break
		}

		_, err = os.Stat(serverSettings.APP_LOCATION + path.PathSeparator + "keys")

		if err != nil {
			err = extensions.MkDir(serverSettings.APP_LOCATION + path.PathSeparator + "keys")
			if err == nil {
				totalSuccesses++
				err := os.Chdir(serverSettings.APP_LOCATION + path.PathSeparator + "keys")
				if err == nil {

					cmd := exec.Command("openssl", "req", "-newkey", "rsa:2048", "-new", "-nodes", "-x509", "-days", "13650", "-subj", "'"+mainCNKeys+"'", "-keyout", "key.pem")
					err = cmd.Start()
					if err != nil {
						logger.Message("Failed to create keys with openssl for application.", logger.RED)
					} else {
						totalSuccesses++
						cmd := exec.Command("openssl", "req", "-new", "-subj", "'"+mainCNKeys+"'", "-key", "key.pem", "-out", "cert.pem")
						err = cmd.Start()
						if err != nil {
							logger.Message("Created to create cert with openssl for application.", logger.RED)
						} else {
							totalSuccesses++
							logger.Message("Created keys for application.", logger.GREEN)
						}

					}
				}
			} else {
				logger.Message("Couldnt create keys dir", logger.RED)
			}
		}

		_, err = os.Stat(serverSettings.APP_LOCATION + path.PathSeparator + "db")

		if err != nil {
			err = extensions.MkDir(serverSettings.APP_LOCATION + path.PathSeparator + "db" + path.PathSeparator + "schemas" + path.PathSeparator + "1.0.0")
			if err == nil {
				totalSuccesses++
				logger.Message("Created db/schemas/1.0.0 to Application.", logger.GREEN)

				err = extensions.MkDir(serverSettings.APP_LOCATION + path.PathSeparator + "db" + path.PathSeparator + "goFiles" + path.PathSeparator + "v1")
				if err == nil {
					totalSuccesses++
					logger.Message("Created goFiles for Application.", logger.GREEN)

					err := extensions.Write(`package model
// Anything go files you put here will be combined into the model package so you can extend whatever you want after the models have been generated
`, serverSettings.APP_LOCATION+path.PathSeparator+"db"+path.PathSeparator+"goFiles"+path.PathSeparator+"v1"+path.PathSeparator+"blankExample.go")
					if err != nil {
						logger.Message("failed to create blankExample.go "+err.Error(), logger.RED)
					} else {
						totalSuccesses++
						logger.Message("Created blankExample.go  for Application.", logger.GREEN)
					}
				} else {
					logger.Message("failed to create goFiles dir: "+err.Error(), logger.RED)
				}
			} else {
				logger.Message("Couldnt create db/schemas/1.0.0 dir", logger.RED)
			}
		}

		_, err = os.Stat(serverSettings.APP_LOCATION + path.PathSeparator + "log")

		if err != nil {
			err = extensions.MkDirRWAll(serverSettings.APP_LOCATION + path.PathSeparator + "log")
			if err == nil {
				totalSuccesses++
				logger.Message("Created log for Application.", logger.GREEN)
			} else {
				logger.Message("Failed to create log dir.", logger.RED)
			}
		}

		_, err = os.Stat(serverSettings.APP_LOCATION + path.PathSeparator + "db" + path.PathSeparator + "bootstrap")

		if err != nil {
			err = extensions.MkDir(serverSettings.APP_LOCATION + path.PathSeparator + "db" + path.PathSeparator + "bootstrap")
			if err == nil {
				totalSuccesses++
				logger.Message("Created /db/bootstrap for Application.", logger.GREEN)
			} else {
				logger.Message("Failed to create /db/bootstrap dir.", logger.RED)
			}
		}

		_, err = os.Stat(serverSettings.APP_LOCATION + path.PathSeparator + "settings")

		if err != nil {
			err = extensions.MkDir(serverSettings.APP_LOCATION + path.PathSeparator + "settings")
			if err == nil {
				totalSuccesses++
				logger.Message("Created settings Application.", logger.GREEN)
				settingsFile := `
// Package settings provides settings for go files to reference paths and other constants to know where files are located.
package settings

import (	
	"log"
	"os"
	"sync"
	"encoding/json"
	"github.com/DanielRenne/GoCore/core/serverSettings"
)
var Version string 

type appSettings struct {
	DeveloperMode                      bool     ` + "`" + `json:"developerMode"` + "`" + ` 
	// TODO: Add more application specific settings here
}

type webConfig struct {
	AppSettings appSettings ` + "`" + `json:"appSettings"` + "`" + ` 
}

type FullWebConfig struct {
	Application serverSettings.Application ` + "`" + `json:"application"` + "`" + ` 
	DbConnections []serverSettings.DbConnection ` + "`" + `json:"dbConnections"` + "`" + ` 
	AppSettings struct {
		DeveloperMode                      bool          ` + "`" + `json:"developerMode"` + "`" + `
		// TODO: Add more application specific settings here
	} ` + "`" + `json:"appSettings"` + "`" + ` 
}
var AppSettings appSettings
var ServerSettings serverSettings.Application
var AppSettingsSync sync.RWMutex

func Initialize() {
	log.Println("Settings Initialized.")
	if serverSettings.APP_LOCATION == "" {
		log.Println("server settings APP_LOCATION is blank.  This is not right!")
		os.Exit(1)
	}
	ServerSettings = serverSettings.WebConfig.Application

	jsonData, err := os.ReadFile(serverSettings.APP_LOCATION + "/webConfig.json")
	if err != nil {
		log.Println("Reading of webConfig.json failed at settings.init():  " + err.Error())
		os.Exit(1)
	}

	var config webConfig

	errUnmarshal := json.Unmarshal(jsonData, &config)
	if errUnmarshal != nil {
		log.Println("Parsing / Unmarshaling of webConfig.json failed:  " + errUnmarshal.Error())
		os.Exit(1)
	}

	AppSettingsSync.Lock()
	AppSettings = config.AppSettings
	AppSettingsSync.Unlock()
}
`
				err := extensions.Write(settingsFile, serverSettings.APP_LOCATION+path.PathSeparator+"settings"+path.PathSeparator+"settings.go")
				if err != nil {
					logger.Message("failed to settings.go: "+err.Error(), logger.RED)
				} else {
					totalSuccesses++
					logger.Message("Created settings/settings.go for Application.", logger.GREEN)
				}
			} else {
				logger.Message("Couldnt settings dir", logger.RED)
			}
		}

		var mainNameFileGo string
		for {
			reader := bufio.NewReader(os.Stdin)
			log.Println("What do you want to call your main package fileName?")
			mainNameFileGo, _ = reader.ReadString('\n')
			mainNameFileGo = strings.Trim(mainNameFileGo, "\n")
			ok := false
			if strings.Index(mainNameFileGo, " ") == -1 {
				ok = true
			} else {
				logger.Message("No spaces please", logger.GREEN)
			}
			if ok {
				if strings.Index(mainNameFileGo, ".go") == -1 {
					mainNameFileGo += ".go"
				}
				break
			}
		}

		_, err = os.Stat(serverSettings.APP_LOCATION + path.PathSeparator + ".gitignore")

		if err != nil {
			err := extensions.Write(`*.idea
*.pyc
nohup.out
debug
/log/*
appModelBuild
`+strings.ReplaceAll(mainNameFileGo, ".go", "")+`
webConfig.json
.DS_Store
.history
*.db
/models/
`, serverSettings.APP_LOCATION+path.PathSeparator+".gitignore")
			if err != nil {
				logger.Message("failed to create .gitignore"+err.Error(), logger.RED)
			} else {
				totalSuccesses++
				logger.Message("Created .gitignore for Application.", logger.GREEN)
			}
		}

		// todo parse go.mod for name
		//module github.com/your_example/repo_project
		var moduleName string
		for {
			reader := bufio.NewReader(os.Stdin)
			logger.Message("What is your go module name?", logger.GREEN)
			moduleName, _ = reader.ReadString('\n')
			moduleName = strings.Trim(moduleName, "\n")
			ok := false
			if strings.Index(moduleName, " ") == -1 {
				ok = true
			} else {
				logger.Message("No spaces please", logger.GREEN)
			}
			if ok {
				break
			}
		}

		addCronJob := "y"
		for {
			reader := bufio.NewReader(os.Stdin)
			log.Println("Do you want to include cron jobs to your main.go? ('y' or 'n')")
			addCronJob, _ = reader.ReadString('\n')
			addCronJob = strings.Trim(addCronJob, "\n")
			if addCronJob == "" {
				addCronJob = "y"
			}
			ok := false
			if addCronJob == "y" || addCronJob == "n" {
				ok = true
			} else {
				logger.Message("Invalid type 'n' or 'y'", logger.RED)
			}
			if ok {
				break
			}
		}

		cronImport := ""
		cronStartCode := ""

		if addCronJob == "y" {
			totalExpectedSuccesses = totalExpectedSuccesses + 2
			cronImport = "\t\"github.com/DanielRenne/GoCore/core/cron\"\n//" + strings.ReplaceAll(mainNameFileGo, ".go", "") + "cron \"github.com/yourusername/package/cron\"\n"
			cronStartCode = "\t//" + strings.ReplaceAll(mainNameFileGo, ".go", "") + "cron.Start()"
			_, err = os.Stat(serverSettings.APP_LOCATION + path.PathSeparator + "cron")

			if err != nil {
				err = extensions.MkDir(serverSettings.APP_LOCATION + path.PathSeparator + "cron")
				if err == nil {
					totalSuccesses++
					logger.Message("Created /cron dir for Application.", logger.GREEN)
					cronFile := `
package cron

import (
	"log"
	"runtime/debug"
	"time"

	//"github.com/DanielRenne/GoCore/core/cron"
)

func Start() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("\n\nPanic Stack: " + string(debug.Stack()))
			time.Sleep(time.Millisecond * 3000)
			Start()
			return
		}
	}()
	// setup cron jobs for your app you can do
	// core.CRON_TOP_OF_HOUR
	// core.CRON_TOP_OF_SECOND
	// And much more, just check the core package for more or send a pull request if you need faster tickers
	// go core.RegisterRecurring(core.CRON_TOP_OF_30_SECONDS, yourCronFunction)
}
`
					err := extensions.Write(cronFile, serverSettings.APP_LOCATION+path.PathSeparator+"cron"+path.PathSeparator+"cron.go")
					if err != nil {
						logger.Message("failed to create cron go file"+err.Error(), logger.RED)
					} else {
						totalSuccesses++
						logger.Message("Created cron/cron.go for Application.", logger.GREEN)
					}
				} else {
					logger.Message("Failed to create /cron dir.", logger.RED)
				}
			}
		}
		if totalExpectedSuccesses != totalSuccesses {
			logger.Message("Something failed in the GoCore app generation, please look above in the logs", logger.RED)
		}
		mainFile := `
package main

import (
	"os"
	"net/http"
	"time"
	"log" 
	"runtime/debug"
	"net"
	"github.com/DanielRenne/GoCore/core/app"
	"github.com/DanielRenne/GoCore/core/ginServer"
	"github.com/DanielRenne/GoCore/core/dbServices"
	"github.com/gin-gonic/gin"
	
	_ "` + moduleName + `/models/v1/model"
	"` + moduleName + `/settings"
	` + cronImport + `
)

func loadApp(c *gin.Context) {
	handleRespondHTML(c, []byte(appIndex(c)), time.Now())
}

func handleRespondHTML(c *gin.Context, data []byte, modTime time.Time) {
	ginServer.RenderHTML(string(data), c)
}

func appIndex(c *gin.Context) (htmlContent string) {
	// Show mongo restart page if it cannot dial to port
	dialer, _ := dbServices.GetMongoDialInfo()
	conn, err := net.Dial("tcp", dialer.Addrs[0])
	if err != nil { 
		htmlContent = ` + "`" + `
			<span style="color: red">
				<h2>
					An error occurred and application cannot run due to mongo database being down.<br/><br/> 
					Error description: ` + "`" + ` + err.Error() + ` + "`" + `<br/><br/>
				</h2> 
			</span>
			` + "`" + `
		return
	}
	conn.Close()
	htmlContent = "<html><body>Hello goCore World!</body></html>"
	return
}
 
func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("\n\nPanic Stack: " + string(debug.Stack()))
			time.Sleep(time.Millisecond * 3000)
			main()
			return
		}
	}()
	
	// By default go core will mount to web a folder matching web in your app for static hosting but for your initial app, we are turning this off, if you need it create a directory called web or whatever you enter below
	app.StaticWebLocation = ""
	app.Init()
	settings.Initialize()
		// Give mongo 2 seconds to connect
	time.Sleep(time.Second * 2) 
	
	// Detect if port is open on mongo, if not do not do anything database wise because that will block until a connection is established.
	// Most likely you will attempt to build out database models
	dialer, _ := dbServices.GetMongoDialInfo()
	conn, err := net.Dial("tcp", dialer.Addrs[0])
	if err != nil {
		ginServer.Router.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "` + "%v" + `", "An error occurred and this application cannot run (due to mongo database services being down).\n\n Error description: "+err.Error())
		})
		go app.RunServer()
		time.Sleep(time.Minute)
		os.Exit(1) // systemd daemon should respawn your main program
	}
	conn.Close()
	
	` + cronStartCode + `
	
	// Put all your application code here:
	
	// Controllers:
	// Add all your ginServer.Router methods preferrably in your own controllers package
	// Move appIndex() handleRespondHTML() and loadApp() to your controllers package
	ginServer.Router.GET("/", loadApp)
	
	// Block and run the server ports
	app.Run()
}`
		err := extensions.Write(mainFile, serverSettings.APP_LOCATION+path.PathSeparator+mainNameFileGo)
		if err != nil {
			logger.Message("failed to create main go file"+err.Error(), logger.RED)
		} else {
			if totalExpectedSuccesses == totalSuccesses {
				msg := "GoCore " + mainNameFileGo + " and other files generated in your module successfully."
				logger.Message(msg+"  Please `go build "+mainNameFileGo+" && ./"+strings.ReplaceAll(mainNameFileGo, ".go", "")+"` to get started running your app for the first time", logger.GREEN)
				go exec.Command("say", msg).Output()
			}
		}
	}
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
		extensions.CopyFolder("/tmp/tools/appFiles"+path, serverSettings.APP_LOCATION+path)
		logger.Message("Created "+path+" in Application.", logger.GREEN)
	}
	return
}

func replacePath(path string, newpath string, newGithubUser string, newProject string) {
	filepath.Walk(serverSettings.APP_LOCATION+path, func(pth string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			utils.ReplaceTokenInFile(pth, "github.com/DanielRenne/goCoreAppTemplate", newpath)
			// we cant just globally replace DanielRenne with the new github username, so we use the special token DanielRenneFolder

			utils.ReplaceTokenInFile(pth, "goCoreAppPath", newProject+"_path")
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

	_, err := os.Stat("/tmp/tools/appFiles")
	if err == nil {
		humanTitle, err := extensions.ReadFile("/tmp/humanTitle")
		if err != nil {
			log.Println("error reading humanTitle")
			os.Exit(1)
		}
		mainCNKeys, err := extensions.ReadFile("/tmp/mainCNKeys")
		if err != nil {
			log.Println("error reading mainCNKeys")
			os.Exit(1)
		}
		username, err := extensions.ReadFile("/tmp/username")
		if err != nil {
			log.Println("error reading username")
			os.Exit(1)
		}
		foundDbType := false
		databaseType, errDatabaseFile := extensions.ReadFile("/tmp/databaseType")
		if errDatabaseFile != nil {
			log.Println("error reading databaseType")
			os.Exit(1)
		} else {
			os.Remove("/tmp/databaseType")
			foundDbType = true
		}
		_, errDatabaseType := os.Stat(serverSettings.APP_LOCATION + "/databaseType")
		if errDatabaseType == nil {
			databaseType, err = extensions.ReadFile(serverSettings.APP_LOCATION + "/databaseType")
			if err != nil {
				log.Println("error reading databaseType local")
				os.Exit(1)
			}
			foundDbType = false
		}
		parts := strings.Split(serverSettings.APP_LOCATION, "/")
		appName := parts[len(parts)-1]
		githubName := string(username)
		project := githubName + "/" + appName
		//First check for the WebConfig.json file
		_, errNoWebConfig := os.Stat(serverSettings.APP_LOCATION + "/webConfig.json")
		if errNoWebConfig != nil {
			if string(databaseType) == "mongo" || string(databaseType) == "" {
				extensions.CopyFile("/tmp/tools/appFiles/webConfig.json", serverSettings.APP_LOCATION+"/webConfig.json")
			} else if string(databaseType) == "bolt" {
				extensions.CopyFile("/tmp/tools/appFiles/webConfig.bolt.json", serverSettings.APP_LOCATION+"/webConfig.json")
			}
			if foundDbType {
				createFile("/databaseType", string(databaseType))
			}
			logger.Message("Copied webConfig.json to Application.", logger.GREEN)
		}

		_, err = os.Stat(serverSettings.APP_LOCATION + "/webConfig.prod.json")
		if err != nil {
			extensions.CopyFile("/tmp/tools/appFiles/webConfig.prod.json", serverSettings.APP_LOCATION+"/webConfig.prod.json")
			logger.Message("Copied webConfig.json to Application.", logger.GREEN)
		}

		_, err = os.Stat(serverSettings.APP_LOCATION + "/webConfig.dev.json")
		if err != nil {
			extensions.CopyFile("/tmp/tools/appFiles/webConfig.dev.json", serverSettings.APP_LOCATION+"/webConfig.dev.json")
			logger.Message("Copied webConfig.json to Application.", logger.GREEN)
		}

		for _, v := range utils.Array("webConfig.prod.json", "webConfig.dev.json", "webConfig.json") {
			id2 := bson.NewObjectId()
			utils.ReplaceTokenInFile(serverSettings.APP_LOCATION+"/"+v, "goCoreProductName", appName+"BaseProduct")
			utils.ReplaceTokenInFile(serverSettings.APP_LOCATION+"/"+v, "goCoreSessionKey", id2.Hex())
		}

		_, err = os.Stat(serverSettings.APP_LOCATION + "/log")
		if err != nil {
			os.MkdirAll(serverSettings.APP_LOCATION+"/log/plugins", 0777)
		}
		var wasCopied bool
		_, err = os.Stat(serverSettings.APP_LOCATION + path.PathSeparator + "keys")

		if err != nil {
			err = extensions.MkDir(serverSettings.APP_LOCATION + path.PathSeparator + "keys")
			if err == nil {
				err := os.Chdir(serverSettings.APP_LOCATION + path.PathSeparator + "keys")
				if err == nil {

					cmd := exec.Command("openssl", "req", "-newkey", "rsa:2048", "-new", "-nodes", "-x509", "-days", "13650", "-subj", "'"+string(mainCNKeys)+"'", "-keyout", "key.pem")
					err = cmd.Start()
					if err != nil {
						logger.Message("Failed to create keys with openssl for application.", logger.RED)
					} else {
						cmd := exec.Command("openssl", "req", "-new", "-subj", "'"+string(mainCNKeys)+"'", "-key", "key.pem", "-out", "cert.pem")
						err = cmd.Start()
						if err != nil {
							logger.Message("Created to create cert with openssl for application.", logger.RED)
						} else {
							logger.Message("Created keys for application.", logger.GREEN)
						}
					}
				}
				err = os.Chdir(serverSettings.APP_LOCATION)
				if err != nil {
					logger.Message("Failed  to change dir to app location.", logger.RED)
				}
			} else {
				logger.Message("Couldnt create keys dir", logger.RED)
			}
		}

		wasCopied = copyFolder("/web")
		if wasCopied {
			utils.ReplaceTokenInFile(serverSettings.APP_LOCATION+"/web/app/watchFile.json", "github.com/DanielRenne/goCoreAppTemplate", project)
			utils.ReplaceTokenInFile(serverSettings.APP_LOCATION+"/web/app/javascript/build-css.sh", "github.com/DanielRenne/goCoreAppTemplate", project)
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

		if errDatabaseFile == nil {
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

		wasCopied = copyFolder("/controllerRegistry")
		if wasCopied {
			replacePath("/controllerRegistry", project, githubName, appName)
		}

		_, err = os.Stat(serverSettings.APP_LOCATION + "/db")

		if err != nil {
			os.MkdirAll(serverSettings.APP_LOCATION+"/db/schemas/1.0.0", 0777)
			os.MkdirAll(serverSettings.APP_LOCATION+"/db/bootstrap", 0777)
			os.MkdirAll(serverSettings.APP_LOCATION+"/db/goFiles/v1", 0777)
			extensions.WriteToFile("Put model class functions and overrides here", serverSettings.APP_LOCATION+"/db/goFiles/v1/.gitkeep", 0777)
			extensions.CopyFolder("/tmp/tools/appFiles/db/schemas", serverSettings.APP_LOCATION+"/db/schemas/1.0.0")
			extensions.CopyFolder("/tmp/tools/appFiles/db/bootstrap", serverSettings.APP_LOCATION+"/db/bootstrap")
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
	
	"github.com/DanielRenne/GoCore/core/dbServices"
	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/GoCore/core/app"
	"github.com/DanielRenne/GoCore/core/ginServer"
	"github.com/DanielRenne/GoCore/core/logger"
	_ "github.com/`+project+`/controllerRegistry"
	"github.com/`+project+`/br"
	"github.com/`+project+`/controllers"
	"github.com/`+project+`/cron"
	_ "github.com/`+project+`/models/v1/model"
	"github.com/`+project+`/sessionFunctions"
	"github.com/`+project+`/settings"
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
	
	app.Initialize(os.Getenv("`+appName+`_path"), "webConfig.json")
	settings.Initialize()
	br.Schedules.UpdateLinuxToGMT()

	dialer, _ := dbServices.GetMongoDialInfo()
	conn, err := net.Dial("tcp", dialer.Addrs[0])
	if err != nil {
		//lastError := err.Error()
		ginServer.Router.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "%v", "An error occurred and velocity cannot run (due to mongo database services being down).\n\n Error description: "+err.Error()+".")
		})
		go app.RunServer()
		time.Sleep(time.Minute)
		os.Exit(1) // systemd daemon should respawn your main program
	}
	conn.Close()

	if settings.AppSettings.DeveloperGoTrace {
		f, err := os.Create(os.Getenv("`+appName+`_path") + "/log/trace.log")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		err = trace.Start(f)
		if err != nil {
			panic(err)
		}
		mgo.SetDebug(true)

		file, _ := os.OpenFile(os.Getenv("`+appName+`_path") + "/log/studioMongo.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)

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
	
	// Blocking - Finally run your web server after starting cron jobs, setting up controllers.
	app.Run()
}`)

		createFile("/.gitignore", `*.idea
*.pyc
db/bootstrap/*/mongoDump
localWebConfig.json
nohup.out
debug
.happypack
/log/*
webConfig.json
/dist
*.upx
docker/dist
/updates/latest
.DS_Store
.history
*.vscode
web/app/npm-debug.log*
web/app/node_modules
*.db
/models/
package-lock.json
`+appName)

		createFile("/README.md", `# `+appName+` [a [GoCore Application](https://github.com/DanielRenne/GoCore/ "GoCore Application")]

Add an elevator description to pitch of what this GoCore web app does here.

## Setting up a development environment for this application ##

`+"`"+`cd /tmp/ && git clone github.com/`+githubName+`/`+appName+` && export `+appName+`_path=/tmp/`+appName+` && cd `+appName+` && bash bin/start_app`+"`"+`

Once your application is up and running login as username admin and password admin and start coding
`)
	}
}
