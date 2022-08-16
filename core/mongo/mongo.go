package mongo

import (
	"log"
	"net"
	"runtime/debug"
	"time"

	"github.com/DanielRenne/GoCore/core/atomicTypes"
	"github.com/DanielRenne/GoCore/core/dbServices"
)

var IsMongoAlive atomicTypes.AtomicBool
var MongoError error

func init() {
	IsMongoAlive.Set(true)
}

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
