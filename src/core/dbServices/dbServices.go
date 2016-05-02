package dbServices

import (
	"core/serverSettings"
	"database/sql"
	"fmt"
	tiedotDriver "github.com/HouzuoGuo/tiedot/db"
	// tiedotError "github.com/HouzuoGuo/tiedot/dberr"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB
var Tiedot *tiedotDriver.DB

func init() {
	fmt.Println("core dbServices initialized.")

	switch serverSettings.WebConfig.DbConnection.Driver {
	case "mysql":
		openSQLDriver()
	case "mssql":
		openSQLDriver()
	case "sqlite3":
		openSQLDriver()
	case "tiedot":
		openTiedot()
	}

}

func openSQLDriver() {
	var err error
	DB, err = sql.Open(serverSettings.WebConfig.DbConnection.Driver, serverSettings.WebConfig.DbConnection.ConnectionString)

	if err != nil {
		fmt.Println("Open connection failed:" + err.Error())
		return
	}

	fmt.Println("Open Database Connections: " + string(DB.Stats().OpenConnections))
}

func openTiedot() {

	myDBDir := "db/" + serverSettings.WebConfig.DbConnection.AppName + "/" + serverSettings.WebConfig.DbConnection.ConnectionString

	// (Create if not exist) open a database
	Tiedot, err := tiedotDriver.OpenDB(myDBDir)
	if err != nil {
		fmt.Println("Failed to create or open tiedot Database at " + myDBDir + ":\n\t" + err.Error())
	}

	fmt.Println("Successfully opened new tiedot DB at " + myDBDir)
	for _, collection := range Tiedot.AllCols() {
		fmt.Println("First tiedot collection:  " + collection)
		break
	}
}

// Return all collection names.
func tieDotAllCols() (ret []string) {
	ret = Tiedot.AllCols()
	return
}
