package dbServices

import (
	"database/sql"
	"fmt"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/asdine/storm"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/fatih/color"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/mgo.v2"
	"os"
	"path"
)

var DB *sql.DB
var BoltDB *storm.DB
var MongoSession *mgo.Session
var MongoDB *mgo.Database
var DatabaseInitialized chan int
var dbInitializedCount int

const (

	//Driver Types
	DATABASE_DRIVER_MYSQL   = "mysql"
	DATABASE_DRIVER_MSSQL   = "mssql"
	DATABASE_DRIVER_BOLTDB  = "boltDB"
	DATABASE_DRIVER_MONGODB = "mongoDB"
)

func init() {
	DatabaseInitialized = make(chan int, 1)
}

func Initialize() {

	fmt.Println("core dbServices initialized.")

	switch serverSettings.WebConfig.DbConnection.Driver {
	case DATABASE_DRIVER_MYSQL:
		openSQLDriver()
	case DATABASE_DRIVER_MSSQL:
		openSQLDriver()
	case DATABASE_DRIVER_BOLTDB:
		openBolt()
	case DATABASE_DRIVER_MONGODB:
		openMongo()
	}
}

func WaitForDatabase() chan int {
	dbInitializedCount++
	return DatabaseInitialized
}

func openSQLDriver() {
	var err error
	DB, err = sql.Open(serverSettings.WebConfig.DbConnection.Driver, serverSettings.WebConfig.DbConnection.ConnectionString)
	notifyDBWaits()
	if err != nil {
		color.Red("Open connection failed:" + err.Error())
		return
	}

	color.Cyan("Open Database Connections: " + string(DB.Stats().OpenConnections))
}

func openBolt() {

	myDBDir := serverSettings.APP_LOCATION + "/db/" + serverSettings.WebConfig.DbConnection.ConnectionString

	os.Mkdir(path.Dir(myDBDir), 0777)

	// (Create if not exist) open a database
	var err error
	BoltDB, err = storm.Open(myDBDir, storm.AutoIncrement())

	if err != nil {
		color.Red("Failed to create or open boltDB Database at " + myDBDir + ":\n\t" + err.Error())
		return
	}
	notifyDBWaits()

	color.Cyan("Successfully opened new bolt DB at " + myDBDir)

}

func openMongo() {

	var err error
	MongoSession, err = mgo.Dial(serverSettings.WebConfig.DbConnection.ConnectionString) // open an connection -> Dial function
	if err != nil {                                                                      //  if you have a
		color.Red("Failed to create or open mongo Database at " + serverSettings.WebConfig.DbConnection.ConnectionString + ":\n\t" + err.Error())
		return
	}

	MongoSession.SetMode(mgo.Monotonic, true) // Optional. Switch the session to a monotonic behavior.

	MongoDB = MongoSession.DB(serverSettings.WebConfig.DbConnection.Database)
	notifyDBWaits()

}

func notifyDBWaits() {
	for i := 0; i < dbInitializedCount; i++ {
		DatabaseInitialized <- 1
	}
	dbInitializedCount = 0
}
