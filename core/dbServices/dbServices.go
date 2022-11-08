// Package dbServices provides a set of extensions for database utilities and ORM generation
package dbServices

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"time"

	"sync"

	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/asdine/storm"
	"github.com/fatih/color"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// DBMutex is a mutex for the database connection
var DBMutex *sync.RWMutex

// BoltDB is the global bolt database connection
var BoltDB *storm.DB

// MongoSession is the global mongo session
var MongoSession *mgo.Session

// MongoDB is the global mongo database connection
var MongoDB *mgo.Database

var mongoDBOverride string
var mongoDBNameOverride string
var hasBoltConnected bool

const (
	//Driver Types
	DATABASE_DRIVER_BOLTDB  = "boltDB"
	DATABASE_DRIVER_MONGODB = "mongoDB"
)

func init() {
	DBMutex = &sync.RWMutex{}
}

// OverrideMongoDBConnection allows you to override the connection string for the mongo database
func OverrideMongoDBConnection(connectionString string, dbName string) {
	mongoDBOverride = connectionString
	mongoDBNameOverride = dbName
}

// ReadMongoDB will read the database connection from memory
func ReadMongoDB() (mdb *mgo.Database) {
	DBMutex.RLock()
	mdb = MongoDB
	DBMutex.RUnlock()
	return mdb
}

// Initialize will initialize the database connection
func Initialize() {
	fmt.Println("core dbServices initialized.")
	if serverSettings.WebConfig.DbConnection.Driver == DATABASE_DRIVER_BOLTDB {
		openBolt()
	} else if serverSettings.WebConfig.DbConnection.Driver == DATABASE_DRIVER_MONGODB {
		go openMongo()
	}
}

func openBolt() {
	if !hasBoltConnected {

		myDBDir := serverSettings.APP_LOCATION + "/db/" + serverSettings.WebConfig.DbConnection.ConnectionString

		os.Mkdir(path.Dir(myDBDir), 0777)

		// (Create if not exist) open a database
		var err error
		DBMutex.Lock()
		BoltDB, err = storm.Open(myDBDir)
		DBMutex.Unlock()

		if err != nil {
			color.Red("Failed to create or open boltDB Database at " + myDBDir + ":\n\t" + err.Error())
			os.Exit(1)
		}
		color.Cyan("Successfully opened new bolt DB at " + myDBDir)
	}
	hasBoltConnected = true

}

// GetMongoDialInfo returns a mgo.DialInfo object based on the current serverSettings
func GetMongoDialInfo() (*mgo.DialInfo, error) {
	connectionString := serverSettings.WebConfig.DbConnection.ConnectionString
	if mongoDBOverride != "" {
		connectionString = mongoDBOverride
	}

	overrideConnectionString := os.Getenv("MGO_CONNECTION_STRING")
	if overrideConnectionString != "" {
		connectionString = overrideConnectionString
	}
	return mgo.ParseURL(connectionString)
}

func openMongo() {

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

	mgoTLSEnabled := os.Getenv("MGO_TLS_ENABLED")

	dialInfo, err := GetMongoDialInfo()
	if err != nil {
		log.Println(err)
	}

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		if serverSettings.WebConfig.DbConnection.EnableTLS || mgoTLSEnabled == "1" {
			tlsConfig := &tls.Config{}
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			if err != nil {
				log.Println(err)
			}
			return conn, err
		} else {
			conn, err := net.Dial("tcp", addr.String())
			if err != nil {
				log.Println(err)
			}
			return conn, err
		}

	}

	// This will block forever until a connection is established
	DBMutex.Lock()
	MongoSession, err = mgo.DialWithInfo(dialInfo) // open an connection -> Dial function
	DBMutex.Unlock()
	connectMongoDB()
}

func connectMongoDB() error {
	DBMutex.Lock()
	MongoSession.SetMode(mgo.Monotonic, true) // Optional. Switch the session to a monotonic behavior.
	MongoSession.SetSyncTimeout(2000 * time.Millisecond)

	dbName := serverSettings.WebConfig.DbConnection.Database
	if mongoDBNameOverride != "" {
		dbName = mongoDBNameOverride
	}

	overrideDBName := os.Getenv("MGO_DB_NAME")
	if overrideDBName != "" {
		dbName = overrideDBName
	}

	MongoDB = MongoSession.DB(dbName)
	color.Green("Mongo Database Connected Successfully.")
	DBMutex.Unlock()
	return nil
}
