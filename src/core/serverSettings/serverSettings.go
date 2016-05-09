package serverSettings

import (
	"encoding/json"
	"fmt"
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

type application struct {
	Name           string        `json:"name"`
	Domain         string        `json:"domain"`
	HttpPort       int           `json:"httpPort"`
	HttpsPort      int           `json:"httpsPort"`
	ReleaseMode    string        `json:"releaseMode"`
	WebServiceOnly bool          `json:"webServiceOnly"`
	HtmlTemplates  htmlTemplates `json:"htmlTemplates"`
}

type webConfigObj struct {
	DbConnections []dbConnection `json:"dbConnections"`
	Application   application    `json:"application"`
	DbConnection  dbConnection
}

var WebConfig webConfigObj

func init() {
	fmt.Println("core serverSettings initialized.")

	jsonData, err := ioutil.ReadFile("webConfig.json")
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
