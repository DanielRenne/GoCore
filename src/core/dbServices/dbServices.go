package dbServices

import (
	"core/serverSettings"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func init() {
	fmt.Println("core dbServices initialized.")

	var err error
	DB, err = sql.Open(serverSettings.WebConfig.DbConnection.Driver, serverSettings.WebConfig.DbConnection.ConnectionString)

	if err != nil {
		fmt.Println("Open connection failed:" + err.Error())
		return
	}

	fmt.Println("Open Database Connections: " + string(DB.Stats().OpenConnections))
}
