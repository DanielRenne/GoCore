// Package mongo contains a helper function for ensuring mongo is up or down
package mongo

import (
	"log"
	"net"
	"runtime/debug"
	"time"

	"github.com/DanielRenne/GoCore/core/atomicTypes"
	"github.com/DanielRenne/GoCore/core/dbServices"
)

// IsMongoAlive is a global atomic bool that can be used to check if mongo is alive or not.
var IsMongoAlive atomicTypes.AtomicBool

// MongoError is a global error that shows the last error that occurred with mongo.
var MongoError error

func init() {
	IsMongoAlive.Set(true)
}

// InitializeDaemonChecker is a function that will check if mongo is up or down every 5 seconds and will invoke your callback function with the status of the connection after 10 seconds of tcp timeouts it will flag as false or every 5 seconds return true to you.
// You have to manage state if it has changed
// Something like
/*
mongo.InitializeDaemonChecker(mongoDaemonCallback)
func mongoDaemonCallback(status bool) {
	type response struct {
		Status bool `json:"status"`
	}
	var json response
	json.Status = status
	PublishWebSocketJSON("Mongo", json)
}
*/
func InitializeDaemonChecker(callback func(bool)) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("\n\nPanic Stack: " + string(debug.Stack()))
				InitializeDaemonChecker(callback)
			}
		}()
		for {
			dialer, _ := dbServices.GetMongoDialInfo()
			conn, err := net.Dial("tcp", dialer.Addrs[0])
			if err != nil {
				MongoError = err
				isMongoReallyDead := true
				// Try for 10 more seconds
				for i := 1; i <= 5; i++ {
					dialer, _ := dbServices.GetMongoDialInfo()
					conn, err := net.Dial("tcp", dialer.Addrs[0])
					if err != nil {
						time.Sleep(time.Second * 2)
						continue
					}
					conn.Close()
					isMongoReallyDead = false
					break
				}
				if isMongoReallyDead {
					IsMongoAlive.Set(false)
					callback(false)
				}
				continue
			}
			MongoError = nil
			IsMongoAlive.Set(true)
			callback(true)
			conn.Close()
			time.Sleep(time.Second * 15)
		}
	}()
}
