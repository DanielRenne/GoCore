package sqliteDBServices

import (
	"core/serverSettings"
	"database/sql"
	"fmt"
	"github.com/fatih/color"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

var DB *sql.DB
var BoltDB *storm.DB

func init() {
	fmt.Println("core dbServices initialized.")

	switch serverSettings.WebConfig.DbConnection.Driver {
	case "sqlite3":
		openSQLDriver()
	}

}

func openSQLDriver() {
	var err error
	DB, err = sql.Open(serverSettings.WebConfig.DbConnection.Driver, serverSettings.WebConfig.DbConnection.ConnectionString)

	if err != nil {
		color.Red("Open connection failed:" + err.Error())
		return
	}

	color.Cyan("Open Database Connections: " + string(DB.Stats().OpenConnections))
}
