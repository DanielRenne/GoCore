package dbServices

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"sync"

	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/asdine/storm"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/fatih/color"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	_ "github.com/go-sql-driver/mysql"
)

var DBMutex *sync.RWMutex
var DB *sql.DB
var BoltDB *storm.DB
var MongoSession *mgo.Session
var MongoSessionAuth *mgo.Session
var MongoDB *mgo.Database
var MongoDBAuth *mgo.Database

var mongoDBOverride string
var mongoDBNameOverride string

const (

	//Driver Types
	DATABASE_DRIVER_MYSQL   = "mysql"
	DATABASE_DRIVER_MSSQL   = "mssql"
	DATABASE_DRIVER_BOLTDB  = "boltDB"
	DATABASE_DRIVER_MONGODB = "mongoDB"
)

func init() {
	DBMutex = &sync.RWMutex{}
}

func OverrideMongoDBConnection(connectionString string, dbName string) {
	mongoDBOverride = connectionString
	mongoDBNameOverride = dbName
}

func ReadMongoDB() (mdb *mgo.Database) {
	DBMutex.RLock()
	mdb = MongoDB
	DBMutex.RUnlock()
	return mdb
}

func ReadMongoDBAuth() (mdb *mgo.Database) {
	DBMutex.RLock()
	mdb = MongoDBAuth
	DBMutex.RUnlock()
	return mdb
}

func Initialize() error {

	fmt.Println("core dbServices initialized.")

	switch serverSettings.WebConfig.DbConnection.Driver {
	case DATABASE_DRIVER_MYSQL:
		return openSQLDriver()
	case DATABASE_DRIVER_MSSQL:
		return openSQLDriver()
	case DATABASE_DRIVER_BOLTDB:
		return openBolt()
	case DATABASE_DRIVER_MONGODB:
		return openMongo()
	}
	return nil
}

func openSQLDriver() error {
	var err error
	DBMutex.Lock()
	DB, err = sql.Open(serverSettings.WebConfig.DbConnection.Driver, serverSettings.WebConfig.DbConnection.ConnectionString)
	DBMutex.Unlock()

	if err != nil {
		color.Red("Open connection failed:" + err.Error())
		return err
	}

	DBMutex.RLock()
	color.Cyan("Open Database Connections: " + string(DB.Stats().OpenConnections))
	DBMutex.RUnlock()
	return nil
}

func openBolt() error {

	myDBDir := serverSettings.APP_LOCATION + "/db/" + serverSettings.WebConfig.DbConnection.ConnectionString

	os.Mkdir(path.Dir(myDBDir), 0777)

	// (Create if not exist) open a database
	var err error
	DBMutex.Lock()
	BoltDB, err = storm.Open(myDBDir)
	DBMutex.Unlock()

	if err != nil {
		color.Red("Failed to create or open boltDB Database at " + myDBDir + ":\n\t" + err.Error())
		return err
	}

	color.Cyan("Successfully opened new bolt DB at " + myDBDir)
	return nil
}

func openMongo() error {

	if serverSettings.WebConfig.DbConnection.Replication.Enabled {
		info := new(mgo.DialInfo)
		info.Direct = true
		info.Timeout = time.Millisecond * 3000

		var addresses []string
		addresses = append(addresses, serverSettings.WebConfig.DbConnection.Replication.Master)
		// for i, _ := range serverSettings.WebConfig.DbConnection.Replication.Slaves {
		// 	slave := serverSettings.WebConfig.DbConnection.Replication.Slaves[i]
		// 	addresses = append(addresses, slave)
		// }
		info.Addrs = addresses

		session, err := mgo.DialWithInfo(info)
		// session, err := mgo.Dial(serverSettings.WebConfig.DbConnection.Replication.SessionConnection)
		if err != nil { //  if you have a
			color.Red("Failed to create or open mongo Database to initialize replicaSet at " + serverSettings.WebConfig.DbConnection.Replication.Master + "\n\t" + err.Error())
		} else {
			time.Sleep(time.Millisecond * 500)
			result := Mongo_Result_Repl_Conf{}
			err = session.DB("admin").Run("replSetGetConfig", &result)
			if err != nil {
				color.Red("Failed to get buildInfo:  " + err.Error())
			} else {

				result.Config.Version = result.Config.Version + 1
				result.Config.Members[0].Host = serverSettings.WebConfig.DbConnection.Replication.Master
				result.Config.Members[0].Priority = 1
				result.Config.Members[0].Votes = 1

				for i, _ := range serverSettings.WebConfig.DbConnection.Replication.Slaves {
					slaveAddress := serverSettings.WebConfig.DbConnection.Replication.Slaves[i]
					if len(result.Config.Members) < i+2 {
						var slave Mongo_Replica_Member
						slave.Id = i + 1
						slave.Priority = i + 1
						slave.Votes = 1
						slave.Host = slaveAddress
						result.Config.Members = append(result.Config.Members, slave)
					} else {
						result.Config.Members[i+1].Host = slaveAddress
						result.Config.Members[i+1].Priority = i + 1
						result.Config.Members[i+1].Votes = 1
					}
				}

				result.Config.Settings.HeartbeatTimeoutSecs = 5

				err = session.DB("admin").Run(bson.D{{"replSetReconfig", result.Config}, {"force", true}}, nil)
				if err != nil {
					color.Red("Failed to replSetReconfig:  " + err.Error())
				}
				log.Println("Successfully initialized replica sets.")
			}
		}
	}

	var err error
	connectionString := serverSettings.WebConfig.DbConnection.ConnectionString
	if mongoDBOverride != "" {
		connectionString = mongoDBOverride
	}

	DBMutex.Lock()
	MongoSession, err = mgo.Dial(connectionString) // open an connection -> Dial function
	DBMutex.Unlock()

	if err != nil { //  if you have a
		color.Red("Failed to create or open mongo Database at " + connectionString + "\n\t" + err.Error())
		return err
	}

	if serverSettings.WebConfig.HasDbAuth && serverSettings.WebConfig.DbAuthConnection.AuthServer {

		DBMutex.Lock()
		MongoSessionAuth, err = mgo.Dial(connectionString) // open an connection -> Dial function
		DBMutex.Unlock()

		if err != nil { //  if you have a
			color.Red("Failed to create or open mongo Database at " + connectionString + "\n\t" + err.Error())
			return err
		}
	}

	return connectMongoDB()
}

func connectMongoDB() error {
	DBMutex.Lock()
	MongoSession.SetMode(mgo.Monotonic, true) // Optional. Switch the session to a monotonic behavior.
	MongoSession.SetSyncTimeout(2000 * time.Millisecond)

	dbName := serverSettings.WebConfig.DbConnection.Database
	if mongoDBNameOverride != "" {
		dbName = mongoDBNameOverride
	}

	MongoDB = MongoSession.DB(dbName)
	color.Green("Mongo Database Connected Successfully.")
	if serverSettings.WebConfig.HasDbAuth {
		MongoSessionAuth.SetMode(mgo.Monotonic, true) // Optional. Switch the session to a monotonic behavior.
		MongoSessionAuth.SetSyncTimeout(2000 * time.Millisecond)
		MongoDBAuth = MongoSession.DB(serverSettings.WebConfig.DbAuthConnection.Database)
		color.Green("Mongo Authentication Database Connected Successfully.")
	}
	DBMutex.Unlock()
	return nil
}
