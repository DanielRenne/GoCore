package serverSettings

import (
	"fmt"
	"github.com/Jeffail/gabs"
)

type dbConnection struct {
	AppName          string
	ConnectionString string
	Driver           string
}

type webConfigObj struct {
	DbConnection dbConnection
}

var WebConfig webConfigObj

func init() {
	fmt.Println("core serverSettings initialized.")

	webConfigJSON, errParse := gabs.ParseJSONFile("webConfig.json")

	if errParse != nil {
		fmt.Println("Error parsing webConfig", errParse.Error())
	}

	appName, ok := webConfigJSON.Path("application.name").Data().(string)
	if ok {

		children, _ := webConfigJSON.S("dbConnections").Children()
		for _, child := range children {

			if child.S("appName").Data().(string) == appName {
				WebConfig.DbConnection.ConnectionString = child.S("connectionString").Data().(string)
				WebConfig.DbConnection.Driver = child.S("driver").Data().(string)
				WebConfig.DbConnection.AppName = child.S("appName").Data().(string)
			}

		}
	}
	return
}
