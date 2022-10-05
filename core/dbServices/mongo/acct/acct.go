// Package acct - Internal Stub
package acct

import (
	"log"

	"github.com/DanielRenne/GoCore/core/dbServices"

	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	"time"
)

type Accounts struct {
}

func init() {
	go func() {

		for {
			dbServices.DBMutex.RLock()
			session := dbServices.MongoSession
			dbServices.DBMutex.RUnlock()
			if session != nil {
				log.Println("Building Indexes for MongoDB collection Accounts:")
				return
			}
			time.Sleep(200 * time.Millisecond)

		}
	}()

}
