package serverSettings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
)

type htmlTemplates struct {
	Enabled         bool   `json:"enabled"`
	Directory       string `json:"directory"`
	DirectoryLevels int    `json:"directoryLevels"`
}

type dbConnection struct {
	ConnectionString string `json:"connectionString"`
	Driver           string `json:"driver"`
}

type license struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type contact struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Email string `json:"email"`
}

type info struct {
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	Contact        contact `json:"contact"`
	License        license `json:"license"`
	TermsOfService string  `json:termsOfService"`
}

type application struct {
	Name           string        `json:"name"`
	Domain         string        `json:"domain"`
	HttpPort       int           `json:"httpPort"`
	HttpsPort      int           `json:"httpsPort"`
	ReleaseMode    string        `json:"releaseMode"`
	WebServiceOnly bool          `json:"webServiceOnly"`
	Info           info          `json:"info"`
	HtmlTemplates  htmlTemplates `json:"htmlTemplates"`
}

type webConfigObj struct {
	DbConnections []dbConnection `json:"dbConnections"`
	Application   application    `json:"application"`
	DbConnection  dbConnection
}

var WebConfig webConfigObj
var APP_LOCATION string
var GOCORE_PATH string
var SWAGGER_UI_PATH string

func Initialize(path string) {

	APP_LOCATION = path
	SWAGGER_UI_PATH = APP_LOCATION + "/web/swagger/dist"
	setGoCorePath()

	fmt.Println("core serverSettings initialized.")

	jsonData, err := ioutil.ReadFile(APP_LOCATION + "/webConfig.json")
	if err != nil {
		fmt.Println("Reading of webConfig.json failed:  " + err.Error())
		return
	}

	errUnmarshal := json.Unmarshal(jsonData, &WebConfig)
	if errUnmarshal != nil {
		fmt.Println("Parsing / Unmarshaling of webConfig.json failed:  " + errUnmarshal.Error())
		return
	}

	for _, dbConnection := range WebConfig.DbConnections {
		WebConfig.DbConnection = dbConnection
		return
	}
}

//Sets the GoCore path for go core packages to reference.
func setGoCorePath() {

	_, filename, _, ok := runtime.Caller(1)

	if ok == true {
		GOCORE_PATH = strings.Replace(filename[strings.Index(filename, "/src"):], "/core/serverSettings/serverSettings.go", "", -1)
	} else {
		GOCORE_PATH = "src/github.com/DanielRenne/GoCore"
	}
}
