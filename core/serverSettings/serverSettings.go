package serverSettings

import (
	"encoding/json"
	"fmt"
	"github.com/DanielRenne/GoCore/core/appGen"
	"io/ioutil"
)

type htmlTemplates struct {
	Enabled         bool   `json:"enabled"`
	Directory       string `json:"directory"`
	DirectoryLevels int    `json:"directoryLevels"`
}

type dbConnection struct {
	AppName          string `json:"appName"`
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

const SwaggerUIPath = appGen.APP_LOCATION + "/web/swagger/dist"

func init() {
	fmt.Println("core serverSettings initialized.")

	jsonData, err := ioutil.ReadFile(appGen.APP_LOCATION + "/webConfig.json")
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
		if dbConnection.AppName == WebConfig.Application.Name {
			WebConfig.DbConnection = dbConnection
			return
		}
	}

	// webConfigJSON, errParse := gabs.ParseJSONFile("webConfig.json")

	// if errParse != nil {
	// 	fmt.Println("Error parsing webConfig", errParse.Error())
	// }

	// appName, ok := webConfigJSON.Path("application.name").Data().(string)
	// if ok {

	// 	children, _ := webConfigJSON.S("dbConnections").Children()
	// 	for _, child := range children {

	// 		if child.S("appName").Data().(string) == appName {
	// 			WebConfig.DbConnection.ConnectionString = child.S("connectionString").Data().(string)
	// 			WebConfig.DbConnection.Driver = child.S("driver").Data().(string)
	// 			WebConfig.DbConnection.AppName = child.S("appName").Data().(string)
	// 		}

	// 	}
	// }
	// return
}
