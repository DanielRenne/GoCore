package dbServices

import (
	"core/serverSettings"
	"database/sql"
	"fmt"
	"github.com/asdine/storm"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/fatih/color"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path"
)

var DB *sql.DB
var BoltDB *storm.DB

func init() {
	fmt.Println("core dbServices initialized.")

	switch serverSettings.WebConfig.DbConnection.Driver {
	case "mysql":
		openSQLDriver()
	case "mssql":
		openSQLDriver()
	case "sqlite3":
		openSQLDriver()
	case "boltDB":
		openBolt()
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

func openBolt() {

	myDBDir := "db/" + serverSettings.WebConfig.DbConnection.AppName + "/" + serverSettings.WebConfig.DbConnection.ConnectionString

	os.Mkdir(path.Dir(myDBDir), 0777)

	// (Create if not exist) open a database
	var err error
	BoltDB, err = storm.Open(myDBDir, storm.AutoIncrement())
	if err != nil {
		color.Red("Failed to create or open boltDB Database at " + myDBDir + ":\n\t" + err.Error())
	}

	color.Cyan("Successfully opened new bolt DB at " + myDBDir)

}
